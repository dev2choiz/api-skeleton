package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/dev2choiz/api-skeleton/internal/config"
)

type Cache interface {
	GetString(ctx context.Context, key string) (string, error)
	SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error
	Exists(ctx context.Context, keys ...string) (int, error)
	Raw() *redis.Client
}

type cacheStruct struct {
	re *redis.Client
}

func New(conf *config.Config) (Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", conf.RedisHost, conf.RedisPort),
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	re := &cacheStruct{rdb}

	return re, err
}

func (c *cacheStruct) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.re.Set(ctx, key, b, ttl).Err()
}

func (c *cacheStruct) GetString(ctx context.Context, key string) (string, error) {
	return c.re.Get(ctx, key).Result()
}

func GetJSON[T any](ctx context.Context, r Cache, key string) (T, error) {
	var zero T

	val, err := r.GetString(ctx, key)
	if err != nil {
		return zero, err
	}

	var out T
	return out, json.Unmarshal([]byte(val), &out)
}

func (c *cacheStruct) Exists(ctx context.Context, keys ...string) (int, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	n, err := c.re.Exists(ctx, keys...).Result()

	return int(n), err
}

func (c *cacheStruct) Raw() *redis.Client {
	return c.re
}
