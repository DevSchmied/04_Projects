# Build stage
FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# --- Create env file inside builder ---
RUN cp internal/config/app.env.example internal/config/app.env

RUN go build -o bookmanager-go ./cmd/main.go

# Runtime stage
FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update && apt-get install -y sqlite3 && rm -rf /var/lib/apt/lists/*

# Copy binary
COPY --from=builder /app/bookmanager-go .

# Copy views
COPY --from=builder /app/internal/view ./internal/view

# Copy config folder including generated app.env
COPY --from=builder /app/internal/config ./internal/config

EXPOSE 8080
CMD ["./bookmanager-go"]
