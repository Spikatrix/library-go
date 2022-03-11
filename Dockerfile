FROM golang:1.17.7
WORKDIR /library-go
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build

FROM busybox:1.35
COPY --from=0 /library-go/library-go ./
EXPOSE 8080
CMD ["./library-go"]
