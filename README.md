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
APP_ENV=
PORT=
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
LOG_LEVEL=
ES_URL=
ES_USER=
ES_PASS=
ES_ENABLED=
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
KIBANA_URL_SA=
VA_PREFIX=8808
VA_EXPIRED_HOURS=24
BASIC_AUTH_USERNAME=
BASIC_AUTH_PASSWORD=
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

Using Postman:

```
http://localhost:9090/api/v1/virtual-accounts/create
```

## Sample Request

```http
POST /api/v1/virtual-accounts
Content-Type: application/json
```

```json
{
    "username": "your_username",
    "password": "your_password",
    "channel": "channel",
    "payload": {
        "customer_id": "CUST-001",
        "customer_name": "Mirna",
        "amount": 150000,
        "description": "Pembayaran Invoice #INV-2024-015",
        "reference_id": "INV-2024-015"
    }
}
```

## Sample Response

```json
{
    "status": "VA_SUCCESS",
    "response_code": "VA-0000",
    "message": "Create VA Success",
    "payload": {
        "id": "ba1ac6bb-87a3-4592-8f22-02b128fd1c9c",
        "va_number": "8808202606260848495802",
        "customer_id": "CUST-001",
        "customer_name": "Mirna",
        "amount": 150000,
        "description": "Pembayaran Invoice #INV-2024-015",
        "reference_id": "INV-2024-015",
        "expired_at": "2026-06-27T08:48:49.415589+07:00",
        "created_at": "2026-06-26T08:48:49.415589+07:00",
        "status": "PENDING"
    }
}
```

## License

This project is intended for portfolio demonstration.
