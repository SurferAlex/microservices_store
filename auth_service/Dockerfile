# Этап 1: Сборка приложения
FROM golang:1.25.0-alpine AS builder

# Метаданные
LABEL maintainer="artielucifer@gmail.com"
LABEL description="Auth Service"

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей (для кэширования)
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем приложение с оптимизацией
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w" \
    -o /app/auth_service \
    ./cmd/main.go

# Этап 2: Финальный образ (минимальный)
FROM alpine:3.21.0

# Метаданные
LABEL maintrainer="artielucifer@gmail.com"
LABEL description="Auth Service"

# Создаем непривилерированного пользователя
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Устанавливаем CA сертификаты (для HTTPS запросов)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем собранный бинарник из builder
COPY --from=builder /app/auth_service .

# Меняем владлельца файлов
RUN chown -R appuser:appuser /app

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./auth_service"]