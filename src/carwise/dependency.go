package carwise

import "time"

type UserRepository interface {
	Create(*User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	UpdatePassword(email, hashedPassword string) error
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

type Services struct {
	UserRepo          UserRepository
	TokenRepo         TokenRepository
	AuxRepo           AuxiliaryRepository
	MailGW            MailGateway
	PasswordResetRepo PasswordResetRepository
}
