package music

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

	rg.POST("/", h.addMusic)
	rg.GET("/", h.getAllMusic)
	rg.GET("/:id", h.getMusicById)
	rg.DELETE("/:id", h.deleteMusicById)
}

func (h handler) addMusic(ctx *gin.Context) {
	var requestModel schemas.MusicRequest
	if err := ctx.BindJSON(&requestModel); err != nil {
		h.logger.Printf("request body in wrong format, error: %s",
			err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong request model",
			Errors:  err.Error(),
		})
		return
	}

	id, err := h.service.AddMusic(models.Music{
		Id:     uuid.New(),
		Name:   requestModel.Name,
		Author: requestModel.Author,
		Url:    requestModel.Url,
	})

	if err != nil {
		h.logger.Printf("could not add music %v, error: %s",
			requestModel, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	ctx.Header("Location", fmt.Sprintf("/api/v1/musics/%v", id))
	ctx.JSON(http.StatusCreated, "")
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

	musics, err := h.service.GetAllMusics(page, size)
	if err != nil {
		h.logger.Printf("error while handling get all musics, error: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: fmt.Sprintf("internal error: %s", err.Error()),
		})

		return
	}

	ctx.JSON(http.StatusOK, musics)
}

func (h handler) getMusicById(ctx *gin.Context) {
	musicIdStr := ctx.Param("id")
	musicId, err := uuid.Parse(musicIdStr)
	if err != nil {
		h.logger.Printf("could not parse music id %v, error: %s",
			musicIdStr, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong music id format",
			Errors:  err.Error(),
		})

		return
	}

	music, err := h.service.GetMusicById(musicId)
	if err != nil {
		h.logger.Printf("could not get music %v, error: %s",
			musicId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, schemas.MusicResponse{Music: music})
}

func (h handler) deleteMusicById(ctx *gin.Context) {
	musicIdStr := ctx.Param("id")
	musicId, err := uuid.Parse(musicIdStr)
	if err != nil {
		h.logger.Printf("could not parse music id %v, error: %s",
			musicIdStr, err.Error())
		ctx.JSON(http.StatusBadRequest, schemas.ValidationErrorResponse{
			Message: "wrong music id format",
			Errors:  err.Error(),
		})

		return
	}

	err = h.service.DeleteMusicById(musicId)
	if err != nil {
		h.logger.Printf("could not delete music %v, error: %s",
			musicId, err.Error())
		ctx.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusNoContent, "")
}
