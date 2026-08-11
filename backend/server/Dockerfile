# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Install dependencies (if we had them)
# COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Production stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]
