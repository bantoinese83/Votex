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

```
votex-template/
├── backend/          # Go API server
├── frontend/         # SvelteKit application
├── docker-compose.yml # Database services
└── package.json      # Root scripts & dependencies
```

## 🛠️ Tech Stack

### Backend
![Go](https://img.shields.io/badge/Go-1.24.2+-00ADD8?logo=go&logoColor=white)
![HTTP](https://img.shields.io/badge/HTTP-Server-lightgrey)
- **Go 1.24.2** - High-performance server language
- **HTTP Server** - Built-in Go HTTP package
- **Port**: 8080

### Frontend
![SvelteKit](https://img.shields.io/badge/SvelteKit-2.22.0-FF3E00?logo=svelte&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?logo=typescript&logoColor=white)
![TailwindCSS](https://img.shields.io/badge/TailwindCSS-4.0-38B2AC?logo=tailwind-css&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-7.0-646CFF?logo=vite&logoColor=white)
- **SvelteKit 2.22.0** - Full-stack web framework
- **TypeScript** - Type-safe development
- **TailwindCSS 4.0** - Utility-first CSS framework
- **Vite 7.0** - Fast build tool and dev server
- **Port**: 5173

### Database & ORM
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?logo=postgresql&logoColor=white)
![Drizzle](https://img.shields.io/badge/Drizzle-ORM-FF6B6B?logo=drizzle&logoColor=white)
![Neon](https://img.shields.io/badge/Neon-Serverless-00D4AA?logo=neon&logoColor=white)
- **PostgreSQL 15** - Primary database
- **Drizzle ORM** - Type-safe database toolkit
- **Neon** - Serverless PostgreSQL client

### Authentication
![Lucia](https://img.shields.io/badge/Lucia-Auth-FF6B6B?logo=lucia&logoColor=white)
![Argon2](https://img.shields.io/badge/Argon2-Password%20Hashing-lightgrey)
- **Lucia Auth** - Lightweight authentication library
- **Argon2** - Password hashing
- **Session management** - Secure session handling

### Internationalization
![Paraglide](https://img.shields.io/badge/Paraglide-i18n-00D4AA?logo=paraglide&logoColor=white)
![Languages](https://img.shields.io/badge/Languages-EN%20%7C%20ES-blue)
- **Paraglide** - Type-safe internationalization
- **Supported languages**: English, Spanish

### Development Tools
![ESLint](https://img.shields.io/badge/ESLint-Code%20Linting-4B32C3?logo=eslint&logoColor=white)
![Prettier](https://img.shields.io/badge/Prettier-Code%20Formatting-F7B93E?logo=prettier&logoColor=white)
![Vitest](https://img.shields.io/badge/Vitest-Testing-6E9F18?logo=vitest&logoColor=white)
![Playwright](https://img.shields.io/badge/Playwright-E2E%20Testing-2EAD96?logo=playwright&logoColor=white)
![Storybook](https://img.shields.io/badge/Storybook-Component%20Dev-FF4785?logo=storybook&logoColor=white)
- **ESLint** - Code linting
- **Prettier** - Code formatting
- **Vitest** - Unit & component testing
- **Playwright** - End-to-end testing
- **Storybook** - Component development & documentation

### Infrastructure
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Caching-DC382D?logo=redis&logoColor=white)
![Concurrently](https://img.shields.io/badge/Concurrently-Multi%20Process-lightgrey)
- **Docker Compose** - Local development services
- **Redis** - Caching & session storage
- **Concurrently** - Run multiple services simultaneously

## 📊 Project Status

![GitHub last commit](https://img.shields.io/github/last-commit/bantoinese83/Votex)
![GitHub repo size](https://img.shields.io/github/repo-size/bantoinese83/Votex)
![GitHub language count](https://img.shields.io/github/languages/count/bantoinese83/Votex)
![GitHub top language](https://img.shields.io/github/languages/top/bantoinese83/Votex)
![GitHub issues](https://img.shields.io/github/issues/bantoinese83/Votex)
![GitHub pull requests](https://img.shields.io/github/issues-pr/bantoinese83/Votex)

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
npm run install:all
```

### 3. Environment Setup
Create a `.env` file in the `frontend` directory:
```bash
cd frontend
cp .env.example .env  # If .env.example exists
```

Add your database URL:
```env
DATABASE_URL="postgresql://user:password@localhost:5432/vortexdb"
```

### 4. Start Database Services
```bash
# Start PostgreSQL and Redis
docker-compose up -d
```

### 5. Initialize Database
```bash
# Push database schema
npm run db:push
```

### 6. Start Development Servers
```bash
# Run both frontend and backend
npm run dev
```

Your application will be available at:
- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080

## 📋 Available Commands

### Root Directory Commands

#### Development
```bash
# Run both frontend and backend
npm run dev
make dev

# Run only frontend
npm run dev:frontend
make dev-frontend

# Run only backend
npm run dev:backend
make dev-backend
```

#### Build & Deploy
```bash
# Build frontend for production
npm run build
make build

# Install all dependencies
npm run install:all
make install
```

#### Testing
```bash
# Run all tests
npm run test
make test

# Run linting
npm run lint
make lint

# Format code
npm run format
make format
```

#### Database Operations
```bash
# Push schema changes
npm run db:push
make db-push

# Open Drizzle Studio
npm run db:studio
make db-studio
```

#### Storybook
```bash
# Start Storybook
npm run storybook
make storybook
```

### Frontend Directory Commands
```bash
cd frontend

# Development
npm run dev          # Start dev server
npm run build        # Build for production
npm run preview      # Preview production build

# Testing
npm run test:unit    # Unit tests
npm run test:e2e     # End-to-end tests
npm run test         # All tests

# Code Quality
npm run check        # Type checking
npm run format       # Format code
npm run lint         # Lint code

# Database
npm run db:push      # Push schema
npm run db:migrate   # Run migrations
npm run db:studio    # Open Drizzle Studio
```

## 🗄️ Database Schema

### Users Table
```sql
CREATE TABLE user (
  id TEXT PRIMARY KEY,
  age INTEGER,
  username TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL
);
```

### Sessions Table
```sql
CREATE TABLE session (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES user(id),
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);
```

## 🌐 API Endpoints

### Backend (Go)
- `GET /` - Health check: "Hello from Go backend!"

### Frontend (SvelteKit)
- `GET /` - Home page
- `GET /demo` - Demo overview
- `GET /demo/lucia` - Authentication demo
- `GET /demo/paraglide` - Internationalization demo

## 🔐 Authentication

The project includes a complete authentication system with:

- **User registration and login**
- **Password hashing with Argon2**
- **Session management**
- **Protected routes**

### Demo Authentication
Visit `/demo/lucia` to see the authentication demo in action.

## 🌍 Internationalization

Built-in support for multiple languages:

- **English** (`en`)
- **Spanish** (`es`)

### Adding Translations
Edit the message files:
- `frontend/messages/en.json`
- `frontend/messages/es.json`

### Demo i18n
Visit `/demo/paraglide` to see internationalization in action.

## 🧪 Testing

### Unit Tests
```bash
npm run test:unit
```

### Component Tests
```bash
npm run test:unit -- --run
```

### End-to-End Tests
```bash
npm run test:e2e
```

### All Tests
```bash
npm run test
```

## 📚 Storybook

Component development and documentation:

```bash
npm run storybook
```

Access Storybook at: http://localhost:6006

## 🐳 Docker Services

### PostgreSQL
- **Port**: 5432
- **Database**: vortexdb
- **Username**: user
- **Password**: password

### Redis
- **Port**: 6379
- **Purpose**: Caching and session storage

### Start Services
```bash
docker-compose up -d
```

### Stop Services
```bash
docker-compose down
```

## 🔧 Development Workflow

### 1. Daily Development
```bash
# Start database services
docker-compose up -d

# Start development servers
npm run dev

# In another terminal, run tests
npm run test

# Format code before committing
npm run format
```

### 2. Database Changes
```bash
# After modifying schema.ts
npm run db:push

# View database in browser
npm run db:studio
```

### 3. Adding New Features
```bash
# Create new SvelteKit routes in frontend/src/routes/
# Add API endpoints in backend/main.go
# Update database schema in frontend/src/lib/server/db/schema.ts
```

## 📁 Project Structure

```
votex-template/
├── backend/
│   ├── main.go          # Go server entry point
│   └── go.mod           # Go dependencies
├── frontend/
│   ├── src/
│   │   ├── lib/
│   │   │   └── server/
│   │   │       ├── auth.ts      # Authentication logic
│   │   │       └── db/
│   │   │           ├── index.ts  # Database connection
│   │   │           └── schema.ts # Database schema
│   │   ├── routes/               # SvelteKit routes
│   │   │   ├── demo/             # Demo pages
│   │   │   │   ├── lucia/        # Auth demo
│   │   │   │   └── paraglide/    # i18n demo
│   │   │   └── +page.svelte      # Home page
│   │   └── stories/              # Storybook stories
│   ├── messages/                 # i18n translations
│   │   ├── en.json
│   │   └── es.json
│   ├── drizzle.config.ts         # Database config
│   ├── svelte.config.js          # SvelteKit config
│   └── vite.config.ts            # Vite config
├── docker-compose.yml            # Database services
├── package.json                  # Root scripts
├── Makefile                      # Convenience commands
└── README.md                     # This file
```

## 🚀 Deployment

### Frontend Deployment
```bash
# Build for production
npm run build

# The built files will be in frontend/build/
```

### Backend Deployment
```bash
# Build Go binary
cd backend
go build -o main main.go
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `npm run test`
5. Format code: `npm run format`
6. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/bantoinese83/Votex/issues)
- **Documentation**: This README and inline code comments
- **Community**: Check the demo pages for examples

---

**Happy coding! 🎉**
