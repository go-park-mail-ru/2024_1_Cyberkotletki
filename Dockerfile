FROM golang:1.22.1-alpine3.19 AS build
RUN apk --no-cache add gcc g++ make
RUN apk add git
WORKDIR /app
COPY . .
RUN GOOS=linux go build -o ./.bin ./cmd/app/main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /app/.bin .
COPY --from=build /app/.env .
# RUN chmod +x /app
EXPOSE 8080
ENTRYPOINT ["./.bin"]