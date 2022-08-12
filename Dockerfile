# Start from golang:1.12-alpine base image
FROM golang:1.17-alpine as builder

WORKDIR /

COPY ./ ./ 

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

EXPOSE 8080

FROM alpine:latest

COPY --from=builder /app /app

CMD [ "./app" ]