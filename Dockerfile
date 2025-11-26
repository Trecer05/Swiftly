# ------------ BUILD STAGE -------------
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/service ./cmd/servers/${SERVICE}

# ------------ RUNTIME STAGE -------------
FROM alpine

# Устанавливаем wget для healthcheck
RUN apk add --no-cache wget

WORKDIR /root

COPY --from=builder /bin/service .
# Важно: копируем миграции
COPY --from=builder /app/internal/repository/migrations ./internal/repository/migrations/

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --retries=5 \
  CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./service"]
