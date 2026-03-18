package user

type User struct {
	Id       int    `json:"id" db:"-"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Login    string `json:"login" db:"-"`
	Password string `json:"password" db:"password"`
	Phone    string `json:"phone" db:"phone"`
}
