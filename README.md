# Payment Gateway Service

This project is a payment gateway service that integrates with multiple payment providers such as ADCB and HSBC. It is built with Go and uses PostgreSQL as its database.

## Table of Contents

- [Application Structure](#application-structure)
- [Prerequisites](#prerequisites)
- [Running the Application Manually](#running-the-application-manually)
- [Running the Application Using Docker](#running-the-application-using-docker)
- [Troubleshooting](#troubleshooting)
- [License](#license)

## Application Structure

The application follows a modular structure to maintain separation of concerns and scalability. Below is an overview of the main directories and files:

\`\`\`
payment-gateway-service/
│
├── cmd/
│ └── main.go # Entry point for the application
│
├── internal/
│ ├── adapters/ # Payment provider adapters
│ ├── database/
│ │ ├── db.go # Database connection and setup
│ │ └── migrations/ # Database migrations
│ ├── services/ # Core business logic
│ └── routes/ # API route handlers
│
├── mock_services/ # Mock services for ADCB and HSBC
│ ├── adcb/
│ │ └── main.go
│ └── hsbc/
│ └── main.go
│
├── configs/ # Configuration files and settings
│ └── config.go
│
├── docs/ # Swagger documentation files
│ ├── doc.go
│ ├── swagger.json
│ └── swagger.yaml
│
├── Dockerfile # Dockerfile for the Go application
├── Dockerfile.adcb # Dockerfile for the ADCB mock service
├── Dockerfile.hsbc # Dockerfile for the HSBC mock service
├── docker-compose.yml # Docker Compose file to orchestrate services
├── .env # Environment variables file
├── go.mod # Go module dependencies
└── README.md # Project documentation
\`\`\`

![Payment Gateway Diagram](./payment_gateway.diagram.png)

## Prerequisites

Before running the application, ensure you have the following installed on your machine:

- **Go** (version 1.16 or later)
- **PostgreSQL** (version 13 or later)
- **Docker** (for running the application using Docker)
- **Docker Compose** (for orchestrating multi-container Docker applications)

## Running the Application Manually

### Step 1: Clone the Repository

\`\`\`bash
git clone https://github.com/erik1996/payment-gateway-service
cd payment-gateway-service
\`\`\`

### Step 2: Set Up the Environment Variables

Create a `.env` file in the root of the project with the following content:

\`\`\`env
DB_HOST=localhost
DB_PORT=5432
DB_USER=username
DB_PASSWORD=yourpassword
DB_NAME=payment_db
PORT=8080

HSBC_USER_ID=1
HSBC_USER_SECRET=1111

ADCB_USER_ID=1
ADCB_USER_SECRET=1111

APP_HOST=http://localhost:8080
\`\`\`

### Step 3: Install Dependencies

Ensure you have Go modules enabled and install the necessary dependencies:

\`\`\`bash
go mod tidy
\`\`\`

### Step 4: Set Up the Database

Start PostgreSQL and create the database:

\`\`\`bash
psql -U erik -c "CREATE DATABASE payment_db;"
\`\`\`

### Step 5: Run Migrations

Run the database migrations to set up the schema:

\`\`\`bash
migrate -path ./internal/database/migrations -database "postgres://username:yourpassword@localhost:5432/payment_db?sslmode=disable" up
\`\`\`

### Step 6: Run the Application

Finally, run the application:

\`\`\`bash
go run cmd/main.go
\`\`\`

The application will be available at `http://localhost:8080`.

Once the application is running, you can find the Swagger documentation at `http://localhost:8080/swagger/index.html`.

## Running the Application Using Docker

### Step 1: Clone the Repository

\`\`\`bash
git clone https://github.com/erik1996/payment-gateway-service
cd payment-gateway-service
\`\`\`

### Step 2: Set Up the Environment Variables

Create a `.env` file in the root of the project with the following content:

\`\`\`env
DB_HOST=db
DB_PORT=5432
DB_USER=username
DB_PASSWORD=yourpassword
DB_NAME=payment_db
PORT=8080

HSBC_USER_ID=1
HSBC_USER_SECRET=1111

ADCB_USER_ID=1
ADCB_USER_SECRET=1111

APP_HOST=http://localhost:8080
\`\`\`

### Step 3: Build and Run the Docker Containers

Use Docker Compose to build and run the containers:

\`\`\`bash
docker-compose up --build
\`\`\`

This will start the following services:

- **PostgreSQL** on port 5432
- **Go Application** on port 8080
- **ADCB Mock Service** on port 8082
- **HSBC Mock Service** on port 8081
- **Migrate** service to run the migrations automatically

### Step 4: Access the Application

Once all services are up and running, you can access the application at `http://localhost:8080`.

You can find the Swagger documentation at `http://localhost:8080/swagger/index.html`.

### Step 5: Stop the Application

To stop the application, press `Ctrl+C` and then run:

\`\`\`bash
docker-compose down
\`\`\`

This will stop and remove all the containers.

## Troubleshooting

- **Migrations Fail on First Run:** If the migrations fail during the first run, ensure that the database service is fully up and healthy. You can run the migration command manually using Docker:

  \`\`\`bash
  docker-compose run migrate
  \`\`\`

- **Database Connection Issues:** Ensure that the PostgreSQL service is running and accessible on the specified port. If you encounter issues, check the logs using:

  \`\`\`bash
  docker logs postgres
  \`\`\`

## License

This project is licensed under the MIT License.
