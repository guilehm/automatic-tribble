package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
	"tribble/db"
	"tribble/models"
	"tribble/settings"

	"github.com/dgrijalva/jwt-go"
)

type SignedDetails struct {
	Email string
	ID    string
	jwt.StandardClaims
}

const accessTokenLifetime = time.Minute * time.Duration(10)
const refreshTokenLifetime = time.Hour * time.Duration(24)

func CheckToken(signedToken string) (claims *SignedDetails, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(settings.JWTSecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("token is expired")
	}
	return claims, nil
}

func generateTokens(email string, userId int) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email: email,
		ID:    strconv.Itoa(userId),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(accessTokenLifetime).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(refreshTokenLifetime).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(settings.JWTSecretKey))
	if err != nil {
		log.Printf("could not create claims. %v\n", err.Error())
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(settings.JWTSecretKey))
	if err != nil {
		log.Printf("could not create claims. %v\n", err.Error())
		return
	}
	return token, refreshToken, err
}

func ValidateToken(w http.ResponseWriter, r *http.Request) {

	var tokens models.Tokens

	if err := json.NewDecoder(r.Body).Decode(&tokens); err != nil {
		HandleApiErrors(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `SELECT id FROM users WHERE refresh_token=$1`

	var user models.User
	if err := db.DB.QueryRow(ctx, sql, tokens.RefreshToken).Scan(&user.ID); err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusNotFound, "")
		return
	}
	response, _ := json.Marshal(struct {
		Ok bool `json:"ok"`
	}{Ok: true})
	_, _ = w.Write(response)
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {

	var tokens models.Tokens

	if err := json.NewDecoder(r.Body).Decode(&tokens); err != nil {
		HandleApiErrors(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `SELECT id, email FROM users WHERE refresh_token=$1`

	var user models.User
	if err := db.DB.QueryRow(ctx, sql, tokens.RefreshToken).Scan(&user.ID, &user.Email); err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusNotFound, "")
		return
	}

	token, refresh, err := generateTokens(user.Email, user.ID)
	if err != nil {
		HandleApiErrors(w, http.StatusInternalServerError, err.Error())
		return
	}

	sql = `UPDATE users SET token=$1, refresh_token=$2 WHERE id=$3`
	_, err = db.DB.Query(context.Background(), sql, token, refresh, user.ID)

	if err != nil {
		HandleApiErrors(w, http.StatusInternalServerError, "could not update tokens")
		return
	}

	response, _ := json.Marshal(struct {
		Token   string `json:"token"`
		Refresh string `json:"refresh"`
	}{token, refresh})
	_, _ = w.Write(response)
}
