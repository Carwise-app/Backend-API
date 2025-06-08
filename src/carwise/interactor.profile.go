package carwise

import (
	"fmt"
	"log"
	"mime/multipart"
	"time"
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

func (i *Interactor) NotifyProfile(request ProfileNotifyRequest) error {
	user, err := i.services.UserRepo.GetByID(request.UserId)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if request.DeviceToken != "" {
		user.DeviceToken = request.DeviceToken
	}

	user.EmailNotify = request.EmailNotify
	user.PushNotify = request.PushNotify
	user.UpdatedAt = time.Now().Unix()

	err = i.services.UserRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update user notify: %w", err)
	}

	return nil
}

func (i *Interactor) GetProfileNotify(request GetProfileNotifyRequest) (*GetProfileNotifyResponse, error) {
	user, err := i.services.UserRepo.GetByID(request.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &GetProfileNotifyResponse{
		DeviceToken: user.DeviceToken,
		EmailNotify: user.EmailNotify,
		PushNotify:  user.PushNotify,
	}, nil
}

func (i *Interactor) DeleteAccount(request DeleteAccountRequest) error {
	if request.Role != 2 || request.UserId != request.ProfileId {
		return fmt.Errorf("you are not authorized to access this resource")
	}

	err := i.services.UserRepo.DeleteUser(request.ProfileId)
	if err != nil {
		log.Println("Error deleting user:", err.Error())
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
