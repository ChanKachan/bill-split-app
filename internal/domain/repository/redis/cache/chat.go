package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/ChanKachan/bill-split-app/repository"
	"github.com/redis/go-redis/v9"
	"time"
)

type ChatCache interface {
	GetMessage(ctx context.Context, chatID string) (string, error)
}

type chatCache struct {
	redisDB repository.RedisDB
}

func NewChatCache(
	redisDB repository.RedisDB,
) ChatCache {
	return &chatCache{
		redisDB: redisDB,
	}
}

func (c *chatCache) GetMessage(ctx context.Context, chatID string) (string, error) {
	response, err := c.redisDB.Get(
		ctx,
		fmt.Sprintf("group_message:%s", chatID),
	)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}

	return response, nil
}

// Перезаписывает полностью запись по ключу group_message:id
func (c *chatCache) SaveMessage(ctx context.Context, chatID string, msg string, timeLiveSecond int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := c.redisDB.Set(
		ctx,
		fmt.Sprintf("group_message:%s", chatID),
		msg,
		time.Duration(timeLiveSecond)*time.Second,
	)
	if err != nil {
		return err
	}

	return nil
}
