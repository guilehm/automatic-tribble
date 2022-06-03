package settings

import "os"

type Username string
type ID string

const (
	U Username = "username"
	I ID       = "id"
)

var JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
