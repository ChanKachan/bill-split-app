package chat

import (
	"fmt"

	"github.com/ChanKachan/bill-split-app/internal/domain/repository/postgres"
	"github.com/ChanKachan/bill-split-app/internal/domain/repository/redis/cache"
)

type ChatService interface {
}

type chatService struct {
	chatCache cache.ChatCache
	chatRepo  postgres.ChatRepository
}

func NewChatService(
	chatCache cache.ChatCache,
	chatRepo postgres.ChatRepository,
) ChatService {
	return &chatService{
		chatCache: chatCache,
		chatRepo:  chatRepo,
	}
}

func (cs *chatService) SendMessage(req RequestSendMessage) error {
	if req.Text == "" || req.UserID == 0 || req.ChatId == 0 {
		return fmt.Errorf("data is empty: %v", req)
	}

	return nil
}
