package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name" validate:"required,gte=3"`
	Email        string    `json:"email" validate:"required,email"`
	DateJoined   time.Time `json:"date_joined" validate:"required"`
	Password     string    `json:"-"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type Player struct {
	UserID    int
	XP        int
	Sprite    string
	PositionX int
	PositionY int
}
