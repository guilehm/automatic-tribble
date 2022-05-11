package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tribble/models"
	"tribble/settings"
)

func CreatePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

}
