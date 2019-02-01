package models

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"github.com/stanlyliao/logger"
)

var (
	r *redis.Client
)

// InitRedis redis initial
func InitRedis() {
	host := viper.GetString("redis.host")
	if host == "" {
		host = "127.0.0.1"
	}
	port := viper.GetInt("redis.port")
	if port == 0 {
		port = 6379
	}

	r = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    fmt.Sprintf("%s:%d", host, port),
		DB:      viper.GetInt("redis.db"),
	})

	if _, err := r.Ping().Result(); err != nil {
		logger.Fatal(err)
	}
}
