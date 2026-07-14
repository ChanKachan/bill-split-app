package chat

import "time"

type CreateMessageRequest struct {
	UserID     int
	Message    string
	ChatID     int
	DateCreate string `example:"02-01-2006 15:04:05"`
	DateUpdate string `example:"02-01-2006 15:04:05"`
}

type GetMessagesResponse struct {
	MessageID  int
	UserID     int
	Message    string
	DateCreate time.Time
	DateUpdate time.Time
	ChatID     int
}
