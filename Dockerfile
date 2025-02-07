# Wybierz obraz Go jako bazę
FROM golang:1.23-alpine

# Ustaw katalog roboczy
WORKDIR /app

# Skopiuj pliki zależności
COPY go.mod go.sum ./
RUN go mod download

# Skopiuj całą aplikację
COPY . .

# Zmień katalog na ten zawierający `main.go`
WORKDIR /app/cmd

# Buduj aplikację
RUN go build -o main .

# Ustaw domyślny punkt wejścia
EXPOSE 8080
CMD ["./main"]
