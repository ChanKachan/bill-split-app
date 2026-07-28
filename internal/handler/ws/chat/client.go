package chat

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ChanKachan/bill-split-app/internal/domain/service/chat"
	"github.com/gorilla/websocket"
)

type client struct {
	conn        *websocket.Conn
	wg          *sync.WaitGroup
	chatService chat.ChatService
	data        chatData
}

func newClient(
	conn *websocket.Conn,
	chatService chat.ChatService,
	chatData chatData,
) *client {
	var wg sync.WaitGroup

	return &client{
		conn:        conn,
		wg:          &wg,
		chatService: chatService,
		data:        chatData,
	}
}

func (c *client) run(ctx context.Context) error {
	errors := make(chan error, 2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := c.reader(ctx); err != nil {
			errors <- fmt.Errorf("reader error: %w", err)
		}
	}()
	go func() {
		if err := c.writer(ctx); err != nil {
			errors <- fmt.Errorf("writer error: %w", err)
		}
	}()

	log.Println("Web socket connected")

	firstError := <-errors
	cancel()

	if err := c.conn.Close(); err != nil {
		log.Printf("close websocket connection: %v", err)
	}

	<-errors
	if err := firstError; err != nil {
		return fmt.Errorf("client client web-socket error: %w", err)
	}

	return nil
}

func (c *client) saveMessageToDB(ctx context.Context, chatID, userID int, msg string) error {
	data := chat.RequestSendMessage{
		ChatId: chatID,
		Text:   msg,
		UserID: userID,
	}

	err := c.chatService.CreateMessage(ctx, data)
	if err != nil {
		return fmt.Errorf("Error send message: %w", err)
	}
	return nil
}

// Закрыть соединение web socket
func (c *client) сlose() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("Web socket connection close error: %w", err)
	}
	return nil
}

func (c *client) reader(ctx context.Context) error {
	defer func() {
		c.wg.Done()

		err := c.сlose()
		if err != nil {
			log.Printf("Error close client: %v", err)
			return
		}
		return
	}()

	if err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil { // todo: подсоединить конфиг для время жизни
		return fmt.Errorf("Error setting read deadline: %w", err)
	}

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err) {
				log.Println("Web socket connection close")
				break
			}
			log.Println("Web socket error read message:", err)
			return fmt.Errorf("Error read messange: %w", err)
		}

		err = c.saveMessageToDB(ctx, c.data.chatID, c.data.userID, string(msg))
		if err != nil {
			log.Println("Web socket error save message:", err)
			return fmt.Errorf("Error save message: %w", err)
		}
		c.data.receive <- msg // todo: пока просто отправим его
		c.pongHandler()
	}
	return nil
}

func (c *client) writer(ctx context.Context) error {
	defer c.wg.Done()

	pingTicker := time.NewTicker(50 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case msg := <-c.data.receive:
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return fmt.Errorf("Error send message: %w", err)
			}
		case <-pingTicker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return fmt.Errorf("Error send ping: %w", err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (c *client) pongHandler() {
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
			return fmt.Errorf("Error setting read deadline: %w", err)
		}
		return nil
	})
	return
}
