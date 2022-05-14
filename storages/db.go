package storages

import "tribble/models"

type DBRepository interface {
	models.UserRepository
	Close()
}

var DB DBRepository
