package config

import (
	"GoGameApp/repository/mysql"
	"GoGameApp/service/authservice"
)

type HTTPServer struct {
	Port int
}

type Config struct {
	HTTPServer HTTPServer
	Auth       authservice.Config
	Mysql      mysql.Config
}
