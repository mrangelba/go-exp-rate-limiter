FROM golang:1.21-alpine AS build 

RUN go install github.com/cosmtrek/air@latest
RUN go install go.uber.org/mock/mockgen@latest

WORKDIR /app

ENTRYPOINT ["/go/bin/air"]