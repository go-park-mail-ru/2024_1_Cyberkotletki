# docker build -t core_service -f docker/core.dockerfile .
# docker run -d -p 8080:8080 --name core core_service

# Этап сборки
FROM golang:1.22-alpine AS build

WORKDIR /src
COPY cmd cmd
COPY internal internal
COPY docs docs
COPY go.mod go.mod
COPY config config
COPY pkg pkg
RUN apk add --no-cache gcc libc-dev
RUN go mod tidy
RUN go build -o core cmd/app/main.go


# --------------------------------------------

# Этап запуска
FROM alpine:latest

WORKDIR /app
COPY --from=build /src/core /app
COPY config.yaml /app
CMD ["./core"]
