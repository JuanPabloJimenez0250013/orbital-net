# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

ARG SIM_SERVICE_PORT=8081

COPY sim-service/go.mod sim-service/go.sum ./
COPY pkg ../pkg

RUN go mod download

COPY sim-service/. .
RUN go build -o sim-service main.go

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /app/sim-service .

ARG SIM_SERVICE_PORT=8081
ENV SIM_SERVICE_PORT=$SIM_SERVICE_PORT

EXPOSE $SIM_SERVICE_PORT

CMD ["./sim-service"]
