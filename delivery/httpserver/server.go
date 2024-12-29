package httpserver

import (
	"GoGameApp/config"
	"GoGameApp/delivery/httpserver/userhandler"
	"GoGameApp/service/authservice"
	"GoGameApp/service/userservice"
	"GoGameApp/validator/uservalidator"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config      config.Config
	userHandler userhandler.Handler
}

func New(config config.Config, authSvc authservice.Service,
	userSvc userservice.Service, userValidator uservalidator.Validator) Server {
	return Server{
		config:      config,
		userHandler: userhandler.New(authSvc, userSvc, userValidator, config.Auth),
	}
}

func (s Server) Serve() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/health", s.healthCheck)

	s.userHandler.SetUserRoutes(e)

	// Start Server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))
}
