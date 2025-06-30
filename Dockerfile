# Start from the official Go image for building the app
FROM golang:1.23-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary statically
RUN CGO_ENABLED=0 GOOS=linux go build -o blog-api ./cmd/main.go

# Use a minimal image for the final container
FROM scratch

# Copy the built binary from the builder stage
COPY --from=builder /app/blog-api /blog-api

# Expose port your app listens on (e.g., 8080)
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/blog-api"]
