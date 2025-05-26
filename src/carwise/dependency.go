package carwise

import (
	"mime/multipart"
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

type BrandRepository interface {
	GetAll() ([]Brand, error)
	GetAllWithDetails() ([]BrandWithDetails, error)
	Create(brand *Brand) error
	Update(brand *Brand) error
	Delete(id string) error
	GetById(id string) (*Brand, error)

	GetAllSeriesByBrandId(brandId string) ([]Series, error)
	CreateSeries(series *Series) error
	UpdateSeries(series *Series) error
	DeleteSeries(id string) error
	GetSeriesById(id string) (*Series, error)

	GetAllModelsBySeriesId(seriesId string) ([]Model, error)
	CreateModel(model *Model) error
	UpdateModel(model *Model) error
	DeleteModel(id string) error
	GetModelById(id string) (*Model, error)
}

type ListingRepository interface {
	CreateListing(listing *Listing) error
	GetListingById(id string) (*Listing, error)
	GetListingBySlug(slug string) (*Listing, error)
	UpdateListing(listing *Listing) error
	DeleteListing(id string) error
	ListListing(filter *ListingFilter) ([]Listing, error)
	CountListing(filter *ListingFilter) (int, error)
}

type ImageRepository interface {
	SaveImage(file *multipart.FileHeader, userId string) (*Image, error)
	GetImageById(id string) (*Image, error)
	DeleteImage(id string) error
}

type MessageRepository interface {
	SaveMessage(message *Message) error
	GetMessagesByUserId(userId string, limit, offset int) ([]Message, error)
	CountMessagesByUserId(userId string) (int, error)
	GetChats(userId string, limit, offset int) ([]Chat, error)
	CountChats(userId string) (int, error)
	ReadMessage(messageId string) error
}

type PredictionRepository interface {
	SaveImagePrediction(prediction *ImagePrediction) error
	GetImagePredictionsByImageId(imageId string) (*ImagePrediction, error)
}

type GoogleAuth interface {
	VerifyIDToken(idToken string) (*GoogleResponse, error)
	GetUrl() string
	Callback(state, code string) (*GoogleResponse, error)
}

type FavoriteRepository interface {
	CreateFavorite(favorite *Favorite) error
	GetFavoriteByUserIdAndListingId(userId, listingId string) (*Favorite, error)
	GetFavoritesByUserId(userId string, limit, offset int) ([]Favorite, error)
	CountFavoritesByUserId(userId string) (int, error)
	DeleteFavorite(favorite *Favorite) error
	IsFavorite(userId, listingId string) bool
}

type Services struct {
	UserRepo          UserRepository
	TokenRepo         TokenRepository
	MailGW            MailGateway
	PasswordResetRepo PasswordResetRepository
	BrandRepo         BrandRepository
	ListingRepo       ListingRepository
	ImageRepo         ImageRepository
	MessageRepo       MessageRepository
	PredictionRepo    PredictionRepository
	GoogleAuth        GoogleAuth
	FavoriteRepo      FavoriteRepository
}
