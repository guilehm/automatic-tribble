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

	"github.com/jackc/pgconn"
)

func CreatePlayer(w http.ResponseWriter, r *http.Request) {

	var player *models.Player

	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "unable to decode request body")
		return
	}

	if validationErr := validate.Struct(player); validationErr != nil {
		log.Println(validationErr.Error())
		HandleApiErrors(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	userId, err := strconv.Atoi(r.Context().Value(settings.I).(string))
	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}
	player.UserID = userId
	// TODO: position cannot be hardcoded
	player.PositionX = 8
	player.PositionY = 8

	// TODO: insert player on database

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	player, err = storages.DB.CreatePlayer(ctx, *player)

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

	response, _ := json.Marshal(player)
	_, _ = w.Write(response)
}

func GetPlayerList(w http.ResponseWriter, r *http.Request) {

	userId, err := strconv.Atoi(r.Context().Value(settings.I).(string))
	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	players, err := storages.DB.GetPlayerList(ctx, userId)

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

	response, err := json.Marshal(players)
	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}
	_, _ = w.Write(response)
}
