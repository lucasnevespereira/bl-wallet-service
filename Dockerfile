FROM golang:1.21 AS builder

WORKDIR /app
ADD . /app

RUN go build -o bl-wallet-service cmd/server/main.go

FROM ubuntu:latest AS launcher
COPY --from=builder /app .
CMD ["./bl-wallet-service"]