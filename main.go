package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"tribble/handlers"
	"tribble/middlewares"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/cors"
)

func main() {
	fmt.Println("hello world")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect to database: %v.", err))
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("successfully connected to database")

	r := mux.NewRouter()
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(r)

	r.HandleFunc("/users/{id}/", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUser(pool, w, r)
	}).Methods("GET")
	r.HandleFunc("/users/{id}/", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateUser(pool, w, r)
	}).Methods("PUT")
	r.HandleFunc("/users/{id}/", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteUser(pool, w, r)
	}).Methods("DELETE")
	
	r.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserList(pool, w, r)
	}).Methods("GET")
	r.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateUser(pool, w, r)
	}).Methods("POST")

	_ = http.ListenAndServe(":"+os.Getenv("PORT"), middlewares.LogRequest(handler))
}
