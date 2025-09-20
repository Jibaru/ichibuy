run:
	@go run cmd/app/main.go

build:
	@swag init -g cmd/app/main.go
	@go build -o bin/app cmd/app/main.go

gen:
	@rm -R internal/infra/persistence/postgres
	@rm -R internal/domain/dao
	@gormless -i internal/domain -o internal/infra/persistence --driver postgres
	@gormless -i internal/domain -o internal/domain/dao --interface

# Database migration commands
migrate-up:
	@echo "Applying database migrations..."
	@goose up

migrate-down:
	@echo "Rolling back last migration..."
	@goose down

migrate-status:
	@echo "Checking migration status..."
	@goose status

migrate-reset:
	@echo "Rolling back all migrations..."
	@goose reset

# Development helpers
dev-setup: migrate-up
	@echo "Development environment setup complete"

.PHONY: run build gen migrate-up migrate-down migrate-status migrate-reset dev-setup