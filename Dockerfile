FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o caravan

FROM alpine:latest

COPY --from=builder caravan /app

WORKDIR /app

EXPOSE 8080

CMD ["./caravan"]