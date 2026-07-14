package chat

type Chat struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	Id         int    `json:"id"`
	Text       string `json:"text"`
	UserId     int    `json:"user_id"`
	ChatId     int    `json:"chat_id"`
	DateCreate string `json:"date_create"`
	DateUpdate string `json:"date_update"`
}
