package models

import "context"

type UserRepository interface {
	GetUser(ctx context.Context, ID int) (*User, error)
	GetUserList(ctx context.Context) ([]*User, error)
	SaveUser(ctx context.Context, user User) (*User, error)
	UpdateUser(ctx context.Context, user User) (*User, error)
	DeleteUser(ctx context.Context, ID int) (*User, error)
}
