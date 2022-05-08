package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"tribble/db"
	"tribble/models"

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

func generateTokens(email string, userId int) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email: email,
		Uid:   strconv.Itoa(userId),
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

func ValidateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tokens models.Tokens

	if err := json.NewDecoder(r.Body).Decode(&tokens); err != nil {
		HandleApiErrors(w, http.StatusBadRequest, err.Error())
		return
	}

	sql := `SELECT id FROM users WHERE refresh_token=$1`
	row := db.DB.QueryRow(context.Background(), sql, tokens.RefreshToken)

	var user models.User
	if err := row.Scan(&user.ID); err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusNotFound, "")
		return
	}
	response, _ := json.Marshal(struct {
		Ok bool `json:"ok"`
	}{Ok: true})
	w.Write(response)
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tokens models.Tokens

	if err := json.NewDecoder(r.Body).Decode(&tokens); err != nil {
		HandleApiErrors(w, http.StatusBadRequest, err.Error())
		return
	}

	sql := `SELECT id, email FROM users WHERE refresh_token=$1`
	row := db.DB.QueryRow(context.Background(), sql, tokens.RefreshToken)

	var user models.User
	if err := row.Scan(&user.ID, &user.Email); err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusNotFound, "")
		return
	}

	token, refresh, err := generateTokens(user.Email, user.ID)
	if err != nil {
		HandleApiErrors(w, http.StatusInternalServerError, err.Error())
		return
	}

	response, _ := json.Marshal(struct {
		Token   string `json:"token"`
		Refresh string `json:"refresh"`
	}{token, refresh})
	w.Write(response)
}
