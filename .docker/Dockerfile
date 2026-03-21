# Builder
FROM golang:1.23.5 AS builder

RUN apt-get update && apt-get install -y ca-certificates git-core ssh

WORKDIR /opt/driver-server

# Конфиг
COPY . .
COPY .docker/config/config.yaml .docker/config/config.yaml ./

# Установка зависимостей
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# swagger
RUN --mount=type=cache,target=/go/pkg/mod \
    go install github.com/swaggo/swag/cmd/swag@latest

# создание swagger-документации
RUN swag init -g cmd/app/main.go    

# Сборка
RUN go build -o application ./cmd/app

# Сделать файл исполняемым
RUN chmod +x application

# Старт приложения
ENTRYPOINT ["./application"]

