FROM golang:1.24-alpine as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o order-service ./cmd/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/order-service .
COPY .env .env

EXPOSE 8081
ENTRYPOINT ["./order-service"]