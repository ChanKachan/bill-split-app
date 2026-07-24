package chat

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	entityChat "github.com/ChanKachan/bill-split-app/internal/domain/entity/chat"
	"github.com/ChanKachan/bill-split-app/internal/domain/repository/postgres/chat"
	"github.com/ChanKachan/bill-split-app/internal/domain/repository/redis/cache"
	"github.com/ChanKachan/bill-split-app/internal/utils"
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
	messages, err := cs.chatCache.GetMessagesFromList(ctx, strconv.Itoa(req.ChatId), 0, 49)
	if err != nil {
		return nil, fmt.Errorf("get chat messages from cache error: %w", err)
	}
	return messages, nil
}

func (cs *chatService) SendMessage(ctx context.Context, req RequestSendMessage) error {
	if req.Text == "" || req.UserID == 0 || req.ChatId == 0 {
		return fmt.Errorf("data is empty: %v", req)
	}

	timeNow := time.Now()
	timeNowStr := timeNow.Format("02-01-2006 15:04:05")

	msgID, err := cs.chatRepo.CreateMessage(ctx, chat.CreateMessageRequest{
		UserID:     req.UserID,
		Message:    req.Text,
		ChatID:     req.ChatId,
		DateCreate: timeNow,
		DateUpdate: timeNow,
	})
	if err != nil {
		return fmt.Errorf("create chat message error: %w", err)
	}

	msgData := entityChat.Message{
		Id:         msgID,
		Text:       req.Text,
		UserId:     req.UserID,
		ChatId:     req.ChatId,
		DateCreate: timeNowStr,
		DateUpdate: timeNowStr,
	}

	err = cs.updateMessageToCache(ctx, msgData)
	if err != nil {
		return fmt.Errorf("update message to cache error: %w", err)
	}

	return nil
}

// Обновляем кэш если он заполнен
// Если элементов в листе меньше 50, то мы добавляем элемент без удаления правого (раннего сообщения)
func (cs *chatService) updateMessageToCache(ctx context.Context, chatData entityChat.Message) error {
	messages, err := cs.chatRepo.GetLastMessages(ctx, chatData.ChatId)
	if err != nil {
		return fmt.Errorf("get last messages from repository error: %w", err)
	}

	err = cs.chatCache.AddMessageOnLeftToList(ctx, strconv.Itoa(chatData.ChatId), utils.ConvertStructsToString(chatData)...)
	if err != nil {
		return fmt.Errorf("add message to cache error: %w", err)
	}

	if len(messages) > 50 {
		err = cs.chatCache.DelOnRightMessageFromList(ctx, strconv.Itoa(chatData.ChatId))
		if err != nil {
			return fmt.Errorf("delete on message from cache error: %w", err)
		}
	}

	return nil
}
