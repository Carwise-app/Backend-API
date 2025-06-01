package carwise

// @model SendMessageRequest
// @Description Request body for sending a message
type SendMessageRequest struct {
	ListingId        string `json:"-"`
	ReceiverId       string `json:"-"`
	Message          string `json:"message"`
	UserId           string `json:"-"`
	Role             int    `json:"-"`
	IsReceiverActive bool   `json:"-"`
}

// @model GetMessagesRequest
// @Description Request body for getting messages
type GetMessagesRequest struct {
	ListingId  string `json:"-"`
	ReceiverId string `json:"-"`
	Limit      int    `json:"-"`
	Page       int    `json:"-"`
	UserId     string `json:"-"`
	Role       int    `json:"-"`
}

// @model GetMessagesResponse
// @Description Response body for getting messages
type GetMessagesResponse struct {
	Messages []MessageInfo `json:"messages"`
	Total    int           `json:"total"`
}

// @model MessageInfo
// @Description Message information
type MessageInfo struct {
	Id        string   `json:"id"`
	Sender    UserInfo `json:"sender"`
	Receiver  UserInfo `json:"receiver"`
	Message   string   `json:"message"`
	Read      bool     `json:"read"`
	CreatedAt int64    `json:"created_at"`
}

// @model GetChatsRequest
// @Description Request body for getting chats
type GetChatsRequest struct {
	Limit  int    `json:"-"`
	Page   int    `json:"-"`
	UserId string `json:"-"`
	Role   int    `json:"-"`
}

// @model GetChatsResponse
// @Description Response body for getting chats
type GetChatsResponse struct {
	Chats []ChatInfo `json:"chats"`
	Total int        `json:"total"`
}

// @model ChatInfo
// @Description Chat information
type ChatInfo struct {
	User            UserInfo        `json:"user"`
	Listing         ListListingInfo `json:"listing"`
	LastMessageTime int64           `json:"last_message_time"`
}
