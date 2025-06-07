package carwise

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

func (i *Interactor) CreateListing(request *CreateListingRequest) (string, error) {
	brand, err := i.services.BrandRepo.GetById(request.BrandId)
	if err != nil {
		return "", err
	}

	series, err := i.services.BrandRepo.GetSeriesById(request.SeriesId)
	if err != nil {
		return "", err
	}

	if series.BrandId != brand.Id {
		return "", errors.New("series not found")
	}

	model, err := i.services.BrandRepo.GetModelById(request.ModelId)
	if err != nil {
		return "", err
	}

	if model.SeriesId != series.Id {
		return "", errors.New("model not found")
	}

	listing := &Listing{
		Id:                uuid.New().String(),
		BrandId:           request.BrandId,
		SeriesId:          request.SeriesId,
		ModelId:           request.ModelId,
		Slug:              makeUniqueSlugWithHash(request.Title),
		Title:             request.Title,
		Description:       request.Description,
		Currency:          request.Currency,
		Price:             request.Price,
		City:              request.City,
		District:          request.District,
		Neighborhood:      request.Neighborhood,
		Images:            request.Images,
		FuelType:          request.DetailInfo.FuelType,
		TransmissionType:  request.DetailInfo.TransmissionType,
		BodyType:          request.DetailInfo.BodyType,
		DriveType:         request.DetailInfo.DriveType,
		EnginePower:       request.DetailInfo.EnginePower,
		EngineVolume:      request.DetailInfo.EngineVolume,
		Kilometers:        request.DetailInfo.Kilometers,
		Year:              request.DetailInfo.Year,
		Color:             request.DetailInfo.Color,
		HeavyDamage:       request.DetailInfo.HeavyDamage,
		FrontBumper:       request.DetailInfo.FrontBumper,
		FrontHood:         request.DetailInfo.FrontHood,
		Roof:              request.DetailInfo.Roof,
		FrontRightDoor:    request.DetailInfo.FrontRightDoor,
		RearRightDoor:     request.DetailInfo.RearRightDoor,
		FrontLeftMudguard: request.DetailInfo.FrontLeftMudguard,
		FrontLeftDoor:     request.DetailInfo.FrontLeftDoor,
		RearLeftDoor:      request.DetailInfo.RearLeftDoor,
		RearLeftMudguard:  request.DetailInfo.RearLeftMudguard,
		RearBumper:        request.DetailInfo.RearBumper,
		CreatedBy:         request.UserId,
		CreatedAt:         time.Now().Unix(),
		UpdatedAt:         time.Now().Unix(),
		Status:            1,
	}

	err = i.services.ListingRepo.CreateListing(listing)
	if err != nil {
		return "", err
	}

	go func() {
		user, err := i.services.UserRepo.GetByID(listing.CreatedBy)
		if err != nil {
			log.Printf("GetByID error for listing %s: %v", listing.Id, err)
			return
		}

		message := fmt.Sprintf("Merhaba %s! İlanınız yayınlandı. İlanınızı görüntülemek için tıklayınız.", user.FirstName)

		err = i.CreatePushNotification(
			SystemMessage,
			"İlanınız Yayınlandı!",
			message,
			map[string]string{
				"listing_id": listing.Id,
			},
			user.Id,
		)
		if err != nil {
			log.Printf("CreatePushNotification error for listing %s: %v", listing.Id, err)
		} else {
			log.Printf("Price drop notification sent successfully for listing %s", listing.Id)
		}
	}()

	return listing.Id, nil
}

