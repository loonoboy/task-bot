# Этап 1: Сборка приложения
FROM golang:1.23.3-alpine as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файл go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости Go
RUN go mod tidy

# Копируем весь код приложения в контейнер
COPY . .

# Компилируем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o telegram-bot .

# Этап 2: Финальный образ
FROM alpine:latest

# Устанавливаем необходимые для работы сертификаты
RUN apk --no-cache add ca-certificates

# Копируем скомпилированное приложение из первого этапа
COPY --from=builder /app/telegram-bot /usr/local/bin/telegram-bot

# Создаём директорию для логов
RUN mkdir -p /var/log/telegram-bot

# Устанавливаем переменные окружения
ENV PATH="/usr/local/bin:${PATH}"

# Команда для запуска приложения
CMD ["telegram-bot"]
