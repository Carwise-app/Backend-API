package carwise

type UserRepository interface {
	Create(*User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
}

type Services struct {
	UserRepo UserRepository
}
