package main

import (
	"GoGameApp/adapter/redis"
	"GoGameApp/config"
	"GoGameApp/delivery/httpserver"
	"GoGameApp/repository/migrator"
	"GoGameApp/repository/mysql"
	"GoGameApp/repository/mysql/mysqlaccesscontrol"
	"GoGameApp/repository/mysql/mysqluser"
	"GoGameApp/repository/redis/redismatching"
	"GoGameApp/service/authorizationservice"
	"GoGameApp/service/authservice"
	"GoGameApp/service/backofficeuserservice"
	"GoGameApp/service/matchingservice"
	"GoGameApp/service/userservice"
	"GoGameApp/validator/matchingvalidator"
	"GoGameApp/validator/uservalidator"
	"fmt"
)

func main() {
	//? notes about reading conifgs
	// 1. load default
	// 2. read file and merge (overwrite)
	// 3. get env and merge (overwrite) -> higher priority
	cfg := config.Load("config.yml")
	fmt.Printf("conf: %+v\n", cfg)

	//TODO: add command for migrations
	migrator := migrator.New(cfg.Mysql)
	migrator.Up()

	authSvc, userSvc, userValidator, backofficeUserSvc, authorizationSvc, matchingSvc, matchingValidator := setupServices(cfg)
	server := httpserver.New(cfg, authSvc, userSvc, userValidator, backofficeUserSvc, authorizationSvc, matchingSvc, matchingValidator)

	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service,
	uservalidator.Validator,
	backofficeuserservice.Service,
	authorizationservice.Service,
	matchingservice.Service,
	matchingvalidator.Validator) {
	authSvc := authservice.New(cfg.Auth)

	mysqlRepo := mysql.New(cfg.Mysql)

	userMysql := mysqluser.New(mysqlRepo)
	userSvc := userservice.New(userMysql, authSvc)
	uV := uservalidator.New(userMysql)

	backofficeUserSvc := backofficeuserservice.New()

	aclMysql := mysqlaccesscontrol.New(mysqlRepo)
	authorizationSvc := authorizationservice.New(aclMysql)

	mV := matchingvalidator.New()
	redisAdapter := redis.New(cfg.Redis)
	matchingRepo := redismatching.New(redisAdapter)
	matchingSvc := matchingservice.New(matchingRepo, cfg.Matching)

	return authSvc, userSvc, uV, backofficeUserSvc, authorizationSvc, matchingSvc, mV
}
