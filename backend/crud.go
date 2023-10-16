package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

type Item struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  string             `json:"name"`
	Price float64            `json:"price"`
}

func validateItem(item *Item) error {
	if item.Name == "" {
		return errors.New("name field is required")
	}
	if item.Price <= 0 {
		return errors.New("price should be greater than 0")
	}
	return nil
}

func RemoveItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database("ecommerce").Collection("items")
	res, err := coll.DeleteOne(ctx, bson.M{"_id": itemID})
	if err != nil {
		http.Error(w, "Error deleting item", http.StatusInternalServerError)
		return
	}
	if res.DeletedCount == 0 {
		http.Error(w, "No item found to delete", http.StatusNotFound)
		return
	}

	log.Println("Deleted item with id ", itemID.String())
	w.WriteHeader(http.StatusOK)
}

// UpdateItem updates an item based on its ID.
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Bad Input", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.Database("ecommerce").Collection("items").
		UpdateOne(ctx, bson.M{"_id": itemID}, bson.M{"$set": item})
	if err != nil || res.MatchedCount == 0 {
		http.Error(w, "Update Failed", http.StatusInternalServerError)
		return
	}

	log.Println("Updated item with id ", itemID.String())
	w.WriteHeader(http.StatusOK)
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	var items []Item

	// Get a new context with a timeout for the MongoDB operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database("ecommerce").Collection("items")
	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to retrieve items", http.StatusInternalServerError)
		log.Printf("Get items error: %v", err)
		return
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var item Item
		if err := cur.Decode(&item); err != nil {
			http.Error(w, "Failed to decode item", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := cur.Err(); err != nil {
		http.Error(w, "Cursor error", http.StatusInternalServerError)
		log.Printf("Cursor error: %v", err)
		return
	}

	jsonResponse, err := json.Marshal(items)
	if err != nil {
		http.Error(w, "Failed to marshal items", http.StatusInternalServerError)
		return
	}

	log.Println("Returning items")

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Fatal("Error while writing response: ", err)
	}
}

func AddItem(w http.ResponseWriter, r *http.Request) {
	var item Item

	// Decode the request body into the item struct
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate the item
	if err := validateItem(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get a new context with a timeout for the MongoDB operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("ecommerce").Collection("items")
	_, err := collection.InsertOne(ctx, item)
	if err != nil {
		log.Printf("Failed to insert record: %v", err)
		http.Error(w, "Failed to save to DB", http.StatusInternalServerError)
		return
	}

	message := []byte(fmt.Sprintf("New item added with ID: %s", item.ID))
	publishToRabbitMQ(message)

	log.Println("Item added ", item)
	w.WriteHeader(http.StatusOK)
}

// GetItemByID retrieves a specific item based on its ID.
func GetItemByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database("ecommerce").Collection("items")

	var item Item
	err = coll.FindOne(ctx, bson.M{"_id": itemID}).Decode(&item)
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	log.Println("Item found ", item)
	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
	}
}
