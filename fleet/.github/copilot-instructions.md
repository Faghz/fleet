# Copilot Instructions for Fleet Management API

## Architecture Overview

This is a **Go HTTP API service with MQTT integration** using clean architecture with clear layer separation:

- **HTTP Transport** (`pkg/transport/http/`) - Fiber HTTP handlers, validation, middleware
- **MQTT Transport** (`pkg/transport/mqtt/`) - Real-time vehicle data handlers using Eclipse Paho MQTT client
- **Services** (`service/`) - Business logic layer (user, auth, vehicle)
- **Repository** (`pkg/repository/`) - Data access via SQLC-generated code + manual cache layer
- **External** (`pkg/external/`) - Infrastructure connections (PostgreSQL, Redis, MQTT broker)

**Key Pattern**: Each service defines its own interface to decouple from concrete repository implementations, enabling easy mocking.

## Code Generation Workflow

This project heavily relies on code generation. **Always regenerate after modifying these sources:**

### Database Layer (SQLC)
- **Source**: SQL queries in `db/queries/*.sql` + schema in `db/schema.sql`
- **Generate**: `make gen-db` (uses Docker, no local sqlc needed)
- **Output**: Generated code goes to `./generated/` directory
  - `generated/*sql.go` - Query implementations
  - `generated/querier.go` - Query interface
  - `generated/models.go` - DB model structs
- **Manual steps after generation**:
  1. Copy repository implementations (`*sql.go`, `querier.go`, `db.go`) to `pkg/repository/`
  2. Copy model structs from `generated/models.go` to `pkg/models/` (if needed)
  3. Update import paths if necessary
- **Config**: `configs/sqlc.yaml` targets PostgreSQL with pgx/v5, outputs to `../generated`

### API Documentation (Swagger)
- **Source**: Annotations in HTTP handler files (`@Summary`, `@Router`, etc.)
- **Generate**: `make docs`
- **Output**: `cmd/server/docs/` (swagger.json, swagger.yaml, docs.go)

### Mocks (gomock)
- **Source**: `//go:generate` directives on service interfaces
- **Generate**: `go generate ./service/...` or `go generate ./...`
- **Example**: `//go:generate mockgen -destination=../../pkg/mocks/user_repository_mock.go -package=mocks github.com/elzestia/fleet/service/user UserRepository`
- **Output**: `pkg/mocks/*_mock.go` files for each service interface
- **Usage**: Service tests import from `pkg/mocks`, see `service/user/user_test.go` for patterns

## Database Management (dbmate + Docker)

**All migrations use Docker** - no local dbmate needed:

```bash
make migrate-create        # Prompts for migration name (no spaces)
make migrate-up            # Apply pending + auto-dump schema
make migrate-status        # Check current state
```

**Migration Format** (single file, up/down):
```sql
-- migrate:up
CREATE TABLE example (...);
-- migrate:down
DROP TABLE IF EXISTS example;
```

**Connection**: Configured via `.env` → `DATABASE_URL` (Postgres format). The Makefile auto-constructs this from individual env vars.

**Schema Dump**: `db/schema.sql` is auto-generated after each migration (version controlled).

**Vehicle Tables**: Core entities include `vehicle` (vehicle metadata) and `vehicle_location` (GPS tracking with timestamps).

## MQTT Integration

**Real-time vehicle tracking** using Eclipse Mosquitto broker:

- **MQTT Transport**: `pkg/transport/mqtt/` handles real-time vehicle location updates
- **Broker**: Mosquitto configured in `docker-compose.yml` (ports 1883 MQTT, 9001 WebSocket)
- **Client**: Paho MQTT client with auto-reconnect and connection handlers
- **Topics**: Fleet-related topics handled in `pkg/transport/mqtt/handler/fleet.go`
- **Configuration**: `MQTTConfig` in `configs/models.go` with broker, credentials, QoS settings

**MQTT Handler Pattern**:
```go
mqttHandler := CreateMqttConsumer(cfg, logger, mqttClient, vehicleService)
// Initializes fleet handlers automatically
```

## Authentication & Session Management

**Token System**: PASETO tokens (v4.public) with session validation
- **Generate keys**: `make generate-auth-key-save` → stores in `.auth-key`
- **Auth flow**:
  1. Login → `service/user/auth.go` validates password (bcrypt)
  2. Generate PASETO token with claims (subject=userID, jti=sessionID, orgID)
  3. Store session in PostgreSQL + Redis cache (`pkg/repository/sesion.cache.go`)
  4. Middleware (`pkg/api/http/handler/midlleware.go`) verifies token + checks session cache/DB

