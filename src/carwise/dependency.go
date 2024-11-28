package carwise

type UserRepository interface {
	Create(*User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
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

type Services struct {
	UserRepo  UserRepository
	TokenRepo TokenRepository
	AuxRepo   AuxiliaryRepository
}
