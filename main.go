package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"tribble/db"
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

	var err error
	db.DB, err = pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect to database: %v.", err))
	}
	defer db.DB.Close()

	err = db.DB.Ping(ctx)
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

	r.HandleFunc("/users/", handlers.GetUserList).Methods("GET")
	r.HandleFunc("/users/", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/users/validate/", handlers.ValidateToken).Methods("POST")
	r.HandleFunc("/users/refresh/", handlers.RefreshToken).Methods("POST")

	r.HandleFunc("/users/{id}/", handlers.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}/", handlers.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}/", handlers.DeleteUser).Methods("DELETE")

	_ = http.ListenAndServe(":"+os.Getenv("PORT"), middlewares.LogRequest(handler))
}
