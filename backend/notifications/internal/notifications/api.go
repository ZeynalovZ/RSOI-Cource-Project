package notifications

import (
	"github.com/Feokrat/music-dating-app/notifications/internal/schemas"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type handler struct {
	s      Service
	logger *log.Logger
}

func RegisterHandlers(rg *gin.RouterGroup, service Service, logger *log.Logger) {
	h := handler{logger: logger, s: service}

	rg.GET("/chats/:user_id", h.GetChatsByUserID)
	rg.GET("/messages/chat/:id", h.GetMessagesByChatId)
	rg.POST("/chats", h.CreateChat)
	rg.POST("/messages", h.CreateMessage)
}

func (h handler) GetChatsByUserID(ctx *gin.Context) {
	userIdStr := ctx.Param("user_id")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userIdStr, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	chats, err := h.s.GetAllChats(userId)
	if err != nil {
		h.logger.Printf("could not get chats for user %v, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	if len(chats) == 0 {
		h.logger.Printf("chats weren't found %v",
			userId)
		ctx.JSON(http.StatusNotFound, schemas.ErrorResponse{
			Message: "Chats weren't found",
		})
		return
	}

	var chatsInfo schemas.ChatsResponse

	for i := 0; i < len(chats); i++ {
		var chat schemas.ChatsModel
		chat.Id = chats[i].Id
		chat.IsRead = true
		chat.UserId1 = chats[i].UserId1
		chat.UserId2 = chats[i].UserId2
		messages, error := h.s.GetAllMessages(chats[i].Id)
		if error != nil {
			h.logger.Printf("Error occured during getting messages of chat %v", chats[i].Id)
		} else {
			if len(messages) != 0 {
				chat.LastMessage = messages[len(messages)-1].Content
			}

		}
		chatsInfo.Chats = append(chatsInfo.Chats, chat)
	}

	ctx.JSON(http.StatusOK, chatsInfo)
}

func (h handler) CreateChat(ctx *gin.Context) {
	userIdStr1 := ctx.Query("user_id1")
	userId1, err := uuid.Parse(userIdStr1)
	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userIdStr1, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	userIdStr2 := ctx.Query("user_id2")
	userId2, err := uuid.Parse(userIdStr2)
	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userIdStr2, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	chatId, err := h.s.CreateChat(userId1, userId2)
	if err != nil {
		h.logger.Printf("could not create chats for user %v and %v, error: %s",
			userId1, userId2, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, chatId)
}

func (h handler) GetMessagesByChatId(ctx *gin.Context) {
	chatIdStr := ctx.Param("id")
	chatId, err := uuid.Parse(chatIdStr)
	if err != nil {
		h.logger.Printf("could not parse chat id %v, error: %s",
			chatIdStr, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong char id format",
			Errors:  err.Error(),
		})

		return
	}

	messages, err := h.s.GetAllMessages(chatId)
	if err != nil {
		h.logger.Printf("could not get chats for user %v, error: %s",
			chatId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, schemas.MessageResponse{
		Messages: messages,
	})
}

func (h handler) CreateMessage(ctx *gin.Context) {
	var messageRequest schemas.MessageRequest
	if err := ctx.BindJSON(&messageRequest); err != nil {
		h.logger.Printf("request body in wrong format, error: %s",
			err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong request model",
			Errors:  err.Error(),
		})
		return
	}

	messageId, err := h.s.CreateMessage(messageRequest.ChatId, messageRequest.UserId, messageRequest.Message)
	if err != nil {
		h.logger.Printf("could not get chats for user %v and chat %v, error: %s",
			messageRequest.UserId, messageRequest.ChatId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, messageId)
}
