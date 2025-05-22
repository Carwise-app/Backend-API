package carwise

import (
	"errors"

	"github.com/google/uuid"
)

func (i *Interactor) GetAllBrands() ([]Brand, error) {
	return i.services.BrandRepo.GetAll()
}

func (i *Interactor) CreateBrand(request *BrandCreateRequest) error {
	if request.Role != 2 {
		return errors.New("unauthorized")
	}

	brand := &Brand{
		Id:      uuid.New().String(),
		ImageId: request.ImageId,
		Name:    request.Name,
	}

	return i.services.BrandRepo.Create(brand)
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

func (i *Interactor) GetAllSeriesByBrandId(brandId string) ([]Series, error) {
	return i.services.BrandRepo.GetAllSeriesByBrandId(brandId)
}

func (i *Interactor) CreateSeries(request *SeriesCreateRequest) error {
	if request.Role != 2 {
		return errors.New("unauthorized")
	}

	_, err := i.services.BrandRepo.GetById(request.BrandId)
	if err != nil {
		return err
	}

	series := &Series{
		Id:      uuid.New().String(),
		BrandId: request.BrandId,
		Name:    request.Name,
	}

	return i.services.BrandRepo.CreateSeries(series)
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

func (i *Interactor) GetAllModelsBySeriesId(seriesId string) ([]Model, error) {
	return i.services.BrandRepo.GetAllModelsBySeriesId(seriesId)
}

func (i *Interactor) CreateModel(request *ModelCreateRequest) error {
	if request.Role != 2 {
		return errors.New("unauthorized")
	}

	_, err := i.services.BrandRepo.GetById(request.BrandId)
	if err != nil {
		return err
	}

	_, err = i.services.BrandRepo.GetSeriesById(request.SeriesId)
	if err != nil {
		return err
	}

	model := &Model{
		Id:       uuid.New().String(),
		BrandId:  request.BrandId,
		SeriesId: request.SeriesId,
		Name:     request.Name,
	}

	return i.services.BrandRepo.CreateModel(model)
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
	return i.services.BrandRepo.GetAllWithDetails()
}
