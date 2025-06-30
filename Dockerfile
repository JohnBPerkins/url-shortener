# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build binary
COPY . .
RUN go build -o shortener .

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/shortener .
COPY --from=builder /app/internal/web ./web/
RUN chmod +x /app/shortener

EXPOSE 50051
EXPOSE 8080
ENTRYPOINT ["./shortener"]