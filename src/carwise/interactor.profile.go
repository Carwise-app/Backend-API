package carwise

import (
	"fmt"
	"mime/multipart"
)

func (i *Interactor) GetProfile(id string) (*ProfileResponse, []string) {
	user, err := i.services.UserRepo.GetByID(id)
	if err != nil {
		return nil, []string{err.Error()}
	}
	return &ProfileResponse{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		ImageUrl:    user.ImageUrl,
		CountryCode: user.CountryCode,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        user.Role,
		Status:      user.Status,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (i *Interactor) EditProfile(userId string, request ProfileEditRequest, avatar *multipart.FileHeader) []string {
	var errors []string
	user, err := i.services.UserRepo.GetByID(userId)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Failed to get user: %v", err))
		return errors
	}

	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.CountryCode = request.CountryCode
	user.PhoneNumber = request.PhoneNumber

	err = i.services.UserRepo.Update(user)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Failed to update user profile: %v", err))
		return errors
	}

	if avatar != nil {
		image, err := i.services.ImageRepo.SaveImage(avatar, userId)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to upload avatar: %v", err))
			return errors
		}

		user.ImageUrl = image.Path
		err = i.services.UserRepo.Update(user)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to update user avatar URL: %v", err))
			return errors
		}
	}

	return nil
}
