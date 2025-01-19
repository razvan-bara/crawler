FROM golang:1.23.4 AS builder

WORKDIR /app

COPY crawler/go.mod crawler/go.sum ./

RUN go mod download

COPY crawler/. .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8081

ENTRYPOINT ["./main"]