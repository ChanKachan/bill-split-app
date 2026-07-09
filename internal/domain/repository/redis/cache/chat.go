package cache

import "github.com/ChanKachan/bill-split-app/repository"

type ChatCache interface{}

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
