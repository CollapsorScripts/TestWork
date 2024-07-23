FROM golang:latest as builder

SHELL ["/bin/bash", "-c"]

# Устанавливаем значение переменной GOARCH внутри Docker контейнера
ENV GOARCH=arm64

# Обновляем
RUN apt-get update -y && apt-get upgrade -y


# Рабочая папка
WORKDIR /go/api
# Копируем необходимые папки в билдер
COPY api .
# Компилируем
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -o ./build/api ./cmd/entrypoint

# Рабочая папка
WORKDIR /go/checker
# Копируем необходимые папки в билдер
COPY checker .
# Компилируем
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -o ./build/checker ./cmd/entrypoint


# Создаем финальный образ
FROM alpine:latest

# Рабочая директория
WORKDIR /app

#Порт для прослушки
ENV PORT=8070

# Копируем исполняемые файлы из предыдущего образа
COPY --from=builder /go/api/build/api .
COPY --from=builder /go/checker/build/checker .

# Устанавливаем права на выполнение (если необходимо)
RUN chmod +x ./api
RUN chmod +x ./checker

# Копируем файл конфигурации
COPY .env .

# Открываем порты
EXPOSE ${PORT}