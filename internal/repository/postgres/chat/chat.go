package chat

import (
	"database/sql"
	"fmt"

	errors "github.com/Trecer05/Swiftly/internal/errors/chat"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
)

func (manager *Manager) GetUserRooms(userId int) (models.ChatRooms, error) {
	var chatRooms models.ChatRooms

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
					cm.user_id
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
					gm.user_id
				FROM group_messages gm
				ORDER BY gm.group_id, gm.sent_at DESC
			),
			combined_rooms AS (
				SELECT 
					pc.chat_id,
					pc.chat_name,
					pc.chat_type,
					lpm.text AS last_message_text,
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
			`, userId)
	if err != nil {
		if err == sql.ErrNoRows { return models.ChatRooms{}, errors.ErrNoRooms }
		return models.ChatRooms{}, err
	}

	for rows.Next() {
		var id int
		var chatName string
		var lastMessageText sql.NullString
		var senderName sql.NullString
		var lastMessageTime sql.NullTime
		var chatType string

		err := rows.Scan(
			&id,
			&chatName,
			&chatType,
			&lastMessageText,
			&senderName,
			&lastMessageTime,
		)
		if err != nil { return models.ChatRooms{}, err }

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
			ID: id,
			Name: chatName,
			Type: typeModel,
			LastMessage: &models.Message{
				Text: lastMessageText.String,
				Time: lastMessageTime.Time,
				Author: models.Client{
					Name: senderName.String,
				},
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
		if err == sql.ErrNoRows { return messages, errors.ErrNoMessages } else { return messages, err }
	}

	for rows.Next() {
		var message models.Message

		err := rows.Scan(
			&message.ID,
			&message.Text,
			&message.Time,
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
		if err == sql.ErrNoRows { return messages, errors.ErrNoMessages } else { return messages, err }
	}

	for rows.Next() {
		var message models.Message

		err := rows.Scan(
			&message.ID,
			&message.Text,
			&message.Time,
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

		messages = append(messages, message)
	}

	return messages, nil
}

func (manager *Manager) SaveMessage(message models.Message, chatType models.DBType) error {
	_, err := manager.Conn.Exec(fmt.Sprintf(`INSERT INTO %s_messages (%s_id, user_id, text) VALUES ($1, $2, $3)`, chatType, chatType), message.ChatID, message.Author.ID, message.Text)
	return err
}