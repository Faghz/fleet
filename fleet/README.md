# Fleet Management API

A Go-based HTTP API service with MQTT integration for real-time fleet management, featuring vehicle tracking, location history, and points of interest management.

## Features

- **Vehicle Management**: CRUD operations for vehicles with metadata storage
- **Real-time Location Tracking**: MQTT-based vehicle location updates with timestamp validation
- **Location History**: Retrieve historical vehicle locations within time ranges
- **Points of Interest**: Manage and query points of interest for fleet operations
- **Authentication & Sessions**: PASETO token-based authentication with session management
- **Clean Architecture**: Layered design with HTTP transport, services, and repository patterns
- **Code Generation**: SQLC for database queries, Swagger for API docs, gomock for testing

## Architecture

The service follows clean architecture principles with clear layer separation:

- **HTTP Transport** (`pkg/transport/http/`): Fiber HTTP handlers, validation, middleware
- **MQTT Transport** (`pkg/transport/mqtt/`): Real-time vehicle data handlers using Eclipse Paho MQTT client
- **Services** (`service/`): Business logic layer (user, auth, vehicle, point_of_interest)
- **Repository** (`pkg/repository/`): Data access via SQLC-generated code + manual cache layer
- **External** (`pkg/external/`): Infrastructure connections (PostgreSQL, Redis, MQTT broker)

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (via Docker)
- Redis (via Docker)
- Mosquitto MQTT Broker (via Docker)

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd fleet
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Start infrastructure services:
   ```bash
   docker-compose up -d postgres redis mosquitto
   ```

4. Run database migrations:
   ```bash
   make migrate-up
   ```

5. Generate authentication keys:
   ```bash
   make generate-auth-key-save
   ```

6. Generate database code:
   ```bash
   make gen-db
   ```

7. Generate mocks (optional, for testing):
   ```bash
   go generate ./...
   ```

## Configuration

The service uses Viper for configuration hierarchy:

1. `.env` file in the project root
2. Environment variables override `.env` values

Key configuration options in `.env`:

```env
# Application
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=fleet
DB_PASSWORD=password
DB_NAME=fleet

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# MQTT
MQTT_BROKER=tcp://localhost:1883
MQTT_CLIENT_ID=fleet-api
MQTT_USERNAME=
MQTT_PASSWORD=

# Authentication
AUTH_PUBLIC_KEY=<generated-public-key>
AUTH_PRIVATE_KEY=<generated-private-key>
```

## Running the Service

### Development Mode (with hot reload)
```bash
make dev
```

### Production Build
```bash
make build
./bin/server
```

### Docker
```bash
docker build -t fleet-api .
docker run -p 8080:8080 fleet-api
```

## Development Commands

```bash
make dev                   # Hot reload via Air
make build                 # Build to bin/server
make test                  # Run tests
make lint                  # golangci-lint
make docs                  # Generate swagger docs
make gen-db                # Generate database code with sqlc
go generate ./...          # Run go generate on all packages (includes mocks)
make migrate-up            # Apply DB migrations
make generate-auth-key-save # Generate and save auth keys
```

## API Documentation

API documentation is generated using Swagger. After running `make docs`, visit:

- Swagger UI: `http://localhost:8080/swagger/index.html`
- Swagger JSON: `http://localhost:8080/swagger/doc.json`

## Testing

Run tests with:
```bash
make test
```

Tests use table-driven patterns with gomock for mocking dependencies.

## Database Management

Migrations are handled via dbmate with Docker:

```bash
make migrate-create        # Create new migration
make migrate-up            # Apply pending migrations
make migrate-status        # Check migration status
```

## MQTT Integration

The service integrates with Mosquitto MQTT broker for real-time vehicle tracking:

- **Topics**: Fleet-related topics handled in `pkg/transport/mqtt/handler/fleet.go`
- **QoS**: Configurable delivery guarantees
- **Connection**: Auto-reconnect with configurable retry intervals

Start MQTT broker:
```bash
docker-compose up mosquitto
```

## Authentication

Uses PASETO v4.public tokens with session validation:

1. Login to receive token
2. Token contains user ID, session ID, and organization ID
3. Sessions stored in PostgreSQL with Redis caching
4. Middleware validates tokens and checks session validity

## Key Patterns

- **Distributed Locking**: Redis-based mutexes for concurrent vehicle location updates
- **Transaction Management**: Proper transaction handling with rollback on errors
- **Timestamp Validation**: Only update locations with newer timestamps
- **Snowflake IDs**: Distributed unique ID generation
- **Email Hashing**: SHA256 hashing for indexed user lookups
