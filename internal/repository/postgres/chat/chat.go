package chat

import (
	"database/sql"
	"fmt"
	"strconv"

	errors "github.com/Trecer05/Swiftly/internal/errors/chat"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
)

func (manager *Manager) GetUserRoomsForStatus(userId int) ([]models.UserRoom, error) {
    var rooms []models.UserRoom
    
    rows, err := manager.Conn.Query(`
        SELECT c.id, 'private' as type
        FROM chats c
        JOIN chat_users cu ON c.id = cu.chat_id
        WHERE cu.user_id = $1`, userId)
    if err != nil {
        return rooms, err
    }
    defer rows.Close()
    
    for rows.Next() {
        var room models.UserRoom
        var roomType string
        if err := rows.Scan(&room.ID, &roomType); err != nil {
            continue
        }
        room.Type = models.ChatType(roomType)
        rooms = append(rooms, room)
    }
    
    groupRows, err := manager.Conn.Query(`
        SELECT g.id, 'group' as type
        FROM groups g
        JOIN group_users gu ON g.id = gu.group_id
        WHERE gu.user_id = $1`, userId)
    if err != nil {
        return rooms, err
    }
    defer groupRows.Close()
    
    for groupRows.Next() {
        var room models.UserRoom
        var roomType string
        if err := groupRows.Scan(&room.ID, &roomType); err != nil {
            continue
        }
        room.Type = models.ChatType(roomType)
        rooms = append(rooms, room)
    }
    
    return rooms, nil
}

func (manager *Manager) GetUserRooms(userId, limit, offset int) (models.ChatRooms, error) {
	var chatRooms models.ChatRooms
	var limitStr *string
	if limit == 0 {
		limitStr = nil
	} else {
		strLimit := strconv.Itoa(limit)
		limitStr = &strLimit
	}

	rows, err := manager.Conn.Query(`WITH user_private_chats AS (
				SELECT 
					c.id AS chat_id,
					c.name AS chat_name,
					'chat' AS chat_type
				FROM 
					chats c
				JOIN 
					chat_users cu ON c.id = cu.chat_id
				WHERE 
					cu.user_id = $1
			),
			last_private_messages AS (
				SELECT DISTINCT ON (cm.chat_id)
					cm.chat_id,
					cm.text,
					cm.sent_at,
					cm.user_id,
					cm.read
				FROM chat_messages cm
				ORDER BY cm.chat_id, cm.sent_at DESC
			),
			user_groups AS (
				SELECT 
					g.id AS chat_id,
					g.name AS chat_name,
					'group' AS chat_type
				FROM 
					groups g
				JOIN 
					group_users gu ON g.id = gu.group_id
				WHERE 
					gu.user_id = $1
			),
			last_group_messages AS (
				SELECT DISTINCT ON (gm.group_id)
					gm.group_id,
					gm.text,
					gm.sent_at,
					gm.user_id,
					gm.read
				FROM group_messages gm
				ORDER BY gm.group_id, gm.sent_at DESC
			),
			combined_rooms AS (
				SELECT 
					pc.chat_id,
					pc.chat_name,
					pc.chat_type,
					lpm.text AS last_message_text,
					lpm.read,
					u.name AS sender_name,
					lpm.sent_at AS last_message_time
				FROM 
					user_private_chats pc
				LEFT JOIN last_private_messages lpm ON lpm.chat_id = pc.chat_id
				LEFT JOIN users u ON lpm.user_id = u.id

				UNION ALL

				SELECT 
					ug.chat_id,
					ug.chat_name,
					ug.chat_type,
					lgm.text AS last_message_text,
					lgm.read,
					u.name AS sender_name,
					lgm.sent_at AS last_message_time
				FROM 
					user_groups ug
				LEFT JOIN last_group_messages lgm ON lgm.group_id = ug.chat_id
				LEFT JOIN users u ON lgm.user_id = u.id
			)

			SELECT *
			FROM combined_rooms
			ORDER BY COALESCE(last_message_time, NOW()) DESC
			LIMIT $2 OFFSET $3
			`, userId, limitStr, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.ChatRooms{}, errors.ErrNoRooms
		}
		return models.ChatRooms{}, err
	}

	for rows.Next() {
		var id int
		var chatName string
		var lastMessageText sql.NullString
		var senderName sql.NullString
		var lastMessageTime sql.NullTime
		var chatType string
		var read bool

		err := rows.Scan(
			&id,
			&chatName,
			&chatType,
			&lastMessageText,
			&read,
			&senderName,
			&lastMessageTime,
		)
		if err != nil {
			return models.ChatRooms{}, err
		}

		var typeModel models.ChatType
		if chatType == "chat" {
			typeModel = models.TypePrivate
		} else {
			if chatType == "group" {
				typeModel = models.TypeGroup
			} else {
				return models.ChatRooms{}, errors.ErrUnknownChatType
			}
		}

		chatRooms.Rooms = append(chatRooms.Rooms, &models.ChatRoom{
			ID:   id,
			Name: chatName,
			Type: typeModel,
			LastMessage: &models.Message{
				Text: &lastMessageText.String,
				Time: lastMessageTime.Time,
				Author: models.Client{
					Name: senderName.String,
				},
				Read: &read,
			},
		})
	}

	return chatRooms, nil
}

