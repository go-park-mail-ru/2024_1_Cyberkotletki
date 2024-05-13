# Этап сборки
FROM golang:1.22-alpine AS build

# Установка Goose
RUN apk add --no-cache git && \
    go install github.com/pressly/goose/cmd/goose@latest

# Этап запуска
FROM alpine:latest

# Установка bash и postgresql-client для проверки соединения
RUN apk add --no-cache bash postgresql-client

# Копирование Goose и файлов миграции из этапа сборки
COPY --from=build /go/bin/goose /usr/local/bin/goose
# Копирование файлов миграции
COPY db/migrations /migrations

WORKDIR /migrations

# Запуск Goose
CMD until psql "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@postgres:5432/kinoskop?sslmode=disable" -c '\q'; do sleep 1; done && \
    goose -dir=/migrations postgres "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@postgres:5432/kinoskop?sslmode=disable" up || exit 1