package library

import (
	"Spikatrix/library-go/pkg/db"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type server struct {
	bookCollection *mongo.Collection
}

func StartLibraryServer() {
	bookCollection, dbClose, err := db.SetupDB(db.DbURI())

	if err != nil {
		log.Fatalf("Failed to setup DB: %+v", err)
	}

	defer dbClose()

	libraryServer := server{
		bookCollection: bookCollection,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", Root).Methods("GET")
	r.HandleFunc("/books", libraryServer.GetBooks).Methods("GET")
	r.HandleFunc("/book/{id}", libraryServer.GetBookByID).Methods("GET")
	r.HandleFunc("/newbook", libraryServer.CreateNewBook).Methods("POST")
	http.Handle("/", r)

	port := "8080"
	log.Println("Server is ready at http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
