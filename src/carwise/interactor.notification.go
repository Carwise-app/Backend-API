package carwise

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

func (i *Interactor) GetNotifications(request *GetNotificationsRequest) (*GetNotificationsResponse, error) {
	notifications, err := i.services.NotificationRepo.GetNotificationsByUserId(request.UserId, request.Limit, (request.Page-1)*request.Limit)
	if err != nil {
		return nil, err
	}

	unreadCount, err := i.services.NotificationRepo.GetUnreadNotificationsCount(request.UserId)
	if err != nil {
		return nil, err
	}

	totalCount, err := i.services.NotificationRepo.GetTotalNotificationsCount(request.UserId)
	if err != nil {
		return nil, err
	}

	if len(notifications) == 0 {
		return &GetNotificationsResponse{
			Notifications: []Notification{},
			Total:         0,
			Unread:        0,
		}, nil
	}

	return &GetNotificationsResponse{
		Notifications: notifications,
		Total:         totalCount,
		Unread:        unreadCount,
	}, nil
}

func (i *Interactor) ReadNotification(request *ReadNotificationRequest) error {
	return i.services.NotificationRepo.ReadNotification(request.NotificationId, request.UserId)
}

func (i *Interactor) DeleteNotification(request *DeleteNotificationRequest) error {
	return i.services.NotificationRepo.DeleteNotification(request.NotificationId, request.UserId)
}

func (i *Interactor) CreateNotification(notification *Notification) error {
	return i.services.NotificationRepo.CreateNotification(notification)
}

func (i *Interactor) PushNotificationToAll(request *PushNotificationRequest) error {
	if request.Role != 2 {
		return errors.New("admin user only can send push notification to all users")
	}

	users, err := i.services.UserRepo.GetAllUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.EmailNotify {
			err := i.services.MailGW.SendEmail(user.Email, request.Title, "notification.html", map[string]interface{}{
				"title":   request.Title,
				"message": request.Message,
				"image":   request.BigImage,
			})
			if err != nil {
				log.Println("Error sending email to user", user.Email, err)
			}
		}
		notificationData := make(map[string]any)
		for k, v := range request.Data {
			notificationData[k] = v
		}
		err := i.CreateNotification(&Notification{
			ID:        uuid.New().String(),
			Title:     request.Title,
			Message:   request.Message,
			Data:      notificationData,
			Status:    SystemMessage,
			CreatedBy: user.Id,
			Read:      false,
			CreatedAt: time.Now().Unix(),
		})
		if err != nil {
			log.Println("Error creating notification", err)
		}
	}

	deviceTokens := []string{}
	for _, user := range users {
		if user.PushNotify {
			deviceTokens = append(deviceTokens, user.DeviceToken)
		}
	}

	err = i.PushNotification(deviceTokens, request.Title, request.Message, request.Data, request.BigImage)
	if err != nil {
		log.Println("Error sending push notification to all users", err)
	}
	return err
}
