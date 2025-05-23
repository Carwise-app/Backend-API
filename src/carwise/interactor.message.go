package carwise

import (
	"log"
	"time"

	"github.com/google/uuid"
)

func (i *Interactor) SendMessage(request *SendMessageRequest) error {
	message := &Message{
		Id:         uuid.New().String(),
		SenderId:   request.UserId,
		ReceiverId: request.ReceiverId,
		Message:    request.Message,
		Read:       false,
		CreatedAt:  time.Now().Unix(),
	}

	return i.services.MessageRepo.SaveMessage(message)
}

func (i *Interactor) GetMessages(request *GetMessagesRequest) (*GetMessagesResponse, error) {
	messages, err := i.services.MessageRepo.GetMessagesByUserId(request.UserId, request.Limit, request.Page)
	if err != nil {
		return nil, err
	}

	total, err := i.services.MessageRepo.CountMessagesByUserId(request.UserId)
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

		chatInfos = append(chatInfos, ChatInfo{
			User: UserInfo{
				Id:          user.Id,
				FirstName:   user.FirstName,
				LastName:    user.LastName,
				Email:       user.Email,
				CountryCode: user.CountryCode,
				PhoneNumber: user.PhoneNumber,
			},
			LastMessageTime: chat.LastMessageTime,
		})
	}

	return &GetChatsResponse{
		Chats: chatInfos,
		Total: total,
	}, nil
}
