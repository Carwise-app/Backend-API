package carwise

import (
	"fmt"
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
	user, err := i.services.UserRepo.GetByEmail(request.Email)
	if err != nil {
		return nil, []string{err.Error()}
	}

	if !comparePasswords(user.PasswordHash, request.Password) {
		return nil, []string{"invalid credentials"}
	}

	return user, nil
}

func (i *Interactor) LogoutUser(token string) []string {

	return nil
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

func (i *Interactor) GetBrands() ([]BrandResponse, error) {
	brands, err := i.services.AuxRepo.GetBrands()
	if err != nil {
		return nil, fmt.Errorf("error fetching brands: %w", err)
	}

	var brandResponses []BrandResponse

	for _, brand := range brands {

		series, err := i.services.AuxRepo.GetSeriesByBrand(brand.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching series for brand %d: %w", brand.ID, err)
		}

		var seriesResponses []SeriesResponse
		for _, s := range series {
	
			models, err := i.services.AuxRepo.GetModelsBySeries(s.ID)
			if err != nil {
				return nil, fmt.Errorf("error fetching models for series %d: %w", s.ID, err)
			}

			var modelResponses []ModelResponse
			for _, model := range models {
				modelResponses = append(modelResponses, ModelResponse{
					Id:    model.ID,
					Name:  model.Name,
					Count: len(models),
				})
			}

	
			seriesResponses = append(seriesResponses, SeriesResponse{
				Id:     s.ID,
				Name:   s.Name,
				Count:  len(models), 
				Models: modelResponses,
			})
		}


		brandResponses = append(brandResponses, BrandResponse{
			Id:     brand.ID,
			Name:   brand.Name,
			Count:  len(series),
			Series: seriesResponses,
		})
	}

	return brandResponses, nil
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
