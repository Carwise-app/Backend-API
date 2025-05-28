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
		Id:      uuid.New().String(),
		ImageId: request.ImageId,
		Name:    request.Name,
	}

	err := i.services.BrandRepo.Create(brand)
	if err != nil {
		return nil, err
	}

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
	brand.ImageId = request.ImageId
	brand.Name = request.Name

	return i.services.BrandRepo.Update(brand)
}

func (i *Interactor) DeleteBrand(request *BrandDeleteRequest) error {
	if request.Role != 2 {
		return errors.New("unauthorized")
	}

	return i.services.BrandRepo.Delete(request.BrandId)
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

	return i.services.BrandRepo.UpdateSeries(series)
}

func (i *Interactor) DeleteSeries(request *SeriesDeleteRequest) error {
	if request.Role != 2 {
		return errors.New("unauthorized")
	}

	return i.services.BrandRepo.DeleteSeries(request.SeriesId)
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

	return i.services.BrandRepo.UpdateModel(model)
}

func (i *Interactor) DeleteModel(request *ModelDeleteRequest) error {
	if request.Role != 2 {
		return errors.New("unauthorized")
	}

	return i.services.BrandRepo.DeleteModel(request.ModelId)
}

func (i *Interactor) GetAllBrandsWithDetails() ([]BrandWithDetails, error) {
	brands, err := i.services.RedisRepo.GetBrandsWithDetails()
	if err != nil {
		if err.Error() == "Redis is not available" {
			log.Printf("Warning: Redis is not available, falling back to database: %v", err)
		} else {
			return nil, fmt.Errorf("failed to get brands from cache: %v", err)
		}
	}

	if brands != nil {
		return brands, nil
	}

	brands, err = i.services.BrandRepo.GetAllWithDetails()
	if err != nil {
		return nil, fmt.Errorf("failed to get brands from database: %v", err)
	}

	if err := i.services.RedisRepo.SetBrandsWithDetails(brands); err != nil {
		log.Printf("Warning: failed to cache brands: %v", err)
	}

	return brands, nil
}
