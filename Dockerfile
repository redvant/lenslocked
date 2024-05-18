FROM node:latest AS tailwind-builder
WORKDIR /tailwind
RUN npm init -y && \
    npm install tailwindcss && \
    npx tailwindcss init
COPY ./templates /templates
COPY ./tailwind/tailwind.config.js /src/tailwind.config.js
COPY ./tailwind/styles.css /src/styles.css
RUN npx tailwindcss -c /src/tailwind.config.js -i /src/styles.css -o /styles.css --minify

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
COPY --from=tailwind-builder /styles.css /app/assets/styles.css
CMD ./server
