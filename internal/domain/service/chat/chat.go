package chat

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ChanKachan/bill-split-app/internal/domain/repository/postgres/chat"
	"github.com/ChanKachan/bill-split-app/internal/domain/repository/redis/cache"
)

type ChatService interface {
	GetChat(ctx context.Context, chatID RequestGetChat) ([]string, error)
	SendMessage(ctx context.Context, req RequestSendMessage) error
}

type chatService struct {
	chatCache cache.ChatCache
	chatRepo  chat.ChatRepository
}

func NewChatService(
	chatCache cache.ChatCache,
	chatRepo chat.ChatRepository,
) ChatService {
	return &chatService{
		chatCache: chatCache,
		chatRepo:  chatRepo,
	}
}

func (cs *chatService) GetChat(ctx context.Context, req RequestGetChat) ([]string, error) {
	if req.ChatId <= 0 {
		return nil, errors.New("chat id must be greater than zero")
	}
	messages, err := cs.chatCache.GetMessagesFromList(ctx, strconv.Itoa(req.ChatId), 0, 50)
	if err != nil {
		return nil, fmt.Errorf("get chat messages from cache error: %w", err)
	}
	return messages, nil
}

func (cs *chatService) SendMessage(ctx context.Context, req RequestSendMessage) error {
	if req.Text == "" || req.UserID == 0 || req.ChatId == 0 {
		return fmt.Errorf("data is empty: %v", req)
	}

	messages, err := cs.GetChat(ctx, RequestGetChat{ChatId: req.ChatId})
	if err != nil {
		return fmt.Errorf("send message error get chat: %w", err)
	}

	if len(messages) == 0 {

	}
	return nil
}
