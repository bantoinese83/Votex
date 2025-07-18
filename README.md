# Votex Template 🚀

[![Go Version](https://img.shields.io/badge/Go-1.24.2+-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![Node.js Version](https://img.shields.io/badge/Node.js-18+-339933?logo=node.js&logoColor=white)](https://nodejs.org/)
[![SvelteKit](https://img.shields.io/badge/SvelteKit-2.22.0-FF3E00?logo=svelte&logoColor=white)](https://kit.svelte.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Required-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

A modern, full-stack web application template built with **Go** backend and **SvelteKit** frontend, featuring authentication, internationalization, database management, and comprehensive development tools.

## 🏗️ Architecture

This template uses a clean, modular architecture to promote separation of concerns and scalability.

- **Backend:** The Go backend is structured using a layered architecture, separating the API handlers, business logic (services), and data access (stores).
- **Frontend:** The SvelteKit frontend is structured to separate UI components, features, API communication, and state management.
- **API:** The frontend and backend communicate via a REST API defined by the `openapi.yaml` specification.

```
votex-template/
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── api/
│   │   ├── service/
│   │   └── store/
│   └── pkg/router/
├── frontend/
│   ├── src/
│   │   └── lib/
│   │       ├── api/
│   │       ├── components/
│   │       ├── features/
│   │       ├── stores/
│   │       └── types/
│   └── ...
├── docker-compose.yml
├── openapi.yaml
���── package.json
```

## 🛠️ Tech Stack

### Backend
![Go](https://img.shields.io/badge/Go-1.24.2+-00ADD8?logo=go&logoColor=white)
- **Go 1.24.2** - High-performance server language
- **Chi** - Lightweight, idiomatic router
- **sqlx** - Extensions to Go's standard database/sql library

### Frontend
![SvelteKit](https://img.shields.io/badge/SvelteKit-2.22.0-FF3E00?logo=svelte&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?logo=typescript&logoColor=white)
![TailwindCSS](https://img.shields.io/badge/TailwindCSS-4.0-38B2AC?logo=tailwind-css&logoColor=white)
- **SvelteKit 2.22.0** - Full-stack web framework
- **TypeScript** - Type-safe development
- **TailwindCSS 4.0** - Utility-first CSS framework

### Database & Cache
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Caching-DC382D?logo=redis&logoColor=white)
- **PostgreSQL 15** - Primary database
- **Redis** - Caching & session storage

## 🚀 Quick Start

### Prerequisites
- **Node.js** 18+
- **Go** 1.24+
- **Docker** & **Docker Compose**
- **Git**

### 1. Clone & Setup
```bash
git clone https://github.com/bantoinese83/Votex.git
cd Votex
```

### 2. Install Dependencies
```bash
# Install all dependencies (root + frontend)
make install
```

### 3. Start Database Services
```bash
# Start PostgreSQL and Redis
docker-compose up -d
```

### 4. Start Development Servers
```bash
# Run both frontend and backend
make dev
```

Your application will be available at:
- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080

## 📋 Available Commands

See the `Makefile` for a full list of available commands.

- `make dev`: Run both frontend and backend
- `make dev-frontend`: Run only frontend
- `make dev-backend`: Run only backend
- `make install`: Install all dependencies
- `make test`: Run frontend tests
- `make lint`: Run frontend linter
- `make format`: Format frontend code

## 🌐 API Endpoints

See the `openapi.yaml` file for a full list of API endpoints.

- `GET /` - Health check

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Format code: `make format`
6. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details

