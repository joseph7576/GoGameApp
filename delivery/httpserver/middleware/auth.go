package middleware

import (
	"GoGameApp/pkg/constant"
	"GoGameApp/service/authservice"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func Auth(service authservice.Service, config authservice.Config) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContextKey: constant.AuthMiddlewareContextKey,
		SigningKey: []byte(config.SignKey),
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			claims, err := service.ParseToken(auth)
			if err != nil {
				return nil, err
			}

			return claims, nil
		},
	})
}
