package user

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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
	Create(user models.User) (uuid.UUID, error)
	GetById(id uuid.UUID) (models.User, error)
	Update(id uuid.UUID, user models.UpdateUserInfo) error
	DeleteById(id uuid.UUID) error
	CreateMusicToUser(models.UserToMusic) error
	CreateUserImage(image models.Image) error
	GetUserImage(id uuid.UUID) (models.Image, error)
	GetAll(page, size int) ([]models.User, error)
	GetRecommendationsForUser(userId uuid.UUID, page int, size int) ([]models.User, error)
	Like(id uuid.UUID, likedId uuid.UUID) (bool, error)
}

const (
	userTable        = "users"
	userToMusicTable = "users_to_music"
	imageTable       = "images"
	likesTable       = "user_likes"
)

func NewRepository(db *sqlx.DB, logger *log.Logger) Repository {
	return repository{
		db:     db,
		logger: logger,
	}
}

func (r repository) Like(id uuid.UUID, likedId uuid.UUID) (bool, error) {
	// todo: add here excepting likes table
	query := fmt.Sprintf("INSERT INTO %s (id, who, from_who)"+
		" values ($1, $2, $3) RETURNING id", likesTable)
	row := r.db.QueryRow(query, uuid.New(), likedId, id)

	if err := row.Scan(&id); err != nil {
		r.logger.Printf("error in db while trying to create user like for user %v, error: %s",
			id, err.Error())
		return false, err
	}

	var like models.UserLikes
	query = fmt.Sprintf(`SELECT * FROM %s WHERE from_who = $1`, likesTable)
	err := r.db.Get(&like, query, likedId)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (r repository) GetRecommendationsForUser(userId uuid.UUID, page int, size int) ([]models.User, error) {
	var users []models.User
	query := fmt.Sprintf("SELECT * FROM %s  WHERE id != $1 LIMIT $2 OFFSET $3 ", userTable)

	err := r.db.Select(&users, query, userId, page*size, page-1)
	if err != nil {
		r.logger.Printf("error in db while trying to get all users, error: %s", err.Error())
		return nil, err
	}

	return users, nil
}

func (r repository) GetUserImage(id uuid.UUID) (models.Image, error) {
	var image models.Image
	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, imageTable)
	fmt.Sprintf(`SELECT * FROM %s WHERE userId = $1`, imageTable)
	err := r.db.Get(&image, query, id)
	if err == sql.ErrNoRows {
		return image, schemas.NotFoundError{Message: fmt.Sprintf("Not found any image of user with id %v", id)}
	}

	return image, err
}

func (r repository) Create(user models.User) (uuid.UUID, error) {
	query := fmt.Sprintf("INSERT INTO %s (id, name, surname, email, phone_number, has_access)"+
		" values ($1, $2, $3, $4, $5, $6) RETURNING id", userTable)

	var id uuid.UUID

	row := r.db.QueryRow(query, user.Id, user.Name, user.Surname, user.Email, user.PhoneNumber, user.HasAccess)

	if err := row.Scan(&id); err != nil {
		r.logger.Printf("error in db while trying to create user %v, error: %s",
			user.Id, err.Error())
		return uuid.Nil, err
	}

	return id, nil
}

func (r repository) GetById(id uuid.UUID) (models.User, error) {
	var user models.User
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, userTable)
	err := r.db.Get(&user, query, id)
	if err == sql.ErrNoRows {
		return user, schemas.NotFoundError{Message: fmt.Sprintf("Not found any user with id %v", id)}
	}

	return user, err
}

func (r repository) Update(id uuid.UUID, user models.UpdateUserInfo) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if user.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *user.Name)
		argId++
	}

	if user.Surname != nil {
		setValues = append(setValues, fmt.Sprintf("surname=$%d", argId))
		args = append(args, *user.Surname)
		argId++
	}

	if user.Email != nil {
		setValues = append(setValues, fmt.Sprintf("email=$%d", argId))
		args = append(args, *user.Email)
		argId++
	}

	if user.PhoneNumber != nil {
		setValues = append(setValues, fmt.Sprintf("phone_number=$%d", argId))
		args = append(args, *user.PhoneNumber)
		argId++
	}

	if user.HasAccess != nil {
		setValues = append(setValues, fmt.Sprintf("has_access=$%d", argId))
		args = append(args, *user.HasAccess)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s u SET %s WHERE u.id = $%v", userTable, setQuery, argId)

	args = append(args, id)

	_, err := r.db.Exec(query, args...)

	return err
}

func (r repository) DeleteById(id uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, userTable)
	_, err := r.db.Exec(query, id)

	return err
}

func (r repository) CreateMusicToUser(userToMusic models.UserToMusic) error {
	query := fmt.Sprintf("INSERT INTO %s (id, user_id, music_id, favourite_level)"+
		" values ($1, $2, $3, $4) RETURNING id", userToMusicTable)

	var id uuid.UUID

	row := r.db.QueryRow(query, userToMusic.Id, userToMusic.UserId, userToMusic.MusicId, userToMusic.FavouriteLevel)

	if err := row.Scan(&id); err != nil {
		r.logger.Printf("error in db while trying to create user to music %v, error: %s",
			userToMusic.Id, err.Error())
		return err
	}

	return nil
}

func (r repository) CreateUserImage(image models.Image) error {
	query := fmt.Sprintf("INSERT INTO %s (id, user_id, image)"+
		" values ($1, $2, $3) RETURNING id", imageTable)

	var id uuid.UUID

	row := r.db.QueryRow(query, image.Id, image.UserId, image.Image)

	if err := row.Scan(&id); err != nil {
		r.logger.Printf("error in db while trying to create user image %v, error: %s",
			image.Id, err.Error())
		return err
	}
	return nil
}

func (r repository) GetAll(page, size int) ([]models.User, error) {
	var users []models.User
	query := fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2", userTable)

	err := r.db.Select(&users, query, page*size, page-1)
	if err != nil {
		r.logger.Printf("error in db while trying to get all users, error: %s", err.Error())
		return nil, err
	}

	return users, nil
}
