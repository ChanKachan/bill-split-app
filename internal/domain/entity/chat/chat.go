package chat

type Chat struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	Id      int    `json:"id"`
	Text    string `json:"text"`
	UserId  int    `json:"user_id"`
	GroupId int    `json:"group_id"`
}
