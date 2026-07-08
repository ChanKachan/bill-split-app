package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ChanKachan/bill-split-app/internal/config"
	"github.com/ChanKachan/bill-split-app/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler interface {
	ConnectionWS(c *gin.Context)
}

type chatHandler struct {
	ws       websocket.Upgrader
	configWS *config.CfgWebSocket
	mutex    sync.RWMutex
}

func NewChatHandler(
	ws websocket.Upgrader,
	configWS *config.CfgWebSocket,
) ChatHandler {
	return &chatHandler{
		ws:       ws,
		configWS: configWS,
	}
}

// Обновление http до web socket
func (ch *chatHandler) ConnectionWS(c *gin.Context) {
	var wg sync.WaitGroup

	//userID, ok := c.Get("userID")
	//if !ok {
	//	log.Println("Connect web socket get user id error")
	//	return
	//}

	conn, err := ch.ws.Upgrade(c.Writer, c.Request, nil)
	defer conn.Close()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

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

	wg.Add(2)
	clientConn := &client{
		conn:    conn,
		send:    make(chan []byte, 256),
		receive: make(chan []byte, 256),
		wg:      &wg,
		id:      1, // todo: нужно получить этот ID
		//id: userID,
	}

	go clientConn.readMessage(ctx)
	go clientConn.writeMessage(ctx)
	log.Println("Web socket connected")

	wg.Wait()

	return
}

// Закрыть соединение web socket
func (c *client) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("Web socket connection close error: %w", err)
	}
	return nil
}

func (c *client) readMessage(ctx context.Context) error {
	defer func() {
		c.wg.Done()

		err := c.Close()
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
			return fmt.Errorf("Error read messange: %w", err)
		}

		c.receive <- msg // todo: пока просто отправим его
		c.pongHandler()
	}
	return nil
}

func (c *client) writeMessage(ctx context.Context) error {
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
