# Tahap 1: Build stage
FROM golang:1.26-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod dan go.sum agar bisa download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy semua file source code
COPY . .

# Build aplikasi menjadi file biner bernama 'main'
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Tahap 2: Final stage (menjalankan biner yang sudah di-build)
FROM alpine:latest

WORKDIR /root/

# Copy hasil build dari tahap 1
COPY --from=builder /app/main .

# Expose port
EXPOSE 8081

# Jalankan aplikasi
CMD ["./main"]