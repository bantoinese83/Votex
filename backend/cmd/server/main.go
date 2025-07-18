package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/user/votex-template/backend/internal/api"
	"github.com/user/votex-template/backend/internal/config"
	"github.com/user/votex-template/backend/internal/middleware"
	"github.com/user/votex-template/backend/internal/service"
	"github.com/user/votex-template/backend/internal/store"
	"github.com/user/votex-template/backend/pkg/logger"
	"github.com/user/votex-template/backend/pkg/router"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup structured logging
	logger.Setup(cfg.LogLevel, cfg.IsDevelopment())
	slog.Info("Starting Votex backend",
		"environment", cfg.Environment,
		"port", cfg.Port,
		"log_level", cfg.LogLevel,
	)

	// Run database migrations
	if err := runMigrations(cfg.DBURL); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Database connection
	db, err := sqlx.Connect("postgres", cfg.DBURL)
	if err != nil {
		slog.Error("Could not connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("Database connected successfully")

	// Initialize store and services
	store := store.New(db)
	authService := service.NewAuthService(store, cfg)

	// Initialize handlers
	authHandler := api.NewAuthHandler(authService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg)

		// Initialize router with middleware
	r := router.New()
	
	// Apply CORS middleware
	r.Use(middleware.CORS(cfg))

	// Health check endpoint
	r.Get("/health", http.HandlerFunc(api.HandleHealthCheck))

	// Auth endpoints
	r.Post("/api/auth/register", http.HandlerFunc(authHandler.Register))
	r.Post("/api/auth/login", http.HandlerFunc(authHandler.Login))
	r.Get("/api/auth/profile", authMiddleware.Authenticate(http.HandlerFunc(authHandler.Profile)).(http.HandlerFunc))

	// Start server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	slog.Info("Go backend server starting", "port", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Could not start server", "error", err)
		os.Exit(1)
	}
}

func runMigrations(dbURL string) error {
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
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
