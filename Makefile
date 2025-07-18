.PHONY: help dev dev-frontend dev-backend build install test lint format clean

# Default target
help:
	@echo "Available commands:"
	@echo "  make dev          - Run both frontend and backend in development mode"
	@echo "  make dev-frontend - Run only frontend in development mode"
	@echo "  make dev-backend  - Run only backend in development mode"
	@echo "  make build        - Build frontend for production"
	@echo "  make install      - Install all dependencies (root + frontend)"
	@echo "  make test         - Run frontend tests"
	@echo "  make lint         - Run frontend linting"
	@echo "  make format       - Format frontend code"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make storybook    - Run Storybook"
	@echo "  make db-push      - Push database schema"
	@echo "  make db-studio    - Open Drizzle Studio"

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

# Build frontend
build:
	@echo "Building frontend..."
	@npm run build

# Install all dependencies
install:
	@echo "Installing all dependencies..."
	@npm run install:all

# Run tests
test:
	@echo "Running tests..."
	@npm run test

# Run linting
lint:
	@echo "Running linting..."
	@npm run lint

# Format code
format:
	@echo "Formatting code..."
	@npm run format

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf frontend/build
	@rm -rf frontend/.svelte-kit
	@rm -rf frontend/node_modules/.vite

# Run Storybook
storybook:
	@echo "Starting Storybook..."
	@npm run storybook

# Push database schema
db-push:
	@echo "Pushing database schema..."
	@npm run db:push

# Open Drizzle Studio
db-studio:
	@echo "Opening Drizzle Studio..."
	@npm run db:studio 