package backofficeuserhandler

import (
	"GoGameApp/delivery/httpserver/middleware"
	"GoGameApp/entity"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetBackofficeUserRoute(e *echo.Echo) {
	backofficeGroup := e.Group("/backoffice/users")

	backofficeGroup.GET("/", h.listUsers,
		middleware.Auth(h.authSvc, h.authConfig),
		middleware.AccessCheck(h.authorizationSvc, entity.UserListPermission))
}
