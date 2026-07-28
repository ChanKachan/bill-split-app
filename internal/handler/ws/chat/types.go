package chat

type chatData struct {
	userID  int
	chatID  int
	send    chan []byte // Сообщение, которое мы хотим отправить (writer)
	receive chan []byte // Сообщение, которое мы получаем (reader)
}
