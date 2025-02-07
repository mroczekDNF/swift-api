# Swift-API

Swift-API is a **RESTful API** designed for handling SWIFT code data. It allows users to fetch, add, update, and delete SWIFT codes while maintaining relationships between headquarters and branches. This project uses PostgreSQL as the database and is fully containerized with Docker for seamless deployment and testing.

---

## Features

- **CRUD operations**: Perform Create, Read, Update, and Delete operations on SWIFT code records.
- **Headquarter/Branch relationships**: Maintain hierarchical relationships between headquarters and their branches.
- **RESTful API design**: All endpoints follow standard RESTful conventions.
- **Dockerized environment**: Simplified deployment with Docker Compose.

---

## Prerequisites

Before running the application, ensure the following tools are installed:

- **Git**: [Install Git](https://git-scm.com/)
- **Docker**: [Install Docker](https://www.docker.com/get-started)

---

## Getting Started

Follow the steps below to set up and run the Swift-API project:

### 1. Clone the repository
Navigate to a directory of your choice and run:

```bash
git clone https://github.com/mroczekDNF/swift-api.git
cd swift-api
```

### 2. Build and run the application
Run the following command to build and start the application:

```bash
docker-compose up --build
```

This command:
- Builds the necessary Docker images.
- Spins up the PostgreSQL database container.
- Starts the Swift-API application on **http://localhost:8080**.

### 3. Verify the API
Visit the following endpoint in your browser or with a tool like Postman:

```bash
http://localhost:8080/v1/swift-codes/BSCHCLR10R6
```

You should receive a response with SWIFT code data, verifying that the API is running correctly.

---

## API Endpoints

### GET: `/v1/swift-codes/{swiftCode}`
- **Description**: Fetches details of a specific SWIFT code.
- **Response**: Includes information about the SWIFT code, associated bank, country, and whether it is a headquarter.

### POST: `/v1/swift-codes`
- **Description**: Adds a new SWIFT code to the database.
- **Request Body**:
  ```json
  {
    "swiftCode": "NEWCODE01XXX",
    "bankName": "New Bank",
    "address": "New Street",
    "countryISO2": "US",
    "countryName": "United States",
    "isHeadquarter": true
  }
  ```

### DELETE: `/v1/swift-codes/{swiftCode}`
- **Description**: Deletes a SWIFT code from the database.

---

## Running Tests

### 1. Set up the test environment
Run the following command to build and start the test environment:

```bash
docker-compose -f docker-compose_test.yml up --build -d
```

This starts a test-specific PostgreSQL database on **port 5433**.

### 2. Run tests
Execute the following command to run integration tests:

```bash
go test ./tests/integration/... -v
```

### 3. Clean up test environment
After running tests, shut down the test environment and remove volumes:

```bash
docker-compose -f docker-compose_test.yml down --volumes
```

---

## Notes

- If you encounter database errors, ensure that the volumes for both the main and test databases are cleaned up before restarting Docker:
  ```bash
  docker-compose down --volumes
  docker-compose -f docker-compose_test.yml down --volumes
  ```

- The API is set to use `localhost:8080` for the main application and `localhost:5433` for the test database.

---

Let me know if you'd like to make further changes or add more details!
