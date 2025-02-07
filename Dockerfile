# Wybierz obraz Go jako bazę
FROM golang:1.23-alpine

# Ustaw katalog roboczy
WORKDIR /app

# Skopiuj pliki zależności
COPY go.mod go.sum ./
RUN go mod download

# Skopiuj całą aplikację
COPY . .

# Upewnij się, że pliki CSV są dostępne we właściwej lokalizacji
# Kopiujemy katalog data do katalogu /app/cmd/data
RUN mkdir -p /app/cmd/data && cp -r /app/data/* /app/cmd/data/

# Zmień katalog na ten zawierający `main.go`
WORKDIR /app/cmd

# Buduj aplikację
RUN go build -o main .

# Ustaw domyślny punkt wejścia
EXPOSE 8080
CMD ["./main"]
