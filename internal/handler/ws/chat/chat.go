package chat

import (
	"encoding/json"
	"fmt"
	"github.com/ChanKachan/bill-split-app/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type ChatHandler interface {
	ConnectionWS(c *gin.Context)
	Close() error
}

type chatHandler struct {
	ws     websocket.Upgrader
	conn   *websocket.Conn
	client *client
	mutex  sync.RWMutex
}

func NewChatHandler(
	ws websocket.Upgrader,
) ChatHandler {
	return &chatHandler{
		ws: ws,
	}
}

// Обновление http до web socket
func (ch *chatHandler) ConnectionWS(c *gin.Context) {
	conn, err := ch.ws.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		if websocket.IsWebSocketUpgrade(c.Request) {
			c.Writer.WriteHeader(http.StatusInternalServerError)

			log.Printf("Ошибка web socket апгрейда: %v", err)

			json.NewEncoder(c.Writer).Encode(types.ResponseError{
				Message: fmt.Sprintf("WebSocket upgrade failed: %v", err),
				Data:    nil,
				Code:    500,
			})
			return
		}
		c.Writer.WriteHeader(http.StatusBadRequest)

		log.Printf(`Connection to web socket error: %w`, err)
		json.NewEncoder(c.Writer).Encode(types.ResponseError{
			Message: fmt.Sprintf("Connection to web socket error: %v", err),
			Data:    nil,
			Code:    400,
		})
		return
	}

	ch.client = &client{
		conn:    conn,
		send:    make(chan []byte, 256),
		receive: make(chan []byte, 256),
		id:      1, // todo: нужно получить этот ID
	}

	go ch.readMessage()
	go ch.writeMessage()
	log.Println("Web socket connected")

	return
}

// Закрыть соединение web socket
func (ch *chatHandler) Close() error {
	if err := ch.conn.Close(); err != nil {
		return fmt.Errorf("Web socket connection close error: %w", err)
	}
	return nil
}

func (ch *chatHandler) readMessage() error {
	for {
		_, msg, err := ch.conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("Error read messange: %w", err)
		}

		ch.client.receive <- msg // todo: пока просто отправим его

	}
}

func (ch *chatHandler) writeMessage() error {
	for {
		select {
		case msg := <-ch.client.receive:
			if err := ch.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return fmt.Errorf("Error send message: %w", err)
			}
		}

	}
}
