package carwise

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Interactor struct {
	services Services
}

func NewInteractor(svcs Services) *Interactor {
	return &Interactor{
		services: svcs,
	}
}

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
		ID:           uuid.New().String(),
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		CountryCode:  request.CountryCode,
		PhoneNumber:  request.PhoneNumber,
		Email:        request.Email,
		PasswordHash: hashedPassword,
		Role:         UserRoleRegular,
		Status:       AccountStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastLogin:    time.Now(),
	}

	err = i.services.UserRepo.Create(user)
	if err != nil {
		return nil, []string{"Failed to create user: " + err.Error()}
	}
	return user, nil
}

func (i *Interactor) LoginUser(request UserLoginRequest) (*User, []string) {

	return nil, nil
}

func (i *Interactor) LogoutUser(token string) []string {

	return nil
}

func (i *Interactor) IsTokenBlackListed(token string) (bool, []string) {

	return false, nil
}

func (i *Interactor) AddTokenBlackList(token string) []string {

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
