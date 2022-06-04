package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
	"tribble/models"
	"tribble/settings"
	"tribble/storages"

	"github.com/dgrijalva/jwt-go"
)

type SignedDetails struct {
	Username string
	ID       string
	jwt.StandardClaims
}

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

func generateTokens(username string, userId int) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Username: username,
		ID:       strconv.Itoa(userId),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(settings.AccessTokenLifetime).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(settings.RefreshTokenLifetime).Unix(),
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

	ok, err := storages.DB.ValidateToken(ctx, tokens.RefreshToken)
	if err != nil || !ok {
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

	user, err := storages.DB.GetUserByRefresh(ctx, tokens.RefreshToken)
	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusNotFound, "")
		return
	}

	token, refresh, err := generateTokens(user.Username, user.ID)
	if err != nil {
		HandleApiErrors(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = storages.DB.UpdateUserTokens(ctx, user.ID, token, refresh)
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
