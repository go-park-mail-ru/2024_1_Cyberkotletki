# docker build -t core_service -f docker/core.dockerfile .
# docker run -d -p 8080:8080 --name core core_service

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
RUN go build -o core cmd/app/main.go


# --------------------------------------------

# Этап запуска
FROM alpine:latest

WORKDIR /app
COPY --from=build /src/core /app
COPY config_static.yaml /app
COPY config_auth.yaml /app
COPY config.yaml /app
CMD ["./core"]
