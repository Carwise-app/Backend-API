package carwise

type Favorite struct {
	UserId    string `json:"user_id"`
	ListingId string `json:"listing_id"`
	CreatedAt int64  `json:"created_at"`
}
