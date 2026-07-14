package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ChanKachan/bill-split-app/repository"
	"github.com/redis/go-redis/v9"
)

type ChatCache interface {
	GetMessage(ctx context.Context, chatID string) (string, error)
	SaveMessage(ctx context.Context, chatID string, msg string, timeLiveSecond int) error
	AddMessageOnLeftToList(ctx context.Context, chatID string, dataMessage ...string) error
	DelOnRightMessageFromList(ctx context.Context, chatID string) error
	GetMessagesFromList(ctx context.Context, chatID string, start, end int64) ([]string, error)
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

func (c *chatCache) AddMessageOnLeftToList(ctx context.Context, chatID string, msg ...string) error {
	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	_, err := c.redisDB.LPush(
		ctx,
		fmt.Sprintf("group_message:%s", chatID),
		msg...,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *chatCache) DelOnRightMessageFromList(ctx context.Context, chatID string) error {
	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	_, err := c.redisDB.RPop(
		ctx,
		fmt.Sprintf("group_message:%s", chatID),
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *chatCache) GetMessagesFromList(ctx context.Context, chatID string, start, end int64) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	messages, err := c.redisDB.LRange(
		ctx,
		fmt.Sprintf("group_message:%s", chatID),
		start,
		end,
	)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
