package db

import (
	"github.com/go-redis/redis/v8"
	"github.com/tuanp/go-mircroservice-boilerplate/pkg/config"

	"log"
	"time"
)

func ConnectRedis(cfgRedis *config.RedisConfig) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfgRedis.Addr,
		DB:           cfgRedis.DB,
		PoolSize:     cfgRedis.PoolSize,
		PoolTimeout:  time.Duration(cfgRedis.PoolTimeout) * time.Second,
		MinIdleConns: cfgRedis.MinIdleConns,
	})

	_, err := redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		log.Fatal("Error pinging redis", err)
	}
	return redisClient
}
