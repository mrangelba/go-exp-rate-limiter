FROM golang:1.22-alpine AS build 

RUN go install github.com/cosmtrek/air@latest

WORKDIR /app

ENTRYPOINT ["/go/bin/air"]