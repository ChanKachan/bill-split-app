package cost

type Cost struct {
	Id          int     `json:"id"`
	UserId      int     `json:"user_id"`
	GroupId     int     `json:"group_id"`
	Description string  `json:"description"`
	Sum         float64 `json:"sum"`
}
