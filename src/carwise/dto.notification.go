package carwise

type GetNotificationsRequest struct {
	UserId string
	Limit  int
	Page   int
}

type GetNotificationsResponse struct {
	Notifications []Notification `json:"notifications"`
	Total         int            `json:"total"`
	Unread        int            `json:"unread"`
}

type DeleteNotificationRequest struct {
	NotificationId string
	UserId         string
}

type ReadNotificationRequest struct {
	NotificationId string
	UserId         string
}
