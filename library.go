package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type book struct {
	Name   string `json:"name"`
	Author string `json:"author"`
}

var bookList []book

func errorHandler(w http.ResponseWriter, errorMsg string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "Something went wrong on our end; sorry about that!")
	log.Printf("%s: %+v", errorMsg, err)
}

func newbook(w http.ResponseWriter, req *http.Request) {
	if reqType := req.Method; reqType != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Add("Allow", "POST")
		fmt.Fprintln(w, "Sorry, only POST requests are supported on this endpoint")
		return
	}

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorHandler(w, "Failed to read request body", err)
		return
	}

	bookItem := book{}
	json.Unmarshal(reqBody, &bookItem)

	if bookItem.Name == "" || bookItem.Author == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid book. "+
			"It should be of the form {\"name\": \"<book name>\", \"author\": \"<author name>\"}")
		return
	}

	bookList = append(bookList, bookItem)
	fmt.Fprintf(w, "Added book %+v\n", bookItem)
}

func books(w http.ResponseWriter, req *http.Request) {
	data, err := json.Marshal(map[string][]book{"books": bookList})
	if err != nil {
		errorHandler(w, "Failed to marshal bookList", err)
		return
	}

	fmt.Fprintln(w, string(data))
}

func root(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintln(w, "<h2>Welcome to the library application!</h2>")
	fmt.Fprintln(w, "<div>Here's what you can do here:</div>")
	fmt.Fprintln(w, "<ul>")
	fmt.Fprintln(w, "	<li>Visit <a href=\"/books\">/books</a> to see all books</li>")
	fmt.Fprintln(w, "	<li>Send a POST request to /newbook to add a new book</li>")
	fmt.Fprintln(w, "</ul>")
}

func main() {
	bookList = make([]book, 0)

	http.HandleFunc("/", root)
	http.HandleFunc("/books", books)
	http.HandleFunc("/newbook", newbook)

	err := http.ListenAndServe(":8080", nil)
	panic(err)
}
