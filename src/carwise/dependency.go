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
	GetAllUsers() ([]User, error)
	GetUsers(page, limit int) ([]User, error)
	CountUsers() (int, error)
}

type RedisRepository interface {
	IsTokenBlackListed(token string) (bool, error)
	AddTokenBlackList(token string) error
	SetBrandsWithDetails(brands []BrandWithDetails) error
	GetBrandsWithDetails() ([]BrandWithDetails, error)
}

type MailGateway interface {
	SendEmail(to, subject, templateName string, data map[string]interface{}) error
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
	ImagesCount() (int, error)
}

type MessageRepository interface {
	SaveMessage(message *Message) error
	GetMessagesByUserId(userId, otherUserId, listingId string, limit, offset int) ([]Message, error)
	CountMessagesByUserId(userId, otherUserId, listingId string) (int, error)
	GetChats(userId string, limit, offset int) ([]Chat, error)
	CountChats(userId string) (int, error)
	ReadMessage(messageId string) error
	MessagesCount() (int, error)
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

	GetFavoritesUsers(listingId string) ([]User, error)
	CountFavorites() (int, error)
}

type PredictRepository interface {
	SavePredict(predict *Predict) error
	GetPredictByUserId(userId string, page, limit int) ([]Predict, error)
	CountPredictByUserId(userId string) (int, error)
	PredictCount() (int, error)
}

type OneSignalRepository interface {
	PushNotification(deviceTokens []string, title, message string, customData map[string]string, bigImage string) error
}

type NotificationRepository interface {
	CreateNotification(notification *Notification) error
	GetNotificationsByUserId(userId string, limit, offset int) ([]Notification, error)
	GetTotalNotificationsCount(userId string) (int, error)
	ReadNotification(id string, userId string) error
	DeleteNotification(id string, userId string) error
	GetUnreadNotificationsCount(userId string) (int, error)
	NotificationsCount() (int, error)
}

type Services struct {
	UserRepo          UserRepository
	RedisRepo         RedisRepository
	MailGW            MailGateway
	PasswordResetRepo PasswordResetRepository
	BrandRepo         BrandRepository
	ListingRepo       ListingRepository
	ImageRepo         ImageRepository
	MessageRepo       MessageRepository
	PredictionRepo    PredictionRepository
	GoogleAuth        GoogleAuth
	FavoriteRepo      FavoriteRepository
	PredictRepo       PredictRepository
	OneSignalRepo     OneSignalRepository
	NotificationRepo  NotificationRepository
}
