FROM golang:1.22.3-alpine3.19 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o ./server ./cmd/server/

FROM alpine:3.19
WORKDIR /app
COPY ./assets ./assets
COPY --from=builder /app/server ./server
CMD ./server
