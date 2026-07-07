# ---- Стадия сборки ----
FROM golang:1.26.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
ENV GOPROXY=https://proxy.golang.org,direct
COPY ./ /app
RUN CGO_ENABLED=0 GOOS=linux go build -o backend ./cmd/main.go

# Запуск
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/backend .
COPY .env .env
ENTRYPOINT [ "./backend" ]