package settings

import "os"

type Email string
type ID string

const (
	E Email = "email"
	I ID    = "id"
)

var JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
