package infra

import (
	"carwise"
	"database/sql"
	"encoding/json"
	"fmt"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository() *NotificationRepository {
	db := ConnectDb()
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) CreateNotification(notification *carwise.Notification) error {
	dataJSON, err := json.Marshal(notification.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal notification data: %w", err)
	}

	query := `
	INSERT INTO notifications (title, message, status, data, created_by)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err = r.db.Exec(query, notification.Title, notification.Message, notification.Status, dataJSON, notification.CreatedBy)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) GetNotificationsByUserId(userId string, limit, offset int) ([]carwise.Notification, error) {
	query := `
	SELECT id, title, message, status, data, read, created_by, created_at 
	FROM notifications 
	WHERE created_by = $1 
	ORDER BY created_at DESC 
	LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []carwise.Notification
	for rows.Next() {
		var notification carwise.Notification
		var dataJSON string

		err := rows.Scan(
			&notification.ID,
			&notification.Title,
			&notification.Message,
			&notification.Status,
			&dataJSON, // Scan JSON as string first
			&notification.Read,
			&notification.CreatedBy,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSON string into map
		if dataJSON != "" {
			if err := json.Unmarshal([]byte(dataJSON), &notification.Data); err != nil {
				return nil, fmt.Errorf("failed to unmarshal data JSON: %w", err)
			}
		}

		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *NotificationRepository) ReadNotification(id string, userId string) error {
	query := `
	UPDATE notifications SET read = true WHERE id = $1 AND created_by = $2
	`

	_, err := r.db.Exec(query, id, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) DeleteNotification(id string, userId string) error {
	query := `
	DELETE FROM notifications WHERE id = $1 AND created_by = $2
	`

	_, err := r.db.Exec(query, id, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) GetUnreadNotificationsCount(userId string) (int, error) {
	query := `
	SELECT COUNT(*) FROM notifications WHERE created_by = $1 AND read = false
	`

	var count int
	err := r.db.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *NotificationRepository) GetTotalNotificationsCount(userId string) (int, error) {
	query := `
	SELECT COUNT(*) FROM notifications WHERE created_by = $1
	`

	var count int
	err := r.db.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
