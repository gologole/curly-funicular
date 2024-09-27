FROM golang:latest

# Копируем файл go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код приложения
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/main.go

# Указываем команду для запуска приложения
CMD ["./main"]
