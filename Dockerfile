FROM golang:1.22 AS builder

LABEL authors="arian"

WORKDIR /app

COPY . .

RUN cd cmd/ && go mod download
RUN cd cmd/ && go build -o main .

CMD ["./cmd/main"]