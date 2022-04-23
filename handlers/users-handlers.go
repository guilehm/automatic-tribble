package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"net/http"
	"time"
	"tribble/models"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

var validate = validator.New()

func GetUser(pool *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	sql := `SELECT id, name, date_joined FROM users WHERE id=$1`
	row := pool.QueryRow(context.Background(), sql, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.DateJoined); err != nil {
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
	w.Write(response)
}

func GetUserList(pool *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sql := `SELECT id, name, date_joined FROM users`
	rows, err := pool.Query(context.Background(), sql)

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
	w.Write(response)
}

func CreateUser(pool *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
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

	var id int

	sql := `INSERT INTO users (name, email, date_joined) VALUES ($1, $2, $3) RETURNING id`
	err := pool.QueryRow(
		context.Background(), sql, user.Name, user.Email, time.Now(),
	).Scan(&id)

	if err != nil {
		log.Println(err.Error())
		// TODO: improve response for unique constraint violated
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}

	response, _ := json.Marshal(struct {
		ID string `json:"id"`
	}{ID: fmt.Sprintf("%d", id)})
	w.Write(response)
}

func UpdateUser(pool *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
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

	sql := `UPDATE users SET name=$2 WHERE id=$1`
	res, err := pool.Exec(context.Background(), sql, id, user.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("PgError: code: %v message: %v", pgErr.Code, pgErr.Message)
			switch pgErr.Code {
			case "23505":
				// unique constraint violated
				HandleApiErrors(w, http.StatusBadRequest, "this name already exists")
				return
			case "22001":
				// value too long for type character
				HandleApiErrors(w, http.StatusBadRequest, "value too long for type character")
				return
			}
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

func DeleteUser(pool *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	sql := `DELETE FROM users where id=$1`
	res, err := pool.Exec(context.Background(), sql, id)
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
