package bcrypt

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHashService struct {
	hashCost int
}

func NewBcryptHashService(hashCost int) *BcryptHashService {
	return &BcryptHashService{hashCost: hashCost}
}

func (b BcryptHashService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.hashCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash), nil
}

func (b BcryptHashService) ValidatePassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
