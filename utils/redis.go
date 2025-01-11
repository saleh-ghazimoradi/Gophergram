package utils

import (
	"github.com/go-redis/redis/v8"
	"github.com/saleh-ghazimoradi/Gophergram/logger"
)

func RedisConnection(addr, pw string, db int) (*redis.Client, error) {
	logger.Logger.Info("Connecting to Redis...")
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	}), nil

}
