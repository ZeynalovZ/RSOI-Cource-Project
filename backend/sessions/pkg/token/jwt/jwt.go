package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTTokenService struct {
	signingKey string
	duration   time.Duration
}

func NewJWTokenService(signingKey string, duration time.Duration) *JWTTokenService {
	return &JWTTokenService{
		signingKey: signingKey,
		duration:   duration,
	}
}

func (t JWTTokenService) GenerateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(t.duration).Unix(),
		Subject:   userId,
	})
	return token.SignedString([]byte(t.signingKey))
}

func (t JWTTokenService) ParseToken(signedToken string) (string, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(t.signingKey), nil
	})

	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("can't extract token claims")
	}
	return claims["sub"].(string), nil
}
