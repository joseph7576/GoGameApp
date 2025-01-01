package matchinghandler

import (
	"GoGameApp/param"
	"GoGameApp/pkg/claim"
	"GoGameApp/pkg/httpmsg"
	"GoGameApp/service/authservice"
	"GoGameApp/service/matchingservice"
	"GoGameApp/validator/matchingvalidator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	authConfig        authservice.Config
	authSvc           authservice.Service
	matchingSvc       matchingservice.Service
	matchingValidator matchingvalidator.Validator
}

func New(authConfig authservice.Config,
	authSvc authservice.Service,
	matchingSvc matchingservice.Service,
	matchingValidator matchingvalidator.Validator) Handler {
	return Handler{
		authConfig:        authConfig,
		authSvc:           authSvc,
		matchingSvc:       matchingSvc,
		matchingValidator: matchingValidator,
	}
}

func (h Handler) AddToWaitingList(c echo.Context) error {
	var req param.AddToWaitingListRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	claims := claim.GetClaimsFromEchoContext(c)
	req.UserID = claims.UserID

	if fieldErrors, err := h.matchingValidator.ValidateAddToWaitingList(req); err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}

	resp, err := h.matchingSvc.AddToWaitingList(req)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)
}
