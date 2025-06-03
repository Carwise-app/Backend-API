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
	log.Printf("CreatePushNotification called - Status: %d, UserID: %s",
		status, userId)

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
		title = "Favorilerinizdeki aracın fiyatı düştü!"
		message = m
		log.Printf("Price drop notification - Title: %s, Message: %s", title, message)
	default:
		log.Printf("Invalid notification type: %d", status)
		return errors.New("geçersiz bildirim türü")
	}

	imageUrl := ""
	if customData["listing_id"] != "" {
		listing, err := i.services.ListingRepo.GetListingById(customData["listing_id"])
		if err != nil {
			log.Printf("Error getting listing by ID %s: %v", customData["listing_id"], err)
			return err
		}

		image := Image{}
		if len(listing.Images) > 0 {
			firstImage, err := i.services.ImageRepo.GetImageById(listing.Images[0])
			if err == nil {
				image = *firstImage
			} else {
				log.Printf("Error getting image for listing %s: %v", customData["listing_id"], err)
			}
		}

		if image.Path != "" {
			imageUrl = "https://carwisegw.yusuftalhaklc.com" + strings.TrimPrefix(image.Path, ".")
			log.Printf("Image URL set for listing %s: %s", customData["listing_id"], imageUrl)
		}
	}

	if status == PriceDropped {
		emails, err := i.services.FavoriteRepo.GetUserEmails(customData["listing_id"])
		if err != nil {
			log.Printf("Error getting user emails for listing %s: %v", customData["listing_id"], err)
			return err
		}
		log.Printf("Found %d users to notify for price drop on listing %s", len(emails), customData["listing_id"])

		link := ""
		if customData["listing_id"] != "" {
			link = fmt.Sprintf("https://carwisegw.yusuftalhaklc.com/listing/%s", customData["listing_id"])
		}

		for _, email := range emails {
			log.Printf("Sending price drop email to: %s", email)
			i.services.MailGW.SendEmail(email, title, "notification", map[string]interface{}{
				"title":   title,
				"message": message,
				"image":   imageUrl,
				"link":    link,
			})
		}

		deviceTokens, err := i.services.FavoriteRepo.GetUserDeviceTokens(customData["listing_id"])
		if err != nil {
			log.Printf("Error getting device tokens for listing %s: %v", customData["listing_id"], err)
			return err
		}
		log.Printf("Found %d device tokens for push notification on listing %s", len(deviceTokens), customData["listing_id"])

		err = i.PushNotification(deviceTokens, title, message, customData, imageUrl)
		if err != nil {
			log.Printf("Error sending push notifications for listing %s: %v", customData["listing_id"], err)
		} else {
			log.Printf("Push notifications sent successfully for listing %s", customData["listing_id"])
		}
		return err
	} else {
		user, err := i.services.UserRepo.GetByID(userId)
		if err != nil {
			log.Printf("Error getting user by ID %s: %v", userId, err)
			return err
		}

		if user.EmailNotify {
			link := ""
			if customData["listing_id"] != "" {
				link = fmt.Sprintf("https://carwisegw.yusuftalhaklc.com/listing/%s", customData["listing_id"])
			}

			log.Printf("Sending email notification to user %s (%s)", userId, user.Email)
			i.services.MailGW.SendEmail(user.Email, title, "notification.html", map[string]interface{}{
				"title":   title,
				"message": message,
				"image":   imageUrl,
				"link":    link,
			})
		} else {
			log.Printf("Email notifications disabled for user %s", userId)
		}

		if user.PushNotify {
			log.Printf("Sending push notification to user %s", userId)
			err = i.PushNotification([]string{user.DeviceToken}, title, message, customData, imageUrl)
			if err != nil {
				log.Printf("Error sending push notification to user %s: %v", userId, err)
			}
			return err
		} else {
			log.Printf("Push notifications disabled for user %s", userId)
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
