package infra

import (
	"carwise"
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository() *MessageRepository {
	db := ConnectDb()
	return &MessageRepository{db: db}
}

func (r *MessageRepository) SaveMessage(message *carwise.Message) error {
	messageId := generateUUID()

	query := `
		INSERT INTO messages (id, sender_id, receiver_id, message, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, messageId, message.SenderId, message.ReceiverId, message.Message, message.CreatedAt)
	if err != nil {
		log.Println("Error saving message:", err)
		return err
	}

	return nil
}

func (r *MessageRepository) GetMessagesBetween(senderId, receiverId string, limit, offset int) ([]carwise.Message, error) {
	query := `
		SELECT id, sender_id, receiver_id, message, created_at
		FROM messages
		WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
		ORDER BY created_at ASC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, senderId, receiverId, limit, offset)
	if err != nil {
		log.Println("Error retrieving messages:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []carwise.Message
	for rows.Next() {
		var msg carwise.Message
		err := rows.Scan(&msg.Id, &msg.SenderId, &msg.ReceiverId, &msg.Message, &msg.CreatedAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, err
	}

	return messages, nil
}

func generateUUID() string {
	return uuid.New().String()
}
