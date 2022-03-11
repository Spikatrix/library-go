# Library API

A simple library application written in Go

## Usage

```bash
git clone git@github.com:Spikatrix/library-go.git
cd library-go
echo "MONGODB_URI=<...>" > .env  # Replace <...> with the MongoDB URI link
echo "MONGODB_URI=<...>" > cmd/library/.env.test  # Replace <...> with the testing MongoDB URI link (DB for testing)
go mod download
go run main.go

# Add a new book
curl http://localhost:8080/newbook -XPOST -H 'Content-Type: application/json' -d '{"name": "Atomic Habits", "author": "James Clear"}'

# List books
curl http://localhost:8080/books

# List a single book
curl http://localhost:8080/book/<bookID>  # Replace <bookID> with the ID of the book
```

If you wish to run this with Docker, run the following:

```bash
docker build . -t spikatrix/library-go
docker run --name library-go -d -p 8080:8080 -e MONGODB_URI=<...> spikatrix/library-go # Replace <...> with the MongoDB URI link
# Use the curl commands shown above to access the library API
```

If you wish to run this in a `kind` kubernetes cluster, run the following:

```bash
cd manifests
kind create cluster --config kind-library-config.yml --name kind
kubectl apply -f library-deployment.yml,mongodb-configmap.yml,mongodb-statefulset.yml
# Wait until all resources are ready. Once ready, use the curl commands shown above to access the library API
```