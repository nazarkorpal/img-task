# Start from golang:1.12-alpine base image
FROM golang:1.19-alpine as builder

WORKDIR /app

COPY . .

RUN apk add gcc g++ vips-dev vips-poppler pkgconf

RUN go mod download

RUN go build -o server ./cmd/main.go

EXPOSE 8080

FROM alpine:latest

COPY --from=builder /app/server /app/server

RUN apk add --update --no-cache --repository http://dl-3.alpinelinux.org/alpine/edge/testing --repository http://dl-3.alpinelinux.org/alpine/edge/main vips-dev vips-poppler

CMD [ "./app/server" ]