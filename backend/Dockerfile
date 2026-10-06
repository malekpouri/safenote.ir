# Build stage (CGO is required by mattn/go-sqlite3)
FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk add --no-cache build-base
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -ldflags="-s -w" -o safenote-backend ./cmd/server

# Run stage
FROM alpine:3.20
WORKDIR /app
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/safenote-backend .

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD wget -q -O /dev/null http://127.0.0.1:8080/health || exit 1

CMD ["./safenote-backend"]
