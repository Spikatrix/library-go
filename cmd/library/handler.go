package library

import (
	"Spikatrix/library-go/pkg/db"
	"Spikatrix/library-go/pkg/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func errorHandler(w http.ResponseWriter, errorMsg string, err error) {
	log.Printf("%s: %+v", errorMsg, err)
	http.Error(w, "Something went wrong on our end; sorry about that!", http.StatusInternalServerError)
}

func (server *server) CreateNewBook(w http.ResponseWriter, req *http.Request) {
	bookItem := models.Book{}
	err := json.NewDecoder(req.Body).Decode(&bookItem)
	if err != nil {
		errorHandler(w, "Failed to decode request body", err)
		return
	}

	bookItem.Name, bookItem.Author = strings.TrimSpace(bookItem.Name), strings.TrimSpace(bookItem.Author)
	if bookItem.Name == "" || bookItem.Author == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid book. "+
			"It should be of the form {\"name\": \"<book name>\", \"author\": \"<author name>\"}")
		return
	}

	bookItem, err = db.CreateNewBook(server.bookCollection, bookItem)
	if err.Error() == db.ErrBookInsertFailed {
		errorHandler(w, "Failed to insert book", err)
		return
	} else if err != nil {
		errorHandler(w, "db.CreateNewBook failed with an unknown error", err)
		return
	}

	fmt.Fprintf(w, "Added book %+v\n", bookItem)
}

func (server *server) GetBookByID(w http.ResponseWriter, req *http.Request) {
	urlQueryID := mux.Vars(req)["id"]
	bookID, err := primitive.ObjectIDFromHex(urlQueryID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid book ID")
		return
	}

	bookItem, err := db.GetBookByID(server.bookCollection, bookID)
	if err.Error() == db.ErrBookNotFound {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No book with id '%s' is available in the database", urlQueryID)
		return
	} else if err.Error() == db.ErrBookDecodeFailed {
		errorHandler(w, "Failed to decode book", err)
		return
	} else if err != nil {
		errorHandler(w, "db.GetBookByID failed with an unknown error", err)
		return
	}

	data, err := json.MarshalIndent(map[string]models.Book{"book": bookItem}, "", "  ")
	if err != nil {
		errorHandler(w, "Failed to marshal bookItem", err)
		return
	}

	fmt.Fprintln(w, string(data))
}

func (server *server) GetBooks(w http.ResponseWriter, req *http.Request) {
	bookList, err := db.GetBooks(server.bookCollection)
	if err.Error() == db.ErrBookQueryFailed {
		errorHandler(w, "Failed to query database", err)
		return
	} else if err.Error() == db.ErrBookDecodeFailed {
		errorHandler(w, "Failed to decode bookList", err)
		return
	} else if err != nil {
		errorHandler(w, "db.GetBooks failed with an unknown error", err)
		return
	}

	data, err := json.MarshalIndent(map[string][]models.Book{"books": bookList}, "", "  ")
	if err != nil {
		errorHandler(w, "Failed to marshal bookList", err)
		return
	}

	fmt.Fprintln(w, string(data))
}

func Root(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	introMessages := []string{
		"<h2>Welcome to the library application!</h2>",
		"<div>Here's what you can do here:</div>",
		"<ul>",
		"	<li>Visit <a href=\"/books\">/books</a> to get all books</li>",
		"	<li>Visit /books/&lt;id&gt; to get a particular book</li>",
		"	<li>Send a POST request to /newbook to add a new book</li>",
		"</ul>",
	}

	for _, introMessage := range introMessages {
		fmt.Fprintln(w, introMessage)
	}
}
