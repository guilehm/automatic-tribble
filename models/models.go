package models

import "time"

type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	DateJoined time.Time `json:"date_joined"`
}
