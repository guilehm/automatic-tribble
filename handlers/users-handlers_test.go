package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
	"tribble/models"
	"tribble/settings"
	"tribble/storages"
	"tribble/storages/postgres"

	"gopkg.in/go-playground/assert.v1"

	"github.com/gorilla/mux"
)

var pg = postgres.GetPostgres()

var frodoEmail = "frodo@gmail.com"
var frodo = &models.User{
	ID:           1,
	Username:     "frodo",
	Email:        &frodoEmail,
	Password:     "password",
	Token:        "token",
	RefreshToken: "refresh",
	DateJoined:   time.Now(),
}

func TestSetup(t *testing.T) {
	storages.DB = pg
	t.Log("setting postgres as default database")
}

func TestCreateUserHandler(t *testing.T) {
	url := "/users/"

	payload, err := json.Marshal(frodo)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/users/", CreateUser).Methods("POST")

	var count int
	sql := `SELECT COUNT(*) FROM users`
	if err = pg.DB.QueryRow(context.Background(), sql).Scan(&count); err != nil {
		t.Fatalf("%s FAILED: could not count users", t.Name())
	}

	assert.Equal(t, count, 0)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("%s FAILED: want %d got %d", t.Name(), http.StatusCreated, status)
	}
	if err = pg.DB.QueryRow(context.Background(), sql).Scan(&count); err != nil {
		t.Fatalf("%s FAILED: could not count users", t.Name())
	}

	assert.Equal(t, count, 1)
}

func TestUpdateUserHandler(t *testing.T) {
	// TODO: test fields that cannot be updated

	url := "/users/"
	newUsername := "FRODO"

	assert.Equal(t, frodo.Username, "frodo")
	frodo.Username = newUsername
	payload, err := json.Marshal(frodo)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	// frodo authentication
	ctx := req.Context()
	ctx = context.WithValue(ctx, settings.E, frodo.Email)
	ctx = context.WithValue(ctx, settings.I, strconv.Itoa(frodo.ID))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/users/", UpdateUser).Methods("PUT")

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("%s FAILED: want %d got %d", t.Name(), http.StatusNoContent, status)
	}

	var username string
	sql := `SELECT username FROM users WHERE username=$1`
	if err = pg.DB.QueryRow(context.Background(), sql, frodo.Username).Scan(&username); err != nil {
		t.Fatalf("%s FAILED: could retrieve user", t.Name())
	}

	assert.Equal(t, username, newUsername)
}

func TestGetUserListHandler(t *testing.T) {

	users := []*models.User{frodo}

	req, err := http.NewRequest("GET", "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUserList)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("%s FAILED: want %d got %d", t.Name(), http.StatusOK, status)
	}

	for _, user := range users {
		// omit fields
		user.Token = ""
		user.RefreshToken = ""
		user.Password = ""
	}

	expected, _ := json.Marshal(users)
	if rr.Body.String() != string(expected) {
		t.Errorf("%s FAILED: want %s got %s", t.Name(), expected, rr.Body.String())
	}
}

func TestGetUserDetailHandler(t *testing.T) {

	url := fmt.Sprintf("/users/%v/", frodo.ID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := mux.NewRouter()
	handler.HandleFunc("/users/{id}/", GetUserDetail).Methods("GET")
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("%s FAILED: want %d got %d", t.Name(), http.StatusOK, status)
	}

	expected, _ := json.Marshal(frodo)
	if rr.Body.String() != string(expected) {
		t.Errorf("%s FAILED: want %s got %s", t.Name(), expected, rr.Body.String())
	}
}

// TODO: create tests to validate unique email
