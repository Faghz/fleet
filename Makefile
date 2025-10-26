include .env
# dbmate connection string format
DATABASE_URL=postgres://$(DATABASE.POSTGRESQL.USER:"%"=%):$(DATABASE.POSTGRESQL.PASSWORD:"%"=%)@$(DATABASE.POSTGRESQL.HOST:"%"=%):$(DATABASE.POSTGRESQL.PORT:"%"=%)/$(DATABASE.POSTGRESQL.DB_NAME:"%"=%)?sslmode=${DATABASE.POSTGRESQL.SSL_MODE:"%"=%}
export
# Build and development commands
.PHONY: build dev test lint clean install-tools generate docs generate-auth-key generate-auth-key-save gen-mocks clean-mocks

build:
	go build -o bin/server ./cmd/server

dev:
	go run github.com/air-verse/air@latest -c ./configs/.air.toml

test:
	go test -v ./...

# Install development tools (simplified approach)
install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest && \
	go install github.com/golang/mock/mockgen@latest && \
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install tools using go generate in tools package  
generate-tools:
	go generate ./tools

# Generate swagger docs
docs:
	swag init -g ./cmd/server/main.go -o ./cmd/server/docs

# Generate database code using Docker (recommended)
gen-db:
	@echo "Generating database code using Docker..."
	docker run --rm -v $(PWD):/src -w /src kjconroy/sqlc:latest generate -f ./configs/sqlc.yaml

# Alternative: Generate database code using Docker
gen-db-docker:
	docker run --rm -v $(PWD):/src -w /src kjconroy/sqlc:latest generate -f ./configs/sqlc.yaml


# Lint code
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/ docs/

# SQL linting
lint-sql:
	@read -p "sql file path (ex: ./dev/db/migrations/20230602080910-migration-file.sql): " FILEPATH; \
  	read -p "dialect (ex: mysql|postgres|etc): " DIALECT; \
	docker run -it --rm -v $(PWD):/sql sqlfluff/sqlfluff:2.1.0 lint $${FILEPATH} -d $${DIALECT}

# fix lint sql file
fix-sql:
	@read -p "sql file path (ex: ./dev/db/migrations/20230602080910-migration-file.sql): " FILEPATH; \
  	read -p "dialect (ex: mysql|postgres|etc): " DIALECT; \
	docker run -it --rm -v $(PWD):/sql sqlfluff/sqlfluff:2.1.0 fix $${FILEPATH} -d $${DIALECT}

# Database migration commands (using dbmate in Docker with host network)
migrate-create:
	@read -p "migration name (do not use space): " NAME \
	&& docker run --rm -it \
		--network host \
		-e DATABASE_URL="$(DATABASE_URL)" \
		-v $(PWD)/db:/db \
		amacneil/dbmate:latest new $${NAME}

migrate-up:
	@docker run --rm \
		--network host \
		-e DATABASE_URL="$(DATABASE_URL)" \
		-v $(PWD)/db:/db \
		amacneil/dbmate:latest --wait up && $(MAKE) migrate-dump

migrate-down:
	@docker run --rm \
		--network host \
		-e DATABASE_URL="$(DATABASE_URL)" \
		-v $(PWD)/db:/db \
		amacneil/dbmate:latest --wait down && $(MAKE) migrate-dump

migrate-rollback:
	@docker run --rm \
		--network host \
		-e DATABASE_URL="$(DATABASE_URL)" \
		-v $(PWD)/db:/db \
		amacneil/dbmate:latest --wait rollback && $(MAKE) migrate-dump

migrate-status:
	@docker run --rm \
		--network host \
		-e DATABASE_URL="$(DATABASE_URL)" \
		-v $(PWD)/db:/db \
		amacneil/dbmate:latest status

migrate-drop:
	@docker run --rm \
		--network host \
		-e DATABASE_URL="$(DATABASE_URL)" \
		-v $(PWD)/db:/db \
		amacneil/dbmate:latest --wait drop

# Dump current database schema (useful for version control)
migrate-dump:
	@docker run --rm \
		--network host \
		-e DATABASE_URL="$(DATABASE_URL)" \
		-v $(PWD)/db:/db \
		amacneil/dbmate:latest dump


# Dump current database schema (useful for version control)
migrate-dump:
	@docker run --rm --network host \
		-e DATABASE_URL="$(DATABASE_URL)" \
		-v $(PWD)/db:/db \
		amacneil/dbmate:latest dump


show-db-conn:
	@echo $(DATABASE_URL)

generate:
	@echo "Running go generate on all packages..."
	go generate ./...

generate-auth-key:
	@echo "Generating auth key..."
	go run ./cmd/generate-auth-key/main.go

generate-auth-key-save:
	@echo "Generating auth key and saving to file..."
	go run ./cmd/generate-auth-key/main.go --save


# Help command
help:
	@echo "Available commands:"
	@echo "  build        			- Build the server binary"
	@echo "  dev          			- Run development server with hot reload"
	@echo "  test         			- Run tests"
	@echo "  lint         			- Run code linter"
	@echo "  install-tools 			- Install development tools"
	@echo "  generate     			- Run go generate on all packages"
	@echo "  generate-auth-key 		- Generate auth key and print to stdout"
	@echo "  generate-auth-key-save - Generate auth key and save to .auth-key file"
	@echo "  seed-admin       		- Seed default admin user to database"
	@echo "  seed-categories  		- Seed default journal categories to database"
	@echo "  seed-transactions		- Seed sample transactions for testing"
	@echo "  docs         			- Generate swagger documentation"
	@echo "  gen-db       			- Generate database code with sqlc"
	@echo "  gen-mocks    			- Generate mocks for all repository interfaces"
	@echo "  clean-mocks  			- Clean generated mock files"
	@echo "  clean        			- Clean build artifacts"
	@echo "  migrate-create     	- Create a new migration"
	@echo "  migrate-up        		- Apply all pending migrations"
	@echo "  migrate-down      		- Rollback the last migration"
	@echo "  migrate-rollback  		- Rollback the last migration"
	@echo "  migrate-status     	- Show migration status"
	@echo "  migrate-drop       	- Drop database"
	@echo "  migrate-wait      		- Wait for database to become available"
	@echo "  migrate-dump      		- Dump current database schema (docker)"
	@echo "  help         			- Show this help message"