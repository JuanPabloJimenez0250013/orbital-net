# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first
COPY sim-service/go.mod sim-service/go.sum ./
COPY pkg ../pkg

# Download dependencies
RUN go mod tidy

# Copy the rest of the microservice
COPY sim-service/. .

# Build the binary
RUN go build -o sim-service main.go

# Final stage
FROM alpine:3.18

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/sim-service .

# Expose the port your service listens on
EXPOSE 8081

CMD ["./sim-service"]
