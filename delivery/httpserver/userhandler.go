package httpserver

import (
	"GoGameApp/service/userservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) userRegister(c echo.Context) error {
	var req userservice.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := s.userSvc.Register(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}
