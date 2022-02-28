package library

import (
	"Spikatrix/library-go/pkg/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetBooks(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()

	dbClose, err := models.SetupDB(models.TestDbURI())
	if err != nil {
		t.Errorf("DB setup failed: %+v", err)
		return
	}
	defer dbClose()

	Books(w, req)

	if statusCode := w.Result().StatusCode; statusCode != http.StatusOK {
		t.Errorf("Get books expected status code %d, got %d. Response: '%s'",
			http.StatusOK, statusCode, w.Body.String())
	}
}

func TestGetBook(t *testing.T) {
	tt := []struct {
		name               string
		bookID             string
		expectedStatusCode int
	}{
		{
			name:               "Invalid bookID",
			bookID:             "abcd",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Book not found",
			bookID:             "000000000000000000000000",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/book/{id}", nil)
			w := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{"id": tc.bookID})

			dbClose, err := models.SetupDB(models.TestDbURI())
			if err != nil {
				t.Errorf("DB setup failed: %+v", err)
				return
			}
			defer dbClose()

			Book(w, req)

			if statusCode := w.Result().StatusCode; statusCode != tc.expectedStatusCode {
				t.Errorf("Get book expected status code %d, got %d (Response: '%s')",
					tc.expectedStatusCode, statusCode, w.Body.String())
			}
		})
	}
}

func TestAddBook(t *testing.T) {
	book := models.Book{Name: "Test book", Author: "Test author"}
	bookJson, err := json.Marshal(book)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/newbook", strings.NewReader(string(bookJson)))
	w := httptest.NewRecorder()

	dbClose, err := models.SetupDB(models.TestDbURI())
	if err != nil {
		t.Errorf("DB setup failed: %+v", err)
		return
	}
	defer dbClose()

	NewBook(w, req)

	if statusCode := w.Result().StatusCode; statusCode != http.StatusOK {
		t.Errorf("New book expected status code %d, got %d. Response: '%s'",
			http.StatusOK, statusCode, w.Body.String())
	}
}
