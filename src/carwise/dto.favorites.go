package carwise

// FavoriteRequest is the request for creating a favorite
// @model FavoriteRequest
// @description The request for creating a favorite
type FavoriteRequest struct {
	UserId    string
	Role      int
	ListingId string
	CreatedAt int64
}

// DeleteFavoriteRequest is the request for deleting a favorite
// @model DeleteFavoriteRequest
// @description The request for deleting a favorite
type DeleteFavoriteRequest struct {
	UserId    string
	Role      int
	ListingId string
}

// GetFavoritesRequest is the request for getting favorites
// @model GetFavoritesRequest
// @description The request for getting favorites
type GetFavoritesRequest struct {
	UserId string
	Role   int
	Limit  int
	Page   int
}
