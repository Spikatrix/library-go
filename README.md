# Library API

A simple library application written in Go

## Usage

```bash
git clone git@github.com:Spikatrix/library-go.git
cd library-go
echo "MONGODB_URI=<...>" > .env  # Replace <...> with the MongoDB URI link
go mod download
go run library.go

# Add a new book
curl http://localhost:8080/newbook -XPOST -H 'Content-Type: application/json' -d '{"name": "Atomic Habits", "author": "James Clear"}'

# List books
curl http://localhost:8080/books

# List a single book
curl http://localhost:8080/book/<bookID>  # Replace <bookID> with the ID of the book
```

If you wish to run this with Docker, run the following in place of `go mod download; go run library.go`:

```bash
docker build . -t spikatrix/library-go
docker run --name library-go -d -p 8080:8080 spikatrix/library-go
```