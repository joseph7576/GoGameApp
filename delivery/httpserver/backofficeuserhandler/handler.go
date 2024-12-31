package backofficeuserhandler

import (
	"GoGameApp/pkg/httpmsg"
	"GoGameApp/service/authorizationservice"
	"GoGameApp/service/authservice"
	"GoGameApp/service/backofficeuserservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	authConfig        authservice.Config
	authSvc           authservice.Service
	authorizationSvc  authorizationservice.Service
	backofficeUserSvc backofficeuserservice.Service
}

func New(authSvc authservice.Service, backofficeUserSvc backofficeuserservice.Service,
	authorizationSvc authorizationservice.Service, authConfig authservice.Config) Handler {
	return Handler{
		authConfig:        authConfig,
		authSvc:           authSvc,
		authorizationSvc:  authorizationSvc,
		backofficeUserSvc: backofficeUserSvc,
	}
}

func (h Handler) listUsers(c echo.Context) error {
	list, err := h.backofficeUserSvc.ListAllUsers()
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": list,
	})
}
