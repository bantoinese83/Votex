package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/user/votex-template/backend/internal/api"
	"github.com/user/votex-template/backend/internal/config"
	"github.com/user/votex-template/backend/internal/middleware"
	"github.com/user/votex-template/backend/internal/service"
	"github.com/user/votex-template/backend/internal/store"
	"github.com/user/votex-template/backend/pkg/router"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set up logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Connect to database
	db, isSQLite, err := store.ConnectDatabase(cfg.DBURL, cfg.SQLitePath)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(cfg, isSQLite); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Initialize store
	storeInstance := store.New(db, isSQLite)

	// Initialize services
	authService := service.NewAuthService(storeInstance, cfg)

	// Initialize handlers
	authHandler := api.NewAuthHandler(authService)
	userHandler := api.NewUserHandler(authService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg)
	rateLimiter := middleware.NewRateLimiter(cfg)

	// Initialize router with middleware
	r := router.New()

	// Apply global middleware
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.CORS(cfg))
	r.Use(rateLimiter.RateLimit)

	// Health check endpoint
	r.Get("/health", http.HandlerFunc(api.HandleHealthCheck))

	// API documentation
	r.Get("/api/docs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Vortex API Documentation</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.0.0/swagger-ui.css" />
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.0.0/swagger-ui-bundle.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: '/openapi.yaml',
                dom_id: '#swagger-ui',
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIBundle.SwaggerUIStandalonePreset
                ],
                layout: "BaseLayout"
            });
        }
    </script>
</body>
</html>
		`))
	}))

	// Serve OpenAPI spec
	r.Get("/openapi.yaml", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		http.ServeFile(w, r, "openapi.yaml")
	}))

	// Auth endpoints
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", http.HandlerFunc(authHandler.Register))
		r.Post("/login", http.HandlerFunc(authHandler.Login))
		r.Post("/password-reset", http.HandlerFunc(authHandler.RequestPasswordReset))
		r.Post("/password-reset/{token}", http.HandlerFunc(authHandler.ResetPassword))

		// Protected auth endpoints
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Get("/profile", http.HandlerFunc(authHandler.Profile))
			r.Put("/profile", http.HandlerFunc(authHandler.UpdateProfile))
			r.Delete("/account", http.HandlerFunc(authHandler.DeleteAccount))
		})
	})

	// User management endpoints
	r.Route("/api/users", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Get("/", http.HandlerFunc(userHandler.ListUsers))
		r.Get("/{id}", http.HandlerFunc(userHandler.GetUser))
		r.Put("/{id}", http.HandlerFunc(userHandler.UpdateUser))
		r.Delete("/{id}", http.HandlerFunc(userHandler.DeleteUser))
	})

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	slog.Info("Go backend server starting",
		"port", cfg.Port,
		"environment", cfg.Environment,
		"database", cfg.DBType,
		"rate_limit", fmt.Sprintf("%d req/min", cfg.RateLimitRequests),
	)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Could not start server", "error", err)
		os.Exit(1)
	}
}

func runMigrations(cfg *config.Config, isSQLite bool) error {
	var sourceURL, databaseURL string

	if isSQLite {
		sourceURL = "file://migrations/sqlite"
		databaseURL = fmt.Sprintf("sqlite3://%s", cfg.SQLitePath)
	} else {
		sourceURL = "file://migrations/postgres"
		databaseURL = cfg.DBURL
	}

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	slog.Info("Database migrations applied successfully")
	return nil
}
