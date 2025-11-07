# Build stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bookmanager-go ./cmd/main.go

# Runtime stage
FROM debian:bookworm-slim
WORKDIR /app
RUN apt-get update && apt-get install -y sqlite3 && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/bookmanager-go .
COPY --from=builder /app/internal/view ./internal/view
EXPOSE 8080
CMD ["./bookmanager-go"]
