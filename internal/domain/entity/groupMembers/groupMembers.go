package groupMembers

type GroupMembers struct {
	UserId  int     `json:"user_id" binding:"required"`
	GroupId int     `json:"group_id" binding:"required"`
	Amount  float32 `json:"amount"`
	Status  string  `json:"status"`
	Del     int     `json:"del"`
}

type UserStatus string

const (
	User  UserStatus = "User"
	Admin UserStatus = "Admin"
)
