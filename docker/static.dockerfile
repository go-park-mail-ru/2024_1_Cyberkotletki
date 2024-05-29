# docker build -t static_service -f docker/static.dockerfile .
# docker run -d -p 8082:8082 --name static static_service

# Этап сборки
FROM golang:1.22-alpine AS build

WORKDIR /src
RUN apk add --no-cache gcc libc-dev libwebp-dev
COPY go.mod go.mod
COPY cmd cmd
COPY internal internal
COPY docs docs
COPY config config
COPY pkg pkg
RUN go mod tidy
RUN go build -o static cmd/static/main.go

# --------------------------------------------

# Этап запуска
FROM alpine:latest

WORKDIR /app
COPY --from=build /src/static /app
COPY config_static.yaml /app

CMD ["./static"]
