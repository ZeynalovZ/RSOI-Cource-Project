package user

import (
	"log"

	"github.com/Feokrat/music-dating-app/users/internal/models"
	"github.com/Feokrat/music-dating-app/users/internal/schemas"
	"github.com/google/uuid"
)

type service struct {
	userRepository Repository
	logger         *log.Logger
}

type Service interface {
	AddUser(user models.User) (uuid.UUID, error)
	GetUserById(id uuid.UUID) (schemas.UserResponse, error)
	DeleteUserById(id uuid.UUID) error
	AddMusicToUser(userToMusic models.UserToMusic) error
	AddImageToUser(image models.Image) error
	UpdateUserInfo(id uuid.UUID, user models.UpdateUserInfo) error
	GetAllUsers(page, size int) (schemas.UsersResponse, error)
	GetUserRecommendations(id uuid.UUID, page, size int) (schemas.UsersResponse, error)
	GetUserImageById(id uuid.UUID) (models.Image, error)
	LikeUser(id uuid.UUID, likedId uuid.UUID) (bool, error)
}

func NewService(repo Repository, logger *log.Logger) Service {
	return service{repo, logger}
}

func (s service) LikeUser(id uuid.UUID, likedId uuid.UUID) (bool, error) {
	isMatch, err := s.userRepository.Like(id, likedId)
	return isMatch, err
}

func (s service) AddUser(user models.User) (uuid.UUID, error) {
	id, err := s.userRepository.Create(user)
	return id, err
}

func (s service) GetUserById(id uuid.UUID) (schemas.UserResponse, error) {
	user, err := s.userRepository.GetById(id)
	if err != nil {
		s.logger.Printf("Error occured during getting user")
		return schemas.UserResponse{}, err
	}
	image, err := s.userRepository.GetUserImage(id)
	if err != nil {
		s.logger.Printf("Error occured during getting user image")
		return schemas.UserResponse{Id: user.Id,
			Name:             user.Name,
			Surname:          user.Surname,
			Description:      "Хочу квас",
			SubscriptionType: 0,
			MusicIds:         []string{"365b8b96-3244-486e-934d-9b020fe6ea72"},
		}, nil
	}
	return schemas.UserResponse{Id: user.Id,
		Name:             user.Name,
		Surname:          user.Surname,
		Description:      "Хочу квас",
		SubscriptionType: 0,
		Image:            image.Image,
		MusicIds:         []string{"365b8b96-3244-486e-934d-9b020fe6ea72"},
	}, err
}

func (s service) GetUserImageById(id uuid.UUID) (models.Image, error) {
	image, err := s.userRepository.GetUserImage(id)
	return image, err
}

func (s service) DeleteUserById(id uuid.UUID) error {
	err := s.userRepository.DeleteById(id)
	return err
}

func (s service) AddMusicToUser(userToMusic models.UserToMusic) error {
	err := s.userRepository.CreateMusicToUser(userToMusic)
	return err
}

func (s service) AddImageToUser(image models.Image) error {
	err := s.userRepository.CreateUserImage(image)
	return err
}

func (s service) UpdateUserInfo(id uuid.UUID, user models.UpdateUserInfo) error {
	err := s.userRepository.Update(id, user)
	return err
}

func (s service) GetAllUsers(page, size int) (schemas.UsersResponse, error) {
	users, err := s.userRepository.GetAll(page, size)
	if err != nil {
		s.logger.Printf("Error occured in getting all users")
		return schemas.UsersResponse{}, err
	}
	var usersResponse schemas.UsersResponse
	for i := 0; i < len(users); i++ {
		usersResponse.Users = append(usersResponse.Users, schemas.UserResponse{Id: users[i].Id,
			Name:             users[i].Name,
			Surname:          users[i].Surname,
			Description:      "Хочу квас",
			SubscriptionType: 1,
			MusicIds:         []string{"365b8b96-3244-486e-934d-9b020fe6ea72"},
		})
	}

	return usersResponse, nil
}

func (s service) GetUserRecommendations(id uuid.UUID, page int, size int) (schemas.UsersResponse, error) {
	users, err := s.userRepository.GetRecommendationsForUser(id, page, size)
	if err != nil {
		s.logger.Printf("Error occured during getting recommendations for user %v", id)
		return schemas.UsersResponse{}, err
	}
	var usersResponse schemas.UsersResponse

	for i := 0; i < len(users); i++ {
		image, err := s.userRepository.GetUserImage(users[i].Id)
		if err != nil {
			s.logger.Printf("Error occured during getting image for user %v", users[i].Id)
			usersResponse.Users = append(usersResponse.Users, schemas.UserResponse{Id: users[i].Id,
				Name:             users[i].Name,
				Surname:          users[i].Surname,
				Description:      "Хочу квас",
				SubscriptionType: 0,
				Image:            "",
				MusicIds:         []string{"365b8b96-3244-486e-934d-9b020fe6ea72"},
			})
		} else {
			usersResponse.Users = append(usersResponse.Users, schemas.UserResponse{Id: users[i].Id,
				Name:             users[i].Name,
				Surname:          users[i].Surname,
				Description:      "Хочу квас",
				SubscriptionType: 0,
				Image:            image.Image,
				MusicIds:         []string{"365b8b96-3244-486e-934d-9b020fe6ea72"},
			})
		}
	}

	return usersResponse, nil
}
