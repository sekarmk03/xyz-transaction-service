FROM golang:1.23.1-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/server

RUN go build -o /app/main .

EXPOSE 50052

CMD ["/app/main"]
