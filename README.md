# Backend

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)

A simple bank backend service built with Go and PostgreSQL.

## Features

- Account management (create, read, update, delete)
- Money transfers between accounts
- Transaction history with entries
- Database transactions with deadlock prevention

## Getting Started

### Prerequisites

- Go 1.21+
- Docker
- golang-migrate

### Setup

```bash
# Start PostgreSQL container
make postgres

# Create database
make createdb

# Run migrations
make migrateup

# Generate SQLC code
make sqlc
```

### Run Tests

```bash
make test
```

## Project Structure

```
├── db/
│   ├── migration/    # SQL migration files
│   ├── query/        # SQL queries for SQLC
│   └── sqlc/         # Generated Go code
├── util/             # Utility functions
└── Makefile
```