**Session Cache Pattern**:
- Key format: `session:{userID}:{sessionID}`
- Check cache first (`GetSessionCache`), fallback to DB, then populate cache
- See `pkg/repository/sesion.cache.go` for implementation

## Configuration Management

**Viper-based hierarchy** (see `configs/loader.go`):
1. `.env` file (root or auto-detected from go.mod)
2. Environment variables override `.env` values
3. Supports both formats: `APP.ENV` or `APP_ENV`

**Nested struct** (`configs/models.go`): `Config` → `AppConfig`, `HttpConfig`, `Database`, `MQTTConfig`, `FunctionConfig` (auth settings)

## Request/Response Patterns

### Validation
- **Custom Validator**: `pkg/api/http/validator.go` + `custom_validator.go`
- **Usage**: Call `inthttp.GetValidator().Validate(&req)` in handlers after parsing request body
- **Tags**: Standard go-playground/validator (`required`, `email`, `min`, `max`, `oneof`)
- **Error formatting**: Auto-translates to `response.FailureError` with pointer paths

### Request Parsing
- **HTTP**: Use `c.BodyParser(&req)` to parse JSON request body
- **Context**: Use `c.UserContext()` to get request context for service calls

### Response Structure
```go
response.BaseResponse{Status, Message, Data}
response.PaginatedResponse{Items, Metadata{Page, Limit, TotalPages}}
```

**Helpers**:
- `response.ResponseJson(c, status, message, data)` - Success (Fiber context)
- `response.GenerateFailure(status, message, details)` - Error
- `response.GenerateBadRequest(...)` - Validation errors

## Testing Conventions

**Table-driven tests with gomock** (see `service/user/user_test.go`):

```go
func createTestUserService(t *testing.T) (*UserService, *mocks.MockUserRepository, *mocks.MockUserRepository) {
    ctrl := gomock.NewController(t)
    mockRepo := mocks.NewMockUserRepository(ctrl)
    // ... return service with mock injected
}

// In test:
mockRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(user, nil)
```

**Run tests**: `make test` or `go test -v ./...`

## Development Commands

```bash
make dev                   # Hot reload via Air (requires .air.toml config)
make build                 # Build to bin/server
make test                  # Run tests
make lint                  # golangci-lint
make docs                  # Generate swagger docs
make gen-db                # Generate database code with sqlc
generate                  # Run go generate on all packages (includes mocks)
make migrate-up            # Apply DB migrations
make generate-auth-key-save # Generate and save auth keys
```

**Docker Services**:
```bash
docker-compose up mosquitto  # Start MQTT broker
docker-compose up rabbitmq   # Start RabbitMQ (if used)
```

## Unique Patterns & Gotchas

### MQTT Message Handling
- **Fleet Handler**: `createFleetHandler()` initializes MQTT topic subscriptions
- **Connection Management**: Auto-reconnect with 10s retry interval, 1min max interval
- **QoS Settings**: Configurable per-message delivery guarantees

### Vehicle Location Tracking
- **Dual Storage**: Vehicle metadata in `vehicle` table, locations in `vehicle_location` with timestamps
- **Latest Location Query**: Complex SQL joining vehicle + latest location by timestamp
- **Location History**: Count locations with `CountVehicleLocationsHistory()`

### Snowflake IDs
- Repository embeds `snowflake.Node` (initialized in `repository.CreateRepository`)
- Generate via `r.GenerateSnowflakeID()` for distributed unique IDs

### Transaction Pattern
```go
tx, _ := repo.BeginTx(ctx)
defer repo.RollbackTx(tx)  // Safe even if committed
querierWithTx := repo.WithTx(tx)
// ... use querierWithTx for DB ops
repo.CommitTx(tx)
```

### Custom Error Handling
- Domain errors in `pkg/api/http/response/error.go` (e.g., `ErrorUserDatabaseUserNotFound`)
- Services return domain errors, handlers map to HTTP status codes

### Email Hashing
- User queries use `email_hash` (SHA256) for indexed lookups
- Pattern in `db/queries/user.sql`: `WHERE email_hash = sqlc.arg('email')`

## File Organization Rules

- **No mixing layers**: HTTP handlers call services (not repos), services call repos
- **MQTT handlers**: Real-time message processing in `pkg/transport/mqtt/handler/`
- **Request/Response DTOs**: Always in `pkg/api/http/request|response/` (HTTP only)
- **Domain models**: `pkg/models/` (separate from DB models in `pkg/repository/`)
- **Generated code workflow**:
  - SQLC generates to `./generated/` first
  - Repository implementations (`*sql.go`, `querier.go`, `db.go`) → copy to `pkg/repository/`
  - Model structs → copy from `generated/models.go` to `pkg/models/` as needed
  - Never edit generated files directly - modify source SQL queries instead
