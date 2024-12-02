package carwise

import (
	"io"
	"time"
)

type UserRepository interface {
	Create(*User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	UpdatePassword(email, hashedPassword string) error
	Update(user *User) error
}

type TokenRepository interface {
	IsTokenBlackListed(token string) (bool, error)
	AddTokenBlackList(token string) error
}

type AuxiliaryRepository interface {
	GetBrands() ([]Brand, error)
	GetSeriesByBrand(brandID int) ([]Series, error)
	GetModelsBySeries(seriesID int) ([]Model, error)
}

type MailGateway interface {
	Send(To string, Body []byte) error
}

type PasswordResetRepository interface {
	SaveResetCode(email, code string, ttl time.Duration) error
	VerifyResetCode(email, code string) (bool, error)
	DeleteResetCode(email string) error
}

type CDNRepository interface {
	SaveUserAvatar(userID string, image io.Reader) (string, error)
}

type CarRepository interface {
	Create(car *Car) error
	GetCars(page, limit, brand_id, series_id, model_id int) ([]Car, error)
	GetByID(id string) (*Car, error)
}

type Services struct {
	UserRepo          UserRepository
	TokenRepo         TokenRepository
	AuxRepo           AuxiliaryRepository
	MailGW            MailGateway
	PasswordResetRepo PasswordResetRepository
	CDNRepo           CDNRepository
	CarRepo           CarRepository
}
