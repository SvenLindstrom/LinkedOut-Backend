FROM golang:1.25.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o linkedOutServer .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/linkedOutServer .
ENV GOOGLE_APPLICATION_CREDENTIALS=./firebase_key.json

CMD ["./linkedOutServer"]
