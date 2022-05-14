package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"tribble/handlers"
	"tribble/middlewares"
	"tribble/storages"
	"tribble/storages/postgres"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	fmt.Println("hello world")

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	var err error
	storages.DB = postgres.GetPostgres()
	defer storages.DB.Close()
	// storages.DB, err = pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v.", err)
	}
	defer storages.DB.Close()

	// err = storages.DB.Ping(ctx)
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
	r.HandleFunc("/users/{id}/", handlers.GetUserDetail).Methods("GET")
	r.HandleFunc("/users/", handlers.CreateUser).Methods("POST")

	r.HandleFunc("/users/", middlewares.Authentication(handlers.UpdateUser)).Methods("PUT")
	r.HandleFunc("/users/", middlewares.Authentication(handlers.DeleteUser)).Methods("DELETE")

	r.HandleFunc("/users/validate/", handlers.ValidateToken).Methods("POST")
	r.HandleFunc("/users/refresh/", handlers.RefreshToken).Methods("POST")
	r.HandleFunc("/users/login/", handlers.Login).Methods("POST")

	r.HandleFunc("/players/", middlewares.Authentication(handlers.CreatePlayer)).Methods("POST")
	r.HandleFunc("/players/", middlewares.Authentication(handlers.GetPlayerList)).Methods("GET")

	_ = http.ListenAndServe(":"+os.Getenv("PORT"), middlewares.LogRequest(middlewares.SetHeaders(handler)))
}
