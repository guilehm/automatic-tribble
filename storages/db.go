package storages

import "tribble/models"

type DBRepository interface {
	models.UserRepository
	models.PlayerRepository
	models.TokenRepository
	Close()
}

var DB DBRepository
