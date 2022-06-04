package settings

import (
	"os"
	"time"
)

type Username string
type ID string

const (
	U Username = "username"
	I ID       = "id"
)

var JWTSecretKey = os.Getenv("JWT_SECRET_KEY")

const AccessTokenLifetime = time.Minute * time.Duration(10)
const RefreshTokenLifetime = time.Hour * time.Duration(24)
