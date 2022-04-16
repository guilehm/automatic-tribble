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
		response, _ := json.Marshal(struct {
			Message string `json:"message"`
		}{Message: "user not found"})
		log.Println(err.Error())

		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		response, _ := json.Marshal(struct {
			Message string `json:"message"`
		}{Message: "could not process response"})
		log.Println(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
		return
	}
	w.Write(response)
}

func CreateUser(pool *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response, _ := json.Marshal(struct {
			Message string `json:"message"`
		}{Message: "unable to decode request body"})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	var id int

	sql := `INSERT INTO users (name, date_joined) VALUES ($1, $2) RETURNING id`
	err = pool.QueryRow(
		context.Background(),
		sql,
		user.Name, time.Now(),
	).Scan(&id)

	if err != nil {
		response, _ := json.Marshal(struct {
			Message string `json:"message"`
		}{Message: err.Error()})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
		return
	}

	response, _ := json.Marshal(struct {
		ID string `json:"id"`
	}{ID: fmt.Sprintf("%d", id)})
	w.Write(response)

}
