package models

import "context"

type UserRepository interface {
	GetUser(ctx context.Context, ID int) (*User, error)
	GetUserList(ctx context.Context) ([]*User, error)
	CreateUser(ctx context.Context, user User) (*User, error)
	UpdateUser(ctx context.Context, user User) (*User, error)
	DeleteUser(ctx context.Context, ID int) error

	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByRefresh(ctx context.Context, refresh string) (*User, error)
	UpdateUserTokens(ctx context.Context, ID int, token, refresh string) error
}

type PlayerRepository interface {
	GetPlayerList(ctx context.Context, ID int) ([]*Player, error)
	CreatePlayer(ctx context.Context, player Player) (*Player, error)
}

type TokenRepository interface {
	ValidateToken(ctx context.Context, refresh string) (bool, error)
}
