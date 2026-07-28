package chat

type RequestSendMessage struct {
	ChatId int
	Text   string
	UserID int
}

type RequestGetChat struct {
	ChatId int
}

type RequestIsChatsMember struct {
	ChatId int
	UserId int
}
