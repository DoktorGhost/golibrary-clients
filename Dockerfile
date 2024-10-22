# Стадия сборки приложения
FROM golang:1.23-alpine as app-builder
RUN apk update && apk add --no-cache curl make git

# Рабочая директория для сборки
WORKDIR /src

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем оставшиеся файлы исходного кода
COPY . .

# Собираем приложение в бинарный файл с именем "app"
RUN go build -o /app/app ./cmd/app

# Финальная стадия для запуска приложения
FROM alpine:latest
RUN apk update && apk add --no-cache curl

# Рабочая директория для финального контейнера
WORKDIR /src

# Копируем скомпилированный бинарник из стадии сборки
COPY --from=app-builder /app/app .
COPY ./migrations /src/migrations

# Указываем, что запускается скомпилированное приложение
CMD ["./app"]