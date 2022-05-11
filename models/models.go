package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name" validate:"required,gte=3"`
	Email        string    `json:"email" validate:"required,email"`
	Password     string    `json:"password" validate:"required,gte=5"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	DateJoined   time.Time `json:"date_joined"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

type Player struct {
	UserID    int
	XP        int
	Sprite    string `validate:"oneof=assassin warrior templar archer mage"`
	PositionX int
	PositionY int
}

type Tokens struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
