package authservice

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uint `json:"user_id"`
}

func (c Claims) Valid() error {
	if c.UserID <= 0 {
		return fmt.Errorf("the user id must be greater than 0")
	}

	return nil
}
