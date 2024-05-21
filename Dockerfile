FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o caravan

EXPOSE 8080

CMD ["./caravan"]