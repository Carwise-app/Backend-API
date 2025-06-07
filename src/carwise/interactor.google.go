package carwise

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (i *Interactor) GoogleAuthUrl() string {
	return i.services.GoogleAuth.GetUrl()
}

func (i *Interactor) GoogleAuthCallback(state, code string) (*User, error) {
	googleAccount, err := i.services.GoogleAuth.Callback(state, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get google account: %v", err)
	}
	account, err := i.services.UserRepo.GetByEmail(googleAccount.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing account: %v", err)
	}
	if account == nil || account.Id == "" {
		now := time.Now().Unix()
		newAccount := User{
			Id:          uuid.New().String(),
			FirstName:   googleAccount.FirstName,
			LastName:    googleAccount.LastName,
			CountryCode: "",
			PhoneNumber: "",
			Email:       googleAccount.Email,
			Password:    "",
			Role:        1,
			Status:      1,
			CreatedAt:   now,
			UpdatedAt:   now,
			LastLogin:   now,
			GoogleId:    googleAccount.Id,
		}
		err := i.services.UserRepo.Create(&newAccount)
		if err != nil {
			return nil, fmt.Errorf("failed to create google account: %v", err)
		}
		return &newAccount, nil
	} else if account.Password != "" {
		return nil, fmt.Errorf("this email is already registered with password. Please login with email and password")
	}

	return account, nil
}

func (i *Interactor) GoogleIdToken(idToken string) (*User, error) {
	googleAccount, err := i.services.GoogleAuth.VerifyIDToken(idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify google id token: %v", err)
	}

	account, err := i.services.UserRepo.GetByEmail(googleAccount.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user account: %v", err)
	}

	if account == nil || account.Id == "" {
		now := time.Now().Unix()
		newAccount := User{
			Id:          uuid.New().String(),
			FirstName:   googleAccount.FirstName,
			LastName:    googleAccount.LastName,
			CountryCode: "",
			PhoneNumber: "",
			Email:       googleAccount.Email,
			Password:    "",
			Role:        1,
			Status:      1,
			CreatedAt:   now,
			UpdatedAt:   now,
			LastLogin:   now,
			GoogleId:    googleAccount.Id,
		}
		err := i.services.UserRepo.Create(&newAccount)
		if err != nil {
			return nil, fmt.Errorf("failed to create google account: %v", err)
		}

		go func() {
			i.services.MailGW.SendEmail(newAccount.Email, "Hoşgeldiniz!", "welcome.html", map[string]interface{}{
				"full_name": newAccount.FirstName + " " + newAccount.LastName,
			})
			i.CreateNotification(
				&Notification{
					ID:        uuid.New().String(),
					Title:     "🎉 Carwise'e Hoş Geldiniz!",
					Status:    SystemMessage,
					Message:   "Hoşgeldin, " + newAccount.FirstName + "! 2. el araç alım-satımı ve araç fiyat tahmini artık çok daha kolay. Akıllı sistemimizle aracınızın gerçek değerini öğrenin, güvenle alım-satım yapın.",
					Read:      false,
					CreatedBy: newAccount.Id,
					CreatedAt: time.Now().Unix(),
				},
			)
		}()

		return &newAccount, nil
	} else if account.Password != "" {
		return nil, fmt.Errorf("this email is already registered with password. Please login with email and password")
	}

	return account, nil
}
