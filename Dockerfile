# Этап 1: сборка бинарного файла
FROM golang:1.24 AS builder

WORKDIR /app

# Копируем файлы зависимостей и загружаем модули
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код и собираем бинарный файл
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Этап 2: финальный минималистичный образ
FROM alpine:3.21

# Устанавливаем CA‑сертификаты (нужно для HTTPS‑запросов)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/main .
COPY --from=builder /app/views ./views
COPY --from=builder /app/static ./static

# Открываем порт, если твоё приложение — веб‑сервер (например, слушает 8080)
EXPOSE 8080

# Команда запуска
CMD ["./main"]