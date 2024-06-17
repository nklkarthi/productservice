package main

import (
	"log"
	"net/http"
	"productservice/db"
	"productservice/handlers"
	"productservice/middleware"

	"github.com/gorilla/mux"
)

func main() {
	db.Init()

	router := mux.NewRouter()

	router.HandleFunc("/api/login", handlers.Login).Methods("POST")
	router.HandleFunc("/api/products", handlers.GetProducts).Methods("GET")
	router.HandleFunc("/api/products", handlers.AddProduct).Methods("POST")
	router.HandleFunc("/api/products/{id:[0-9]+}", handlers.GetProduct).Methods("GET")

	router.Use(middleware.JwtAuthentication) // JWT middleware

	log.Fatal(http.ListenAndServe("localhost:8000", router))
}
