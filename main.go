package main

import (
	"GoGameApp/config"
	"GoGameApp/delivery/httpserver"
	"GoGameApp/repository/mysql"
	"GoGameApp/service/authservice"
	"GoGameApp/service/userservice"
	"time"
)

const (
	JWTSignKey                 = "very_secret_key"
	AccessTokenSubject         = "at"
	RefreshTokenSubject        = "rt"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
)

func main() {

	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8080},
		Auth: authservice.Config{
			SignKey:               JWTSignKey,
			AccessSubject:         AccessTokenSubject,
			RefreshSubject:        RefreshTokenSubject,
			AccessExpirationTime:  AccessTokenExpireDuration,
			RefreshExpirationTime: RefreshTokenExpireDuration,
		},
		Mysql: mysql.Config{
			User:                 "root",
			Passwd:               "root",
			Net:                  "tcp",
			Addr:                 "localhost:3306",
			DBName:               "gameapp_local",
			AllowNativePasswords: true,
		},
	}

	authSvc, userSvc := setupServices(cfg)
	server := httpserver.New(cfg, authSvc, userSvc)

	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.Auth)
	mysqlRepo := mysql.New(cfg.Mysql)
	userSvc := userservice.New(mysqlRepo, authSvc)

	return authSvc, userSvc
}
