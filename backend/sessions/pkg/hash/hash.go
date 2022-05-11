package hash

type HashService interface {
	HashPassword(password string) (string, error)
	ValidatePassword(password string, hashedPassword string) error
}
