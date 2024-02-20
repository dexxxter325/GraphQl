FROM golang:1.22-alpine3.19 AS builder

RUN go version

COPY ./ /GRAPHQL

WORKDIR /GRAPHQL

RUN go mod download
RUN go build -o ./.bin/api ./server/server.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /GRAPHQL/.bin/api .

EXPOSE 82

CMD ["./api"]