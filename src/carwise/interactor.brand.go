package carwise

import (
	"errors"

	"github.com/google/uuid"
)

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
