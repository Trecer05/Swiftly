package chat

import (
	"database/sql"

	models "github.com/Trecer05/Swiftly/internal/model/chat"
	errors "github.com/Trecer05/Swiftly/internal/errors/chat"
)

func (manager *Manager) GetUserRooms(userId, limit, offset int) (models.ChatRooms, error) {
	var chatRooms models.ChatRooms

	rows, err := manager.Conn.Query(`WITH user_chats AS (
			SELECT 
				c.id AS chat_id,
				c.name AS chat_name
			FROM 
				chats c
			JOIN 
				chat_users cu ON c.id = cu.chat_id
			WHERE 
				cu.user_id = $1
			ORDER BY 
				c.id
			LIMIT $2 OFFSET $3
		)
		SELECT 
			uc.chat_id,
			uc.chat_name,
			cm.text AS last_message_text,
			u.name AS sender_name,
			cm.sent_at AS last_message_time
		FROM 
			user_chats uc
		LEFT JOIN LATERAL (
			SELECT 
				cm.text, 
				cm.sent_at, 
				cm.user_id
			FROM 
				chat_messages cm
			WHERE 
				cm.chat_id = uc.chat_id
			ORDER BY 
				cm.sent_at DESC
			LIMIT 1
		) cm ON true
		LEFT JOIN 
			users u ON cm.user_id = u.id
		ORDER BY 
			COALESCE(cm.sent_at, (SELECT MIN(created_at) FROM chats)) DESC`, userId, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows { return models.ChatRooms{}, errors.ErrNoChats }
		return models.ChatRooms{}, err
	}

	for rows.Next() {
		var id int
		var chatName string
		var lastMessageText sql.NullString
		var senderName sql.NullString
		var lastMessageTime sql.NullTime

		err := rows.Scan(
			&id,
			&chatName,
			&lastMessageText,
			&senderName,
			&lastMessageTime,
		)
		if err != nil { return models.ChatRooms{}, err }

		chatRooms.Rooms = append(chatRooms.Rooms, &models.ChatRoom{
			ID: id,
			Name: chatName,
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