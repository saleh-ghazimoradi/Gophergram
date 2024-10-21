package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_modles"
	"time"
)

const UserExpiration = time.Minute

type Cacher interface {
	Get(ctx context.Context, id int64) (*service_modles.Users, error)
	Set(ctx context.Context, user *service_modles.Users) error
}

type cacheRepo struct {
	client *redis.Client
}

func (c *cacheRepo) Get(ctx context.Context, id int64) (*service_modles.Users, error) {
	cahcKey := fmt.Sprintf("user-%v", id)
	data, err := c.client.Get(ctx, cahcKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user service_modles.Users
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (c *cacheRepo) Set(ctx context.Context, user *service_modles.Users) error {
	cacheKey := fmt.Sprintf("user-%v", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.client.SetEX(ctx, cacheKey, data, UserExpiration).Err()
}

func NewCacheRepo(client *redis.Client) Cacher {
	return &cacheRepo{
		client: client,
	}
}
