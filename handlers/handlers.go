package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
	"time"
	"tribble/models"
)

func CreateUser(pool *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response, _ := json.Marshal(struct {
			Message string `json:"message"`
		}{Message: "Unable to decode request body"})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}


	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int

	sql := `INSERT INTO users (name, date_joined) VALUES ($1, $2) RETURNING id`
	err = pool.QueryRow(
		ctx,
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
