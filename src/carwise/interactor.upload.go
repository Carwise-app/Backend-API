package carwise

import (
	"errors"
)

func (i *Interactor) UploadImage(request *UploadImageRequest) (*Image, error) {
	if request.File == nil {
		return nil, errors.New("no file provided")
	}

	if request.File.Size > 5*1024*1024 {
		return nil, errors.New("file size exceeds 5MB limit")
	}

	contentType := request.File.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		return nil, errors.New("only jpeg, png and gif images are allowed")
	}

	image, err := i.services.ImageRepo.SaveImage(request.File, request.UserId)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (i *Interactor) DeleteImage(request *DeleteImageRequest) error {
	image, err := i.services.ImageRepo.GetImageById(request.ImageId)
	if err != nil {
		return err
	}

	if request.Role != 2 && image.CreatedBy != request.UserId {
		return errors.New("unauthorized to delete this image")
	}

	return i.services.ImageRepo.DeleteImage(request.ImageId)
}
