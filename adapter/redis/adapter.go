package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

// adapters are that modules that talks to external services
// repo is an adapter

type Config struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`
}

type Adapter struct {
	client *redis.Client
}

func New(config Config) Adapter {
	return Adapter{client: redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})}
}

func (a Adapter) Client() *redis.Client {
	return a.client
}
