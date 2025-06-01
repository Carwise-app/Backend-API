package carwise

import (
	"log"
	"time"

	"github.com/google/uuid"
)

func (i *Interactor) SendMessage(request *SendMessageRequest) error {
	listing, err := i.services.ListingRepo.GetListingById(request.ListingId)
	if err != nil {
		return err
	}

	if listing.CreatedBy != request.UserId {
		err = i.CreatePushNotification(
			ChatMessage,
			"",
			"",
			map[string]string{
				"listing_id": listing.Id,
				"user_id":    request.UserId,
			},
			listing.CreatedBy,
		)
		if err != nil {
			log.Println("CreatePushNotification error: ", err)
		}
	}

	message := &Message{
		Id:         uuid.New().String(),
		SenderId:   request.UserId,
		ReceiverId: request.ReceiverId,
		Message:    request.Message,
		Read:       false,
		CreatedAt:  time.Now().Unix(),
		ListingId:  listing.Id,
	}

	return i.services.MessageRepo.SaveMessage(message)
}

func (i *Interactor) GetMessages(request *GetMessagesRequest) (*GetMessagesResponse, error) {
	messages, err := i.services.MessageRepo.GetMessagesByUserId(request.UserId, request.ReceiverId, request.ListingId, request.Limit, request.Page)
	if err != nil {
		return nil, err
	}

	total, err := i.services.MessageRepo.CountMessagesByUserId(request.UserId, request.ReceiverId, request.ListingId)
	if err != nil {
		return nil, err
	}

	var MessageInfos []MessageInfo
	for _, message := range messages {
		sender, err := i.services.UserRepo.GetByID(message.SenderId)
		if err != nil {
			return nil, err
		}

		receiver, err := i.services.UserRepo.GetByID(message.ReceiverId)
		if err != nil {
			return nil, err
		}

		if sender.Id != request.UserId {
			if !message.Read {
				i.services.MessageRepo.ReadMessage(message.Id)
				message.Read = true
			}
		}

		MessageInfos = append(MessageInfos, MessageInfo{
			Id: message.Id,
			Sender: UserInfo{
				Id:          sender.Id,
				FirstName:   sender.FirstName,
				LastName:    sender.LastName,
				Email:       sender.Email,
				CountryCode: sender.CountryCode,
				PhoneNumber: sender.PhoneNumber,
			},
			Receiver: UserInfo{
				Id:          receiver.Id,
				FirstName:   receiver.FirstName,
				LastName:    receiver.LastName,
				Email:       receiver.Email,
				CountryCode: receiver.CountryCode,
				PhoneNumber: receiver.PhoneNumber,
			},
			Message:   message.Message,
			Read:      message.Read,
			CreatedAt: message.CreatedAt,
		})
	}

	return &GetMessagesResponse{
		Messages: MessageInfos,
		Total:    total,
	}, nil
}

func (i *Interactor) GetChats(request *GetChatsRequest) (*GetChatsResponse, error) {
	chats, err := i.services.MessageRepo.GetChats(request.UserId, request.Limit, request.Page)
	if err != nil {
		return nil, err
	}

	total, err := i.services.MessageRepo.CountChats(request.UserId)
	if err != nil {
		return nil, err
	}

	log.Println(chats)

	chatInfos := make([]ChatInfo, 0, len(chats))
	for _, chat := range chats {
		user, err := i.services.UserRepo.GetByID(chat.OtherUserId)
		if err != nil {
			return nil, err
		}

		listing, err := i.services.ListingRepo.GetListingById(chat.ListingId)
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
		image := Image{}

		if len(listing.Images) > 0 {
			firstImage, err := i.services.ImageRepo.GetImageById(listing.Images[0])
			if err == nil {
				image = *firstImage
			}
		}

		isFavorite := i.services.FavoriteRepo.IsFavorite(request.UserId, listing.Id)

		chatInfos = append(chatInfos, ChatInfo{
			User: UserInfo{
				Id:          user.Id,
				FirstName:   user.FirstName,
				LastName:    user.LastName,
				Email:       user.Email,
				CountryCode: user.CountryCode,
				PhoneNumber: user.PhoneNumber,
			},
			Listing: ListListingInfo{
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
			},
			LastMessageTime: chat.LastMessageTime,
		})
	}

	return &GetChatsResponse{
		Chats: chatInfos,
		Total: total,
	}, nil
}
