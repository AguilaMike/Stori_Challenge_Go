# Stori Challenge

This project implements the Stori Software Engineer Technical Challenge.

## Prerequisites

- Docker
- Docker Compose

## Running the Project

1. Clone the repository:
   ```
   git clone https://github.com/AguilaMike/Stori_Challenge_Go
   cd Stori_Challenge_Go
   ```
2. Create a `.env` file in the project root and configure your environment variables.
    ```
    DB_DRIVER=postgres
    DB_URL=postgres://user:password@db:5432/stori?sslmode=disable
    ELASTICSEARCH_URL=http://elasticsearch:9200
    NATS_URL=nats://nats:4222
    API_PORT=8080
    GRPC_PORT=50051
    SMTP_HOST=smtp.example.com
    SMTP_PORT=587
    SMTP_USER=your_email@example.com
    SMTP_PASSWORD=your_email_password
    ```
3. Build and run the project using Docker Compose:
    ```
    docker-compose up -d --build
    ```
4. The API will be available at `http://localhost:8080` and the gRPC server at `localhost:50051`.

## Processing CSV Files

To process a transaction CSV file:

1. Place your CSV file in the directory `data/input`.
2. Run the following command:
   ```
   curl -X POST http://localhost:8080/process-csv -d '{"filename": "tu_archivo.csv", "account_id": "uuid_de_la_cuenta"}'
   ```

## Running Migrations

Migrations are automatically run when the application starts. To run them manually:
    ```
    ./scripts/run_migrations.sh
    ```

## Project Structure

- `cmd/`: Contains the main application entry points.
- `internal/`: Contains the core application logic.
  - `account/`: Account domain.
  - `transaction/`: Transactions domain.
  - `common/`: Shared packages (Database, configuration, etc.).
- `pkg/`: Contains shared packages and third-party integrations.
- `api/`: Contains API handlers and gRPC service implementations.
- `scripts/`: Contains database migration files.
- `sqlc`: Contains configuration files by sqlc framwework database
- `web/`: Contains web-related files (HTML, JS, etc.).

## License

This project is licensed under the MIT License.

## Technologies Used

- Go 1.22.2
- PostgreSQL
- Elasticsearch
- NATS
- gRPC
- Docker


### Comands to update code
sqlc generate
migrate -source file://scripts/migrations -database postgres://sa:@dmin1234@localhost:5432/stori?sslmode=disable up
protoc --go_out=. --go-grpc_out=. pkg/proto/*.proto
