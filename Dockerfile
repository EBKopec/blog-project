# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o blog-api ./cmd/main.go

# Stage 2: Minimal image
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/blog-api /app/blog-api

RUN chmod +x /app/blog-api
RUN ls -la /app

EXPOSE 8080

CMD ["/app/blog-api"]
