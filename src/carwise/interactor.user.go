package carwise

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
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
	go func() {
		i.services.MailGW.SendEmail(request.Email, "Hoşgeldiniz!", "welcome.html", map[string]interface{}{
			"full_name": request.FirstName + " " + request.LastName,
		})
		i.CreateNotification(
			&Notification{
				ID:        uuid.New().String(),
				Title:     "🎉 Carwise'e Hoş Geldiniz!",
				Status:    SystemMessage,
				Message:   "Hoşgeldin, " + user.FirstName + "! 2. el araç alım-satımı ve araç fiyat tahmini artık çok daha kolay. Akıllı sistemimizle aracınızın gerçek değerini öğrenin, güvenle alım-satım yapın.",
				Read:      false,
				CreatedBy: user.Id,
				CreatedAt: time.Now().Unix(),
			},
		)
	}()

	return user, nil
}

func (i *Interactor) LoginUser(request UserLoginRequest) (*User, []string) {
	user, err := i.services.UserRepo.GetByEmail(request.Email)
	if err != nil {
		return nil, []string{err.Error()}
	}

	if user == nil {
		return nil, []string{"invalid credentials"}
	}

	if !comparePasswords(user.Password, request.Password) {
		return nil, []string{"invalid credentials"}
	}

	if user.GoogleId != "" {
		return nil, []string{"Please login with Google"}
	}

	if user.Status == 2 {
		return nil, []string{"Your account is not active. Please contact support."}
	}

	return user, nil
}

func (i *Interactor) IsTokenBlackListed(token string) (bool, []string) {
	isBlacklisted, err := i.services.RedisRepo.IsTokenBlackListed(token)
	if err != nil {
		return false, []string{"Failed to check token blacklist: " + err.Error()}
	}

	return isBlacklisted, nil
}

func (i *Interactor) AddTokenBlackList(token string) []string {
	err := i.services.RedisRepo.AddTokenBlackList(token)
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

	resetLink := fmt.Sprintf("https://carwise.yusuftalhaklc.com/kokpit/reset-password?token=%s&email=%s", token, request.Email)
	go i.services.MailGW.SendEmail(request.Email, "Password Reset Request", "reset_password.html", map[string]interface{}{
		"title":   "Şifre Sıfırlama İsteği",
		"message": "Hesabınızla ilişkili şifreyi sıfırlama talebi aldık. Eğer bu talebi siz yaptıysanız, şifrenizi sıfırlamak için aşağıdaki bağlantıya tıklayın:",
		"link":    resetLink,
	})

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

func (i *Interactor) GetUserById(request GetUserByIdRequest) (*UserInfo, []string) {
	user, err := i.services.UserRepo.GetByID(request.Id)
	if err != nil {
		return nil, []string{"An unexpected error occurred. Please try again later."}
	}
	return &UserInfo{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		CountryCode: user.CountryCode,
		PhoneNumber: user.PhoneNumber,
	}, nil
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

func (i *Interactor) UpdateUserRole(request UpdateUserRoleRequest) []string {
	// Check if the requesting user is an admin (role = 2)
	if request.AdminRole != 2 {
		return []string{"Only admin users can update user roles"}
	}

	// Validate role value
	if request.Role < 1 || request.Role > 2 {
		return []string{"Invalid role value. Role must be 1 (normal user) or 2 (admin)"}
	}

	// Check if the user exists
	_, err := i.services.UserRepo.GetByID(request.UserId)
	if err != nil {
		return []string{"User not found"}
	}

	// Prevent admin from changing their own role
	if request.UserId == request.AdminUserId {
		return []string{"Cannot change your own role"}
	}

	// Update the user's role
	err = i.services.UserRepo.UpdateUserRole(request.UserId, request.Role)
	if err != nil {
		return []string{"Failed to update user role: " + err.Error()}
	}

	return nil
}
