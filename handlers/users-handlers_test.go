package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"tribble/models"
	"tribble/storages"
	"tribble/storages/postgres"

	"gopkg.in/go-playground/assert.v1"

	"github.com/gorilla/mux"
)

var frodo = &models.User{
	ID:           1,
	Name:         "frodo",
	Email:        "frodo@gmail.com",
	Password:     "password",
	Token:        "token",
	RefreshToken: "pass",
	DateJoined:   time.Now(),
}

var pg = postgres.GetPostgres()

func TestSetup(t *testing.T) {
	storages.DB = pg
	t.Log("setting postgres as default database")
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
		t.Errorf("%s FAILED: want %v got %v", t.Name(), http.StatusCreated, status)
	}
	if err = pg.DB.QueryRow(context.Background(), sql).Scan(&count); err != nil {
		t.Fatalf("%s FAILED: could not count users", t.Name())
	}

	assert.Equal(t, count, 1)
}
