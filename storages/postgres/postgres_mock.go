package postgres

import (
	"context"
	"tribble/models"
)

type PGMock struct {
	Users []*models.User
}

func (p PGMock) Close() {
}

func (p PGMock) GetUser(ctx context.Context, ID int) (*models.User, error) {
	for _, user := range p.Users {
		if user.ID == ID {
			return user, nil
		}
	}
	return nil, nil
}

func (p PGMock) GetUserList(ctx context.Context) ([]*models.User, error) {
	return p.Users, nil
}

func (p PGMock) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	return nil, nil
}

func (p PGMock) UpdateUser(ctx context.Context, user models.User) (*models.User, error) {
	return nil, nil
}

func (p PGMock) DeleteUser(ctx context.Context, ID int) error {
	return nil
}

func (p PGMock) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, nil
}

func (p PGMock) GetUserByRefresh(ctx context.Context, refresh string) (*models.User, error) {
	return nil, nil
}

func (p PGMock) UpdateUserTokens(ctx context.Context, ID int, token, refresh string) error {
	return nil
}

func (p PGMock) GetPlayerList(ctx context.Context, ID int) ([]*models.Player, error) {
	return nil, nil
}

func (p PGMock) CreatePlayer(ctx context.Context, player models.Player) (*models.Player, error) {
	return nil, nil
}

func (p PGMock) ValidateToken(ctx context.Context, refresh string) (bool, error) {
	return true, nil
}
