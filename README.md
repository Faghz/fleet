# Fleet Services - Docker Compose Setup

This repository contains a multi-service application for fleet management, including a main API service, geofence consumer, and location sync worker. The services are orchestrated using Docker Compose.

## Prerequisites

- Docker (version 20.10 or later)
- Docker Compose (version 2.0 or later)

## Services Overview

The `docker-compose.yml` file defines the following services:

- **postgres**: PostgreSQL database (port 5432)
- **redis**: Redis cache (port 6379)
- **rabbitmq**: RabbitMQ message broker (ports 5672, 15672)
- **mqtt**: Eclipse Mosquitto MQTT broker (ports 1883, 9001)
- **fleet**: Main fleet API service (port 8080)
- **geofence-consumer**: Service for processing geofence alerts
- **location-sync-worker**: Worker for synchronizing vehicle locations via MQTT

## Getting Started

1. **Clone the repository** (if not already done):
   ```bash
   git clone <repository-url>
   cd fleet
   ```

2. **Start the services**:
   ```bash
   docker-compose up -d
   ```

   This will build the custom services (`fleet`, `geofence-consumer`, `location-sync-worker`) and start all containers in detached mode.

3. **Check the status**:
   ```bash
   docker-compose ps
   ```

   Ensure all services are healthy. The services have health checks configured, so they may take a few minutes to start.

## Accessing Services

- **Fleet API**: http://localhost:8080
- **RabbitMQ Management UI**: http://localhost:15672 (username: guest, password: guest)
- **PostgreSQL**: localhost:5432 (user: fleet_user, password: fleet_password, database: trans)
- **Redis**: localhost:6379
- **MQTT**: localhost:1883

## Environment Configuration

The services use environment variables for configuration. These are defined in the `docker-compose.yml` file. For production deployments, consider using a `.env` file or external configuration management.

## Stopping Services

To stop all services:
```bash
docker-compose down
```

To stop and remove volumes (this will delete data):
```bash
docker-compose down -v
```

## Logs

View logs for all services:
```bash
docker-compose logs -f
```

View logs for a specific service:
```bash
docker-compose logs -f <service-name>
```

## Development

For development, you can rebuild and restart services after code changes:
```bash
docker-compose up --build -d <service-name>
```

## Troubleshooting

- Ensure Docker and Docker Compose are installed and running.
- Check that ports 5432, 6379, 5672, 15672, 1883, 9001, and 8080 are not in use by other applications.
- If services fail to start, check logs with `docker-compose logs <service-name>`.
- For database issues, you may need to reset volumes: `docker-compose down -v && docker-compose up -d`.

## Additional Notes

- The `mosquitto.conf` file is mounted into the MQTT container for configuration.
- Database migrations and seeds are handled automatically by the services.
- All services are connected via the `fleet-network` bridge network.