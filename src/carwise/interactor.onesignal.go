package carwise

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

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
	imageUrl := ""
	if customData["listing_id"] != "" {
		listing, err := i.services.ListingRepo.GetListingById(customData["listing_id"])
		if err != nil {
			return err
		}

		image := Image{}
		if len(listing.Images) > 0 {
			firstImage, err := i.services.ImageRepo.GetImageById(listing.Images[0])
			if err == nil {
				image = *firstImage
			}
		}

		if image.Path != "" {
			imageUrl = "https://carwisegw.yusuftalhaklc.com" + strings.TrimPrefix(image.Path, ".")
		}
	}

	if status == PriceDropped {
		emails, err := i.services.FavoriteRepo.GetUserEmails(customData["listing_id"])
		if err != nil {
			return err
		}

		link := ""
		if customData["listing_id"] != "" {
			link = fmt.Sprintf("https://carwisegw.yusuftalhaklc.com/listing/%s", customData["listing_id"])
		}

		for _, email := range emails {
			i.services.MailGW.SendEmail(email, title, "notification", map[string]interface{}{
				"title":   title,
				"message": message,
				"image":   imageUrl,
				"link":    link,
			})
		}

		deviceTokens, err := i.services.FavoriteRepo.GetUserDeviceTokens(customData["listing_id"])
		if err != nil {
			return err
		}
		return i.PushNotification(deviceTokens, title, message, customData, imageUrl)
	} else {
		user, err := i.services.UserRepo.GetByID(userId)
		if err != nil {
			return err
		}

		if user.EmailNotify {
			link := ""
			if customData["listing_id"] != "" {
				link = fmt.Sprintf("https://carwisegw.yusuftalhaklc.com/listing/%s", customData["listing_id"])
			}

			i.services.MailGW.SendEmail(user.Email, title, "notification.html", map[string]interface{}{
				"title":   title,
				"message": message,
				"image":   imageUrl,
				"link":    link,
			})
		}

		if user.PushNotify {
			return i.PushNotification([]string{user.DeviceToken}, title, message, customData, imageUrl)
		}
	}
	return nil
}

func (i *Interactor) PushNotification(deviceTokens []string, title, message string, customData map[string]string, bigImage string) error {
	return i.services.OneSignalRepo.PushNotification(deviceTokens, title, message, customData, bigImage)
}

func (i *Interactor) PushNotificationToAll(request *PushNotificationRequest) error {
	if request.Role != 2 {
		return errors.New("admin user only can send push notification to all users")
	}

	emails, err := i.services.UserRepo.GetAllEmails()
	if err != nil {
		return err
	}

	for _, email := range emails {
		err := i.services.MailGW.SendEmail(email, request.Title, "notification.html", map[string]interface{}{
			"title":   request.Title,
			"message": request.Message,
			"image":   request.BigImage,
		})
		if err != nil {
			log.Println("Error sending email to user", email, err)
		}
	}

	deviceTokens, err := i.services.UserRepo.GetAllDeviceTokens()
	if err != nil {
		return err
	}
	return i.PushNotification(deviceTokens, request.Title, request.Message, request.Data, request.BigImage)
}
