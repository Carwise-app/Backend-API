package carwise

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strings"
	"time"

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
	}

	err = i.services.ListingRepo.CreateListing(listing)
	if err != nil {
		return "", err
	}

	return listing.Id, nil
}

func (i *Interactor) GetListingById(idOrSlug string) (*GetListingResponse, error) {
	var listing *Listing
	var err error

	if isUUID(idOrSlug) {
		listing, err = i.services.ListingRepo.GetListingById(idOrSlug)
	} else {
		listing, err = i.services.ListingRepo.GetListingBySlug(idOrSlug)
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

	response := &GetListingResponse{
		Id:           listing.Id,
		Slug:         listing.Slug,
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
		Images:       []Image{},
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

func makeUniqueSlugWithHash(title string) string {
	baseSlug := slug.Make(title)
	data := fmt.Sprintf("%s%d", title, time.Now().UnixNano())
	hash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	return fmt.Sprintf("%s-%s", baseSlug, hash[:8])
}

func isUUID(str string) bool {
	return len(str) == 36 && strings.Count(str, "-") == 4
}
