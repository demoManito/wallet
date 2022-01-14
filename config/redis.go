package config

import (
	"github.com/go-redis/redis/v8"
)

func InitRedis(conf *Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password, // no password set
		DB:       conf.DB,       // use default DB
		PoolSize: conf.PoolSize,
	})
}
