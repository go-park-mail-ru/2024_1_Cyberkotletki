# docker build -t auth_service -f docker/auth.dockerfile .
# docker run -d -p 8081:8081 --name auth auth_service

# Этап сборки
FROM golang:1.22-alpine AS build

WORKDIR /src
RUN apk add --no-cache gcc libc-dev
COPY go.mod go.mod
COPY cmd cmd
COPY internal internal
COPY docs docs
COPY config config
COPY pkg pkg
RUN go mod tidy
RUN go build -o auth cmd/auth/main.go

# --------------------------------------------

# Этап запуска
FROM alpine:latest

WORKDIR /app
COPY --from=build /src/auth /app
COPY config_auth.yaml /app

CMD ["./auth"]