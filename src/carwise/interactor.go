package carwise

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"mime/multipart"
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
			Logo:   brand.Logo,
			Name:   brand.Name,
			Count:  len(series),
			Series: seriesResponses,
		})
	}

	return brandResponses, nil
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

func (i *Interactor) GetProfile(id string) (*ProfileResponse, []string) {
	user, err := i.services.UserRepo.GetByID(id)
	if err != nil {
		return nil, []string{err.Error()}
	}
	return &ProfileResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		ImageUrl:    user.ImageUrl,
		CountryCode: user.CountryCode,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        user.Role,
		Status:      user.Status,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (i *Interactor) EditProfile(userId string, request ProfileEditRequest, avatar *multipart.FileHeader) []string {
	var errors []string
	user, err := i.services.UserRepo.GetByID(userId)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Failed to get user: %v", err))
		return errors
	}

	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.CountryCode = request.CountryCode
	user.PhoneNumber = request.PhoneNumber

	err = i.services.UserRepo.Update(user)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Failed to update user profile: %v", err))
		return errors
	}

	if avatar != nil {
		file, err := avatar.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to open avatar file: %v", err))
			return errors
		}
		defer file.Close()

		avatarURL, err := i.services.CDNRepo.SaveUserAvatar(userId, file)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to upload avatar: %v", err))
			return errors
		}

		user.ImageUrl = avatarURL
		err = i.services.UserRepo.Update(user)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to update user avatar URL: %v", err))
			return errors
		}
	}

	return nil
}

func (i *Interactor) CreateCar(userId string, request CarCreateRequest) []string {
	request.OwnerId = userId
	request.ID = uuid.New().String()
	request.ListingDate = time.Now()
	var err error
	request.ListingNumber, err = generateSecureListingNumber(10)
	if err != nil {
		return []string{err.Error()}
	}

	//for _, v := range request.Images {
	//	// v byte -> image
	//}

	err = i.services.CarRepo.Create(request.ToCar())
	if err != nil {
		return []string{err.Error()}
	}

	return nil
}

func (i *Interactor) ListCars(page, limit, brand_id, series_id, model_id int) ([]ListCarResponse, []string) {
	cars, err := i.services.CarRepo.GetCars(page, limit, brand_id, series_id, model_id)
	if err != nil {
		return nil, []string{"failed to fetch cars"}
	}

	brands, err := i.GetBrands()
	if err != nil {
		return nil, []string{"failed to fetch brands"}
	}

	brandMap := make(map[int]BrandResponse)
	seriesMap := make(map[int]string)
	modelMap := make(map[int]string)
	for _, brand := range brands {
		brandMap[brand.Id] = brand
		for _, series := range brand.Series {
			seriesMap[series.Id] = series.Name
			for _, models := range series.Models {
				modelMap[models.Id] = models.Name
			}
		}
	}
	var response []ListCarResponse
	for _, v := range cars {
		response = append(response, ListCarResponse{
			Id:          v.ID,
			Thumbnail:   "",
			Currency:    v.Currency,
			Price:       v.Price,
			Brand:       brandMap[v.BrandId].Name,
			Series:      seriesMap[v.SeriesId],
			Model:       modelMap[v.ModelId],
			Title:       v.Title,
			Year:        v.Year,
			Mileage:     v.Mileage,
			ListingDate: v.ListingDate,
			City:        v.City,
			District:    v.District,
		})
	}

	return response, nil
}

func (i *Interactor) GetCarDetail(id string) (*CarDetailResponse, []string) {
	car, err := i.services.CarRepo.GetByID(id)
	if err != nil {
		return nil, []string{"failed to fetch cars"}
	}
	owner, err := i.services.UserRepo.GetByID(car.OwnerId)
	if err != nil {
		return nil, []string{"Error fetching user by"}
	}

	ownerResponse := OwnerResponse{
		Id:          owner.ID,
		FirstName:   owner.FirstName,
		LastName:    owner.LastName,
		CountryCode: owner.CountryCode,
		PhoneNumber: owner.PhoneNumber,
		CreatedAt:   owner.CreatedAt,
	}

	brands, err := i.GetBrands()
	if err != nil {
		return nil, []string{"failed to fetch brands"}
	}

	brandMap := make(map[int]BrandResponse)
	seriesMap := make(map[int]string)
	modelMap := make(map[int]string)
	for _, brand := range brands {
		brandMap[brand.Id] = brand
		for _, series := range brand.Series {
			seriesMap[series.Id] = series.Name
			for _, models := range series.Models {
				modelMap[models.Id] = models.Name
			}
		}
	}

	carDetailResponse := &CarDetailResponse{
		ID:                car.ID,
		Owner:             ownerResponse,
		Title:             car.Title,
		Description:       car.Description,
		Currency:          car.Currency,
		Price:             car.Price,
		City:              car.City,
		District:          car.District,
		Neighborhood:      car.Neighborhood,
		ListingNumber:     car.ListingNumber,
		ListingDate:       car.ListingDate,
		Brand:             brandMap[car.BrandId].Name, 
		Series:            seriesMap[car.SeriesId],   
		Model:             modelMap[car.ModelId],     
		Year:              car.Year,
		FuelType:          car.FuelType,
		Transmission:      car.Transmission,
		Mileage:           car.Mileage,
		BodyType:          car.BodyType,
		EnginePower:       car.EnginePower,
		EngineVolume:      car.EngineVolume,
		DriveType:         car.DriveType,
		Color:             car.Color,
		Warranty:          car.Warranty,
		HeavyDamage:       car.HeavyDamage,
		SellerType:        car.SellerType,
		TradeOption:       car.TradeOption,
		FrontBumper:       car.FrontBumper,
		FrontHood:         car.FrontHood,
		Roof:              car.Roof,
		FrontRightDoor:    car.FrontRightDoor,
		RearRightDoor:     car.RearRightDoor,
		FrontLeftMudguard: car.FrontLeftMudguard,
		FrontLeftDoor:     car.FrontLeftDoor,
		RearLeftDoor:      car.RearLeftDoor,
		RearLeftMudguard:  car.RearLeftMudguard,
		RearBumper:        car.RearBumper,
		Images:            []string{},
	}

	return carDetailResponse, nil
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

func generateToken(length int) (string, error) {
	token := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(token), nil
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
