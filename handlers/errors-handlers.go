package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgconn"
)

func HandleApiErrors(w http.ResponseWriter, status int, message string) {
	if message == "" {
		message = http.StatusText(status)
	}

	response, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{message})
	w.WriteHeader(status)
	w.Write(response)
}

func HandleDatabaseErrors(w http.ResponseWriter, pgErr *pgconn.PgError) {
	log.Printf("PgError: code: %v message: %v", pgErr.Code, pgErr.Message)
	switch pgErr.Code {
	case "23505":
		// unique constraint violated
		HandleApiErrors(w, http.StatusBadRequest, pgErr.Message)
		return
	case "22001":
		// value too long for type character
		HandleApiErrors(w, http.StatusBadRequest, "value too long for type character")
		return
	default:
		HandleApiErrors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}
