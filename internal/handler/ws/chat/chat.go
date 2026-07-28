package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	ws          *websocket.Upgrader
	configWS    *config.CfgWebSocket
	chatService chat.ChatService
	mutex       sync.RWMutex
}

func NewChatHandler(
	ws *websocket.Upgrader,
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
	// Получить данные
	rawUserID, exists := c.Get("userID")
	if !exists {
		log.Println("Connect web socket get user userID error")
		return
	}

	userID, ok := rawUserID.(int)
	if !ok {
		log.Println("Connect web socket get user ID error because type isn't int")
		json.NewEncoder(c.Writer).Encode(types.ResponseError{
			Message: "Connect web socket get user ID error because type isn't int",
			Data:    nil,
			Code:    http.StatusInternalServerError,
		})
		return
	}

	dataURI := c.Param("chatID")
	chatID, err := strconv.Atoi(dataURI)
	if err != nil {
		log.Println("Connect web socket get user ID error:", err)
		json.NewEncoder(c.Writer).Encode(types.ResponseError{
			Message: fmt.Sprintf("Connect web socket get user ID error: %v", err),
			Data:    nil,
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Проверка на доступ к чату
	ch.chatService.IsChatsMember()

	// Обновляем до сокета
	conn, err := ch.upgradeConnection(c)
	if err != nil {
		json.NewEncoder(c.Writer).Encode(types.ResponseError{
			Message: fmt.Sprintf("WebSocket upgrade failed: %v", err),
			Data:    nil,
			Code:    http.StatusInternalServerError,
		})
		return
	}
	defer conn.Close()

	clientConn := newClient(
		conn,
		ch.chatService,
		chatData{
			userID:  userID,
			chatID:  chatID,
			send:    make(chan []byte, 256),
			receive: make(chan []byte, 256),
		},
	)

	err = clientConn.run(c.Request.Context())
	if err != nil {
		log.Printf("WebSocket client failed: %v", err)
		return
	}
	return
}

func (ch *chatHandler) upgradeConnection(c *gin.Context) (*websocket.Conn, error) {
	conn, err := ch.ws.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		//if websocket.IsWebSocketUpgrade(c.Request) {
		//	c.Writer.WriteHeader(http.StatusInternalServerError)
		//
		//	log.Printf("Ошибка web socket апгрейда: %v", err)
		//
		//	json.NewEncoder(c.Writer).Encode(types.ResponseError{
		//		Message: fmt.Sprintf("WebSocket upgrade failed: %v", err),
		//		Data:    nil,
		//		Code:    500,
		//	})
		//	return nil, err
		//}
		//c.Writer.WriteHeader(http.StatusBadRequest)
		//
		//log.Printf(`Connection to web socket error: %v`, err)
		//json.NewEncoder(c.Writer).Encode(types.ResponseError{
		//	Message: fmt.Sprintf("Connection to web socket error: %v", err),
		//	Data:    nil,
		//	Code:    400,
		//})
		//return
		return nil, fmt.Errorf("upgrade websocket connection: %w", err)
	}

	return conn, nil
}
