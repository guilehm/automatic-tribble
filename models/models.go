package models

import "time"

type User struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	DateJoined time.Time `json:"date_joined"`
}
