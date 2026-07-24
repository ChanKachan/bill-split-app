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
	HGet(ctx context.Context, key string, field string) (string, error)
	HSet(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	HDel(ctx context.Context, key string, field ...string) error
	LPush(ctx context.Context, key string, values ...string) (int64, error)
	LPop(ctx context.Context, key string) (string, error)
	LLen(ctx context.Context, key string) (int64, error)
	LRange(ctx context.Context, key string, start, end int64) ([]string, error)
	RPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	RPop(ctx context.Context, key string) (string, error)
	LTrim(ctx context.Context, key string, start, end int64) error
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

func (r *redisDB) HGet(ctx context.Context, key string, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

func (r *redisDB) HSet(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.HSet(ctx, key, value, expiration).Err()
}

func (r *redisDB) HDel(ctx context.Context, key string, field ...string) error {
	return r.client.HDel(ctx, key, field...).Err()
}

func (r *redisDB) LPush(ctx context.Context, key string, values ...string) (int64, error) {
	return r.client.LPush(ctx, key, values).Result()
}

func (r *redisDB) LPop(ctx context.Context, key string) (string, error) {
	return r.client.LPop(ctx, key).Result()
}

func (r *redisDB) LLen(ctx context.Context, key string) (int64, error) {
	return r.client.LLen(ctx, key).Result()
}

func (r *redisDB) LRange(ctx context.Context, key string, start, end int64) ([]string, error) {
	return r.client.LRange(ctx, key, start, end).Result()
}

func (r *redisDB) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return r.client.RPush(ctx, key, values).Result()
}

func (r *redisDB) RPop(ctx context.Context, key string) (string, error) {
	return r.client.RPop(ctx, key).Result()
}

func (r *redisDB) LTrim(ctx context.Context, key string, start, end int64) error {
	return r.client.LTrim(ctx, key, start, end).Err()
}
