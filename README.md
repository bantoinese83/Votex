# VORTEX Stack Starter Template

This repository contains a starter template for the VORTEX stack.

### VORTEX Stack Components
| Layer           | Technology                  | Rationale/Benefit                               |
| --------------- | --------------------------- | ----------------------------------------------- |
| **Frontend**    | SvelteKit                   | Fastest loads, minimal runtime                  |
| **Styling**     | Tailwind CSS                | Efficient, scalable CSS                         |
| **Backend**     | Go (+ HTMX)                 | Maximum API speed and concurrency               |
| **Database**    | PostgreSQL                  | Reliable, highly optimized queries              |
| **Cache**       | Redis                       | Near-instant data, pub/sub capability           |
| **Serverless/Edge** | Vercel Edge/Cloudflare      | Compute close to users, lowest latency          |
| **Option**      | PlanetScale (MySQL)         | Serverless, horizontally scalable DB            |

---

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Go](https://go.dev/doc/install)
- [Node.js](https://nodejs.org/en/download/)

### 1. Backend (Go, PostgreSQL, Redis)

The backend consists of a Go server and Docker containers for PostgreSQL and Redis.

1.  **Start services:**
    ```bash
    docker-compose up -d
    ```

2.  **Run the Go server:**
    ```bash
    cd backend
    go run main.go
    ```
    The backend server will be running on `http://localhost:8080`.

### 2. Frontend (SvelteKit + Tailwind CSS)

You need to create the SvelteKit project inside the `frontend` directory.

1.  **Create a new SvelteKit app:**
    ```bash
    # Make sure you are in the root of the votex-template directory
    npm create svelte@latest frontend
    ```
    - When prompted, choose:
        - **App template:** Skeleton project
        - **Type checking with TypeScript:** Yes
        - **ESLint, Prettier, Playwright, Vitest:** Add them as you see fit.
    
2. **Install Tailwind CSS**
   Follow the official SvelteKit installation guide for Tailwind CSS:
   [https://tailwindcss.com/docs/guides/sveltekit](https://tailwindcss.com/docs/guides/sveltekit)

3.  **Run the frontend dev server:**
    ```bash
    cd frontend
    npm install
    npm run dev
    ```
    The frontend will be available at `http://localhost:5173`.
