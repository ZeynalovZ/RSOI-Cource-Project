package music

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Feokrat/music-dating-app/users/internal/models"
	"github.com/Feokrat/music-dating-app/users/internal/schemas"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db     *sqlx.DB
	logger *log.Logger
}

type Repository interface {
	Create(music models.Music) (uuid.UUID, error)
	GetById(id uuid.UUID) (models.Music, error)
	DeleteById(id uuid.UUID) error
	GetAll(page, size int) ([]models.Music, error)
}

const (
	musicTable = "musics"
)

func NewRepository(db *sqlx.DB, logger *log.Logger) Repository {
	return repository{
		db:     db,
		logger: logger,
	}
}

func (r repository) Create(music models.Music) (uuid.UUID, error) {
	query := fmt.Sprintf("INSERT INTO %s (id, name, author, url)"+
		" values ($1, $2, $3, $4) RETURNING id", musicTable)

	var id uuid.UUID

	row := r.db.QueryRow(query, music.Id, music.Name, music.Author, music.Url)

	if err := row.Scan(&id); err != nil {
		r.logger.Printf("error in db while trying to create music %v, error: %s",
			music.Id, err.Error())
		return uuid.Nil, err
	}

	return id, nil
}

func (r repository) GetById(id uuid.UUID) (models.Music, error) {
	var music models.Music
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, musicTable)
	err := r.db.Get(&music, query, id)
	if err == sql.ErrNoRows {
		return music, schemas.NotFoundError{Message: fmt.Sprintf("Not found any music with id %d", id)}
	}

	return music, err
}

func (r repository) DeleteById(id uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, musicTable)
	_, err := r.db.Exec(query, id)

	return err
}

func (r repository) GetAll(page, size int) ([]models.Music, error) {
	var musics []models.Music
	query := fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2", musicTable)

	err := r.db.Select(&musics, query, page*size, page-1)
	if err != nil {
		r.logger.Printf("error in db while trying to get all musics, error: %s", err.Error())
		return nil, err
	}

	return musics, nil
}
