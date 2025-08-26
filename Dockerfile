FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mailing ./cmd/mailing

FROM debian:bookworm-slim

WORKDIR /app

# для работы с TLS (SMTP/HTTPS)
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/mailing .
COPY --from=builder /app/db ./db
COPY --from=builder /app/.env.example .env

CMD ["./mailing"]