func (i *Interactor) GetListingById(request *GetListingRequest) (*GetListingResponse, error) {
	var listing *Listing
	var err error

	if isUUID(request.Id) {
		listing, err = i.services.ListingRepo.GetListingById(request.Id)
	} else {
		listing, err = i.services.ListingRepo.GetListingBySlug(request.Id)
	}

	if err != nil {
		return nil, err
	}

	createdBy, err := i.services.UserRepo.GetByID(listing.CreatedBy)
	if err != nil {
		return nil, err
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
	images := make([]Image, 0, len(listing.Images))
	for _, imageId := range listing.Images {
		image, err := i.services.ImageRepo.GetImageById(imageId)
		if err != nil {
			return nil, err
		}
		images = append(images, *image)

	}

	isFavorite := i.services.FavoriteRepo.IsFavorite(request.UserId, listing.Id)

	response := &GetListingResponse{
		Id:           listing.Id,
		Slug:         listing.Slug,
		Status:       listing.Status,
		Brand:        *brand,
		Series:       *series,
		Model:        *model,
		Title:        listing.Title,
		Description:  listing.Description,
		Currency:     listing.Currency,
		Price:        listing.Price,
		City:         listing.City,
		District:     listing.District,
		Neighborhood: listing.Neighborhood,
		Images:       images,
		DetailInfo: ListingDetailInfo{
			FuelType:          listing.FuelType,
			TransmissionType:  listing.TransmissionType,
			BodyType:          listing.BodyType,
			DriveType:         listing.DriveType,
			EnginePower:       listing.EnginePower,
			EngineVolume:      listing.EngineVolume,
			Kilometers:        listing.Kilometers,
			Year:              listing.Year,
			Color:             listing.Color,
			HeavyDamage:       listing.HeavyDamage,
			FrontBumper:       listing.FrontBumper,
			FrontHood:         listing.FrontHood,
			Roof:              listing.Roof,
			FrontRightDoor:    listing.FrontRightDoor,
			RearRightDoor:     listing.RearRightDoor,
			FrontLeftMudguard: listing.FrontLeftMudguard,
			FrontLeftDoor:     listing.FrontLeftDoor,
			RearLeftDoor:      listing.RearLeftDoor,
			RearLeftMudguard:  listing.RearLeftMudguard,
			RearBumper:        listing.RearBumper,
		},
		IsFavorite: isFavorite,
		CreatedBy: UserInfo{
			Id:          createdBy.Id,
			FirstName:   createdBy.FirstName,
			LastName:    createdBy.LastName,
			Email:       createdBy.Email,
			CountryCode: createdBy.CountryCode,
			PhoneNumber: createdBy.PhoneNumber,
		},
		CreatedAt: listing.CreatedAt,
		UpdatedAt: listing.UpdatedAt,
	}

	return response, nil
}

func (i *Interactor) UpdateListing(request *UpdateListingRequest) error {
	listing, err := i.services.ListingRepo.GetListingById(request.Id)
	if err != nil {
		return err
	}

	if listing.Status == 2 || listing.Status == 3 {
		return errors.New("listing is not active")
	}

	if request.Role != 2 {
		if listing.CreatedBy != request.UserId {
			return errors.New("unauthorized")
		}
	}

	oldPrice := listing.Price

	listing.BrandId = request.BrandId
	listing.SeriesId = request.SeriesId
	listing.ModelId = request.ModelId
	listing.Title = request.Title
	listing.Description = request.Description
	listing.Currency = request.Currency
	listing.Price = request.Price
	listing.City = request.City
	listing.District = request.District
	listing.Neighborhood = request.Neighborhood
	listing.Images = request.Images
	listing.FuelType = request.DetailInfo.FuelType
	listing.TransmissionType = request.DetailInfo.TransmissionType
	listing.BodyType = request.DetailInfo.BodyType
	listing.DriveType = request.DetailInfo.DriveType
	listing.EnginePower = request.DetailInfo.EnginePower
	listing.EngineVolume = request.DetailInfo.EngineVolume
	listing.Kilometers = request.DetailInfo.Kilometers
	listing.Year = request.DetailInfo.Year
	listing.Color = request.DetailInfo.Color
	listing.HeavyDamage = request.DetailInfo.HeavyDamage
	listing.FrontBumper = request.DetailInfo.FrontBumper
	listing.FrontHood = request.DetailInfo.FrontHood
	listing.Roof = request.DetailInfo.Roof
	listing.FrontRightDoor = request.DetailInfo.FrontRightDoor
	listing.RearRightDoor = request.DetailInfo.RearRightDoor
	listing.FrontLeftMudguard = request.DetailInfo.FrontLeftMudguard
	listing.FrontLeftDoor = request.DetailInfo.FrontLeftDoor
	listing.RearLeftDoor = request.DetailInfo.RearLeftDoor
	listing.RearLeftMudguard = request.DetailInfo.RearLeftMudguard
	listing.RearBumper = request.DetailInfo.RearBumper
	listing.UpdatedAt = time.Now().Unix()

	err = i.services.ListingRepo.UpdateListing(listing)
	if err != nil {
		return err
	}

	log.Printf("Comparing prices - New: %d, Old: %d", request.Price, oldPrice)

	if request.Price < oldPrice {
		log.Printf("Price drop condition met!")
		message := fmt.Sprintf("🎉 Harika haber! Takip ettiğiniz araçta fiyat düşüşü var.\n\nÖnceki fiyat: %s TL\nYeni fiyat: %s TL\n\n💰 %s TL tasarruf edebilirsiniz!",
			formatPrice(oldPrice),
			formatPrice(request.Price),
			formatPrice(oldPrice-request.Price),
		)

		log.Printf("Price drop detected - Listing ID: %s, Old Price: %d, New Price: %d, Difference: %d",
			listing.Id, oldPrice, request.Price, oldPrice-request.Price)

		go func() {
			err = i.CreatePushNotification(
				PriceDropped,
				"",
				message,
				map[string]string{
					"listing_id": listing.Id,
				},
				"",
			)
			if err != nil {
				log.Printf("CreatePushNotification error for listing %s: %v", listing.Id, err)
			} else {
				log.Printf("Price drop notification sent successfully for listing %s", listing.Id)
			}
		}()
	} else {
		go func() {
			user, err := i.services.UserRepo.GetByID(listing.CreatedBy)
			if err != nil {
				log.Printf("GetByID error for listing %s: %v", listing.Id, err)
				return
			}

			message := fmt.Sprintf("Merhaba %s! İlanınız güncellendi. İlanınızı görüntülemek için tıklayınız.", user.FirstName)

			i.CreatePushNotification(
				SystemMessage,
				"İlanınız Güncellendi!",
				message,
				map[string]string{
					"listing_id": listing.Id,
				},
				user.Id,
			)
		}()
	}

	return nil
}

func (i *Interactor) DeleteListing(request *DeleteListingRequest) error {
	listing, err := i.services.ListingRepo.GetListingById(request.Id)
	if err != nil {
		return err
	}

	if request.Role != 2 {
		if listing.CreatedBy != request.UserId {
			return errors.New("unauthorized")
		}
	}

	err = i.services.ListingRepo.DeleteListing(listing.Id)
	if err != nil {
		return err
	}

	go func() {
		user, err := i.services.UserRepo.GetByID(listing.CreatedBy)
		if err != nil {
			log.Printf("GetByID error for listing %s: %v", listing.Id, err)
			return
		}

		message := fmt.Sprintf("Merhaba %s! İlanınız silindi.", user.FirstName)

		i.CreatePushNotification(
			SystemMessage,
			"İlanınız Silindi!",
			message,
			map[string]string{},
			user.Id,
		)
	}()

	return nil
}

func (i *Interactor) UpdateListingStatus(request *UpdateListingStatusRequest) error {
	listing, err := i.services.ListingRepo.GetListingById(request.Id)
	if err != nil {
		return err
	}

	if listing.Status == 3 {
		return errors.New("listing has been closed by admin, no modifications allowed")
	}

	if request.Role != 2 {
		if listing.CreatedBy != request.UserId {
			return errors.New("unauthorized")
		}
	}

	listing.Status = request.Status
	listing.UpdatedAt = time.Now().Unix()

	err = i.services.ListingRepo.UpdateListing(listing)
	if err != nil {
		return err
	}

	if listing.Status == 1 {
		go func() {
			user, err := i.services.UserRepo.GetByID(listing.CreatedBy)
			if err != nil {
				log.Printf("GetByID error for listing %s: %v", listing.Id, err)
				return
			}

			message := fmt.Sprintf("Merhaba %s! İlanınız aktif hale getirildi. İlanınızı görüntülemek için tıklayınız.", user.FirstName)

			i.CreatePushNotification(
				SystemMessage,
				"İlanınız Aktif Hale Getirildi!",
				message,
				map[string]string{
					"listing_id": listing.Id,
				},
				user.Id,
			)
		}()
	}

	if listing.Status == 2 {
		go func() {
			user, err := i.services.UserRepo.GetByID(listing.CreatedBy)
			if err != nil {
				log.Printf("GetByID error for listing %s: %v", listing.Id, err)
				return
			}

			message := fmt.Sprintf("Merhaba %s! İlanınız satıldı olarak işaretlendi. İlanınızı görüntülemek için tıklayınız.", user.FirstName)

			i.CreatePushNotification(
				SystemMessage,
				"İlanınız Satıldı Olarak İşaretlendi!",
				message,
				map[string]string{
					"listing_id": listing.Id,
				},
				user.Id,
			)
		}()
	}

	return nil
}

func (i *Interactor) ListListing(request *ListListingRequest) (*ListListingResponse, error) {
	listings, err := i.services.ListingRepo.ListListing(&request.Filter)
	if err != nil {
		return nil, err
	}

	total, err := i.services.ListingRepo.CountListing(&request.Filter)
	if err != nil {
		return nil, err
	}

	listingsInfo := make([]ListListingInfo, 0, len(listings))
	for _, listing := range listings {
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

		isFavorite := i.services.FavoriteRepo.IsFavorite(request.Filter.UserId, listing.Id)

		listingsInfo = append(listingsInfo, ListListingInfo{
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
			IsFavorite:   isFavorite,
		})
	}

	response := &ListListingResponse{
		Listings: listingsInfo,
		Total:    total,
	}

	return response, nil
}

func makeUniqueSlugWithHash(title string) string {
	baseSlug := slug.Make(title)
	data := fmt.Sprintf("%s%d", title, time.Now().UnixNano())
	hash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	return fmt.Sprintf("%s-%s", baseSlug, hash[:8])
}

func isUUID(str string) bool {
	return len(str) == 36 && strings.Count(str, "-") == 4
}

func formatPrice(price int) string {
	return humanize.Comma(int64(price))
}
