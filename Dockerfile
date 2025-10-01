# Stage 1: сборка
FROM golang:1.25 AS builder
WORKDIR /build

COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/bot

# Stage 2: минимальный образ
FROM scratch
WORKDIR /app

COPY --from=builder /build/bot .
COPY --from=builder /build/configs ./configs

EXPOSE 8080

# В ENV передаем имя файла конфигурации
# В локальном запуске можно не задавать ENV, будет values_local.yaml
ENTRYPOINT ["./bot"]