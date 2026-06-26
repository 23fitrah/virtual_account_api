# Virtual Account API

RESTful API for managing Virtual Account transactions, built with Golang using a clean architecture approach. This project provides endpoints to create virtual accounts, retrieve account details, process payments, and monitor transaction status.

## Features

* Create Virtual Account
* Get Virtual Account Status
* Get List Virtual Account
* Process Payment
* Transaction Payment History

## Tech Stack

* Go (Golang)
* Gin Framework
* SQL Server
* Elasticsearch
* Redis
* GORM
* Docker Container
* Unit Test

## Project Structure

```
.
├── cmd/
├── config/
├── constants/
├── internal/
│   ├── handler/
|   ├── injector/
|   ├── middleware/
|   ├── providers/
|   ├── repositories/
│   ├── routes/
│   ├── services/
│   ├── validations/
├── logs/
├── models/
├── resources/
├── tests/
├── utils/
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

## API Endpoints

| Method | Endpoint                               | Description                     |
| ------ | --------------------------------------- | -------------------------------- |
| POST   | `/api/v1/virtual-accounts/create`             | Create a new Virtual Account    |
| POST   | `/api/v1/virtual-accounts/:va_number`        | Get Virtual Account Status information |
| POST   | `/api/v1/virtual-accounts?page=1&limit=10&customer_id=&status=`        | Get List Virtual Account |
| POST   | `/api/v1/virtual-accounts/payments`     | Process payment                 |
| POST   | `/api/v1/payments/history?page=1&limit=10&status=`                 | Get transaction history         |

## Getting Started

### Clone Repository

```bash
git clone https://github.com/yourusername/virtual-account-api.git
cd virtual-account-api
```

### Install Dependencies

```bash
go mod tidy
```

### Configure Environment

Create a `.env` file.

```env
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_NAME=virtual_account
DB_USER=postgres
DB_PASSWORD=password
```

### Run Application

```bash
go run cmd/main.go
```

or

```bash
docker-compose up --build
```

## Running Tests

```bash
go test ./...
```

To check test coverage:

```bash
go test ./... -cover
```

## API Documentation

If Swagger is enabled:

```
http://localhost:8080/swagger/index.html
```

## Sample Request

```http
POST /api/v1/virtual-accounts
Content-Type: application/json
```

```json
{
  "customer_name": "John Doe",
  "bank_code": "014",
  "amount": 250000,
  "expired_at": "2026-12-31T23:59:59Z"
}
```

## Sample Response

```json
{
  "virtual_account_number": "0141234567890",
  "customer_name": "John Doe",
  "amount": 250000,
  "status": "ACTIVE",
  "expired_at": "2026-12-31T23:59:59Z"
}
```

## License

This project is intended for learning purposes and portfolio demonstration.
