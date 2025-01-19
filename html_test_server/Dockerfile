FROM golang:1.23.4 AS builder

WORKDIR /app

COPY html_test_server/go.mod html_test_server/go.sum ./

RUN go mod download

COPY html_test_server/. .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

ENTRYPOINT ["./main"]