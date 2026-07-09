package chat

type RequestSendMessage struct {
	ChatId int
	Text   string
	UserID int
}
