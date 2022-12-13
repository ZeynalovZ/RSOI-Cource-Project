package gateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Feokrat/music-dating-app/gateway/internal/config"
	"github.com/Feokrat/music-dating-app/gateway/internal/models"
	"github.com/Feokrat/music-dating-app/gateway/internal/schemas"
	"github.com/google/uuid"
)

type UsersService interface {
	AddUser(request schemas.UserRequest) (uuid.UUID, error)
	UpdateUserInfo(id uuid.UUID, user models.UpdateUserInfo) (int, error)
	GetUserById(id uuid.UUID) (schemas.UserResponse, int, error)
	DeleteUserById(id uuid.UUID) (int, error)
	GetAllUsers(page, size int) (schemas.UsersResponse, int, error)
	GetAllMusics(page, size int) (schemas.MusicsResponse, int, error)
	GetUserRecommendations(id uuid.UUID) (schemas.UsersResponse, int, error)
	GetUserImage(userId uuid.UUID) (schemas.UserImageResponse, int, error)
	LikeUser(whoLikedId uuid.UUID, whomLikedId uuid.UUID) (schemas.LikeResponse, int, error)
	CreateChatForMatch(whoLikedId uuid.UUID, whomLikedId uuid.UUID) (uuid.UUID, int, error)
}

type usersService struct {
	config config.ServicesConfig
	client *http.Client
	logger *log.Logger
}

func NewUsersService(cfg config.ServicesConfig, logger *log.Logger) UsersService {
	return usersService{cfg, http.DefaultClient, logger}
}

func (s usersService) CreateChatForMatch(whoLikedId uuid.UUID, whomLikedId uuid.UUID) (uuid.UUID, int, error) {
	likeUserUrl := s.config.NotificationService + "/api/v1/chats" + fmt.Sprintf("?user_id1=%v&user_id2=%v", whoLikedId, whomLikedId)
	s.logger.Print(likeUserUrl)
	req, err := http.NewRequest("POST", likeUserUrl, nil)
	if err != nil {
		s.logger.Printf("could not create request, error: %s", err.Error())
		return uuid.UUID{}, 0, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Printf("could not get like info, error: %s", err.Error())
		return uuid.UUID{}, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Printf("could not read response body, error: %s",
			err.Error())
		return uuid.UUID{}, 0, err
	}

	var id uuid.UUID
	err = json.Unmarshal(body, &id)
	if err != nil {
		s.logger.Printf("could not unmarshal response body, error: %s", err.Error())
		return uuid.UUID{}, 0, err
	}

	return id, resp.StatusCode, nil
}

func (s usersService) LikeUser(whoLikedId uuid.UUID, whomLikedId uuid.UUID) (schemas.LikeResponse, int, error) {
	likeUserUrl := s.config.UserService + "/api/v1/users/" + fmt.Sprintf("like/%v?liked=%v", whoLikedId, whomLikedId)
	s.logger.Print(likeUserUrl)
	req, err := http.NewRequest("POST", likeUserUrl, nil)
	if err != nil {
		s.logger.Printf("could not create request, error: %s", err.Error())
		return schemas.LikeResponse{}, 0, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Printf("could not get like info, error: %s", err.Error())
		return schemas.LikeResponse{}, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Printf("could not read response body, error: %s",
			err.Error())
		return schemas.LikeResponse{}, 0, err
	}

	var like schemas.LikeResponse
	err = json.Unmarshal(body, &like)
	if err != nil {
		s.logger.Printf("could not unmarshal response body, error: %s", err.Error())
		return schemas.LikeResponse{}, 0, err
	}

	return like, resp.StatusCode, nil
}

func (s usersService) UpdateUserInfo(id uuid.UUID, user models.UpdateUserInfo) (int, error) {
	userServiceUrl := s.config.UserService + fmt.Sprintf("/api/v1/users/%v", id)
	s.logger.Print(userServiceUrl)

	var userBytes bytes.Buffer
	err := json.NewEncoder(&userBytes).Encode(user)
	if err != nil {
		s.logger.Printf("could not convert to io read user, error: %s", err.Error())
		return 0, err
	}

	req, err := http.NewRequest("PUT", userServiceUrl, &userBytes)
	if err != nil {
		s.logger.Printf("could not create request, error: %s", err.Error())
		return 0, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Printf("could not get user info, error: %s", err.Error())
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, errors.New("error occurred during updating user info")
	}

	return 0, nil
}

func (s usersService) GetUserImage(userId uuid.UUID) (schemas.UserImageResponse, int, error) {
	getUserImageByidUrl := s.config.UserService + "/api/v1/users" + fmt.Sprintf("/%v/image", userId)
	s.logger.Print(getUserImageByidUrl)
	req, err := http.NewRequest("GET", getUserImageByidUrl, nil)
	if err != nil {
		s.logger.Printf("could not create request, error: %s", err.Error())
		return schemas.UserImageResponse{}, 0, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Printf("could not get user image info, error: %s", err.Error())
		return schemas.UserImageResponse{}, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Printf("could not read response body, error: %s",
			err.Error())
		return schemas.UserImageResponse{}, 0, err
	}

	var userImage schemas.UserImageResponse
	err = json.Unmarshal(body, &userImage)
	if err != nil {
		s.logger.Printf("could not unmarshal response body, error: %s", err.Error())
		return schemas.UserImageResponse{}, 0, err
	}

	return userImage, resp.StatusCode, nil
}

