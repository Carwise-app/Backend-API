package infra

import (
	"carwise"
	"database/sql"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository() *NotificationRepository {
	db := ConnectDb()
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) CreateNotification(notification *carwise.Notification) error {
	query := `
	INSERT INTO notifications (title, message, status, data, created_by)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, notification.Title, notification.Message, notification.Status, notification.Data, notification.CreatedBy)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) GetNotificationsByUserId(userId string, limit, offset int) ([]carwise.Notification, error) {
	query := `
	SELECT id, title, message, status, data, created_at FROM notifications WHERE created_by = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []carwise.Notification
	for rows.Next() {
		var notification carwise.Notification
		err := rows.Scan(&notification.ID, &notification.Title, &notification.Message, &notification.Status, &notification.Data, &notification.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
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
