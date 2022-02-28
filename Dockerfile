FROM golang:1.17.7

WORKDIR /library-go

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
EXPOSE 8080
CMD ["go", "run", "main.go"]
