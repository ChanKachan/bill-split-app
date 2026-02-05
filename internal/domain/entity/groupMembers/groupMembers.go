package groupMembers

type GroupMembers struct {
	UserId  int     `json:"user_id"`
	GroupId string  `json:"group_id"`
	Amount  float32 `json:"amount"`
	Status  string  `json:"status"`
	Del     int     `json:"del"`
}
