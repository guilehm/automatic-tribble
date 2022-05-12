package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"tribble/db"

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

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `SELECT id, name, date_joined FROM users WHERE id=$1`

	var user models.User
	if err := db.DB.QueryRow(ctx, sql, id).Scan(&user.ID, &user.Name, &user.DateJoined); err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusNotFound, "")
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
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `SELECT id, name, date_joined FROM users`
	rows, err := db.DB.Query(ctx, sql)

	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}

	userList := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Name, &user.DateJoined)
		if err != nil {
			HandleApiErrors(w, http.StatusInternalServerError, "")
			return
		}
		userList = append(userList, user)
	}

	response, err := json.Marshal(userList)
	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}
	_, _ = w.Write(response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
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
	var id int

	sql := `INSERT INTO users (name, email, date_joined, password, token, refresh_token) 
			VALUES ($1, $2, $3, $4, $5, $6) 
			RETURNING id`
	err = db.DB.QueryRow(
		ctx,
		sql,
		user.Name,
		user.Email,
		user.DateJoined,
		user.Password,
		user.Token,
		user.RefreshToken,
	).Scan(&id)

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

	response, _ := json.Marshal(struct {
		ID string `json:"id"`
	}{ID: fmt.Sprintf("%d", id)})
	_, _ = w.Write(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusBadRequest, "")
		return
	}

	if validationErr := validate.StructPartial(user, user.Name); validationErr != nil {
		log.Println(validationErr.Error())
		HandleApiErrors(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `UPDATE users SET name=$2 WHERE id=$1`
	res, err := db.DB.Exec(ctx, sql, id, user.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			HandleDatabaseErrors(w, pgErr)
			return
		}
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}

	if rowsAffected := res.RowsAffected(); rowsAffected == 0 {
		HandleApiErrors(w, http.StatusNotFound, "")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sql := `DELETE FROM users where id=$1`
	res, err := db.DB.Exec(ctx, sql, id)
	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}
	if rowsAffected := res.RowsAffected(); rowsAffected == 0 {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusNotFound, "")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

	var id int
	var password string
	sql := `SELECT id, password FROM users WHERE email=$1`

	if err := db.DB.QueryRow(ctx, sql, userLogin.Email).Scan(&id, &password); err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusNotFound, "")
		return
	}

	ok := verifyPassword(password, userLogin.Password)
	if !ok {
		HandleApiErrors(w, http.StatusBadRequest, "invalid password")
		return
	}

	token, refresh, err := generateTokens(userLogin.Email, id)
	if err != nil {
		HandleApiErrors(w, http.StatusInternalServerError, "could not update tokens")
		return
	}

	sql = `UPDATE users SET token=$1, refresh_token=$2 WHERE id=$3`
	_, err = db.DB.Exec(ctx, sql, token, refresh, id)

	if err != nil {
		HandleApiErrors(w, http.StatusInternalServerError, "could not update tokens")
		return
	}

	response, _ := json.Marshal(struct {
		Id      int    `json:"id"`
		Token   string `json:"token"`
		Refresh string `json:"refresh_token"`
	}{id, token, refresh})
	_, _ = w.Write(response)

}
