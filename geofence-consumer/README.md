# Geofence Consumer Service

A Go-based microservice that consumes geofence events from RabbitMQ for the Fleet Management API. This service processes real-time vehicle location data to detect when vehicles reach points of interest.

## Overview

The Geofence Consumer Service is part of the Fleet Management API ecosystem. It listens for geofence events published via RabbitMQ and processes them using the GeoFence service. The service integrates with PostgreSQL for data persistence and Redis for caching.

## Architecture

- **Transport Layer**: RabbitMQ consumer for real-time event processing
- **Service Layer**: Business logic for geofence event handling
- **External Dependencies**: PostgreSQL, Redis, and RabbitMQ connections
- **Logging**: Structured logging with Zap

## Features

- Real-time geofence event consumption from RabbitMQ
- Processing of "Reached Nearest Point of Interest" events
- Health check endpoints
- Configurable RabbitMQ consumer settings
- Graceful shutdown handling

## Prerequisites

- Go 1.24.4 or later
- PostgreSQL database
- Redis server
- RabbitMQ server
- Docker (for code generation and optional services)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/elzestia/fleet.git
   cd geofence-consumer
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables (copy and modify `.env.example` to `.env`):
   ```bash
   cp .env.example .env
   ```

## Configuration

The service uses Viper for configuration management. Key configuration sections:

- **App**: Application name, environment, log level
- **Database**: PostgreSQL and Redis connection settings
- **RabbitMQ**: Broker connection and consumer configuration

Environment variables override `.env` file values. Use the format `APP.ENV` or `APP_ENV`.

## Running the Service

### Development Mode
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
docker build -t geofence-consumer .
docker run geofence-consumer
```

## Development Commands

- `make build` - Build the service
- `make dev` - Run in development mode with hot reload
- `make test` - Run tests
- `make lint` - Lint the code
- `make docs` - Generate Swagger documentation
- `make gen-db` - Generate database code using SQLC
- `make install-tools` - Install development tools

## API Documentation

API documentation is generated using Swagger. After running `make docs`, access the docs at `/swagger/index.html` when the service is running.

## Testing

Run tests with:
```bash
make test
```

## Health Checks

The service includes health check endpoints to verify database and Redis connectivity.

## Event Processing

The service consumes `ReachedNearestPointOfInterestEvent` messages from RabbitMQ and processes them through the GeoFence service. Events are logged and can be extended for additional business logic.

## Contributing

1. Follow the existing code structure and patterns
2. Add tests for new functionality
3. Update documentation as needed
4. Ensure code passes linting
