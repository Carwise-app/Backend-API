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
			last_login 
		FROM users 
		WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&user.Id,
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
			last_login 
		FROM users 
		WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(
		&user.Id,
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
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
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
			last_login
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
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
            updated_at = $6
        WHERE id = $7`

	_, err := r.db.Exec(query, user.FirstName, user.LastName, user.ImageUrl, user.CountryCode, user.PhoneNumber, user.UpdatedAt, user.Id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
