package handlers

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SecretKey = os.Getenv("JWT_SECRET_KEY")

type SignedDetails struct {
	Email string
	Uid   string
	jwt.StandardClaims
}

const accessTokenLifetime = time.Minute * time.Duration(10)
const refreshTokenLifetime = time.Hour * time.Duration(24)

func generateTokens(email string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email: email,
		Uid:   uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(accessTokenLifetime).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(refreshTokenLifetime).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SecretKey))
	if err != nil {
		log.Printf("could not create claims. %v\n", err.Error())
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SecretKey))
	if err != nil {
		log.Printf("could not create claims. %v\n", err.Error())
		return
	}
	return token, refreshToken, err
}
