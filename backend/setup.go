package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func setupMongoDB() {
	// Set up MongoDB connection
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb://mongo-statefulset-0.mongo-service,mongo-statefulset-1.mongo-service," +
		"mongo-statefulset-2.mongo-service/?replicaSet=yourReplicaSetName\n").SetServerAPIOptions(serverAPI)

	var err error
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
}

func closeMongoDB() {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	cancel()
}
