// Package main provides the entry point for the tmux agent API server.
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"tmux-client/internal/api"
	"tmux-client/internal/audit"
	"tmux-client/internal/config"
	"tmux-client/internal/tools/command"
	"tmux-client/internal/tools/file"
	"tmux-client/internal/tools/human"
	"tmux-client/internal/tools/plan"
	"tmux-client/internal/tools/tmux"
)

// Version is set at build time.
var Version = "0.0.1-beta"

// @title tmux-client API
// @version 0.0.1-beta
// @description API for managing tmux sessions and windows
// @host localhost:8081
// @BasePath /
func main() {
	// Parse command line flags
	configPath := flag.String("config", "", "Path to configuration file")
	showVersion := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("tmux-agent %s\n", Version)
		os.Exit(0)
	}

	// Initialize logger
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Str("service", "tmux-agent").Logger()

	logger.Info().Str("version", Version).Msg("Starting tmux agent API server")

	// Load configuration
	cfg, err := config.LoadOrDefault(*configPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize audit logger
	auditLogger, err := audit.NewLogger(&cfg.Audit)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize audit logger")
	}
	defer func() {
		if err := auditLogger.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close audit logger")
		}
	}()

	// Initialize tools
	tmuxTool, err := tmux.NewTool(&cfg.Tmux)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to initialize tmux tool (tmux may not be installed)")
		tmuxTool = nil
	}

	fileTool, err := file.NewTool(&cfg.File)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize file tool")
	}

	commandTool, err := command.NewTool(&cfg.Command)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize command tool")
	}

	humanTool, err := human.NewTool(&cfg.Human)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize human approval tool")
	}
	defer func() {
		if err := humanTool.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close human tool")
		}
	}()

	planTool, err := plan.NewTool(&cfg.Plan)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize plan tool")
	}
	defer func() {
		if err := planTool.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close plan tool")
		}
	}()

	// Initialize middleware
	middleware := api.NewMiddleware(&cfg.Security, logger)

	// Initialize API handler
	handler := api.NewHandler(
		cfg,
		logger,
		auditLogger,
		tmuxTool,
		fileTool,
		commandTool,
		humanTool,
		planTool,
		Version,
	)

	// Setup router
	r := chi.NewRouter()

	// Apply middleware stack
	r.Use(middleware.Recovery)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.SecurityHeaders)
	// r.Use(middleware.CORS)
	r.Use(middleware.MaxBodySize)
	// r.Use(middleware.IPWhitelist)
	r.Use(middleware.RateLimit)
	// r.Use(middleware.Auth)

	// Register API routes
	handler.RegisterRoutes(r)

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	serverErrors := make(chan error, 1)
	go func() {
		logger.Info().
			Str("address", cfg.Server.Address()).
			Bool("tls", cfg.Server.EnableTLS).
			Msg("HTTP server listening")

		var err error
		if cfg.Server.EnableTLS {
			err = server.ListenAndServeTLS(cfg.Server.TLSCertFile, cfg.Server.TLSKeyFile)
		} else {
			err = server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Wait for interrupt signal or server error
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Fatal().Err(err).Msg("Server error")

	case sig := <-shutdown:
		logger.Info().Str("signal", sig.String()).Msg("Shutdown signal received")

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer cancel()

		// Attempt graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			logger.Error().Err(err).Msg("Graceful shutdown failed, forcing close")
			server.Close()
		}

		logger.Info().Msg("Server stopped")
	}
}
