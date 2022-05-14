package models

import "context"

type UserRepository interface {
	GetUser(ctx context.Context, ID int) (*User, error)
	GetUserList(ctx context.Context) ([]*User, error)
	CreateUser(ctx context.Context, user User) (*User, error)
	UpdateUser(ctx context.Context, user User) (*User, error)
	DeleteUser(ctx context.Context, ID int) error

	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUserTokens(ctx context.Context, ID int, token, refresh string) error
}
