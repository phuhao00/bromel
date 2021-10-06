package xredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/phuhao00/bromel"
)

type Config struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type Client struct {
	*bromel.BaseComponent
}

func NewClient(ctx context.Context, config *Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password, // no password set
		DB:       config.DB,       // use default DB
	})
	result := rdb.Ping(ctx)
	if result.Err() != nil {
		panic(result.Err())
	}
	return rdb
}
