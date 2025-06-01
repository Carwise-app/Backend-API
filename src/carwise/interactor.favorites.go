package carwise

import (
	"log"
	"time"
)

func (i *Interactor) CreateFavorite(request *FavoriteRequest) error {
	listing, err := i.services.ListingRepo.GetListingById(request.ListingId)
	if err != nil {
		return err
	}

	Favorite := &Favorite{
		UserId:    request.UserId,
		ListingId: listing.Id,
		CreatedAt: time.Now().Unix(),
	}

	if listing.CreatedBy != request.UserId {
		go func() {
			err = i.CreatePushNotification(
				FavoriteAdded,
				"",
				"",
				map[string]string{
					"listing_id": listing.Id,
					"user_id":    request.UserId,
				},
				listing.CreatedBy,
			)
			if err != nil {
				log.Println("CreatePushNotification error: ", err)
			}
		}()
	}

	return i.services.FavoriteRepo.CreateFavorite(Favorite)
}

func (i *Interactor) DeleteFavorite(request *DeleteFavoriteRequest) error {
	favorite, err := i.services.FavoriteRepo.GetFavoriteByUserIdAndListingId(request.UserId, request.ListingId)
	if err != nil {
		return err
	}

	return i.services.FavoriteRepo.DeleteFavorite(favorite)
}

func (i *Interactor) GetFavorites(request *GetFavoritesRequest) (*ListListingResponse, error) {
	offset := (request.Page - 1) * request.Limit
	favorites, err := i.services.FavoriteRepo.GetFavoritesByUserId(request.UserId, request.Limit, offset)
	if err != nil {
		log.Println("GetFavorites error: ", err)
		return nil, err
	}
	total, err := i.services.FavoriteRepo.CountFavoritesByUserId(request.UserId)
	if err != nil {
		log.Println("GetFavorites error: ", err)
		return nil, err
	}

	favoriteItems := make([]ListListingInfo, 0, len(favorites))
	for _, favorite := range favorites {
		listing, err := i.services.ListingRepo.GetListingById(favorite.ListingId)
		if err != nil {
			continue
		}

		brand, err := i.services.BrandRepo.GetById(listing.BrandId)
		if err != nil {
			return nil, err
		}

		series, err := i.services.BrandRepo.GetSeriesById(listing.SeriesId)
		if err != nil {
			return nil, err
		}

		model, err := i.services.BrandRepo.GetModelById(listing.ModelId)
		if err != nil {
			return nil, err
		}

		image := Image{}
		if len(listing.Images) > 0 {
			firstImage, err := i.services.ImageRepo.GetImageById(listing.Images[0])
			if err == nil {
				image = *firstImage
			}
		}

		favoriteItems = append(favoriteItems, ListListingInfo{
			Id:           listing.Id,
			Slug:         listing.Slug,
			Status:       listing.Status,
			Brand:        *brand,
			Series:       *series,
			Model:        *model,
			Title:        listing.Title,
			Currency:     listing.Currency,
			Price:        listing.Price,
			City:         listing.City,
			District:     listing.District,
			Neighborhood: listing.Neighborhood,
			Image:        image,
			CreatedAt:    listing.CreatedAt,
		})
	}

	return &ListListingResponse{
		Listings: favoriteItems,
		Total:    total,
	}, nil

}
