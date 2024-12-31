package claim

import (
	"GoGameApp/config"
	"GoGameApp/service/authservice"

	"github.com/labstack/echo/v4"
)

func GetClaimsFromEchoContext(c echo.Context) *authservice.Claims {
	return c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)
}
