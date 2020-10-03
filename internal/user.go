package internal

import "time"

type User struct {
	Id          string     `json:"id" bson:"_id"`
	Username    string     `json:"username" bson:"username"`
	Active      bool       `json:"active" validate:"active"`
	Locked      bool       `json:"locked"`
	CreatedDate *time.Time `json:"createdDate"`
}
