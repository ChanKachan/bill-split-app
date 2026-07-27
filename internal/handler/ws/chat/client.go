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
	send        chan []byte // Сообщение, которое мы хотим отправить (writer)
	receive     chan []byte // Сообщение, которое мы получаем (reader)
	wg          *sync.WaitGroup
	chatService chat.ChatService
	userID      int
}

func newClient(
	conn *websocket.Conn,
	send, receive chan []byte,
	chatService chat.ChatService,
	userID int,
) *client {
	var wg sync.WaitGroup

	return &client{
		conn:        conn,
		send:        send,
		receive:     receive,
		wg:          &wg,
		chatService: chatService,
		userID:      userID,
	}
}

func (c *client) run(ctx context.Context) {
	c.wg.Add(2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// todo: добавить обработку ошибки
	go c.reader(ctx)
	go c.writer(ctx)

	log.Println("Web socket connected")
	c.wg.Wait()
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

		err = c.saveMessageToDB(ctx, 1, c.userID, string(msg)) // todo: стоят замоканные данные
		if err != nil {
			log.Println("Web socket error save message:", err)
			return fmt.Errorf("Error save message: %w", err)
		}
		c.receive <- msg // todo: пока просто отправим его
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
		case msg := <-c.receive:
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
	return nil
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
