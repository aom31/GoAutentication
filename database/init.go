package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DbInstance() *mongo.Client {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MongoDb := os.Getenv("MONGODB_URL")
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = mongoClient.Connect(context)
	if err != nil {
		log.Fatal(err)
	}

	return mongoClient
}

var Client *mongo.Client = DbInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	return collection
}
