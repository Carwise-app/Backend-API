package carwise

type Message struct {
	Id         string
	ListingId  string
	SenderId   string
	ReceiverId string
	Message    string
	Read       bool
	CreatedAt  int64
}

type Chat struct {
	ListingId       string
	OtherUserId     string
	LastMessageTime int64
}
