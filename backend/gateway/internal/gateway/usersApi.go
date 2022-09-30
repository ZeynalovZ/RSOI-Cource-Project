package gateway

import (
	"fmt"
	"github.com/Feokrat/music-dating-app/gateway/internal/TokenValidator"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Feokrat/music-dating-app/gateway/internal/schemas"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handler struct {
	service           UsersService
	logger            *log.Logger
	validationService TokenValidator.ValidationService
}

func RegisterUsersHandlers(rg *gin.RouterGroup, service UsersService, validationService TokenValidator.ValidationService, logger *log.Logger) {
	h := handler{service, logger, validationService}

	rg.PUT("/users/:id", h.updateUserById)
	rg.GET("/users", h.getUserById)
	rg.DELETE("/users/:id", h.deleteUserById)
	rg.GET("/users/list", h.getAllUsers)
	rg.GET("/musics", h.getAllMusic)
	rg.GET("recommendation-list", h.getUserRecommendations)
	rg.POST("/users/like/:id", h.LikeUser)
	rg.POST("/users/dislike", h.DislikeUser)
}

func (h handler) LikeUser(ctx *gin.Context) {

	reqToken := ctx.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	userId, err := h.validationService.Validate(reqToken)

	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	likedUserIdStr := ctx.Param("id")
	likedId, err := uuid.Parse(likedUserIdStr)
	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			likedUserIdStr, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	liked, code, err := h.service.LikeUser(userId, likedId)
	if err != nil {
		h.logger.Printf("could not create user %v like for, error: %s",
			userId, likedId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	if liked.IsMatch == true {
		chatId1, code, err := h.service.CreateChatForMatch(userId, likedId)
		if err != nil {
			h.logger.Printf("Error occured during creating new chat %v", code)
		}
		h.logger.Printf("new chats %v", chatId1)
	}

	ctx.JSON(code, liked)

}

func (h handler) DislikeUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, schemas.IdResponse{})
}

func (h handler) createEmptyUser(ctx *gin.Context) {
	var requestModel = schemas.UserRequest{}
	id, err := h.service.AddUser(requestModel)
	if err != nil {
		h.logger.Printf("could not add user %v, error: %s",
			requestModel, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	ctx.Header("Location", fmt.Sprintf("/api/v1/users/%v", id))
	ctx.JSON(http.StatusCreated, id)
}

func (h handler) updateUserById(ctx *gin.Context) {
	var requestModel schemas.UserRequest
	if err := ctx.BindJSON(&requestModel); err != nil {
		h.logger.Printf("request body in wrong format, error: %s",
			err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong request model",
			Errors:  err.Error(),
		})
		return
	}

	//id, err := h.service.AddUser(models.User{
	//	Id:          uuid.New(),
	//	Name:        requestModel.Name,
	//	Surname:     requestModel.Surname,
	//	Email:       requestModel.Email,
	//	PhoneNumber: requestModel.PhoneNumber,
	//	HasAccess:   requestModel.HasAccess,
	//})
	id := 1

	//if err != nil {
	//	h.logger.Printf("could not add user %v, error: %s",
	//		requestModel, err.Error())
	//	ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
	//		Message: err.Error(),
	//	})
	//	return
	//}

	ctx.Header("Location", fmt.Sprintf("/api/v1/users/%v", id))
	ctx.JSON(http.StatusCreated, "")
}

func (h handler) getUserById(ctx *gin.Context) {
	reqToken := ctx.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	userId, err := h.validationService.Validate(reqToken)

	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	user, code, err := h.service.GetUserById(userId)
	if err != nil {
		h.logger.Printf("could not get user %v, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(code, user)
}

func (h handler) deleteUserById(ctx *gin.Context) {
	userIdStr := ctx.Param("id")
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

	code, err := h.service.DeleteUserById(userId)
	if err != nil {
		h.logger.Printf("could not get user %v, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(code, "")
}

func (h handler) getAllUsers(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		h.logger.Printf("could not convert page param to int")
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "page param is not int",
		})

		return
	}

	sizeStr := ctx.Query("size")
	if sizeStr == "" {
		sizeStr = "10000"
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		h.logger.Printf("could not convert size param to int")
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "size param is not int",
		})

		return
	}

	users, code, err := h.service.GetAllUsers(page, size)
	if err != nil {
		h.logger.Printf("could not get users, error: %s",
			err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(code, users)
}

func (h handler) getAllMusic(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		h.logger.Printf("could not convert page param to int")
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "page param is not int",
		})

		return
	}

	sizeStr := ctx.Query("size")
	if sizeStr == "" {
		sizeStr = "10000"
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		h.logger.Printf("could not convert size param to int")
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "size param is not int",
		})

		return
	}

	musics, code, err := h.service.GetAllMusics(page, size)
	if err != nil {
		h.logger.Printf("could not get musics, error: %s",
			err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(code, musics)
}

func (h handler) getUserRecommendations(ctx *gin.Context) {

	reqToken := ctx.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	userId, err := h.validationService.Validate(reqToken)

	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	users, code, err := h.service.GetUserRecommendations(userId)
	if err != nil {
		h.logger.Printf("could not get recommendation list for id %v, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: fmt.Sprintf("could not get recommendation list for id %v", userId),
			Errors:  err.Error(),
		})

		return
	}

	ctx.JSON(code, users)
}
