package main

import (
	"GoGameApp/config"
	"GoGameApp/delivery/httpserver"
	"GoGameApp/repository/migrator"
	"GoGameApp/repository/mysql"
	"GoGameApp/service/authservice"
	"GoGameApp/service/userservice"
	"GoGameApp/validator/uservalidator"
	"fmt"
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
	//? notes about reading conifgs
	// 1. load default
	// 2. read file and merge (overwrite)
	// 3. get env and merge (overwrite) -> higher priority
	conf := config.Load("config.yml")
	fmt.Printf("conf: %+v\n", conf)

	//TODO: merge conf with cfg
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

	//TODO: add command for migrations
	migrator := migrator.New(cfg.Mysql)
	migrator.Up()

	authSvc, userSvc, userValidator := setupServices(cfg)
	server := httpserver.New(cfg, authSvc, userSvc, userValidator)

	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator) {
	authSvc := authservice.New(cfg.Auth)
	mysqlRepo := mysql.New(cfg.Mysql)
	userSvc := userservice.New(mysqlRepo, authSvc)
	uV := uservalidator.New(mysqlRepo)

	return authSvc, userSvc, uV
}
