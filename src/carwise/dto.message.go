package carwise

type SendMessageRequest struct {
	ReceiverId string `json:"-"`
	Message    string `json:"message"`
	UserId     string `json:"-"`
	Role       int    `json:"-"`
}

type GetMessagesRequest struct {
	ReceiverId string `json:"-"`
	Limit      int    `json:"-"`
	Page       int    `json:"-"`
	UserId     string `json:"-"`
	Role       int    `json:"-"`
}

type GetMessagesResponse struct {
	Messages []MessageInfo `json:"messages"`
	Total    int           `json:"total"`
}

type MessageInfo struct {
	Id        string   `json:"id"`
	Sender    UserInfo `json:"sender"`
	Receiver  UserInfo `json:"receiver"`
	Message   string   `json:"message"`
	Read      bool     `json:"read"`
	CreatedAt int64    `json:"created_at"`
}

type GetChatsRequest struct {
	Limit  int    `json:"-"`
	Page   int    `json:"-"`
	UserId string `json:"-"`
	Role   int    `json:"-"`
}

type GetChatsResponse struct {
	Chats []ChatInfo `json:"chats"`
	Total int        `json:"total"`
}

type ChatInfo struct {
	User      UserInfo `json:"user"`
	LastMessageTime int64    `json:"last_message_time"`
}
