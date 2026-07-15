package chat

import "time"

type CreateMessageRequest struct {
	UserID     int
	Message    string
	ChatID     int
	DateCreate time.Time
	DateUpdate time.Time
}

type GetMessagesResponse struct {
	MessageID  int
	UserID     int
	Message    string
	DateCreate time.Time
	DateUpdate time.Time
	ChatID     int
}
