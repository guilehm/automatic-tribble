package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"
	"tribble/settings"
	"tribble/storages"

	"github.com/jackc/pgx/v4"

	"golang.org/x/crypto/bcrypt"

	"net/http"
	"tribble/models"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/jackc/pgconn"
)

var (
	validate = validator.New()
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func verifyPassword(userPassword string, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	return err == nil
}

func GetUserDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := storages.DB.GetUser(ctx, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandleDatabaseErrors(w, pgErr)
			return
		}
		if err == pgx.ErrNoRows {
			HandleApiErrors(w, http.StatusNotFound, "")
			return
		}
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}
	_, _ = w.Write(response)
}

func GetUserList(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	users, err := storages.DB.GetUserList(ctx)

	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}
	response, err := json.Marshal(users)
	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}
	_, _ = w.Write(response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var user *models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "unable to decode request body")
		return
	}

	if validationErr := validate.Struct(user); validationErr != nil {
		log.Println(validationErr.Error())
		HandleApiErrors(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	// TODO: validate unique email

	password, err := hashPassword(user.Password)
	if err != nil {
		HandleApiErrors(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	token, refresh, err := generateTokens(user.Email, user.ID)
	if err != nil {
		HandleApiErrors(w, http.StatusInternalServerError, "could not generate tokens")
		return
	}

	user.Password = password
	user.Token = token
	user.RefreshToken = refresh
	user.DateJoined = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err = storages.DB.CreateUser(ctx, *user)

	if err != nil {
		log.Println(err.Error())
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandleDatabaseErrors(w, pgErr)
			return
		}
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}

	response, _ := json.Marshal(user)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	userId, err := strconv.Atoi(r.Context().Value(settings.I).(string))
	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}

	var user models.User
	user.ID = userId
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusBadRequest, "")
		return
	}

	if validationErr := validate.StructPartial(user, user.Username); validationErr != nil {
		log.Println(validationErr.Error())
		HandleApiErrors(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// TODO: validate unique email

	_, err = storages.DB.UpdateUser(ctx, user)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandleDatabaseErrors(w, pgErr)
			return
		}
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	userId, err := strconv.Atoi(r.Context().Value(settings.I).(string))
	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = storages.DB.DeleteUser(ctx, userId)
	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusNotFound, "")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func Login(w http.ResponseWriter, r *http.Request) {

	var userLogin models.UserLogin

	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		HandleApiErrors(w, http.StatusBadRequest, "")
		return
	}
	if err := validate.Struct(userLogin); err != nil {
		HandleApiErrors(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := storages.DB.GetUserByEmail(ctx, userLogin.Email)
	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusNotFound, "")
		return
	}

	ok := verifyPassword(user.Password, userLogin.Password)
	if !ok {
		HandleApiErrors(w, http.StatusBadRequest, "invalid password")
		return
	}

	token, refresh, err := generateTokens(userLogin.Email, user.ID)
	if err != nil {
		HandleApiErrors(w, http.StatusInternalServerError, "could not update tokens")
		return
	}

	if err = storages.DB.UpdateUserTokens(ctx, user.ID, token, refresh); err != nil {
		HandleApiErrors(w, http.StatusInternalServerError, "could not update tokens")
		return
	}

	if err != nil {
		HandleApiErrors(w, http.StatusInternalServerError, "could not update tokens")
		return
	}

	response, _ := json.Marshal(struct {
		Id      int    `json:"id"`
		Token   string `json:"token"`
		Refresh string `json:"refresh_token"`
	}{user.ID, token, refresh})
	_, _ = w.Write(response)
}
