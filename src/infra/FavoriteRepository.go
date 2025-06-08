package infra

import (
	"carwise"
	"database/sql"
)

type FavoriteRepository struct {
	db *sql.DB
}

func NewFavoriteRepository() *FavoriteRepository {
	db := ConnectDb()
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) CreateFavorite(favorite *carwise.Favorite) error {
	query := `
	INSERT INTO favorites (user_id, listing_id)
	VALUES ($1, $2)
	`
	_, err := r.db.Exec(query, favorite.UserId, favorite.ListingId)
	return err
}

func (r *FavoriteRepository) GetFavoriteByUserIdAndListingId(userId, listingId string) (*carwise.Favorite, error) {
	query := `
	SELECT user_id, listing_id, created_at
	FROM favorites
	WHERE user_id = $1 AND listing_id = $2
	LIMIT 1
	`
	var favorite carwise.Favorite
	err := r.db.QueryRow(query, userId, listingId).Scan(&favorite.UserId, &favorite.ListingId, &favorite.CreatedAt)
	return &favorite, err
}

func (r *FavoriteRepository) GetFavoritesByUserId(userId string, limit, offset int) ([]carwise.Favorite, error) {
	query := `
	SELECT user_id, listing_id, created_at
	FROM favorites
	WHERE user_id = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	favorites := []carwise.Favorite{}
	for rows.Next() {
		var favorite carwise.Favorite
		err := rows.Scan(&favorite.UserId, &favorite.ListingId, &favorite.CreatedAt)
		if err != nil {
			return nil, err
		}
		favorites = append(favorites, favorite)
	}
	return favorites, nil
}

func (r *FavoriteRepository) CountFavoritesByUserId(userId string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM favorites
		WHERE user_id = $1
	`
	var count int
	err := r.db.QueryRow(query, userId).Scan(&count)
	return count, err
}

func (r *FavoriteRepository) DeleteFavorite(favorite *carwise.Favorite) error {
	query := `
	DELETE FROM favorites
	WHERE user_id = $1 AND listing_id = $2
	`
	_, err := r.db.Exec(query, favorite.UserId, favorite.ListingId)
	return err
}

func (r *FavoriteRepository) IsFavorite(userId, listingId string) bool {
	query := `
	SELECT EXISTS (
		SELECT 1
		FROM favorites
		WHERE user_id = $1 AND listing_id = $2
	)
	`
	var exists bool
	err := r.db.QueryRow(query, userId, listingId).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (r *FavoriteRepository) GetFavoritesUsers(listingId string) ([]carwise.User, error) {
	query := `
	SELECT u.id, u.email, u.device_token, u.email_notify, u.push_notify
	FROM favorites f
	INNER JOIN users u ON f.user_id = u.id
	WHERE f.listing_id = $1
	`
	rows, err := r.db.Query(query, listingId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []carwise.User{}
	for rows.Next() {
		var user carwise.User
		err := rows.Scan(&user.Id, &user.Email, &user.DeviceToken, &user.EmailNotify, &user.PushNotify)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *FavoriteRepository) CountFavorites() (int, error) {
	query := `
		SELECT COUNT(*) FROM favorites
	`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}
