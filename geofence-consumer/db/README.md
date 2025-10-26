# Database Management

This directory contains all database-related files for the project, including migrations, queries, and schema documentation.

## Directory Structure

```
db/
├── migrations/           # Database migration files (dbmate format)
├── queries/             # Raw SQL query files (used by sqlc)
│   ├── auth.sql        # Authentication-related queries
│   ├── session.sql     # Session management queries
│   └── user.sql        # User management queries
├── schema.sql          # Current database schema (auto-generated)
└── README.md          # This file
```

## Migration System

This project uses [dbmate](https://github.com/amacneil/dbmate) for database migrations.

### Prerequisites

Before using the migration system, ensure you have the following installed:

#### 1. Install dbmate
```bash
# Install dbmate via Go
go install github.com/amacneil/dbmate/v2@latest

# Or install via the project's Makefile
make install-tools
```

### Configuration

dbmate is configured via:
- `.env` file: Contains `DATABASE_URL` for database connection
- `.dbmate` file: Contains dbmate-specific configuration (migrations directory, schema file location)
- `Makefile`: Contains convenient migration commands

### Available Commands

#### Migration Management
```bash
make migrate-create          # Create a new migration
make migrate-up              # Apply all pending migrations + dump schema
make migrate-down            # Rollback the last migration + dump schema  
make migrate-rollback        # Rollback the last migration + dump schema
make migrate-status          # Show migration status
make migrate-drop            # Drop the entire database
make migrate-wait            # Wait for database to become available
```

#### Schema Management
```bash
make migrate-dump            # Dump schema using dbmate (may not work with Docker)
```

#### Development Tools
```bash
make install-tools           # Install dbmate and other development tools
make show-db-conn           # Display current database connection string
```

## Migration File Format

Migrations are stored in `./db/migrations/` and follow the dbmate format with both up and down operations in a single file:

```sql
-- migrate:up
CREATE TABLE example (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_example_name ON example (name);

-- migrate:down
DROP TABLE IF EXISTS example;
```

### Migration Naming Convention
- Format: `YYYYMMDDHHMMSS_migration_name.sql`
- Use underscores instead of hyphens in migration names
- Use descriptive names: `create_user_table`, `add_email_index`, etc.

### Creating a New Migration

1. **Using Make command:**
   ```bash
   make migrate-create
   # You'll be prompted to enter a migration name (no spaces)
   ```

2. **Manual creation:**
   ```bash
   dbmate new create_customer_table
   ```

This creates a new file in `./db/migrations/` with the current timestamp.

## Schema Management

### Auto-Generated Schema
The `schema.sql` file is automatically generated and updated after each migration operation. This file:
- Contains the complete current database structure
- Should be committed to version control
- Helps track schema changes over time
- Is used for database documentation

### Manual Schema Dump
If you need to manually regenerate the schema file:
```bash
make migrate-dump_manual
```

## Docker MySQL Integration

This project is configured to work with MySQL running in Docker containers:

### Connection Method
- Uses TCP connections instead of Unix sockets
- Configured via `DATABASE_URL` with `protocol=tcp` parameter
- Works with both local Docker and remote MySQL instances

### Troubleshooting Docker Issues

**If migrations fail to connect:**
1. Ensure MySQL container is running: `docker ps`
2. Check port mapping: `docker port <container_name>`
3. Verify connection: `make migrate-status`

**If schema dump fails:**
- Use `make migrate-dump_manual` instead of `make migrate-dump`
- The manual command uses `mysqldump` with proper TCP parameters

## Development Workflow

### Initial Setup
1. Install prerequisites (dbmate, mysql-client)
2. Start MySQL container
3. Run initial migrations: `make migrate-up`

### Adding New Features
1. Create migration: `make migrate-create`
2. Edit the generated migration file
3. Apply migration: `make migrate-up`
4. Commit both the migration file and updated `schema.sql`

### Rolling Back Changes
```bash
make migrate-down    # Rollback last migration
make migrate-status  # Check current state
```

## Files Integration

### sqlc Integration
The `db/queries/` directory contains SQL files used by [sqlc](https://sqlc.dev/) to generate Go code:
- `auth.sql` - Authentication-related queries
- `user.sql` - User management queries  
- `sequence.sql` - Sequence/ID generation queries
- and more added as needed

### Configuration Files
- `configs/sqlc.yaml` - sqlc configuration
- `.env` - Database connection parameters
- `.dbmate` - dbmate-specific settings

## Environment Variables

Required environment variables in `.env`:
```ini
# Database Connection
DATABASE.MYSQL.HOST=localhost
DATABASE.MYSQL.PORT=3306
DATABASE.MYSQL.USER=root
DATABASE.MYSQL.PASSWORD=your-password
DATABASE.MYSQL.DB_NAME=your-database

# dbmate URL (auto-constructed from above or set manually)
DATABASE_URL=mysql://user:password@host:port/database?protocol=tcp
```

## Common Issues and Solutions

### Migration Fails with "Can't connect to MySQL server"
- **Cause**: MySQL container not running or wrong connection parameters
- **Solution**: Check Docker containers and verify `.env` configuration

### Schema dump produces empty file
- **Cause**: `dbmate dump` doesn't work well with Docker MySQL
- **Solution**: Use `make migrate-dump_manual` instead

### Migration file has wrong format
- **Cause**: Missing `-- migrate:up` and `-- migrate:down` sections
- **Solution**: Ensure both sections are present in migration files

### Permission denied on mysqldump
- **Cause**: mysql-client not installed or not in PATH
- **Solution**: Install MySQL client tools following the instructions above

## Best Practices

1. **Always test migrations:**
   - Test both up and down migrations
   - Use `migrate-status` to verify state

2. **Keep migrations atomic:**
   - One logical change per migration
   - Use transactions where appropriate

3. **Version control:**
   - Commit migration files and `schema.sql` together
   - Never edit applied migrations

4. **Backup before major changes:**
   ```bash
   make migrate-dump_manual  # Create schema backup
   ```

5. **Use descriptive names:**
   - Good: `create_user_email_index`
   - Bad: `update_table`, `fix_bug`
