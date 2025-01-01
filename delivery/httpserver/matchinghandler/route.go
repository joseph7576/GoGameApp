package matchinghandler

import (
	"GoGameApp/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	matchingGroup := e.Group("/matching")

	matchingGroup.POST("/add-to-waiting-list", h.AddToWaitingList,
		middleware.Auth(h.authSvc, h.authConfig))
}
