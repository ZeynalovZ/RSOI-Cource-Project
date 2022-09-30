package user

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Feokrat/music-dating-app/users/internal/models"
	"github.com/Feokrat/music-dating-app/users/internal/schemas"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handler struct {
	service Service
	logger  *log.Logger
}

func RegisterHandlers(rg *gin.RouterGroup, service Service, logger *log.Logger) {
	h := handler{service, logger}

	rg.POST("/", h.addUser)
	rg.GET("/:id", h.getUserById)
	rg.DELETE("/:id", h.deleteUserById)
	rg.POST("/add-music", h.addMusic)
	rg.POST("/add-image", h.addImage)
	rg.GET("/:id/image", h.getImage)
	rg.PUT("/:id", h.updateUserById)
	rg.GET("/list", h.getAllUsers)
	rg.GET("/recommendation-list/:id", h.getUserRecommendations)
	rg.POST("/like/:id", h.likeUser)
}

func (h handler) likeUser(ctx *gin.Context) {
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

	likedIdStr := ctx.Query("liked")
	likedId, err := uuid.Parse(likedIdStr)
	if err != nil {
		h.logger.Printf("could not parse user id %v, error: %s",
			userIdStr, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong user id format",
			Errors:  err.Error(),
		})

		return
	}

	isMatch, err := h.service.LikeUser(userId, likedId)
	if err != nil {
		h.logger.Printf("could not like %v user, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, schemas.LikeResponse{IsMatch: isMatch})
}

func (h handler) getImage(ctx *gin.Context) {
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

	image, err := h.service.GetUserImageById(userId)
	if err != nil {

		if (err == schemas.NotFoundError{}) {
			ctx.JSON(http.StatusNotFound, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		h.logger.Printf("could not get user %v image, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	h.logger.Print(image)
	ctx.JSON(http.StatusOK, schemas.UserImageResponse{Image: image.Image})
}

func (h handler) addUser(ctx *gin.Context) {
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

	id, err := h.service.AddUser(models.User{
		Id:          uuid.New(),
		Name:        requestModel.Name,
		Surname:     requestModel.Surname,
		Email:       requestModel.Email,
		PhoneNumber: requestModel.PhoneNumber,
		HasAccess:   requestModel.HasAccess,
	})

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

func (h handler) getUserById(ctx *gin.Context) {
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

	user, err := h.service.GetUserById(userId)
	if err != nil {
		h.logger.Printf("could not get user %v, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	h.logger.Print(user)
	ctx.JSON(http.StatusOK, user)
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

	err = h.service.DeleteUserById(userId)
	if err != nil {
		h.logger.Printf("could not delete user %v, error: %s",
			userId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusNoContent, "")
}

func (h handler) addMusic(ctx *gin.Context) {
	var requestModel schemas.UserToMusicRequest
	if err := ctx.BindJSON(&requestModel); err != nil {
		h.logger.Printf("request body in wrong format, error: %s",
			err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong request model",
			Errors:  err.Error(),
		})
		return
	}

	err := h.service.AddMusicToUser(models.UserToMusic{
		Id:             uuid.New(),
		UserId:         requestModel.UserId,
		MusicId:        requestModel.MusicId,
		FavouriteLevel: requestModel.FavouriteLevel,
	})

	if err != nil {
		h.logger.Printf("could not add music %v to user %v, error: %s",
			requestModel.MusicId, requestModel.UserId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})
		return
	}
}

func (h handler) addImage(ctx *gin.Context) {
	var requestModel schemas.ImageRequest
	if err := ctx.BindJSON(&requestModel); err != nil {
		h.logger.Printf("request body in wrong format, error: %s",
			err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong request model",
			Errors:  err.Error(),
		})
		return
	}

	err := h.service.AddImageToUser(models.Image{
		Id:     uuid.New(),
		UserId: requestModel.UserId,
		Image:  requestModel.Image,
	})

	if err != nil {
		h.logger.Printf("could not add image to user %v, error: %s",
			requestModel.UserId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})
		return
	}
}

func (h handler) updateUserById(ctx *gin.Context) {
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

	t, err := h.service.GetUserById(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	fmt.Printf("%v", t)

	//if user == (schemas.UserResponse{}) {
	//	ctx.JSON(http.StatusNotFound, schemas.ErrorResponse{
	//		Message: fmt.Sprintf("Not found user with id %d\n", userId),
	//	})
	//}

	var requestModel schemas.UpdateRequest
	if err := ctx.BindJSON(&requestModel); err != nil {
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "Wrong user model.",
			Errors:  err.Error(),
		})
		return
	}

	err = h.service.UpdateUserInfo(userId, models.UpdateUserInfo{
		Name:        requestModel.Name,
		Surname:     requestModel.Surname,
		Email:       requestModel.Email,
		PhoneNumber: requestModel.PhoneNumber,
		HasAccess:   requestModel.HasAccess,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, nil)
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

	users, err := h.service.GetAllUsers(page, size)
	if err != nil {
		h.logger.Printf("error while handling get all users, error: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: fmt.Sprintf("internal error: %s", err.Error()),
		})

		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (h handler) getUserRecommendations(ctx *gin.Context) {
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
	users, err := h.service.GetUserRecommendations(userId, 1, 100)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: fmt.Sprintf("Couldn't get user recommendations for user %v", userId),
			Errors:  err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, users)
}
