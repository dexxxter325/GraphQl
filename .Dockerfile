FROM golang:1.22-alpine3.19 AS builder

RUN go version

COPY ./ /GRAPHQL

WORKDIR GRAPHQL

RUN go mod download && go get -u ./...
RUN go build -o ./bin/api ./server.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /CRUD_API/bin/api .

EXPOSE 82

CMD ["./api"]