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

### GET: `/v1/swift-codes/{swiftCode}`
- **Description**: Fetches details of SWIFT code data for a specific country.

### POST: `/v1/swift-codes`
- **Description**: Adds a new SWIFT code to the database.

### DELETE: `/v1/swift-codes/{swiftCode}`
- **Description**: Deletes a SWIFT code from the database.

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

## Architecture & Implementation Details

### Data Validation and Error Handling
- The API includes robust validation mechanisms to ensure the correctness of input data:
  - **SWIFT codes** must be between 8 and 11 characters.
  - **Country ISO2 codes** must be exactly 2 uppercase letters.
  - Essential fields such as `bank_name`, `country_name`, and `is_headquarter` are mandatory.
- Addresses, while not critical, are gracefully handled. If the address is missing or empty, it defaults to `"UNKNOWN"` to align with business logic, where the address is less significant than other fields.

### Headquarter and Branch Relationships
- The database tracks hierarchical relationships between headquarters and branches using the `headquarter_id` column.
- Each branch is linked to its headquarters through this column, which is indexed for efficient querying.
- Headquarters are uniquely identified by their SWIFT codes, which typically end with "XXX".

### Performance Optimizations
- The addition of the `headquarter_id` column significantly improves query performance, especially when retrieving branches associated with a specific headquarters.
- Indexed relationships ensure that queries for branches are fast and scalable, even as the dataset grows.

### Robustness Against Errors
- The new functionality for managing headquarter-branch relationships has been implemented with safeguards to maintain data integrity:
  - When a headquarter is deleted, the `headquarter_id` field of its associated branches is automatically set to `NULL`.
  - When new headquarters or branches are added, the system ensures that they are appropriately linked or assigned to maintain consistent relationships.
  - Comprehensive testing ensures the functionality is error-free and resilient.

### Database Schema
- The `swift_codes` table includes the following columns:
  - `id`: Primary key.
  - `swift_code`: Unique identifier for each record.
  - `bank_name`, `address`, `country_iso2`, `country_name`: Key fields containing bank information.
  - `is_headquarter`: Boolean field indicating whether the SWIFT code belongs to a headquarters.
  - `headquarter_id`: Nullable foreign key linking a branch to its headquarters.

### Scalability and Extendability
- The architecture is designed to accommodate future features, such as:
  - Additional metadata for SWIFT codes.
  - Enhanced reporting or more complex relationships between records.

### Business Logic and Practical Decisions
- The decision to handle unknown addresses as `"UNKNOWN"` reflects a practical understanding of business needs, emphasizing the importance of fields like SWIFT code and bank name over address details.
- By focusing on critical data, the system ensures reliability and usability even in scenarios with incomplete information.

