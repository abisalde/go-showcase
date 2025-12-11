
FROM golang:1.25-alpine AS build


WORKDIR /app

RUN apk add --no-cache git gcc musl-dev bash

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY /configs/air.toml .air.toml

COPY . .


ENV AIR_WORKSPACE=/app \
    ENVIRONMENT=development




RUN mkdir -p /app/tmp



CMD ["air", "-c", "/app/.air.toml"]