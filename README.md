# SSO Service

🚀 A modular and secure Single Sign-On (SSO) microservice written in Go, with gRPC support, MongoDB as a data store, and JWT-based authentication.

## Features

- 🔐 Secure user authentication with JWT tokens
- 📡 gRPC API for service-to-service communication
- 💾 MongoDB storage
- 📁 Configuration via YAML (`config/local.yaml`)
- 🪵 Structured logging with `log/slog`
- ♻️ Graceful shutdown with signal handling
- ✅ Modular architecture following Clean Architecture principles

## Project Structure

```
sso/
├── cmd/
│   └── sso/                # Entry point (main.go)
├── config/                 # Configuration files (YAML)
│   └── local.yaml
├── internal/
│   ├── app/                # Composition root / DI container
│   ├── services/           # Business logic (e.g., AuthService)
│   ├── storage/            # MongoDB repository layer
│   ├── domain/             # Domain models (User, App, etc.)
│   └── jwt/                # JWT utility functions
├── go.mod / go.sum         # Dependencies
```

## Configuration

All settings are loaded from `config/local.yaml`. Example:

```yaml
env: "local"
grpc:
  port: 50051
mongo:
  address: "mongodb://localhost:27017"
token_ttl: "15m"
```

## Getting Started

### Prerequisites

- Go 1.20+
- Running MongoDB instance
- Protobuf compiler (`protoc`) for gRPC

### Installation

```bash
git clone https://github.com/cosmowake/sso.git
cd sso
go mod tidy
```

### Run the Service

```bash
go run cmd/sso/main.go --config=./config/local.yaml
```

### Build Binary

```bash
go build -o sso ./cmd/sso
./sso --config=./config/local.yaml
```

## Testing

Use `grpcurl` or Postman with gRPC support to test the endpoints.
