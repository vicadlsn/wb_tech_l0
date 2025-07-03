FROM golang:1.24-alpine

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /migrate

COPY migrations ./migrations

ENTRYPOINT ["goose", "-dir", "./migrations"]