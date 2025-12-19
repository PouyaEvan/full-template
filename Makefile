.PHONY: migrate-up migrate-down run-backend run-frontend

# Database Migrations
migrate-up:
	migrate -path backend/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" -verbose up

migrate-down:
	migrate -path backend/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable" -verbose down

# Swagger
swag:
	cd backend && swag init -g cmd/api/swagger.go -o docs

# Development
run-backend:
	cd backend && go run cmd/api/main.go cmd/api/swagger.go

run-frontend:
	cd frontend && npm run dev

# Docker
docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down
