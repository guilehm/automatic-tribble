package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"tribble/models"
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

	// player.UserID = r.Context().Value(settings.I).(int)
	// TODO: position cannot be hardcoded
	player.PositionX = 8
	player.PositionY = 8

}
