package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/ChanKachan/bill-split-app/internal/config"
	"github.com/ChanKachan/bill-split-app/internal/domain/service/chat"
	"github.com/ChanKachan/bill-split-app/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler interface {
	ConnectionWS(c *gin.Context)
}

type chatHandler struct {
	ws          websocket.Upgrader
	configWS    *config.CfgWebSocket
	chatService chat.ChatService
	mutex       sync.RWMutex
}

func NewChatHandler(
	ws websocket.Upgrader,
	configWS *config.CfgWebSocket,
	chatService chat.ChatService,
) ChatHandler {
	return &chatHandler{
		ws:          ws,
		configWS:    configWS,
		chatService: chatService,
	}
}

// Обновление http до web socket
func (ch *chatHandler) ConnectionWS(c *gin.Context) {
	var wg sync.WaitGroup

	rawUserID, exists := c.Get("userID")
	if !exists {
		log.Println("Connect web socket get user userID error")
		return
	}

	userID, ok := rawUserID.(int)
	if !ok {
		log.Println("Connect web socket get user ID error because type isn't int")
		return
	}

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

	defer conn.Close()

	wg.Add(2)

	clientConn := newClient(
		conn,
		make(chan []byte, 256),
		make(chan []byte, 256),
		&wg,
		ch.chatService,
		userID,
	)

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	go clientConn.reader(ctx)
	go clientConn.writer(ctx)
	log.Println("Web socket connected")

	wg.Wait()

	return
}
