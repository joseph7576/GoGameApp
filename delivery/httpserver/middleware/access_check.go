package middleware

import (
	"GoGameApp/entity"
	"GoGameApp/pkg/claim"
	"GoGameApp/pkg/errmsg"
	"GoGameApp/service/authorizationservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AccessCheck(service authorizationservice.Service, permissions ...entity.PermissionTitle) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := claim.GetClaimsFromEchoContext(c)
			isAllowed, err := service.CheckAccess(claims.UserID, claims.Role, permissions...)
			if err != nil {
				//TODO: LOG unexpected error
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errmsg.ErrMsgSomethingWentWrong,
				})
			}

			if !isAllowed {
				return c.JSON(http.StatusForbidden, echo.Map{
					"message": errmsg.ErrMsgUserNotAllowed,
				})
			}

			return next(c)
		}
	}
}
