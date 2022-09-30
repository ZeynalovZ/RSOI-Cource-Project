package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Feokrat/music-dating-app/gateway/internal/config"
	"github.com/Feokrat/music-dating-app/gateway/internal/schemas"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

type NotificationService interface {
	GetAllChatsByUserId(userId uuid.UUID) (schemas.ChatsNotiResponse, int, error)
	GetMessagesByChatId(chatId uuid.UUID) (schemas.MessageNotiResponse, int, error)
	CreateMessageForChat(request schemas.MessageRequest) (uuid.UUID, int, error)
}

type notificationService struct {
	config config.ServicesConfig
	client *http.Client
	logger *log.Logger
}

func NewNotificationService(cfg config.ServicesConfig, logger *log.Logger) NotificationService {
	return notificationService{cfg, http.DefaultClient, logger}
}

func (s notificationService) CreateMessageForChat(request schemas.MessageRequest) (uuid.UUID, int, error) {
	messagesUrl := s.config.NotificationService + "/api/v1/messages/"
	s.logger.Print(messagesUrl)

	var messageBytes bytes.Buffer
	err := json.NewEncoder(&messageBytes).Encode(request)
	if err != nil {
		s.logger.Printf("could not convert to io read messages, error: %s", err.Error())
		return uuid.UUID{}, 0, err
	}

	req, err := http.NewRequest("POST", messagesUrl, &messageBytes)
	if err != nil {
		s.logger.Printf("could not create request, error: %s", err.Error())
		return uuid.UUID{}, 0, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Printf("could not get messages info, error: %s", err.Error())
		return uuid.UUID{}, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Printf("could not read response body, error: %s",
			err.Error())
		return uuid.UUID{}, 0, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return uuid.UUID{}, 0, err
	}

	var messageId uuid.UUID
	err = json.Unmarshal(body, &messageId)
	if err != nil {
		s.logger.Printf("could not unmarshal response body, error: %s", err.Error())
		return uuid.UUID{}, 0, err
	}

	return messageId, resp.StatusCode, nil
}

func (s notificationService) GetMessagesByChatId(chatId uuid.UUID) (schemas.MessageNotiResponse, int, error) {
	chatsUrl := s.config.NotificationService + "/api/v1/messages/chat/" + fmt.Sprintf("%v", chatId)
	s.logger.Print(chatsUrl)
	req, err := http.NewRequest("GET", chatsUrl, nil)
	if err != nil {
		s.logger.Printf("could not create request, error: %s", err.Error())
		return schemas.MessageNotiResponse{}, 0, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Printf("could not get chats info, error: %s", err.Error())
		return schemas.MessageNotiResponse{}, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Printf("could not read response body, error: %s",
			err.Error())
		return schemas.MessageNotiResponse{}, 0, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return schemas.MessageNotiResponse{}, resp.StatusCode, nil
	}

	var messages schemas.MessageNotiResponse
	err = json.Unmarshal(body, &messages)
	if err != nil {
		s.logger.Printf("could not unmarshal response body, error: %s", err.Error())
		return schemas.MessageNotiResponse{}, 0, err
	}

	return messages, resp.StatusCode, nil
}

func (s notificationService) GetAllChatsByUserId(userId uuid.UUID) (schemas.ChatsNotiResponse, int, error) {
	chatsUrl := s.config.NotificationService + "/api/v1/chats/" + fmt.Sprintf("%v", userId)
	s.logger.Print(chatsUrl)
	req, err := http.NewRequest("GET", chatsUrl, nil)
	if err != nil {
		s.logger.Printf("could not create request, error: %s", err.Error())
		return schemas.ChatsNotiResponse{}, 0, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Printf("could not get chats info, error: %s", err.Error())
		return schemas.ChatsNotiResponse{}, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Printf("could not read response body, error: %s",
			err.Error())
		return schemas.ChatsNotiResponse{}, 0, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return schemas.ChatsNotiResponse{}, resp.StatusCode, nil
	}

	var chats schemas.ChatsNotiResponse
	err = json.Unmarshal(body, &chats)
	if err != nil {
		s.logger.Printf("could not unmarshal response body, error: %s", err.Error())
		return schemas.ChatsNotiResponse{}, 0, err
	}

	return chats, resp.StatusCode, nil
}
