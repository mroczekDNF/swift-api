# Użyj oficjalnego obrazu Go jako podstawy
FROM golang:1.20

# Ustaw katalog roboczy w kontenerze
WORKDIR /app

# Skopiuj pliki projektu do kontenera
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Przejdź do katalogu, w którym znajduje się główny plik aplikacji
WORKDIR /app/cmd

# Buduj aplikację
RUN go build -o /app/main .

# Ustawienie punktu wejścia
CMD ["/app/main"]

# Otwórz port aplikacji
EXPOSE 8080
