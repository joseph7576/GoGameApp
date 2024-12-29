package userhandler

import (
	"GoGameApp/param"
	"GoGameApp/pkg/constant"
	"GoGameApp/pkg/httpmsg"
	"GoGameApp/service/authservice"
	"GoGameApp/service/userservice"
	"GoGameApp/validator/uservalidator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	authConfig    authservice.Config
	authSvc       authservice.Service
	userSvc       userservice.Service
	userValidator uservalidator.Validator
}

func New(authSvc authservice.Service, userSvc userservice.Service,
	userValidator uservalidator.Validator, authConfig authservice.Config) Handler {
	return Handler{
		authConfig:    authConfig,
		authSvc:       authSvc,
		userSvc:       userSvc,
		userValidator: userValidator,
	}
}

func (h Handler) userRegister(c echo.Context) error {
	var req param.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if fieldErrors, err := h.userValidator.ValidateRegisterRequest(req); err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}

	resp, err := h.userSvc.Register(req)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h Handler) userLogin(c echo.Context) error {
	var req param.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if fieldErrors, err := h.userValidator.ValidateLoginRequest(req); err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}

	resp, err := h.userSvc.Login(req)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)
}

func (h Handler) userProfile(c echo.Context) error {
	claims := getClaims(c)

	req := param.ProfileRequest{UserID: claims.UserID}
	resp, err := h.userSvc.Profile(req)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)
}

func getClaims(c echo.Context) *authservice.Claims {
	return c.Get(constant.AuthMiddlewareContextKey).(*authservice.Claims)
}
