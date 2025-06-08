package carwise

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func (i *Interactor) CreateBrand(request *BrandCreateRequest) (*Brand, error) {
	if request.Role != 2 {
		return nil, errors.New("unauthorized")
	}

	brand := &Brand{
		Id:        uuid.New().String(),
		ImagePath: request.ImagePath,
		Name:      request.Name,
	}

	err := i.services.BrandRepo.Create(brand)
	if err != nil {
		return nil, err
	}

	go i.BrandCache()

	return brand, nil
}

func (i *Interactor) UpdateBrand(request *BrandUpdateRequest) error {
	if request.Role != 2 {
		return errors.New("unauthorized")
	}

	brand, err := i.services.BrandRepo.GetById(request.BrandId)
	if err != nil {
		return err
	}
	brand.ImagePath = request.ImagePath
	brand.Name = request.Name

	err = i.services.BrandRepo.Update(brand)
	if err != nil {
		return err
	}

	go i.BrandCache()

	return nil
}

func (i *Interactor) DeleteBrand(request *BrandDeleteRequest) error {
	if request.Role != 2 {
		return errors.New("unauthorized")
	}

	err := i.services.BrandRepo.Delete(request.BrandId)
	if err != nil {
		return err
	}

	go i.BrandCache()

	return nil
}

func (i *Interactor) CreateSeries(request *SeriesCreateRequest) (*Series, error) {
	if request.Role != 2 {
		return nil, errors.New("unauthorized")
	}

	_, err := i.services.BrandRepo.GetById(request.BrandId)
	if err != nil {
		return nil, err
	}

	series := &Series{
		Id:      uuid.New().String(),
		BrandId: request.BrandId,
		Name:    request.Name,
	}

	err = i.services.BrandRepo.CreateSeries(series)
	if err != nil {
		return nil, err
	}

	go i.BrandCache()

	return series, nil
}

func (i *Interactor) UpdateSeries(request *SeriesUpdateRequest) error {
	if request.Role != 2 {
		return errors.New("unauthorized")
	}

	series, err := i.services.BrandRepo.GetSeriesById(request.SeriesId)
	if err != nil {
		return err
	}
	series.Name = request.Name

	err = i.services.BrandRepo.UpdateSeries(series)
	if err != nil {
		return err
	}

	go i.BrandCache()

	return nil
}

func (i *Interactor) DeleteSeries(request *SeriesDeleteRequest) error {
	if request.Role != 2 {
		return errors.New("unauthorized")
	}

	err := i.services.BrandRepo.DeleteSeries(request.SeriesId)
	if err != nil {
		return err
	}

	go i.BrandCache()

	return nil
}

func (i *Interactor) CreateModel(request *ModelCreateRequest) (*Model, error) {
	if request.Role != 2 {
		return nil, errors.New("unauthorized")
	}

	_, err := i.services.BrandRepo.GetById(request.BrandId)
	if err != nil {
		return nil, err
	}

	_, err = i.services.BrandRepo.GetSeriesById(request.SeriesId)
	if err != nil {
		return nil, err
	}

	model := &Model{
		Id:       uuid.New().String(),
		BrandId:  request.BrandId,
		SeriesId: request.SeriesId,
		Name:     request.Name,
	}

	err = i.services.BrandRepo.CreateModel(model)
	if err != nil {
		return nil, err
	}

	go i.BrandCache()

	return model, nil
}

func (i *Interactor) UpdateModel(request *ModelUpdateRequest) error {
	if request.Role != 2 {
		return errors.New("unauthorized")
	}

	model, err := i.services.BrandRepo.GetModelById(request.ModelId)
	if err != nil {
		return err
	}
	model.Name = request.Name

	err = i.services.BrandRepo.UpdateModel(model)
	if err != nil {
		return err
	}

	go i.BrandCache()

	return nil
}

func (i *Interactor) DeleteModel(request *ModelDeleteRequest) error {
	if request.Role != 2 {
		return errors.New("unauthorized")
	}
	err := i.services.BrandRepo.DeleteModel(request.ModelId)
	if err != nil {
		return err
	}

	go i.BrandCache()

	return nil
}

func (i *Interactor) GetAllBrandsWithDetails() ([]BrandWithDetails, error) {
	brands, err := i.services.RedisRepo.GetBrandsWithDetails()
	if err == nil && brands != nil {
		return brands, nil
	}
	if err != nil {
		log.Printf("Warning: Redis is not available, falling back to database: %v", err)
	}

	brands, err = i.services.BrandRepo.GetAllWithDetails()
	if err != nil {
		return nil, fmt.Errorf("failed to get brands from database: %v", err)
	}

	go func(data []BrandWithDetails) {
		if err := i.services.RedisRepo.SetBrandsWithDetails(data); err != nil {
			log.Printf("Warning: failed to cache brands: %v", err)
		}
	}(brands)

	return brands, nil
}

func (i *Interactor) BrandCache() {
	brands, err := i.services.BrandRepo.GetAllWithDetails()
	if err != nil {
		log.Printf("Warning: failed to get brands from database: %v", err)
	}

	if err := i.services.RedisRepo.SetBrandsWithDetails(brands); err != nil {
		log.Printf("Warning: failed to cache brands: %v", err)
	}
}
