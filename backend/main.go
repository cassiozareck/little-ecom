package main

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

var client *mongo.Client
var cancel context.CancelFunc

func main() {
	setupMongoDB()
	defer closeMongoDB()
	log.Println("Mongo setup done!")

	err := setupRabbitMQ()
	if err != nil {
		log.Fatal("Error setting rabbitmq: ", err)
	}

	r := mux.NewRouter()
	r.Use(corsMiddleware)
	r.Use(loggingMiddleware)

	r.HandleFunc("/item/{id:[a-f\\d]{24}}", GetItemByID).Methods("GET")
	r.HandleFunc("/item", AddItem).Methods("POST")
	r.HandleFunc("/item/{id:[a-f\\d]{24}}", RemoveItem).Methods("DELETE")
	r.HandleFunc("/item/{id:[a-f\\d]{24}}", UpdateItem).Methods("PUT")
	r.HandleFunc("/items/{owner}", GetItemsByOwner).Methods("GET")
	r.HandleFunc("/items", GetItems).Methods("GET")
	r.HandleFunc("/buy/{id:[a-f\\d]{24}}", BuyItem).Methods("POST")

	log.Println("Routers done!")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
