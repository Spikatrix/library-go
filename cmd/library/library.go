package library

import (
	"Spikatrix/library-go/pkg/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func StartLibraryServer() {
	dbClose, err := models.SetupDB(models.DbURI())
	if err != nil {
		panic(err)
	}

	defer dbClose()

	r := mux.NewRouter()
	r.HandleFunc("/", Root).Methods("GET")
	r.HandleFunc("/books", Books).Methods("GET")
	r.HandleFunc("/book/{id}", Book).Methods("GET")
	r.HandleFunc("/newbook", NewBook).Methods("POST")
	http.Handle("/", r)

	const port = "8080"
	log.Println("Server is ready at http://localhost:" + port)
	panic(http.ListenAndServe(":"+port, nil))
}
