package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"virsh-sandbox/internal/ansible"
	"virsh-sandbox/internal/libvirt"
	"virsh-sandbox/internal/rest"
	"virsh-sandbox/internal/store"
	postgresStore "virsh-sandbox/internal/store/postgres"
	"virsh-sandbox/internal/vm"
)

// @title virsh-sandbox API
// @version 0.0.1-beta
// @description API for managing virtual machine sandboxes using libvirt
// @BasePath /

// @tag.name Sandbox
// @tag.description Sandbox lifecycle management - create, start, run commands, snapshot, and destroy sandboxes

// @tag.name VMs
// @tag.description Virtual machine listing and information

// @tag.name Ansible
// @tag.description Ansible playbook job management

// @tag.name Health
// @tag.description Health check endpoints
func main() {
	// Context with OS signal cancellation
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Logging setup
	logger := setupLogger()
	slog.SetDefault(logger)

	// Read configuration from environment
	apiAddr := getenv("API_HTTP_ADDR", ":8080")
	dbURL := getenv("DATABASE_URL", "postgresql://virsh_sandbox:virsh_sandbox@postgres:5432/virsh_sandbox")
	network := getenv("LIBVIRT_NETWORK", "default")
	libvirtURI := getenv("LIBVIRT_URI", "qemu:///system")

	defaultVCPUs := atoiDefault(getenv("DEFAULT_VCPUS", "2"), 2)
	defaultMemMB := atoiDefault(getenv("DEFAULT_MEMORY_MB", "2048"), 2048)
	cmdTimeout := durationFromSecondsEnv("COMMAND_TIMEOUT_SEC", 600)              // 10m default
	ipDiscoveryTimeout := durationFromSecondsEnv("IP_DISCOVERY_TIMEOUT_SEC", 120) // 2m default

	// Ansible configuration
	ansibleInventoryPath := getenv("ANSIBLE_INVENTORY_PATH", "/ansible/inventory")
	ansibleImage := getenv("ANSIBLE_IMAGE", "ansible-sandbox")
	ansiblePlaybooks := strings.Split(getenv("ANSIBLE_ALLOWED_PLAYBOOKS", "ping.yml"), ",")

	logger.Info("starting virsh-sandbox API",
		"addr", apiAddr,
		"db", dbURL,
		"network", network,
		"default_vcpus", defaultVCPUs,
		"default_memory_mb", defaultMemMB,
		"command_timeout", cmdTimeout.String(),
		"ip_discovery_timeout", ipDiscoveryTimeout.String(),
	)

	st, err := postgresStore.New(ctx, store.Config{
		DatabaseURL:     dbURL,
		MaxOpenConns:    16,
		MaxIdleConns:    8,
		ConnMaxLifetime: time.Hour,
		AutoMigrate:     true,
		ReadOnly:        false,
	})
	if err != nil {
		logger.Error("failed to initialize store", "error", err)
		os.Exit(1)
	}
	defer func() {
		if cerr := st.Close(); cerr != nil {
			logger.Error("failed to close store", "error", cerr)
		}
	}()

	// Initialize libvirt manager from environment
	lvMgr := libvirt.NewFromEnv()

	// Initialize domain manager for direct libvirt queries
	domainMgr := libvirt.NewDomainManager(libvirtURI)

	// Initialize VM service
	vmSvc := vm.NewService(lvMgr, st, vm.Config{
		Network:            network,
		DefaultVCPUs:       defaultVCPUs,
		DefaultMemoryMB:    defaultMemMB,
		CommandTimeout:     cmdTimeout,
		IPDiscoveryTimeout: ipDiscoveryTimeout,
	})

	// Initialize Ansible runner
	ansibleRunner := ansible.NewRunner(ansibleInventoryPath, ansibleImage, ansiblePlaybooks)

	// REST server setup
	restSrv := rest.NewServer(vmSvc, domainMgr, ansibleRunner)

	// Build http.Server so we can gracefully shutdown
	httpSrv := &http.Server{
		Addr:              apiAddr,
		Handler:           restSrv.Router, // use the chi router directly for graceful shutdowns
		ReadHeaderTimeout: 15 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Start HTTP server
	serverErrCh := make(chan error, 1)
	go func() {
		logger.Info("http server listening", "addr", apiAddr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
		}
	}()

	// Wait for signal or server error
	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-serverErrCh:
		logger.Error("server error", "error", err)
	}

	// Attempt graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		logger.Error("http server graceful shutdown failed", "error", err)
		_ = httpSrv.Close()
	} else {
		logger.Info("http server shut down gracefully")
	}
}

// setupLogger configures slog with level and format from environment.
func setupLogger() *slog.Logger {
	level := slog.LevelInfo
	switch strings.ToLower(getenv("LOG_LEVEL", "info")) {
	case "debug":
		level = slog.LevelDebug
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	jsonFmt := strings.ToLower(getenv("LOG_FORMAT", "text")) == "json"

	var handler slog.Handler
	if jsonFmt {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}
	return slog.New(handler)
}

// getenv returns the value of the environment variable k or def if not set.
func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// atoiDefault parses s as int, returning def if empty or invalid.
func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}

// durationFromSecondsEnv reads an environment variable name as seconds and returns a duration.
// If missing or invalid, returns the defaultSeconds value.
func durationFromSecondsEnv(envName string, defaultSeconds int) time.Duration {
	raw := os.Getenv(envName)
	if raw == "" {
		return time.Duration(defaultSeconds) * time.Second
	}
	// Support plain int seconds or Golang duration format
	if d, err := time.ParseDuration(raw); err == nil {
		return d
	}
	sec, err := strconv.Atoi(raw)
	if err != nil {
		log.Printf("invalid %s=%q, falling back to default %ds", envName, raw, defaultSeconds)
		return time.Duration(defaultSeconds) * time.Second
	}
	return time.Duration(sec) * time.Second
}
