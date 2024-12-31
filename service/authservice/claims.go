package authservice

import (
	"GoGameApp/entity"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uint        `json:"user_id"`
	Role   entity.Role `json:"role"`
}

func (c Claims) Valid() error {
	if c.UserID <= 0 {
		return fmt.Errorf("the user id must be greater than 0")
	}

	return nil
}
