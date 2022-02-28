package models

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var BookCollection *mongo.Collection

const dbName = "library"
const collectionName = "books"

type Book struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name   string             `bson:"name" json:"name"`
	Author string             `bson:"author" json:"author"`
}

func SetupDB(dbURI string) (func(), error) {
	var err error

	if dbURI == "" {
		err = fmt.Errorf("MongoDB URI is not specified; have you set the MONGODB_URI environment variable?")
		return nil, err
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURI))
	if err != nil {
		return nil, err
	}

	dbClose := func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}

	BookCollection = client.Database(dbName).Collection(collectionName)

	return dbClose, err
}

func DbURI() string {
	godotenv.Load() // Loads .env file data
	return os.Getenv("MONGODB_URI")
}

func TestDbURI() string {
	godotenv.Load(".env.test") // Loads .env.test file data
	return os.Getenv("MONGODB_URI")
}