func (manager *Manager) GetChatMessages(chatId, limit, offset int) ([]models.Message, error) {
	var messages []models.Message

	rows, err := manager.Conn.Query(`SELECT 
				cm.id AS message_id,
				cm.text AS message_text,
				cm.sent_at AS sent_time,
				cm.read AS read,
				cm.edited AS edited,
				u.id AS user_id,
				u.name AS user_name
			FROM 
				chat_messages cm
			JOIN 
				users u ON cm.user_id = u.id
			WHERE 
				cm.chat_id = $1
			ORDER BY 
				cm.sent_at DESC
			LIMIT $2 OFFSET $3`, chatId, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return messages, errors.ErrNoMessages
		} else {
			return messages, err
		}
	}

	for rows.Next() {
		var message models.Message

		err := rows.Scan(
			&message.ID,
			&message.Text,
			&message.Time,
			&message.Read,
			&message.Edited,
			&message.Author.ID,
			&message.Author.Name,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return messages, errors.ErrNoMessages
			} else {
				return messages, err
			}
		}

		fileURLs, err := manager.getChatMessageFiles(message.ID)
		if err != nil {
			fmt.Printf("Ошибка загрузки файлов для сообщения %d: %v\n", message.ID, err)
		} else {
			message.FileUrls = fileURLs
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (manager *Manager) GetGroupMessages(groupId, limit, offset int) ([]models.Message, error) {
	var messages []models.Message

	rows, err := manager.Conn.Query(`SELECT 
				gm.id AS message_id,
				gm.text AS message_text,
				gm.sent_at AS sent_time,
				gm.read AS read,
				gm.edited AS edited,
				u.id AS user_id,
				u.name AS user_name
			FROM 
				group_messages gm
			JOIN 
				users u ON gm.user_id = u.id
			WHERE 
				gm.chat_id = $1
			ORDER BY 
				gm.sent_at DESC
			LIMIT $2 OFFSET $3`, groupId, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return messages, errors.ErrNoMessages
		} else {
			return messages, err
		}
	}

	for rows.Next() {
		var message models.Message

		err := rows.Scan(
			&message.ID,
			&message.Text,
			&message.Time,
			&message.Read,
			&message.Edited,
			&message.Author.ID,
			&message.Author.Name,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return messages, errors.ErrNoMessages
			} else {
				return messages, err
			}
		}

		fileURLs, err := manager.getGroupMessageFiles(message.ID)
		if err != nil {
			fmt.Printf("Ошибка загрузки файлов для сообщения %d: %v\n", message.ID, err)
		} else {
			message.FileUrls = fileURLs
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (manager *Manager) SaveMessage(message models.Message, chatType models.DBType) error {
	tx, err := manager.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var messageID int
	var insertQuery string

	if chatType == models.DBChat {
		insertQuery = `INSERT INTO chat_messages (chat_id, user_id, text) VALUES ($1, $2, $3) RETURNING id`
	} else {
		insertQuery = `INSERT INTO group_messages (group_id, user_id, text) VALUES ($1, $2, $3) RETURNING id`
	}

	err = tx.QueryRow(insertQuery, message.ChatID, message.Author.ID, message.Text).Scan(&messageID)
	if err != nil {
		return err
	}

	if len(message.FileUrls) > 0 {
		for _, fileURL := range message.FileUrls {
			var fileID int
			var fileInsertQuery string

			if chatType == models.DBChat {
				fileInsertQuery = `INSERT INTO chat_file_urls (chat_id, file_url) VALUES ($1, $2) RETURNING id`
			} else {
				fileInsertQuery = `INSERT INTO group_file_urls (group_id, file_url) VALUES ($1, $2) RETURNING id`
			}

			err = tx.QueryRow(fileInsertQuery, message.ChatID, fileURL).Scan(&fileID)
			if err != nil {
				return err
			}

			var linkInsertQuery string
			if chatType == models.DBChat {
				linkInsertQuery = `INSERT INTO chat_messages_file_urls (chat_message_id, chat_file_id) VALUES ($1, $2)`
			} else {
				linkInsertQuery = `INSERT INTO group_messages_file_urls (group_message_id, group_file_id) VALUES ($1, $2)`
			}

			_, err = tx.Exec(linkInsertQuery, messageID, fileID)
			if err != nil {
				return err
			}
		}
	}

	if message.FileURL != nil && *message.FileURL != "" {
		var fileID int
		var fileInsertQuery string

		if chatType == models.DBChat {
			fileInsertQuery = `INSERT INTO chat_file_urls (chat_id, file_url) VALUES ($1, $2) RETURNING id`
		} else {
			fileInsertQuery = `INSERT INTO group_file_urls (group_id, file_url) VALUES ($1, $2) RETURNING id`
		}

		err = tx.QueryRow(fileInsertQuery, message.ChatID, *message.FileURL).Scan(&fileID)
		if err != nil {
			return err
		}

		var linkInsertQuery string
		if chatType == models.DBChat {
			linkInsertQuery = `INSERT INTO chat_messages_file_urls (chat_message_id, chat_file_id) VALUES ($1, $2)`
		} else {
			linkInsertQuery = `INSERT INTO group_messages_file_urls (group_message_id, group_file_id) VALUES ($1, $2)`
		}

		_, err = tx.Exec(linkInsertQuery, messageID, fileID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (manager *Manager) DeleteMessage(chatId, messageId int, chatType models.DBType) error {
	_, err := manager.Conn.Exec(fmt.Sprintf(`
		DELETE FROM %s_messages
		WHERE id = $1 and %s_id = $2
	`, chatType, chatType), messageId, chatId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrNoMessages
		}
		return err
	}

	return nil
}

func (manager *Manager) UpdateMessage(chatId, messageId int, text *string, chatType models.DBType) error {
	_, err := manager.Conn.Exec(fmt.Sprintf(`
		UPDATE %s_messages
		SET text = $1, edited = true
		WHERE id = $2 and %s_id = $3
		`, chatType, chatType), text, messageId, chatId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrNoMessages
		}
		return err
	}

	return nil
}

func (manager *Manager) ValidateOwnerId(groupId int, userId int) (bool, error) {
	var id int

	err := manager.Conn.QueryRow("SELECT owner_id FROM groups WHERE id = $1", groupId).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.ErrNoGroupFound
		}
		return false, err
	}

	return id == userId, err
}

func (manager *Manager) UpdateChatMessageStatus(messageId int, read bool) error {
	if _, err := manager.Conn.Exec(`
		UPDATE chat_messages
		SET read = $1
		WHERE id = $2`, read, messageId); err != nil {
		return err
	}

	return nil
}

func (manager *Manager) UpdateGroupMessageStatus(messageId int, read bool) error {
	if _, err := manager.Conn.Exec(`
		UPDATE group_messages
		SET read = $1
		WHERE id = $2`, read, messageId); err != nil {
		return err
	}

	return nil
}

func (manager *Manager) getChatMessageFiles(messageID int) ([]string, error) {
	var fileURLs []string

	rows, err := manager.Conn.Query(`
		SELECT cf.file_url 
		FROM chat_file_urls cf
		JOIN chat_messages_file_urls cmf ON cf.id = cmf.chat_file_id
		WHERE cmf.chat_message_id = $1
		ORDER BY cf.id`, messageID)

	if err != nil {
		return fileURLs, err
	}
	defer rows.Close()

	for rows.Next() {
		var fileURL string
		if err := rows.Scan(&fileURL); err != nil {
			return fileURLs, err
		}
		fileURLs = append(fileURLs, fileURL)
	}

	return fileURLs, nil
}

func (manager *Manager) getGroupMessageFiles(messageID int) ([]string, error) {
	var fileURLs []string

	rows, err := manager.Conn.Query(`
		SELECT gf.file_url 
		FROM group_file_urls gf
		JOIN group_messages_file_urls gmf ON gf.id = gmf.group_file_id
		WHERE gmf.group_message_id = $1
		ORDER BY gf.id`, messageID)

	if err != nil {
		return fileURLs, err
	}
	defer rows.Close()

	for rows.Next() {
		var fileURL string
		if err := rows.Scan(&fileURL); err != nil {
			return fileURLs, err
		}
		fileURLs = append(fileURLs, fileURL)
	}

	return fileURLs, nil
}

func (manager *Manager) UpdateGroup(id int, req models.GroupEdit) error {
	switch {
	case req.Name != nil && req.Description != nil:
		if _, err := manager.Conn.Exec(`
			UPDATE groups
			SET name = $1, description = $2
			WHERE id = $3`, req.Name, req.Description, id); err != nil {
				if err == sql.ErrNoRows {
					return errors.ErrNoGroupFound
				}

				return err
			}
		return nil
	case req.Name != nil:
		if _, err := manager.Conn.Exec(`
			UPDATE groups
			SET name = $1
			WHERE id = $2`, req.Name, id); err != nil {
				if err == sql.ErrNoRows {
					return errors.ErrNoGroupFound
				}

				return err
			}
	case req.Description != nil:
		if _, err := manager.Conn.Exec(`
			UPDATE groups
			SET description = $1
			WHERE id = $2`, req.Description, id); err != nil {
				if err == sql.ErrNoRows {
					return errors.ErrNoGroupFound
				}

				return err
			}
	}
	
	return nil
}
