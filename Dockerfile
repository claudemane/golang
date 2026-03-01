# ── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && \
    go build -o app ./cmd/api/main.go

# ── Stage 2: Run (tiny image) ─────────────────────────────────────────────────
FROM alpine

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

ENTRYPOINT ["./app"]
