package infra

import (
	"carwise"
	"database/sql"
	"errors"
	"fmt"
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
			password_hash, 
			role, 
			status, 
			created_at, 
			updated_at, 
			last_login 
		FROM users 
		WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.ImageUrl,
		&user.CountryCode,
		&user.PhoneNumber,
		&user.Email,
		&user.PasswordHash,
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
			password_hash, 
			role, 
			status, 
			created_at, 
			updated_at, 
			last_login 
		FROM users 
		WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.ImageUrl,
		&user.CountryCode,
		&user.PhoneNumber,
		&user.Email,
		&user.PasswordHash,
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
			password_hash, 
			role, 
			status, 
			created_at, 
			updated_at, 
			last_login
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)`
	_, err := r.db.Exec(query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.ImageUrl,
		user.CountryCode,
		user.PhoneNumber,
		user.Email,
		user.PasswordHash,
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
			password_hash = $1, 
			updated_at = NOW() 
		WHERE email = $2`

	result, err := r.db.Exec(query, hashedPassword, email)
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
            updated_at = NOW()
        WHERE id = $6`

	_, err := r.db.Exec(query, user.FirstName, user.LastName, user.ImageUrl, user.CountryCode, user.PhoneNumber, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
