package carwise

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
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
		customData["image"] = imageUrl
	}

	if status == PriceDropped {
		users, err := i.services.FavoriteRepo.GetFavoritesUsers(customData["listing_id"])
		if err != nil {
			log.Printf("Error getting user emails for listing %s: %v", customData["listing_id"], err)
			return err
		}

		link := ""
		if customData["listing_id"] != "" {
			link = fmt.Sprintf("https://carwisegw.yusuftalhaklc.com/listing/%s", customData["listing_id"])
		}

		for _, user := range users {
			if user.EmailNotify {
				i.services.MailGW.SendEmail(user.Email, title, "notification.html", map[string]interface{}{
					"title":   title,
					"message": message,
					"image":   imageUrl,
					"link":    link,
				})
			}
			notificationData := make(map[string]any)
			for k, v := range customData {
				notificationData[k] = v
			}
			i.CreateNotification(&Notification{
				ID:        uuid.New().String(),
				Title:     title,
				Message:   message,
				Data:      notificationData,
				Status:    status,
				CreatedBy: user.Id,
				Read:      false,
				CreatedAt: time.Now().Unix(),
			})
		}

		deviceTokens := []string{}
		for _, user := range users {
			if user.PushNotify {
				deviceTokens = append(deviceTokens, user.DeviceToken)
			}
		}

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

		notificationData := make(map[string]any)
		for k, v := range customData {
			notificationData[k] = v
		}

		i.CreateNotification(&Notification{
			ID:        uuid.New().String(),
			Title:     title,
			Message:   message,
			Data:      notificationData,
			Status:    status,
			CreatedBy: userId,
			Read:      false,
			CreatedAt: time.Now().Unix(),
		})
	}
	return nil
}

func (i *Interactor) PushNotification(deviceTokens []string, title, message string, customData map[string]string, bigImage string) error {
	return i.services.OneSignalRepo.PushNotification(deviceTokens, title, message, customData, bigImage)
}
