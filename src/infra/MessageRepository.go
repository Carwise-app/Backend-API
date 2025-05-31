package infra

import (
	"carwise"
	"database/sql"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{db: ConnectDb()}
}

func (r *MessageRepository) SaveMessage(message *carwise.Message) error {
	query := `
		INSERT INTO messages (id, sender_id, receiver_id, message, read, created_at, listing_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, message.Id, message.SenderId, message.ReceiverId, message.Message, message.Read, message.CreatedAt, message.ListingId)
	return err
}

func (r *MessageRepository) GetMessagesByUserId(userId, otherUserId, listingId string, limit, page int) ([]carwise.Message, error) {
	query := `
		SELECT id, sender_id, receiver_id, message, read, created_at, listing_id
		FROM messages 
		WHERE (sender_id = $1 OR receiver_id = $1) 
		AND listing_id = $2
		AND (sender_id = $3 OR receiver_id = $3)
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`

	offset := limit * (page - 1)
	rows, err := r.db.Query(query, userId, listingId, otherUserId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []carwise.Message
	for rows.Next() {
		var message carwise.Message
		err := rows.Scan(&message.Id, &message.SenderId, &message.ReceiverId, &message.Message, &message.Read, &message.CreatedAt, &message.ListingId)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (r *MessageRepository) CountMessagesByUserId(userId, otherUserId, listingId string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM messages 
		WHERE (sender_id = $1 OR receiver_id = $1) 
		AND listing_id = $2
		AND (sender_id = $3 OR receiver_id = $3)
	`

	var count int
	err := r.db.QueryRow(query, userId, listingId, otherUserId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *MessageRepository) GetChats(userId string, limit, page int) ([]carwise.Chat, error) {
	query := `
		WITH chat_messages AS (
			SELECT 
				CASE 
					WHEN sender_id = $1 THEN receiver_id 
					ELSE sender_id 
				END AS other_user_id,
				created_at,
				listing_id,
				ROW_NUMBER() OVER (
					PARTITION BY 
						CASE 
							WHEN sender_id = $1 THEN receiver_id 
							ELSE sender_id 
						END,
						listing_id
					ORDER BY created_at DESC
				) as rn
			FROM messages 
			WHERE sender_id = $1 OR receiver_id = $1
		)
		SELECT other_user_id, created_at as last_message_time, listing_id
		FROM chat_messages
		WHERE rn = 1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userId, limit, limit*(page-1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []carwise.Chat
	for rows.Next() {
		var chat carwise.Chat
		err := rows.Scan(&chat.OtherUserId, &chat.LastMessageTime, &chat.ListingId)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (r *MessageRepository) CountChats(userId string) (int, error) {
	query := `
		SELECT COUNT(DISTINCT (other_user_id, listing_id))
		FROM (
			SELECT 
				CASE 
					WHEN sender_id = $1 THEN receiver_id 
					ELSE sender_id 
				END AS other_user_id,
				listing_id
			FROM messages 
			WHERE sender_id = $1 OR receiver_id = $1
		) AS chat_pairs
	`
	var count int
	err := r.db.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *MessageRepository) ReadMessage(messageId string) error {
	query := `
		UPDATE messages
		SET read = true
		WHERE id = $1
	`
	_, err := r.db.Exec(query, messageId)
	return err
}