func (s usersService) AddUser(user schemas.UserRequest) (uuid.UUID, error) {
	userServiceUrl := s.config.UserService + fmt.Sprintf("/api/v1/users")
	s.logger.Print(userServiceUrl)

	var userBytes bytes.Buffer
	err := json.NewEncoder(&userBytes).Encode(user)
	if err != nil {
		s.logger.Printf("could not convert to io read user, error: %s", err.Error())
		return uuid.UUID{}, err
	}

	req, err := http.NewRequest("POST", userServiceUrl, &userBytes)
	if err != nil {
		s.logger.Printf("could not create request, error: %s", err.Error())
		return uuid.UUID{}, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Printf("could not get user info, error: %s", err.Error())
		return uuid.UUID{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Printf("could not read response body, error: %s",
			err.Error())
		return uuid.UUID{}, err
	}

	var id uuid.UUID
	err = json.Unmarshal(body, &id)
	if err != nil {
		s.logger.Printf("could not unmarshal response body, error: %s", err.Error())
		return uuid.UUID{}, err
	}

	return id, nil
}

func (s usersService) GetUserById(id uuid.UUID) (schemas.UserResponse, int, error) {
	getUserByidUrl := s.config.UserService + "/api/v1/users" + fmt.Sprintf("/%v", id)
	s.logger.Print(getUserByidUrl)
	req, err := http.NewRequest("GET", getUserByidUrl, nil)
	if err != nil {
		s.logger.Printf("could not create request, error: %s", err.Error())
		return schemas.UserResponse{}, 0, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Printf("could not get user info, error: %s", err.Error())
		return schemas.UserResponse{}, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Printf("could not read response body, error: %s",
			err.Error())
		return schemas.UserResponse{}, 0, err
	}

	var user schemas.UserResponse
	err = json.Unmarshal(body, &user)
	if err != nil {
		s.logger.Printf("could not unmarshal response body, error: %s", err.Error())
		return schemas.UserResponse{}, 0, err
	}

	return user, resp.StatusCode, nil
}

func (s usersService) DeleteUserById(id uuid.UUID) (int, error) {
	deleteUserByIdUrl := s.config.UserService + fmt.Sprintf("/%v", id)
	req, err := http.NewRequest("DELETE", deleteUserByIdUrl, nil)
	if err != nil {
		s.logger.Printf("could not create request, error: %s", err.Error())
		return 0, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Printf("could not get user info, error: %s", err.Error())
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func (s usersService) GetAllUsers(page, size int) (schemas.UsersResponse, int, error) {
	getUsersUrl := s.config.UserService + fmt.Sprintf("/list?page=%v&size=%v", page, size)
	req, err := http.NewRequest("GET", getUsersUrl, nil)
	if err != nil {
		s.logger.Printf("could not create request, error: %s", err.Error())
		return schemas.UsersResponse{}, 0, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Printf("could not get users, error: %s", err.Error())
		return schemas.UsersResponse{}, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Printf("could not read response body, error: %s",
			err.Error())
		return schemas.UsersResponse{}, 0, err
	}

	var users schemas.UsersResponse
	err = json.Unmarshal(body, &users)
	if err != nil {
		s.logger.Printf("could not unmarshal response body, error: %s", err.Error())
		return schemas.UsersResponse{}, 0, err
	}

	return users, resp.StatusCode, nil
}

func (s usersService) GetAllMusics(page, size int) (schemas.MusicsResponse, int, error) {
	getMusicsUrl := s.config.MusicService + fmt.Sprintf("?page=%v&size=%v", page, size)
	req, err := http.NewRequest("GET", getMusicsUrl, nil)
	if err != nil {
		s.logger.Printf("could not create request, error: %s", err.Error())
		return schemas.MusicsResponse{}, 0, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Printf("could not get musics, error: %s", err.Error())
		return schemas.MusicsResponse{}, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Printf("could not read response body, error: %s",
			err.Error())
		return schemas.MusicsResponse{}, 0, err
	}

	var musics schemas.MusicsResponse
	err = json.Unmarshal(body, &musics)
	if err != nil {
		s.logger.Printf("could not unmarshal response body, error: %s", err.Error())
		return schemas.MusicsResponse{}, 0, err
	}

	return musics, resp.StatusCode, nil
}

func (s usersService) GetUserRecommendations(id uuid.UUID) (schemas.UsersResponse, int, error) {
	getRecommendationsUrl := s.config.UserService + "/api/v1/users" + fmt.Sprintf("/recommendation-list/%v", id)
	req, err := http.NewRequest("GET", getRecommendationsUrl, nil)
	if err != nil {
		s.logger.Printf("could not create request, error: %s", err.Error())
		return schemas.UsersResponse{}, 0, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Printf("could not get recommendations, error: %s", err.Error())
		return schemas.UsersResponse{}, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Printf("could not read response body, error: %s",
			err.Error())
		return schemas.UsersResponse{}, 0, err
	}

	var users schemas.UsersResponse
	err = json.Unmarshal(body, &users)
	if err != nil {
		s.logger.Printf("could not unmarshal response body, error: %s", err.Error())
		return schemas.UsersResponse{}, 0, err
	}

	return users, resp.StatusCode, nil
}
