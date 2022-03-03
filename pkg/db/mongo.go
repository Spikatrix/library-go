package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName         = "library"
	collectionName = "books"
)

func SetupDB(dbURI string) (*mongo.Collection, func(), error) {
	var err error

	if dbURI == "" {
		err = fmt.Errorf("MongoDB URI is not specified; have you set the MONGODB_URI environment variable?")
		return nil, nil, err
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURI))
	if err != nil {
		return nil, nil, err
	}

	dbClose := func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}

	bookCollection := initBookCollection(client)

	return bookCollection, dbClose, err
}

func initBookCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(dbName).Collection(collectionName)
}

func DbURI() string {
	godotenv.Load() // Loads .env file data
	return os.Getenv("MONGODB_URI")
}

func TestDbURI() string {
	godotenv.Load(".env.test") // Loads .env.test file data
	return os.Getenv("MONGODB_URI")
}
