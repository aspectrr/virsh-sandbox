// Package api provides HTTP middleware for the tmux agent API.
package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"tmux-client/internal/config"
	"tmux-client/internal/types"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const (
	// RequestIDKey is the context key for the request ID.
	RequestIDKey contextKey = "request_id"
	// StartTimeKey is the context key for the request start time.
	StartTimeKey contextKey = "start_time"
)

// Middleware provides HTTP middleware functions.
type Middleware struct {
	config      *config.SecurityConfig
	logger      zerolog.Logger
	rateLimiter *RateLimiter
}

// NewMiddleware creates a new middleware instance.
func NewMiddleware(cfg *config.SecurityConfig, logger zerolog.Logger) *Middleware {
	var limiter *RateLimiter
	if cfg.RateLimitPerMin > 0 {
		limiter = NewRateLimiter(cfg.RateLimitPerMin, time.Minute)
	}

	return &Middleware{
		config:      cfg,
		logger:      logger,
		rateLimiter: limiter,
	}
}

// RequestID adds a unique request ID to each request.
func (m *Middleware) RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add to response header
		w.Header().Set("X-Request-ID", requestID)

		// Add to context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		ctx = context.WithValue(ctx, StartTimeKey, time.Now())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logger logs HTTP requests.
func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process request
		next.ServeHTTP(wrapped, r)

		// Log after request completes
		duration := time.Since(start)

		requestID := ""
		if id, ok := r.Context().Value(RequestIDKey).(string); ok {
			requestID = id
		}

		m.logger.Info().
			Str("request_id", requestID).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Int("status", wrapped.statusCode).
			Dur("duration", duration).
			Int("bytes", wrapped.bytesWritten).
			Msg("HTTP request")
	})
}

// Auth validates API key authentication.
func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth if not enabled
		if !m.config.EnableAuth {
			next.ServeHTTP(w, r)
			return
		}

		// Get API key from header
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			// Also check Authorization header with Bearer token
			auth := r.Header.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				apiKey = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		if apiKey == "" {
			m.writeError(w, r, http.StatusUnauthorized, types.ErrCodeForbidden, "API key required", "")
			return
		}

		// Validate API key
		valid := false
		for _, key := range m.config.APIKeys {
			if apiKey == key {
				valid = true
				break
			}
		}

		if !valid {
			m.writeError(w, r, http.StatusUnauthorized, types.ErrCodeForbidden, "Invalid API key", "")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// IPWhitelist restricts access to allowed IPs.
func (m *Middleware) IPWhitelist(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip if no IPs configured
		if len(m.config.AllowedIPs) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// Get client IP
		clientIP := getClientIP(r)

		// Check if IP is allowed
		allowed := false
		for _, ip := range m.config.AllowedIPs {
			if clientIP == ip {
				allowed = true
				break
			}
			// Check if it's a wildcard match
			if ip == "*" {
				allowed = true
				break
			}
		}

		if !allowed {
			m.logger.Warn().
				Str("client_ip", clientIP).
				Msg("Access denied: IP not in whitelist")
			m.writeError(w, r, http.StatusForbidden, types.ErrCodeForbidden, "Access denied", "IP not allowed")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimit applies rate limiting per client IP.
func (m *Middleware) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip if rate limiting is disabled
		if m.rateLimiter == nil {
			next.ServeHTTP(w, r)
			return
		}

		clientIP := getClientIP(r)

		if !m.rateLimiter.Allow(clientIP) {
			m.logger.Warn().
				Str("client_ip", clientIP).
				Msg("Rate limit exceeded")

			w.Header().Set("Retry-After", "60")
			m.writeError(w, r, http.StatusTooManyRequests, "RATE_LIMITED", "Rate limit exceeded", "Try again later")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// MaxBodySize limits the request body size.
func (m *Middleware) MaxBodySize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.config.MaxRequestSize > 0 {
			r.Body = http.MaxBytesReader(w, r.Body, m.config.MaxRequestSize)
		}
		next.ServeHTTP(w, r)
	})
}

// Recovery recovers from panics and returns a 500 error.
func (m *Middleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := ""
				if id, ok := r.Context().Value(RequestIDKey).(string); ok {
					requestID = id
				}

				m.logger.Error().
					Str("request_id", requestID).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Interface("error", err).
					Msg("Panic recovered")

				m.writeError(w, r, http.StatusInternalServerError, types.ErrCodeInternal, "Internal server error", "")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// CORS adds Cross-Origin Resource Sharing headers.
func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-API-Key, X-Request-ID")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ContentType sets the default content type for API responses.
func (m *Middleware) ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// SecurityHeaders adds security-related HTTP headers.
func (m *Middleware) SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// writeError writes an error response.
func (m *Middleware) writeError(w http.ResponseWriter, r *http.Request, status int, code, message, details string) {
	requestID := ""
	if id, ok := r.Context().Value(RequestIDKey).(string); ok {
		requestID = id
	}

	resp := types.APIResponse{
		Success: false,
		Error: &types.APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now().UTC(),
		RequestID: requestID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"success":false,"error":{"code":"%s","message":"%s","details":"%s"},"timestamp":"%s","request_id":"%s"}`,
		resp.Error.Code, resp.Error.Message, resp.Error.Details, resp.Timestamp.Format(time.RFC3339), resp.RequestID)
}

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

// getClientIP extracts the client IP from the request.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP in the chain
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	// Handle IPv6 with brackets
	ip = strings.TrimPrefix(ip, "[")
	ip = strings.TrimSuffix(ip, "]")

	return ip
}

// RateLimiter implements a token bucket rate limiter per client.
type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*rateLimitClient
	limit    int
	window   time.Duration
	cleanupC chan struct{}
}

type rateLimitClient struct {
	tokens    int
	lastReset time.Time
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients:  make(map[string]*rateLimitClient),
		limit:    limit,
		window:   window,
		cleanupC: make(chan struct{}),
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Allow checks if a request from the given client is allowed.
func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	client, exists := rl.clients[clientID]
	if !exists {
		rl.clients[clientID] = &rateLimitClient{
			tokens:    rl.limit - 1, // Use one token for this request
			lastReset: now,
		}
		return true
	}

	// Reset tokens if window has passed
	if now.Sub(client.lastReset) >= rl.window {
		client.tokens = rl.limit - 1
		client.lastReset = now
		return true
	}

	// Check if tokens are available
	if client.tokens > 0 {
		client.tokens--
		return true
	}

	return false
}

// cleanup periodically removes stale client entries.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			cutoff := time.Now().Add(-2 * rl.window)
			for clientID, client := range rl.clients {
				if client.lastReset.Before(cutoff) {
					delete(rl.clients, clientID)
				}
			}
			rl.mu.Unlock()

		case <-rl.cleanupC:
			return
		}
	}
}

// Close stops the rate limiter's cleanup goroutine.
func (rl *RateLimiter) Close() {
	close(rl.cleanupC)
}

// GetRequestID extracts the request ID from the context.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// GetStartTime extracts the request start time from the context.
func GetStartTime(ctx context.Context) time.Time {
	if t, ok := ctx.Value(StartTimeKey).(time.Time); ok {
		return t
	}
	return time.Time{}
}
