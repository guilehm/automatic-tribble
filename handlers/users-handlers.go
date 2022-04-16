package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"time"
	"tribble/models"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func GetUser(pool *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	sql := `SELECT id, name, date_joined FROM users WHERE id=$1`
	row := pool.QueryRow(context.Background(), sql, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.DateJoined)

	if err != nil {
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

	var userList []models.User
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
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "unable to decode request body")
		return
	}

	var id int

	sql := `INSERT INTO users (name, date_joined) VALUES ($1, $2) RETURNING id`
	err = pool.QueryRow(
		context.Background(), sql, user.Name, time.Now(),
	).Scan(&id)

	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}

	response, _ := json.Marshal(struct {
		ID string `json:"id"`
	}{ID: fmt.Sprintf("%d", id)})
	w.Write(response)
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
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusNotFound, "")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
