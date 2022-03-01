package library

import (
	"Spikatrix/library-go/pkg/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func StartLibraryServer() {
	dbClose, err := models.SetupDB(models.DbURI())

	// TODO: panic is not the best way to handle the error, try using any logger
	if err != nil {
		panic(err)
	}

	defer dbClose()

	r := mux.NewRouter()
	r.HandleFunc("/", Root).Methods("GET")
	// TODO: handler names can be improved, these seem a bit confusing
	r.HandleFunc("/books", Books).Methods("GET")
	r.HandleFunc("/book/{id}", Book).Methods("GET")
	r.HandleFunc("/newbook", NewBook).Methods("POST")
	http.Handle("/", r)

	// TODO: inside a func, const is irrelevant, a simple var would suffice. Here, the port can be const or better declare it as env var
	const port = "8080"
	log.Println("Server is ready at http://localhost:" + port)
	panic(http.ListenAndServe(":"+port, nil))
}
