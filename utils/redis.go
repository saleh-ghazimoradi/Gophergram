package utils

import "github.com/go-redis/redis/v8"

func RedisConnection(addr, pw string, db int) (*redis.Client, error) {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	}), nil
}
