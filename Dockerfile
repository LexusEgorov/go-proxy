###########################
#1 Сборочный этап (builder)
###########################
FROM golang:1.24-alpine AS builder
LABEL stage=builder

#1.1 Рабочая директория внутри контейнера
WORKDIR /app

#1.2 Копируем файлы с зависимостями и скачиваем модули
COPY go.mod go.sum ./
RUN go mod download

#1.3 Копируем исходники приложения
COPY . .

#1.4 Собираем приложение
RUN go build -o proxy ./cmd/proxy

#########################
#2 Финальный легкий образ
#########################

#2.1
FROM alpine:latest

#2.2 Рабочая директория внутри контейнера
WORKDIR /app

#2.3 Копируем собранное приложение
COPY --from=builder /app/proxy ./
#2.4 Копируем файлы конфигурации
COPY --from=builder /app/configs ./configs

#2.5 Экспонируем порт, который слушает приложение
EXPOSE 8081

#2.6 Запуск
ENTRYPOINT [ "./proxy" ]