package schemas

import (
	"github.com/Feokrat/music-dating-app/gateway/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type Error500response struct {
	Message string
	Code    int
}

type Success200response struct {
	Message string
	Code    int
}

type dataResponse struct {
	Data interface{} `json:"data"`
}

type UserChatResponse struct {
	Id    uuid.UUID `json:"id"`
	Image string    `json:"image"`
	Name  string    `json:"name"`
}

type ChatResponse struct {
	Id          uuid.UUID        `json:"id"`
	User        UserChatResponse `json:"user"`
	LastMessage string           `json:"lastMessage"`
	IsRead      bool             `json:"isRead"`
}

type LikeResponse struct {
	IsMatch bool
}

type ChatsNotiResponse struct {
	Chats []ChatsNotiModel
}

type ChatsNotiModel struct {
	Id          uuid.UUID `json:"id"`
	LastMessage string    `json:"lastMessage"`
	IsRead      bool      `json:"isRead"`
	UserId1     uuid.UUID `json:"UserId1"`
	UserId2     uuid.UUID `json:"UserId2"`
}

type ChatsResponse struct {
	Chats []ChatsModel `json:"chats"`
}

type MessageNotiModel struct {
	Id            uuid.UUID `json:"id"`
	CreatorUserId uuid.UUID `json:"creator_user_id"`
	ChatId        uuid.UUID `json:"chat_id"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"created_at"`
	ParentMessage uuid.UUID `json:"parent_message"`
}

type MessageNotiResponse struct {
	Messages []MessageNotiModel `json:"messages"`
}

type MessageFrontRequest struct {
	UserId  uuid.UUID `json:"userId"`
	ChatId  uuid.UUID `json:"chatId"`
	Message string    `json:"message"`
}

type MessageModel struct {
	UserId uuid.UUID `json:"userId"`
	Text   string    `json:"text"`
}

type MessageRequest struct {
	UserId  uuid.UUID `json:"user_id"`
	ChatId  uuid.UUID `json:"chat_id"`
	Message string    `json:"message"`
}

type MessageResponse struct {
	Messages  []MessageModel `json:"messages"`
	CreatedAt string         `json:"createdAt"`
}

type ChatsUser struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Image string    `json:"image"`
}

type ChatsModel struct {
	Id          uuid.UUID `json:"id"`
	User        ChatsUser `json:"user"`
	LastMessage string    `json:"lastMessage"`
	IsRead      bool      `json:"isRead"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type idResponse struct {
	ID interface{} `json:"id"`
}

type IdResponse struct {
	ID uuid.UUID `json:"id"`
}

type messageResponse struct {
	Message string `json:"message"`
}

type UserRequest struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	HasAccess   bool   `json:"hasAccess"`
}

func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, messageResponse{message})
}

func RespondWithData(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, dataResponse{data})
}

func RespondWithId(c *gin.Context, statusCode int, id int) {
	c.JSON(statusCode, idResponse{id})
}

func RespondWithToken(c *gin.Context, statusCode int, token string) {
	c.JSON(statusCode, idResponse{token})
}

type ValidationErrorResponse struct {
	Message string `json:"message"`
	Errors  string `json:"errors"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Surname          string    `json:"surname"`
	Description      string    `json:"description"`
	Image            string    `json:"image"`
	SubscriptionType int       `json:"subscriptionType"`
	MusicIds         []string  `json:"musicIds"`
}

type UserResponseDto struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Surname          string    `json:"surname"`
	Description      string    `json:"description"`
	Image            string    `json:"image"`
	SubscriptionType int       `json:"subscriptionType"`
	MusicIds         []string  `json:"musicIds"`
}

type UsersResponse struct {
	Users []UserResponse `json:"users"`
}

type MusicsResponse struct {
	Musics []models.Music `json:"musics"`
}

type UserImageResponse struct {
	Image string `json:"image"`
}
