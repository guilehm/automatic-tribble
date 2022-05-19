package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"tribble/models"

	"github.com/gorilla/mux"
)

var frodo = &models.User{
	ID:           1,
	Name:         "frodo",
	Email:        "frodo@gmail.com",
	Password:     "pass",
	Token:        "token",
	RefreshToken: "pass",
	DateJoined:   time.Now(),
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
		t.Errorf("%s FAILED: want %v got %v", t.Name(), http.StatusOK, status)
	}

	expected, _ := json.Marshal(users)
	if rr.Body.String() != string(expected) {
		t.Errorf("%s FAILED: want %v got %v", t.Name(), expected, rr.Body.String())
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
		t.Errorf("%s FAILED: want %v got %v", t.Name(), http.StatusOK, status)
	}

	expected, _ := json.Marshal(frodo)
	if rr.Body.String() != string(expected) {
		t.Errorf("%s FAILED: want %v got %v", t.Name(), expected, rr.Body.String())
	}
}
