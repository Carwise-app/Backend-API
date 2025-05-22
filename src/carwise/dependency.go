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

type MessageRepository interface {
	SaveMessage(message *Message) error
	GetMessagesBetween(senderId, receiverId string, limit, offset int) ([]Message, error)
}

type Services struct {
	UserRepo          UserRepository
	TokenRepo         TokenRepository
	MailGW            MailGateway
	PasswordResetRepo PasswordResetRepository
	CDNRepo           CDNRepository
	MessageRepo       MessageRepository
}
