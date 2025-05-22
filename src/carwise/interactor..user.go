package carwise

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (i *Interactor) CreateUser(request UserCreateRequest) (*User, []string) {
	existingUser, _ := i.services.UserRepo.GetByEmail(request.Email)
	if existingUser != nil {
		return nil, []string{"Email is already in use."}
	}

	hashedPassword, err := hashPassword(request.Password)
	if err != nil {
		return nil, []string{"Failed to hash password."}
	}
	user := &User{
		Id:          uuid.New().String(),
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		CountryCode: request.CountryCode,
		PhoneNumber: request.PhoneNumber,
		Email:       request.Email,
		Password:    hashedPassword,
		Role:        1,
		Status:      1,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		LastLogin:   time.Now().Unix(),
	}

	err = i.services.UserRepo.Create(user)
	if err != nil {
		return nil, []string{"Failed to create user: " + err.Error()}
	}
	return user, nil
}

func (i *Interactor) LoginUser(request UserLoginRequest) (*User, []string) {
	user, err := i.services.UserRepo.GetByEmail(request.Email)
	if err != nil {
		return nil, []string{err.Error()}
	}

	if !comparePasswords(user.Password, request.Password) {
		return nil, []string{"invalid credentials"}
	}

	return user, nil
}

func (i *Interactor) IsTokenBlackListed(token string) (bool, []string) {
	isBlacklisted, err := i.services.TokenRepo.IsTokenBlackListed(token)
	if err != nil {
		return false, []string{"Failed to check token blacklist: " + err.Error()}
	}

	return isBlacklisted, nil
}

func (i *Interactor) AddTokenBlackList(token string) []string {
	err := i.services.TokenRepo.AddTokenBlackList(token)
	if err != nil {
		return []string{"Failed to add token to blacklist: " + err.Error()}
	}

	return nil
}

func (i *Interactor) ResetPasswordRequest(request ResetPasswordRequest) []string {
	existingUser, err := i.services.UserRepo.GetByEmail(request.Email)
	if err != nil {
		log.Printf("Error fetching user by email: %v\n", err)
		return []string{"An unexpected error occurred. Please try again later."}
	}

	if existingUser == nil {
		return []string{"No account found with this email."}
	}
	token, err := generateToken(40)
	if err != nil {
		log.Printf("Error generate password reset token: %v\n", err)
		return []string{"An unexpected error occurred. Please try again later."}
	}
	err = i.services.PasswordResetRepo.SaveResetCode(request.Email, token, 5*24*time.Hour)
	if err != nil {
		fmt.Printf("Failed to save reset code: %v\n", err)
	}

	resetLink := fmt.Sprintf("http://localhost:3000/reset-password?token=%s&email=%s", token, request.Email)
	emailBody := fmt.Sprintf(`From: Carwise <app.carwise@gmail.com>
Subject: Password Reset Request
Dear User,
We received a request to reset the password associated with your account. If you made this request, please click the link below to reset your password:

%s

This link will expire in 5 days. If you did not request a password reset, you can safely ignore this email.

Best regards,
Carwise Team`, resetLink)

	err = i.services.MailGW.Send(request.Email, []byte(emailBody))
	if err != nil {
		log.Printf("Error send password reset email: %v\n", err)
		return []string{"An unexpected error occurred. Please try again later."}
	}

	return nil
}

func (i *Interactor) ChangePassword(request ChangePasswordRequest, token, email string) []string {
	verify, err := i.services.PasswordResetRepo.VerifyResetCode(email, token)
	if err != nil {
		log.Printf("Error verifying reset token: %v\n", err)
		return []string{"An unexpected error occurred. Please try again later."}
	}
	if !verify {
		return []string{"Invalid or expired password reset token."}
	}

	hashedPassword, err := hashPassword(request.Password)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err)
		return []string{"An unexpected error occurred. Please try again later."}
	}

	err = i.services.UserRepo.UpdatePassword(email, hashedPassword)
	if err != nil {
		log.Printf("Error updating password: %v\n", err)
		return []string{"An unexpected error occurred. Please try again later."}
	}

	err = i.services.PasswordResetRepo.DeleteResetCode(email)
	if err != nil {
		log.Printf("Error deleting reset token: %v\n", err)
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePasswords(passwordHash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) == nil
}

func generateToken(size int) (string, error) {
	if size <= 0 {
		return "", fmt.Errorf("invalid size for token generation")
	}

	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	return hex.EncodeToString(buf), nil

}

func generateSecureListingNumber(length int) (string, error) {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randPart := make([]rune, length)

	for i := range randPart {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		randPart[i] = letters[idx.Int64()]
	}

	return string(randPart), nil
}
