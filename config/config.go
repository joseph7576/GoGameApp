package config

import (
	"GoGameApp/adapter/redis"
	"GoGameApp/repository/mysql"
	"GoGameApp/service/authservice"
	"GoGameApp/service/matchingservice"
)

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	HTTPServer HTTPServer             `koanf:"http_server"`
	Auth       authservice.Config     `koanf:"auth"`
	Mysql      mysql.Config           `koanf:"mysql"`
	Matching   matchingservice.Config `koanf:"matching_service"`
	Redis      redis.Config           `koanf:"redis"`
}
