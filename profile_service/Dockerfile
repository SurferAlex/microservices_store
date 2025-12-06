# Этап 1: Сборка приложения
FROM golang:1.25.0-alpine AS builder

# Метаданные
LABEL maintainer="artielucifer@gmail.com"
LABEL description="Profile Service"

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей (для кэширования)
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w" \
    -o /app/profile_service \
    ./cmd/main.go

# Этап 2: Финальный образ (минимальный)
FROM alpine:3.21.0

# Метаданные
LABEL maintainer="artielucifer@gmail.com"
LABEL descriptions="Profile Service"

# Создаем непривилигированного пользователя
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Устанавливаем CA сертификаты (для HTTPS запросов)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем собранный бинарник из builder
COPY --from=builder /app/profile_service .

# Меняем владлеьца файлов
RUN chown -R appuser:appuser /app

# Переключаемся на непривелигированного пользователя
USER appuser

# Открываем порт
EXPOSE 8081

# Запускаем приложение
CMD ["./profile_service"]