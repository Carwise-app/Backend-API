package carwise

func (i *Interactor) GetNotifications(request *GetNotificationsRequest) (*GetNotificationsResponse, error) {
	notifications, err := i.services.NotificationRepo.GetNotificationsByUserId(request.UserId, request.Limit, request.Page)
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
