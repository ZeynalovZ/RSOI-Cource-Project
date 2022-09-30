package notifications

import (
	"github.com/Feokrat/music-dating-app/notifications/internal/models"
	"github.com/google/uuid"
	"log"
)

type service struct {
	_chatRepository            ChatRepository
	_messageRepository         MessageRepository
	_messageStatusesRepository MessageStatusesRepository
	logger                     *log.Logger
}

type Service interface {
	GetAllChats(userId uuid.UUID) ([]models.Chats, error)
	GetAllMessages(chatId uuid.UUID) ([]models.Messages, error)
	CreateMessage(chatId uuid.UUID, userId uuid.UUID, message string) (uuid.UUID, error)
	CreateChat(userId1 uuid.UUID, userId2 uuid.UUID) (uuid.UUID, error)
}

func NewChatService(logger *log.Logger, chatr ChatRepository, messager MessageRepository, messagesr MessageStatusesRepository) Service {
	return service{chatr,
		messager,
		messagesr,
		logger}
}

func (s service) GetAllChats(userId uuid.UUID) ([]models.Chats, error) {
	chats, err := s._chatRepository.GetAllChatsByUserId(userId)
	if err != nil {
		s.logger.Printf("Error occured during getting chats for user %v", userId)
		return nil, err
	}

	return chats, err
}

func (s service) GetAllMessages(chatId uuid.UUID) ([]models.Messages, error) {
	return s._messageRepository.GetAllMessages(chatId)
}

func (s service) CreateMessage(chatId uuid.UUID, userId uuid.UUID, message string) (uuid.UUID, error) {
	return s._messageRepository.CreateMessage(message, chatId, userId)
}

func (s service) CreateChat(userId1 uuid.UUID, userId2 uuid.UUID) (uuid.UUID, error) {
	return s._chatRepository.CreateChat(userId1, userId2)
}
