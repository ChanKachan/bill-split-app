package chat

import "github.com/gorilla/websocket"

type client struct {
	conn    *websocket.Conn
	send    chan []byte // Сообщение, которое мы хотим отправить (writer)
	receive chan []byte // Сообщение, которое мы получаем (reader)
	id      int
}
