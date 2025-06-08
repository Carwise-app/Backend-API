package carwise

import "errors"

type Interactor struct {
	services Services
}

func NewInteractor(svcs Services) *Interactor {
	return &Interactor{
		services: svcs,
	}
}

func (i *Interactor) Count(req *CountRequest) (*CountResponse, error) {
	if req.Role == 2 {
		return nil, errors.New("you are not authorized to access this resource")
	}

	userCount, err := i.services.UserRepo.CountUsers()
	if err != nil {
		return nil, err
	}

	imageCount, err := i.services.ImageRepo.ImagesCount()
	if err != nil {
		return nil, err
	}

	messageCount, err := i.services.MessageRepo.MessagesCount()
	if err != nil {
		return nil, err
	}

	predictCount, err := i.services.PredictRepo.PredictCount()
	if err != nil {
		return nil, err
	}

	notificationCount, err := i.services.NotificationRepo.NotificationsCount()
	if err != nil {
		return nil, err
	}

	favoriteCount, err := i.services.FavoriteRepo.CountFavorites()
	if err != nil {
		return nil, err
	}

	listingCount, err := i.services.ListingRepo.CountListing(&ListingFilter{})
	if err != nil {
		return nil, err
	}

	return &CountResponse{
		UserCount:         userCount,
		ImageCount:        imageCount,
		MessageCount:      messageCount,
		PredictCount:      predictCount,
		NotificationCount: notificationCount,
		FavoriteCount:     favoriteCount,
		ListingCount:      listingCount,
	}, nil
}

func (i *Interactor) GetUsers(req *GetUsersRequest) (*GetUsersResponse, error) {
	if req.Role == 2 {
		return nil, errors.New("you are not authorized to access this resource")
	}

	users, err := i.services.UserRepo.GetUsers(req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	total, err := i.services.UserRepo.CountUsers()
	if err != nil {
		return nil, err
	}

	return &GetUsersResponse{
		Users: users,
		Total: total,
	}, nil
}
