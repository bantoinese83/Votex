.PHONY: help dev dev-frontend dev-backend build build-frontend build-backend install test test-frontend test-backend test-integration lint lint-frontend lint-backend format format-frontend format-backend clean storybook db-push db-studio docker-build docker-up docker-down docker-dev docker-logs migrate-up migrate-down

# Default target
help:
	@echo "Available commands:"
	@echo "  make dev          - Run both frontend and backend in development mode"
	@echo "  make dev-frontend - Run only frontend in development mode"
	@echo "  make dev-backend  - Run only backend in development mode"
	@echo "  make build        - Build both frontend and backend for production"
	@echo "  make build-frontend - Build frontend for production"
	@echo "  make build-backend - Build backend for production"
	@echo "  make install      - Install all dependencies (root + frontend + backend)"
	@echo "  make test         - Run all tests (frontend + backend)"
	@echo "  make test-frontend - Run frontend tests"
	@echo "  make test-backend - Run backend tests"
	@echo "  make test-integration - Run integration tests"
	@echo "  make lint         - Run all linting (frontend + backend)"
	@echo "  make lint-frontend - Run frontend linting"
	@echo "  make lint-backend - Run backend linting"
	@echo "  make format       - Format all code (frontend + backend)"
	@echo "  make format-frontend - Format frontend code"
	@echo "  make format-backend - Format backend code"
	@echo "  make clean        - Clean build artifacts and containers"
	@echo "  make storybook    - Run Storybook"
	@echo "  make db-push      - Push database schema"
	@echo "  make db-studio    - Open Drizzle Studio"
	@echo "  make docker-build - Build Docker images"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make docker-dev   - Start development containers"
	@echo "  make docker-logs  - View Docker logs"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"

# Run both frontend and backend
dev:
	@echo "Starting both frontend and backend..."
	@npm run dev

# Run only frontend
dev-frontend:
	@echo "Starting frontend development server..."
	@npm run dev:frontend

# Run only backend
dev-backend:
	@echo "Starting backend server..."
	@npm run dev:backend

# Build both frontend and backend
build:
	@echo "Building both frontend and backend..."
	@npm run build

# Build frontend only
build-frontend:
	@echo "Building frontend..."
	@npm run build:frontend

# Build backend only
build-backend:
	@echo "Building backend..."
	@npm run build:backend

# Install all dependencies
install:
	@echo "Installing all dependencies..."
	@npm run install:all

# Run all tests
test:
	@echo "Running all tests..."
	@npm run test

# Run frontend tests
test-frontend:
	@echo "Running frontend tests..."
	@npm run test:frontend

# Run backend tests
test-backend:
	@echo "Running backend tests..."
	@npm run test:backend

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@npm run test:integration

# Run all linting
lint:
	@echo "Running all linting..."
	@npm run lint

# Run frontend linting
lint-frontend:
	@echo "Running frontend linting..."
	@npm run lint:frontend

# Run backend linting
lint-backend:
	@echo "Running backend linting..."
	@npm run lint:backend

# Format all code
format:
	@echo "Formatting all code..."
	@npm run format

# Format frontend code
format-frontend:
	@echo "Formatting frontend code..."
	@npm run format:frontend

# Format backend code
format-backend:
	@echo "Formatting backend code..."
	@npm run format:backend

# Clean build artifacts and containers
clean:
	@echo "Cleaning build artifacts and containers..."
	@npm run clean

# Run Storybook
storybook:
	@echo "Starting Storybook..."
	@npm run storybook

# Database operations
db-push:
	@echo "Pushing database schema..."
	@npm run db:push

db-studio:
	@echo "Opening Drizzle Studio..."
	@npm run db:studio

# Docker operations
docker-build:
	@echo "Building Docker images..."
	@npm run docker:build

docker-up:
	@echo "Starting Docker containers..."
	@npm run docker:up

docker-down:
	@echo "Stopping Docker containers..."
	@npm run docker:down

docker-dev:
	@echo "Starting development containers..."
	@npm run docker:dev

docker-logs:
	@echo "Viewing Docker logs..."
	@npm run docker:logs

# Migration operations
migrate-up:
	@echo "Running database migrations..."
	@npm run migrate:up

migrate-down:
	@echo "Rolling back database migrations..."
	@npm run migrate:down 