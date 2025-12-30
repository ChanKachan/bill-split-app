package group

import "time"

type Group struct {
	Id        int       `json:"id"`
	CreateAt  time.Time `json:"create_at"`
	DateStart time.Time `json:"date_start"`
	DateEnd   time.Time `json:"date_end"`
}
