package library

import (
	"Spikatrix/library-go/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ErrorHandler(w http.ResponseWriter, errorMsg string, err error) {
	log.Printf("%s: %+v", errorMsg, err)
	http.Error(w, "Something went wrong on our end; sorry about that!", http.StatusInternalServerError)
}

func NewBook(w http.ResponseWriter, req *http.Request) {
	bookItem := models.Book{}
	err := json.NewDecoder(req.Body).Decode(&bookItem)
	if err != nil {
		ErrorHandler(w, "Failed to decode request body", err)
		return
	}

	bookItem.Name, bookItem.Author = strings.TrimSpace(bookItem.Name), strings.TrimSpace(bookItem.Author)
	if bookItem.Name == "" || bookItem.Author == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid book. "+
			"It should be of the form {\"name\": \"<book name>\", \"author\": \"<author name>\"}")
		return
	}

	result, err := models.BookCollection.InsertOne(context.TODO(), bookItem)
	if err != nil {
		ErrorHandler(w, "Failed to insert book", err)
		return
	}

	bookItem.ID = result.InsertedID.(primitive.ObjectID)
	fmt.Fprintf(w, "Added book %+v\n", bookItem)
}

func Book(w http.ResponseWriter, req *http.Request) {
	bookID, err := primitive.ObjectIDFromHex(mux.Vars(req)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid book ID")
		return
	}

	var bookItem models.Book
	err = models.BookCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: bookID}}).Decode(&bookItem)
	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No book with id '%s' is available in the database", bookID)
		return
	} else if err != nil {
		ErrorHandler(w, "Failed to decode book", err)
		return
	}

	data, err := json.MarshalIndent(map[string]models.Book{"book": bookItem}, "", "  ")
	if err != nil {
		ErrorHandler(w, "Failed to marshal bookItem", err)
		return
	}

	fmt.Fprintln(w, string(data))
}

func Books(w http.ResponseWriter, req *http.Request) {
	cursor, err := models.BookCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		ErrorHandler(w, "Failed to query database", err)
		return
	}

	var bookList []models.Book
	if err = cursor.All(context.TODO(), &bookList); err != nil {
		ErrorHandler(w, "Failed to decode bookList", err)
		return
	}
	if bookList == nil {
		bookList = []models.Book{}
	}

	data, err := json.MarshalIndent(map[string][]models.Book{"books": bookList}, "", "  ")
	if err != nil {
		ErrorHandler(w, "Failed to marshal bookList", err)
		return
	}

	fmt.Fprintln(w, string(data))
}

func Root(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintln(w, "<h2>Welcome to the library application!</h2>")
	fmt.Fprintln(w, "<div>Here's what you can do here:</div>")
	fmt.Fprintln(w, "<ul>")
	fmt.Fprintln(w, "	<li>Visit <a href=\"/books\">/books</a> to get all books</li>")
	fmt.Fprintln(w, "	<li>Visit /books/&lt;id&gt; to get a particular book</li>")
	fmt.Fprintln(w, "	<li>Send a POST request to /newbook to add a new book</li>")
	fmt.Fprintln(w, "</ul>")
}
