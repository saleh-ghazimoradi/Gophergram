package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/saleh-ghazimoradi/Gophergram/internal/service/service_models"
	"time"
)

const UserExpiration = time.Minute

type CacheRepository interface {
	Get(ctx context.Context, id int64) (*service_models.User, error)
	Set(ctx context.Context, user *service_models.User) error
}

type cacheRepository struct {
	client *redis.Client
}

func (c *cacheRepository) Get(ctx context.Context, id int64) (*service_models.User, error) {
	cacheKey := fmt.Sprintf("user-%v", id)
	data, err := c.client.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	user := &service_models.User{}
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (c *cacheRepository) Set(ctx context.Context, user *service_models.User) error {
	cacheKey := fmt.Sprintf("user-%v", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.client.SetEX(ctx, cacheKey, data, UserExpiration).Err()
}

func NewCacheRepository(client *redis.Client) CacheRepository {
	return &cacheRepository{
		client: client,
	}
}
