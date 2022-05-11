package token

type TokenService interface {
	GenerateToken(userId string) (string, error)
	ParseToken(signedToken string) (string, error)
}
