package authservice

import (
	"GoGameApp/entity"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	SignKey               string
	AccessExpirationTime  time.Duration
	RefreshExpirationTime time.Duration
	AccessSubject         string
	RefreshSubject        string
}

type Service struct {
	config Config
}

func New(config Config) Service {
	return Service{config: config}
}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.config.AccessSubject, s.config.AccessExpirationTime)
}

func (s Service) CreateRefreshToken(user entity.User) (string, error) {
	return s.createToken(user.ID, s.config.RefreshSubject, s.config.RefreshExpirationTime)
}

func (s Service) ParseToken(tokenString string) (*Claims, error) {
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.SignKey), nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	} else {
		return nil, err
	}
}

func (s Service) createToken(userID uint, subject string, expireDuration time.Duration) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   subject,
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.config.SignKey))
}
