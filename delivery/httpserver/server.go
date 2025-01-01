package httpserver

import (
	"GoGameApp/config"
	"GoGameApp/delivery/httpserver/backofficeuserhandler"
	"GoGameApp/delivery/httpserver/matchinghandler"
	"GoGameApp/delivery/httpserver/userhandler"
	"GoGameApp/service/authorizationservice"
	"GoGameApp/service/authservice"
	"GoGameApp/service/backofficeuserservice"
	"GoGameApp/service/matchingservice"
	"GoGameApp/service/userservice"
	"GoGameApp/validator/matchingvalidator"
	"GoGameApp/validator/uservalidator"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config                config.Config
	userHandler           userhandler.Handler
	backofficeUserHandler backofficeuserhandler.Handler
	matchingHandler       matchinghandler.Handler
}

func New(config config.Config, authSvc authservice.Service,
	userSvc userservice.Service, userValidator uservalidator.Validator,
	backofficeUserSvc backofficeuserservice.Service,
	authorizationSvc authorizationservice.Service,
	matchingSvc matchingservice.Service,
	matchingValidator matchingvalidator.Validator) Server {
	return Server{
		config:                config,
		userHandler:           userhandler.New(authSvc, userSvc, userValidator, config.Auth),
		backofficeUserHandler: backofficeuserhandler.New(authSvc, backofficeUserSvc, authorizationSvc, config.Auth),
		matchingHandler:       matchinghandler.New(config.Auth, authSvc, matchingSvc, matchingValidator),
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

	s.userHandler.SetRoutes(e)
	s.backofficeUserHandler.SetRoutes(e)
	s.matchingHandler.SetRoutes(e)

	// Start Server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))
}
