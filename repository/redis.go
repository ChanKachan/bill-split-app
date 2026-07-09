package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDB interface {
	RedisClose() error
	Ping(ctx context.Context) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
}

type redisDB struct {
	client *redis.Client
}

func NewRedisDB(redisOption *redis.Options) RedisDB {
	return &redisDB{
		client: redis.NewClient(
			redisOption,
		),
	}
}

func (r *redisDB) RedisClose() error {
	return r.client.Close()
}

func (r *redisDB) Ping(ctx context.Context) (string, error) {
	return r.client.Ping(ctx).Result()
}

func (r *redisDB) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisDB) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisDB) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *redisDB) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}
