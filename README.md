# Swift API

## **Overview**
Swift API is a RESTful service for handling SWIFT codes, providing functionalities to retrieve, store, and manage SWIFT banking data.

---

## **Prerequisites**
Ensure you have the following dependencies installed:

- **Go 1.19+** → [Download & Install](https://go.dev/dl/)
- **Docker & Docker Compose** (Optional for database) → [Install Docker](https://docs.docker.com/get-docker/)
- **PostgreSQL 13+** (If running manually)
- **Git** (For cloning the repository)

---

## **Installation**
### **Clone the Repository**
```sh
git clone https://github.com/mroczekDNF/swift-api.git
cd swift-api
```

### **Set Up Environment Variables**
Create a `.env` file in the root directory:

```sh
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=swift_db
```

Alternatively, export them manually:

```sh
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=your_db_user
export DB_PASSWORD=your_db_password
export DB_NAME=swift_db
```

---

## **Running the Application**

### **1. Run PostgreSQL Database**
#### **Option 1: Using Docker (Recommended)**
```sh
docker-compose up -d
```
Or manually start a PostgreSQL container:
```sh
docker run --name swift-postgres -e POSTGRES_USER=your_db_user -e POSTGRES_PASSWORD=your_db_password -e POSTGRES_DB=swift_db -p 5432:5432 -d postgres:13
```

#### **Option 2: Using a Local PostgreSQL Installation**
```sh
psql -U your_db_user -c "CREATE DATABASE swift_db;"
```

### **2. Initialize the Database**
```sh
go run cmd/migrate.go
```

### **3. Build and Run the Application**
```sh
go run cmd/main.go
```
Or build the binary:
```sh
go build -o swift-api cmd/main.go
./swift-api
```

---

## **Testing**
### **Run Unit Tests**
```sh
go test ./internal/repositories
```

### **Run Integration Tests**
```sh
go test ./internal/integration
```

---

## **Using the API**
### **Check API Health**
```sh
curl http://localhost:8080/health
```

### **Fetch SWIFT Code Details**
```sh
curl http://localhost:8080/v1/swift-codes/BANKUS33XXX
```

### **Fetch SWIFT Codes by Country**
```sh
curl http://localhost:8080/v1/swift-codes/country/US
```

### **Add a New SWIFT Code**
```sh
curl -X POST http://localhost:8080/v1/swift-codes -H "Content-Type: application/json" -d '{
  "swiftCode": "BANKUS44XXX",
  "bankName": "New Bank USA",
  "countryISO2": "US",
  "countryName": "United States",
  "isHeadquarter": true
}'
```

### **Delete a SWIFT Code**
```sh
curl -X DELETE http://localhost:8080/v1/swift-codes/BANKUS44XXX
```

---

## **Stopping the Application**
To stop the application, press `CTRL + C`.

If using Docker:
```sh
docker-compose down
```
Or stop the PostgreSQL container:
```sh
docker stop swift-postgres && docker rm swift-postgres
```

---

Now your **Swift API** application should be running successfully! 🚀
