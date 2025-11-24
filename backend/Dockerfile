# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN apk add --no-cache build-base
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o safenote-backend ./cmd/server/main.go

# Run Stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/safenote-backend .

# Install sqlite runtime dependencies if needed (alpine usually needs none for static go, but CGO might need libc)
# Using CGO_ENABLED=1 for go-sqlite3 usually requires gcc in builder, and libc in runner.
# Let's ensure we have ca-certificates.
RUN apk --no-cache add ca-certificates

EXPOSE 8080

CMD ["./safenote-backend"]
