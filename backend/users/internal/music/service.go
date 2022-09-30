package music

import (
	"log"

	"github.com/Feokrat/music-dating-app/users/internal/models"
	"github.com/Feokrat/music-dating-app/users/internal/schemas"
	"github.com/google/uuid"
)

type service struct {
	musicRepository Repository
	logger          *log.Logger
}

type Service interface {
	AddMusic(music models.Music) (uuid.UUID, error)
	GetMusicById(id uuid.UUID) (models.Music, error)
	DeleteMusicById(id uuid.UUID) error
	GetAllMusics(page, size int) (schemas.MusicsResponse, error)
	GetUserRecommendations(id uuid.UUID, page, size int) (schemas.UsersResponse, error)
}

func NewService(repo Repository, logger *log.Logger) Service {
	return service{repo, logger}
}

func (s service) AddMusic(music models.Music) (uuid.UUID, error) {
	id, err := s.musicRepository.Create(music)
	return id, err
}

func (s service) GetMusicById(id uuid.UUID) (models.Music, error) {
	user, err := s.musicRepository.GetById(id)
	return user, err
}

func (s service) DeleteMusicById(id uuid.UUID) error {
	err := s.musicRepository.DeleteById(id)
	return err
}

func (s service) GetAllMusics(page, size int) (schemas.MusicsResponse, error) {
	musics, err := s.musicRepository.GetAll(page, size)
	return schemas.MusicsResponse{Musics: musics}, err
}

func (s service) GetUserRecommendations(id uuid.UUID, page, size int) (schemas.UsersResponse, error) {
	return schemas.UsersResponse{}, nil
}
