# Simple Bank Backend

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![gRPC](https://img.shields.io/badge/gRPC-244c5a?style=for-the-badge&logo=google&logoColor=white)
![Kubernetes](https://img.shields.io/badge/Kubernetes-326CE5?style=for-the-badge&logo=kubernetes&logoColor=white)
![AWS](https://img.shields.io/badge/AWS-232F3E?style=for-the-badge&logo=amazon-aws&logoColor=white)

A production-ready banking backend service built with Go, featuring both REST and gRPC APIs, deployed on AWS EKS.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              CLIENTS                                         │
│                                                                              │
│    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐                 │
│    │   Web App    │    │  Mobile App  │    │  gRPC Client │                 │
│    │   (REST)     │    │   (REST)     │    │   (Evans)    │                 │
│    └──────┬───────┘    └──────┬───────┘    └──────┬───────┘                 │
│           │                   │                   │                          │
└───────────┼───────────────────┼───────────────────┼──────────────────────────┘
            │                   │                   │
            ▼                   ▼                   ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         AWS CLOUD (EKS)                                      │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                    Kubernetes Cluster                                   │ │
│  │                                                                         │ │
│  │   ┌─────────────────┐         ┌─────────────────┐                      │ │
│  │   │  Load Balancer  │         │  Ingress/NLB    │                      │ │
│  │   │    :8080        │         │    :9090        │                      │ │
│  │   └────────┬────────┘         └────────┬────────┘                      │ │
│  │            │                           │                                │ │
│  │            ▼                           ▼                                │ │
│  │   ┌────────────────────────────────────────────────┐                   │ │
│  │   │              SimpleBank Service                │                   │ │
│  │   │  ┌──────────────────┬──────────────────────┐  │                   │ │
│  │   │  │  HTTP Gateway    │    gRPC Server       │  │                   │ │
│  │   │  │  (Port 8080)     │    (Port 9090)       │  │                   │ │
│  │   │  │                  │                      │  │                   │ │
│  │   │  │  - REST API      │  - Native gRPC      │  │                   │ │
│  │   │  │  - Swagger UI    │  - Reflection       │  │                   │ │
│  │   │  │  - JSON          │  - Protobuf         │  │                   │ │
│  │   │  └──────────────────┴──────────────────────┘  │                   │ │
│  │   │                      │                         │                   │ │
│  │   │              ┌───────┴───────┐                │                   │ │
│  │   │              │  gapi Package │                │                   │ │
│  │   │              │  - Handlers   │                │                   │ │
│  │   │              │  - Logger     │                │                   │ │
│  │   │              │  - Auth       │                │                   │ │
│  │   │              └───────────────┘                │                   │ │
│  │   └────────────────────────────────────────────────┘                   │ │
│  │                            │                                            │ │
│  └────────────────────────────┼────────────────────────────────────────────┘ │
│                               │                                              │
│                               ▼                                              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                        AWS RDS (PostgreSQL)                            │ │
│  │   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  │ │
│  │   │   users     │  │  accounts   │  │  transfers  │  │  sessions   │  │ │
│  │   └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘  │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  ┌────────────────────┐  ┌────────────────────┐                             │
│  │  AWS ECR           │  │  AWS Secrets Mgr   │                             │
│  │  (Docker Images)   │  │  (DB Credentials)  │                             │
│  └────────────────────┘  └────────────────────┘                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Request Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          HTTP Request Flow                                   │
│                                                                              │
│  curl POST /v1/create_user                                                  │
│       │                                                                      │
│       ▼                                                                      │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐  │
│  │ HTTP Logger │───▶│ gRPC-Gateway│───▶│  Validator  │───▶│   Handler   │  │
│  │ Middleware  │    │   (Mux)     │    │             │    │             │  │
│  └─────────────┘    └─────────────┘    └─────────────┘    └──────┬──────┘  │
│                                                                   │         │
│                                                                   ▼         │
│                                                            ┌─────────────┐  │
│                                                            │  DB Store   │  │
│                                                            │  (SQLC)     │  │
│                                                            └─────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                          gRPC Request Flow                                   │
│                                                                              │
│  grpcurl pb.Simplebank/CreateUser                                           │
│       │                                                                      │
│       ▼                                                                      │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐  │
│  │ gRPC Logger │───▶│ Interceptor │───▶│  Validator  │───▶│   Handler   │  │
│  │             │    │   Chain     │    │             │    │             │  │
│  └─────────────┘    └─────────────┘    └─────────────┘    └──────┬──────┘  │
│                                                                   │         │
│                                                                   ▼         │
│                                                            ┌─────────────┐  │
│                                                            │  DB Store   │  │
│                                                            │  (SQLC)     │  │
│                                                            └─────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Authentication Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Authentication Flow                                  │
│                                                                              │
│   ┌──────────┐                                           ┌──────────────┐   │
│   │  Client  │                                           │   Database   │   │
│   └────┬─────┘                                           └──────┬───────┘   │
│        │                                                        │           │
│        │  1. POST /v1/login_user                               │           │
│        │     {username, password}                               │           │
│        │ ──────────────────────────────────────▶                │           │
│        │                                                        │           │
│        │                    2. Verify password                  │           │
│        │                    ◀──────────────────────────────────▶│           │
│        │                                                        │           │
│        │  3. Return tokens                                      │           │
│        │     {access_token (15m), refresh_token (24h)}         │           │
│        │ ◀──────────────────────────────────────                │           │
│        │                                                        │           │
│        │  4. POST /v1/update_user                              │           │
│        │     Authorization: Bearer <access_token>               │           │
│        │ ──────────────────────────────────────▶                │           │
│        │                                                        │           │
│        │                    5. Verify token (PASETO)            │           │
│        │                    6. Check user permission            │           │
│        │                    7. Execute request                  │           │
│        │                                                        │           │
│        │  8. Return response                                    │           │
│        │ ◀──────────────────────────────────────                │           │
│        │                                                        │           │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Database Schema

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            Database Schema                                   │
│                                                                              │
│  ┌──────────────────┐         ┌──────────────────┐                          │
│  │      users       │         │     sessions     │                          │
│  ├──────────────────┤         ├──────────────────┤                          │
│  │ username (PK)    │◀────────│ username (FK)    │                          │
│  │ hashed_password  │         │ id (PK/UUID)     │                          │
│  │ full_name        │         │ refresh_token    │                          │
│  │ email (UNIQUE)   │         │ user_agent       │                          │
│  │ password_changed │         │ client_ip        │                          │
│  │ created_at       │         │ is_blocked       │                          │
│  └────────┬─────────┘         │ expires_at       │                          │
│           │                   │ created_at       │                          │
│           │                   └──────────────────┘                          │
│           │                                                                  │
│           ▼                                                                  │
│  ┌──────────────────┐                                                       │
│  │    accounts      │                                                       │
│  ├──────────────────┤         ┌──────────────────┐                          │
│  │ id (PK)          │◀────────│     entries      │                          │
│  │ owner (FK)       │         ├──────────────────┤                          │
│  │ balance          │         │ id (PK)          │                          │
│  │ currency         │         │ account_id (FK)  │                          │
│  │ created_at       │         │ amount (+/-)     │                          │
│  └────────┬─────────┘         │ created_at       │                          │
│           │                   └──────────────────┘                          │
│           │                                                                  │
│           ▼                                                                  │
│  ┌──────────────────┐                                                       │
│  │    transfers     │                                                       │
│  ├──────────────────┤                                                       │
│  │ id (PK)          │                                                       │
│  │ from_account (FK)│                                                       │
│  │ to_account (FK)  │                                                       │
│  │ amount           │                                                       │
│  │ created_at       │                                                       │
│  └──────────────────┘                                                       │
└─────────────────────────────────────────────────────────────────────────────┘
```

## CI/CD Pipeline

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          CI/CD Pipeline (GitHub Actions)                     │
│                                                                              │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐   │
│  │  Push   │───▶│  Test   │───▶│  Build  │───▶│  Push   │───▶│ Deploy  │   │
│  │  Code   │    │  (Go)   │    │ Docker  │    │  ECR    │    │  EKS    │   │
│  └─────────┘    └─────────┘    └─────────┘    └─────────┘    └─────────┘   │
│                                                                              │
│  Trigger:       Run:           Build:         Push to:       Apply:         │
│  - push main    - go test      - Dockerfile   - AWS ECR      - kubectl      │
│  - PR           - lint         - multi-stage  - tag: latest  - deployment   │
│                                                               - service      │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Features

- User registration and authentication (PASETO tokens)
- Account management (create, read, update, delete)
- Money transfers between accounts with deadlock prevention
- Dual API: REST (gRPC-Gateway) + Native gRPC
- Swagger UI documentation
- Structured logging (zerolog)
- Input validation
- Database migrations (auto-run on startup)
- Docker containerization
- Kubernetes deployment
- CI/CD with GitHub Actions

## Tech Stack

| Category | Technology |
|----------|------------|
| Language | Go 1.23 |
| Database | PostgreSQL (AWS RDS) |
| ORM | SQLC (type-safe SQL) |
| Auth | PASETO tokens |
| API | gRPC + gRPC-Gateway |
| Docs | Swagger/OpenAPI |
| Logging | Zerolog |
| Testing | Testify + GoMock |
| Container | Docker |
| Orchestration | Kubernetes (AWS EKS) |
| CI/CD | GitHub Actions |
| Registry | AWS ECR |
| Secrets | AWS Secrets Manager |

## Project Structure

```
backend/
├── main.go                 # Entry point - runs gRPC & HTTP servers
├── api/                    # Gin HTTP API (legacy REST)
│   ├── server.go
│   ├── user.go
│   ├── account.go
│   ├── transfer.go
│   └── middleware.go
├── gapi/                   # gRPC API
│   ├── server.go           # gRPC server setup
│   ├── rpc_create_user.go
│   ├── rpc_login_user.go
│   ├── rpc_update_user.go
│   ├── authorization.go    # Token verification
│   ├── logger.go           # gRPC & HTTP logging
│   ├── metadata.go         # Extract client info
│   └── converter.go        # DB → Proto conversion
├── pb/                     # Generated protobuf code
├── proto/                  # Proto definitions
│   ├── service_simple_bank.proto
│   ├── rpc_create_user.proto
│   ├── rpc_login_user.proto
│   ├── rpc_update_user.proto
│   └── user.proto
├── db/
│   ├── migration/          # SQL migrations
│   ├── query/              # SQLC queries
│   ├── sqlc/               # Generated Go code
│   └── mock/               # Mock store for tests
├── token/                  # JWT & PASETO token makers
├── util/                   # Config, password, random
├── val/                    # Input validation
├── doc/
│   ├── swagger/            # Swagger UI + JSON
│   └── statik/             # Embedded static files
├── eks/                    # Kubernetes manifests
│   ├── deployment.yaml
│   ├── service.yaml
│   └── ingress.yaml
├── .github/workflows/      # CI/CD pipelines
│   ├── test.yml
│   └── deploy.yml
├── Dockerfile
├── docker-compose.yaml
└── Makefile
```

## Getting Started

### Prerequisites

- Go 1.23+
- Docker
- PostgreSQL
- protoc (Protocol Buffer compiler)
- golang-migrate
- sqlc
- mockgen
- grpcurl or Evans (gRPC testing)

### Local Setup

```bash
# Clone the repo
git clone https://github.com/satvikmpatil/simplebank.git
cd simplebank

# Start PostgreSQL
make postgres
make createdb

# Run migrations
make migrateup

# Generate code
make sqlc
make proto
make mock

# Run server
make server
```

### Environment Variables

Create `app.env`:
```env
ENVIRONMENT=dev
DB_DRIVER=postgres
DB_SOURCE=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
MIGRATION_URL=file://db/migration
HTTP_SERVER_ADDRESS=0.0.0.0:8080
GRPC_SERVER_ADDRESS=0.0.0.0:9090
TOKEN_SYMMETRIC_KEY=12345678901234567890123456789012
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=24h
```

## API Endpoints

### gRPC-Gateway (REST) - Port 8080

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/v1/create_user` | Create new user | No |
| POST | `/v1/login_user` | Login, get tokens | No |
| PATCH | `/v1/update_user` | Update user info | Yes |
| GET | `/swagger/` | Swagger UI docs | No |

### gRPC - Port 9090

| Service | Method | Description |
|---------|--------|-------------|
| Simplebank | CreateUser | Create new user |
| Simplebank | LoginUser | Login user |
| Simplebank | UpdateUser | Update user |

## Usage Examples

### REST API (curl)

```bash
# Create user
curl -X POST http://localhost:8080/v1/create_user \
  -H "Content-Type: application/json" \
  -d '{"username":"john","full_name":"John Doe","email":"john@example.com","password":"secret123"}'

# Login
curl -X POST http://localhost:8080/v1/login_user \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"secret123"}'

# Update user (with token)
curl -X PATCH http://localhost:8080/v1/update_user \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{"username":"john","full_name":"John Smith"}'
```

### gRPC (grpcurl)

```bash
# Create user
grpcurl -plaintext -d '{"username":"john","full_name":"John Doe","email":"john@example.com","password":"secret123"}' \
  localhost:9090 pb.Simplebank/CreateUser

# Login
grpcurl -plaintext -d '{"username":"john","password":"secret123"}' \
  localhost:9090 pb.Simplebank/LoginUser
```

### Evans (Interactive gRPC Client)

```bash
make evans

# In Evans REPL:
package pb
service Simplebank
call CreateUser
```

## Kubernetes Deployment

### Deploy to EKS

```bash
# Configure kubectl
aws eks update-kubeconfig --name simple_bank --region ap-south-1

# Apply manifests
kubectl apply -f eks/deployment.yaml
kubectl apply -f eks/service.yaml

# Check status
kubectl get pods
kubectl get services
```

### Kubernetes Architecture

```yaml
# deployment.yaml
replicas: 2
image: <aws-account>.dkr.ecr.<region>.amazonaws.com/simplebank:latest
ports:
  - containerPort: 8080  # HTTP
  - containerPort: 9090  # gRPC
```

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make server` | Run the server |
| `make proto` | Generate protobuf & swagger |
| `make sqlc` | Generate DB code |
| `make mock` | Generate mocks |
| `make test` | Run all tests |
| `make evans` | Start gRPC client |
| `make postgres` | Create PostgreSQL container |
| `make startdb` | Start database |
| `make stopdb` | Stop database |
| `make createdb` | Create database |
| `make dropdb` | Drop database |
| `make migrateup` | Run migrations |
| `make migratedown` | Rollback migrations |
| `make docker-build` | Build Docker image |

## Testing

```bash
# Run all tests
make test

# Run specific test
go test -v ./api -run TestCreateUserAPI
```

## License

MIT License

## Author

**Satvik Patil**
- GitHub: [@satvikmpatil](https://github.com/satvikmpatil)
- Email: 7satvikpatil@gmail.com
