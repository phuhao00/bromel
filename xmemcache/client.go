package xmemcache

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/phuhao00/bromel"
)

var (
	_ bromel.Component = (*Client)(nil)
)

type Client struct {
	*bromel.BaseComponent
}

type Config struct {
	ServerAddress string
}

func NewClient(config *Config) *memcache.Client {
	client := memcache.New(config.ServerAddress)
	err := client.Ping()
	if err != nil {
		panic(err)
	}
	return client
}
