package carwise

import "errors"

var (
	SystemMessage = 1
	ChatMessage   = 2
	FavoriteAdded = 3
	PriceDropped  = 4
)

func (i *Interactor) CreatePushNotification(
	status int,
	t,
	m string,
	customData map[string]string,
	userId string,
	largeIcon string,
) error {
	var title string
	var message string

	switch status {
	case SystemMessage:
		title = t
		message = m
	case ChatMessage:
		title = "Yeni Mesaj"
		message = "Aracınızla ilgili yeni bir mesajınız var."
	case FavoriteAdded:
		title = "Aracınız Favorilere Eklendi"
		message = "Bir kullanıcı aracınızı favorilerine ekledi."
	case PriceDropped:
		title = "Fiyat Güncellemesi"
		message = "Favorilerinizdeki aracın fiyatı düştü!"
	default:
		return errors.New("geçersiz bildirim türü")
	}

	if status == PriceDropped {
		deviceTokens, err := i.services.FavoriteRepo.GetUserDeviceTokens(customData["listing_id"])
		if err != nil {
			return err
		}
		return i.PushNotification(deviceTokens, title, message, customData, largeIcon)
	} else {
		user, err := i.services.UserRepo.GetByID(userId)
		if err != nil {
			return err
		}
		if user.PushNotify {
			return i.PushNotification([]string{user.DeviceToken}, title, message, customData, largeIcon)
		}
	}
	return nil
}

func (i *Interactor) PushNotification(deviceTokens []string, title, message string, customData map[string]string, largeIcon string) error {
	return i.services.OneSignalRepo.PushNotification(deviceTokens, title, message, customData, largeIcon)
}

func (i *Interactor) PushNotificationToAll(request *PushNotificationRequest) error {
	if request.Role != 2 {
		return errors.New("admin user only can send push notification to all users")
	}
	deviceTokens, err := i.services.UserRepo.GetAllDeviceTokens()
	if err != nil {
		return err
	}
	return i.PushNotification(deviceTokens, request.Title, request.Message, request.Data, request.LargeIcon)
}
