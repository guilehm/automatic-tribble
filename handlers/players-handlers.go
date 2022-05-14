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

	"github.com/jackc/pgconn"
)

func CreatePlayer(w http.ResponseWriter, r *http.Request) {

	var player models.Player

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

	sql := `INSERT INTO players (user_id, xp, sprite, position_x, position_y)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id`

	var playerID int
	err = db.DB.QueryRow(
		ctx,
		sql,
		player.UserID,
		player.XP,
		player.Sprite,
		player.PositionX,
		player.PositionY,
	).Scan(&playerID)

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
		Id int `json:"id"`
	}{playerID})
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

	sql := `SELECT name, xp, sprite, position_x, position_y FROM players WHERE user_id=$1`
	rows, err := db.DB.Query(ctx, sql, userId)

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

	players := make([]models.Player, 0)
	for rows.Next() {
		var player models.Player
		err = rows.Scan(
			&player.Name,
			&player.XP,
			&player.Sprite,
			&player.PositionX,
			&player.PositionY,
		)
		if err != nil {
			log.Printf("could not scan player: %v", err.Error())
			HandleApiErrors(w, http.StatusInternalServerError, "")
			return
		}
		players = append(players, player)
	}

	response, err := json.Marshal(players)
	if err != nil {
		log.Println(err.Error())
		HandleApiErrors(w, http.StatusInternalServerError, "")
		return
	}
	_, _ = w.Write(response)
}
