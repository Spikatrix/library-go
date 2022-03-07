package db

import (
	"Spikatrix/library-go/pkg/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetBookByID(bookCollection *mongo.Collection, bookID primitive.ObjectID) (models.Book, error) {
	var bookItem models.Book
	err := bookCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: bookID}}).Decode(&bookItem)
	if err == mongo.ErrNoDocuments {
		return bookItem, errors.New(ErrBookNotFound)
	} else if err != nil {
		return bookItem, errors.New(ErrBookDecodeFailed)
	}

	return bookItem, nil
}

func GetBooks(bookCollection *mongo.Collection) ([]models.Book, error) {
	var bookList []models.Book

	cursor, err := bookCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return bookList, errors.New(ErrBookQueryFailed)
	}

	if err := cursor.All(context.TODO(), &bookList); err != nil {
		return bookList, errors.New(ErrBookDecodeFailed)
	}

	if bookList == nil {
		// Book list is empty, set it to an empty list instead of nil
		bookList = []models.Book{}
	}

	return bookList, nil
}

func CreateNewBook(bookCollection *mongo.Collection, bookItem models.Book) (models.Book, error) {
	result, err := bookCollection.InsertOne(context.TODO(), bookItem)
	if err != nil {
		return bookItem, errors.New(ErrBookInsertFailed)
	}

	bookItem.ID = result.InsertedID.(primitive.ObjectID)

	return bookItem, nil
}
