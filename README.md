# Simple Bank Backend

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-00ADD8?style=for-the-badge&logo=go&logoColor=white)

A simple bank backend service built with Go, PostgreSQL, and Gin framework.

## Features

- User registration and authentication (PASETO tokens)
- Account management (create, read, update, delete)
- Money transfers between accounts
- Transaction history with entries
- Database transactions with deadlock prevention
- Auth middleware for protected routes

## Tech Stack

- **Language**: Go
- **Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: SQLC (type-safe SQL)
- **Auth**: PASETO tokens
- **Testing**: Testify + GoMock
- **Container**: Docker

## Getting Started

### Prerequisites

- Go 1.21+
- Docker
- golang-migrate
- sqlc
- mockgen

### Setup

```bash
# Start PostgreSQL container
make startdb

# Create database (first time only)
make createdb

# Run migrations
make migrateup

# Generate SQLC code
make sqlc

# Generate mocks for testing
make mock
```

### Run Server

```bash
make server
```

Server runs at `http://localhost:8080`

### Run Tests

```bash
make test
```

## API Endpoints

### Public Routes

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/users` | Create new user |
| POST | `/users/login` | Login user |

### Protected Routes (requires Bearer token)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/accounts` | Create account |
| GET | `/accounts/:id` | Get account by ID |
| GET | `/accounts` | List accounts |
| PUT | `/accounts/:id` | Update account |
| DELETE | `/accounts/:id` | Delete account |
| POST | `/transfers` | Create transfer |

## Example Usage

```bash
# Create user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"secret123","full_name":"John Doe","email":"john@example.com"}'

# Login
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"secret123"}'

# Create account (with token)
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{"owner":"john","currency":"USD"}'
```

## Project Structure

```
├── api/              # HTTP handlers and middleware
├── db/
│   ├── migration/    # SQL migration files
│   ├── mock/         # Mock store for testing
│   ├── query/        # SQL queries for SQLC
│   └── sqlc/         # Generated Go code
├── token/            # PASETO token maker
├── util/             # Utility functions
├── app.env           # Environment config
├── main.go           # Entry point
└── Makefile          # Build commands
```

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make postgres` | Create new PostgreSQL container |
| `make startdb` | Start existing container |
| `make stopdb` | Stop container |
| `make createdb` | Create database |
| `make dropdb` | Drop database |
| `make migrateup` | Run all migrations |
| `make migratedown` | Rollback all migrations |
| `make sqlc` | Generate SQLC code |
| `make mock` | Generate mock store |
| `make test` | Run tests |
| `make server` | Start server |
