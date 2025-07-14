.PHONY: help lint lint-go lint-web format format-web build test test-go test-web clean dev-up dev-down dev-logs dev-restart db-reset

help:
	@echo "CoScribe - Makefile Commands"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint: lint-go lint-web 

lint-go:
	@docker run --rm -v $(PWD):/app -w /app golang:1.22-alpine sh -c "go mod tidy && go vet ./..."

lint-web:
	@docker run --rm -v $(PWD)/web:/app -w /app node:18-alpine sh -c "npm install --silent && npm run lint"

format-web:
	@docker run --rm -v $(PWD)/web:/app -w /app node:18-alpine sh -c "npm install --silent && npm run format"


test: test-go test-web

test-go:
	@docker run --rm -v $(PWD):/app -w /app golang:1.22-alpine sh -c "go mod tidy && go test -v ./..."

test-web:
	@docker run --rm -v $(PWD)/web:/app -w /app node:18-alpine sh -c "npm install --silent && npm test -- --watchAll=false --verbose"

build:
	@docker compose build

dev-up:
	@docker compose up --build -d

dev-down:
	@docker compose down

dev-restart: dev-down dev-up

dev-logs:
	@docker compose logs -f

dev-status:
	@docker compose ps

db-reset:
	@docker compose down postgres
	@docker volume rm coscribe_postgres_data 2>/dev/null || true
	@docker compose up postgres -d

db-logs:
	@docker compose logs -f postgres

clean-cache:
	@docker builder prune -f

install-deps:
	@docker run --rm -v $(PWD):/app -w /app golang:1.22-alpine go mod tidy
	@docker run --rm -v $(PWD)/web:/app -w /app node:18-alpine npm install

ci: lint test

pre-commit: format-web lint test

quick-start: clean dev-up

health-check:
	@curl -f http://localhost:8080/health 2>/dev/null && echo "✅ Backend is healthy" || echo "❌ Backend is unhealthy"
	@curl -f http://localhost:3000 2>/dev/null && echo "✅ Frontend is accessible" || echo "❌ Frontend is not accessible"

dev-web:
	@cd web && npm start

dev-full:
	@docker compose up --build -d
	@sleep 5
	@cd web && BROWSER=none npm start

dev-full-stop:
	@docker compose down
	@pkill -f "react-scripts start" || true