FROM golang:1.22-alpine AS build 

WORKDIR /app
COPY ./ . 

RUN apk --no-cache add ca-certificates git
RUN CGO_ENABLED=0 GOOS=linux go build -C cmd/server/ -o /app_bin

FROM scratch
COPY --from=build /app_bin /app_bin
COPY --from=build /app/configs/.env.example /.env
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080
ENTRYPOINT ["/app_bin"]