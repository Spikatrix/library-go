# Library API

A simple library application written in Go

## Usage

```
git clone git@github.com:Spikatrix/library-go.git
cd library-go
go mod download
echo "MONGODB_URI=<...>" > .env  # Replace <...> with the MongoDB URI link
go run library.go

# Add a new book
curl http://localhost:8080/newbook -XPOST -H 'Content-Type: application/json' -d '{"name": "Atomic Habits", "author": "James Clear"}'

# List books
curl http://localhost:8080/books
```
