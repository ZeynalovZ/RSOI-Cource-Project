package notifications

import (
	"github.com/Feokrat/music-dating-app/gateway/internal/TokenValidator"
	"github.com/Feokrat/music-dating-app/gateway/internal/gateway"
	"github.com/Feokrat/music-dating-app/gateway/internal/schemas"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
)

type handler struct {
	service     NotificationService
	validator   TokenValidator.ValidationService
	userService gateway.UsersService
	logger      *log.Logger
}

func RegisterChatHandlers(rg *gin.RouterGroup, service NotificationService,
	validator TokenValidator.ValidationService, usersService gateway.UsersService, logger *log.Logger) {
	h := handler{service, validator, usersService, logger}

	rg.GET("/", h.GetAllChats)
	rg.GET("/:id", h.GetChatById)
	rg.POST("/sendMessage", h.CreateMessageInChat)
}

func (h handler) GetAllChats(ctx *gin.Context) {
	reqToken := ctx.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	userId, err := h.validator.Validate(reqToken)

	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	chats, code, err := h.service.GetAllChatsByUserId(userId)
	if err != nil {
		h.logger.Printf("Error occurred during getting all chats in gateway. Error: %v", err.Error())
		if code == http.StatusNotFound {
			ctx.JSON(http.StatusNotFound, schemas.ErrorResponse{Message: err.Error()})
		}

		ctx.JSON(http.StatusInternalServerError, schemas.Error500response{Message: err.Error()})
	}

	if code == http.StatusNotFound {
		ctx.JSON(code, schemas.ChatsResponse{})
	}

	var chatsResponse schemas.ChatsResponse
	for i := 0; i < len(chats.Chats); i++ {
		var chatModel schemas.ChatsModel
		chatModel.Id = chats.Chats[i].Id
		chatModel.LastMessage = chats.Chats[i].LastMessage
		chatModel.IsRead = chats.Chats[i].IsRead
		var UserID uuid.UUID
		if chats.Chats[i].UserId2 == userId {
			UserID = chats.Chats[i].UserId1
		} else {
			UserID = chats.Chats[i].UserId2
		}
		user, code, err := h.userService.GetUserById(UserID)
		if err != nil {
			if code == http.StatusNotFound {
				h.logger.Printf("Image for user %v was not found", userId)
			}
			h.logger.Printf("Error occured during getting image for user %v ", userId)
		} else {
			chatModel.User.Image = user.Image
			chatModel.User.Name = user.Name
			chatModel.User.Id = user.Id
		}
		chatsResponse.Chats = append(chatsResponse.Chats, chatModel)
	}

	ctx.JSON(code, chatsResponse)
}

func (h handler) GetChatById(ctx *gin.Context) {
	reqToken := ctx.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	userId, err := h.validator.Validate(reqToken)

	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	chatIdStr := ctx.Param("id")
	chatId, err := uuid.Parse(chatIdStr)
	if err != nil {
		h.logger.Printf("could not parse chat id %v, error: %s",
			chatIdStr, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	messages, code, err := h.service.GetMessagesByChatId(chatId)
	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			chatIdStr, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	if len(messages.Messages) == 0 {
		ctx.JSON(http.StatusOK, schemas.MessageResponse{Messages: []schemas.MessageModel{}})
		return
	}

	var messagesResponse schemas.MessageResponse
	for i := 0; i < len(messages.Messages); i++ {
		var messageModel schemas.MessageModel
		messageModel.UserId = messages.Messages[i].CreatorUserId
		messageModel.Text = messages.Messages[i].Content
		messagesResponse.Messages = append(messagesResponse.Messages, messageModel)
	}
	messagesResponse.CreatedAt = messages.Messages[0].CreatedAt.String()
	ctx.JSON(code, messagesResponse)
}

func (h handler) CreateMessageInChat(ctx *gin.Context) {
	reqToken := ctx.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	userId, err := h.validator.Validate(reqToken)

	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	var messageFrontRequest schemas.MessageFrontRequest
	if err := ctx.BindJSON(&messageFrontRequest); err != nil {
		h.logger.Printf("request body in wrong format, error: %s",
			err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong request model",
			Errors:  err.Error(),
		})
		return
	}
	var messageRequest schemas.MessageRequest
	messageRequest.Message = messageFrontRequest.Message
	messageRequest.ChatId = messageFrontRequest.ChatId
	messageRequest.UserId = userId

	messageId, code, err := h.service.CreateMessageForChat(messageRequest)
	if err != nil {
		h.logger.Printf("Error occurred during getting all chats in gateway. Error: %v", err.Error())
		if code == http.StatusNotFound {
			ctx.JSON(http.StatusNotFound, schemas.ErrorResponse{Message: err.Error()})
		}

		ctx.JSON(http.StatusInternalServerError, schemas.Error500response{Message: err.Error()})
	}

	ctx.JSON(code, messageId)
}
