package infra

import (
	"carwise"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	database := ConnectDb()
	return &UserRepository{db: database}
}

func (r *UserRepository) GetByID(id string) (*carwise.User, error) {
	user := &carwise.User{}
	query := `
		SELECT 
			id, 
			google_id,
			first_name, 
			last_name, 
			image_url, 
			country_code, 
			phone_number, 
			email, 
			password, 
			role, 
			status, 
			created_at, 
			updated_at, 
			last_login,
			device_token,
			email_notify,
			push_notify
		FROM users 
		WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&user.Id,
		&user.GoogleId,
		&user.FirstName,
		&user.LastName,
		&user.ImageUrl,
		&user.CountryCode,
		&user.PhoneNumber,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
		&user.DeviceToken,
		&user.EmailNotify,
		&user.PushNotify,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to query user by ID: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*carwise.User, error) {
	user := &carwise.User{}
	query := `
		SELECT 
			id, 
			google_id,
			first_name, 
			last_name, 
			image_url, 
			country_code, 
			phone_number, 
			email, 
			password, 
			role, 
			status, 
			created_at, 
			updated_at, 
			last_login,
			device_token,
			email_notify,
			push_notify
		FROM users 
		WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(
		&user.Id,
		&user.GoogleId,
		&user.FirstName,
		&user.LastName,
		&user.ImageUrl,
		&user.CountryCode,
		&user.PhoneNumber,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
		&user.DeviceToken,
		&user.EmailNotify,
		&user.PushNotify,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query user by Email: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Create(user *carwise.User) error {
	query := `
		INSERT INTO users (
			id, 
			first_name, 
			last_name, 
			image_url, 
			country_code, 
			phone_number, 
			email, 
			password, 
			role, 
			status, 
			created_at, 
			updated_at, 
			last_login,
			google_id,
			device_token,
			email_notify,
			push_notify
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)`
	_, err := r.db.Exec(query,
		user.Id,
		user.FirstName,
		user.LastName,
		user.ImageUrl,
		user.CountryCode,
		user.PhoneNumber,
		user.Email,
		user.Password,
		user.Role,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
		user.LastLogin,
		user.GoogleId,
		user.DeviceToken,
		user.EmailNotify,
		user.PushNotify,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdatePassword(email, hashedPassword string) error {
	query := `
		UPDATE users 
		SET 
			password = $1, 
			updated_at = $2
		WHERE email = $3`

	result, err := r.db.Exec(query, hashedPassword, time.Now().Unix(), email)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no user found with the provided email")
	}

	return nil
}

func (r *UserRepository) Update(user *carwise.User) error {
	query := `
        UPDATE users 
        SET 
            first_name = $1,
            last_name = $2,
            image_url = $3,
			country_code = $4,
			phone_number = $5,
			device_token = $6,
			email_notify = $7,
			push_notify = $8,
            updated_at = $9
        WHERE id = $10`

	_, err := r.db.Exec(query, user.FirstName, user.LastName, user.ImageUrl,
		user.CountryCode, user.PhoneNumber, user.DeviceToken, user.EmailNotify,
		user.PushNotify, user.UpdatedAt, user.Id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetAllUsers() ([]carwise.User, error) {
	query := `
		SELECT id, email, device_token, email_notify, push_notify
		FROM users
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	users := []carwise.User{}
	for rows.Next() {
		var user carwise.User
		err := rows.Scan(&user.Id, &user.Email, &user.DeviceToken, &user.EmailNotify, &user.PushNotify)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetUsers(page, limit int) ([]carwise.User, error) {
	query := `
		SELECT 
			id,
			google_id,
			first_name,
			last_name,
			image_url,
			country_code,
			phone_number,
			email,
			role,
			status,
			created_at,
			updated_at,
			last_login
		FROM users
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, (page-1)*limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	users := []carwise.User{}
	for rows.Next() {
		var user carwise.User
		err := rows.Scan(
			&user.Id, &user.GoogleId, &user.FirstName, &user.LastName, &user.ImageUrl,
			&user.CountryCode, &user.PhoneNumber, &user.Email, &user.Role,
			&user.Status, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) CountUsers() (int, error) {
	query := `
		SELECT COUNT(*) FROM users
	`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

func (r *UserRepository) DeleteUser(id string) error {
	queryListingDelete := `
		DELETE FROM listings WHERE created_by = $1
	`
	_, err := r.db.Exec(queryListingDelete, id)
	if err != nil {
		return fmt.Errorf("failed to delete listings: %w", err)
	}

	queryImageDelete := `
		DELETE FROM images WHERE created_by = $1
	`
	_, err = r.db.Exec(queryImageDelete, id)
	if err != nil {
		return fmt.Errorf("failed to delete images: %w", err)
	}

	queryMessageDelete := `
		DELETE FROM messages WHERE sender_id = $1
		OR receiver_id = $1
	`
	_, err = r.db.Exec(queryMessageDelete, id)
	if err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}

	queryFavoriteDelete := `
		DELETE FROM favorites WHERE user_id = $1
	`
	_, err = r.db.Exec(queryFavoriteDelete, id)
	if err != nil {
		return fmt.Errorf("failed to delete favorites: %w", err)
	}

	queryPredictDelete := `
		DELETE FROM predicts WHERE created_by = $1
	`
	_, err = r.db.Exec(queryPredictDelete, id)
	if err != nil {
		return fmt.Errorf("failed to delete predicts: %w", err)
	}

	queryNotificationDelete := `
		DELETE FROM notifications WHERE created_by = $1
	`
	_, err = r.db.Exec(queryNotificationDelete, id)
	if err != nil {
		return fmt.Errorf("failed to delete notifications: %w", err)
	}

	queryUserDelete := `
		DELETE FROM users WHERE id = $1
	`
	_, err = r.db.Exec(queryUserDelete, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *UserRepository) UpdateUserRole(userId string, role int) error {
	query := `
		UPDATE users 
		SET 
			role = $1,
			updated_at = $2
		WHERE id = $3`

	result, err := r.db.Exec(query, role, time.Now().Unix(), userId)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no user found with the provided ID")
	}

	return nil
}
