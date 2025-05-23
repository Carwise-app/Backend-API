package carwise

type Message struct {
	Id         string
	SenderId   string
	ReceiverId string
	Message    string
	Read       bool
	CreatedAt  int64
}

type Chat struct {
	OtherUserId     string
	LastMessageTime int64
}
