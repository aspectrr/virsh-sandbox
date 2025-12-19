// Package audit provides append-only audit logging for all tool operations.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"tmux-client/internal/config"
	"tmux-client/internal/types"
)

// Logger provides thread-safe, append-only audit logging.
type Logger struct {
	mu       sync.Mutex
	file     *os.File
	encoder  *json.Encoder
	config   *config.AuditConfig
	entries  []types.AuditEntry // In-memory buffer for queries
	maxInMem int
}

// NewLogger creates a new audit logger.
func NewLogger(cfg *config.AuditConfig) (*Logger, error) {
	if !cfg.Enabled {
		return &Logger{
			config:   cfg,
			entries:  make([]types.AuditEntry, 0),
			maxInMem: 10000,
		}, nil
	}

	// Ensure directory exists
	dir := filepath.Dir(cfg.LogFile)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return nil, fmt.Errorf("failed to create audit log directory: %w", err)
		}
	}

	// Open file in append-only mode
	file, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open audit log file: %w", err)
	}

	logger := &Logger{
		file:     file,
		encoder:  json.NewEncoder(file),
		config:   cfg,
		entries:  make([]types.AuditEntry, 0),
		maxInMem: 10000,
	}

	return logger, nil
}

// Log writes an audit entry to the log file and in-memory buffer.
func (l *Logger) Log(entry types.AuditEntry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Sanitize sensitive fields if configured
	if l.config != nil && len(l.config.SensitiveFields) > 0 {
		entry = l.sanitizeEntry(entry)
	}

	// Add to in-memory buffer
	l.entries = append(l.entries, entry)
	if len(l.entries) > l.maxInMem {
		// Remove oldest entries
		l.entries = l.entries[len(l.entries)-l.maxInMem:]
	}

	// Write to file if enabled
	if l.file != nil {
		if err := l.encoder.Encode(entry); err != nil {
			return fmt.Errorf("failed to write audit entry: %w", err)
		}
		// Sync to ensure durability
		if err := l.file.Sync(); err != nil {
			return fmt.Errorf("failed to sync audit log: %w", err)
		}
	}

	return nil
}

// LogToolCall is a convenience method for logging tool calls.
func (l *Logger) LogToolCall(
	requestID string,
	tool string,
	action string,
	arguments any,
	result any,
	apiErr *types.APIError,
	duration time.Duration,
	clientIP string,
	userAgent string,
) error {
	entry := types.AuditEntry{
		Timestamp:  time.Now().UTC(),
		RequestID:  requestID,
		Tool:       tool,
		Action:     action,
		DurationMs: duration.Milliseconds(),
		ClientIP:   clientIP,
		UserAgent:  userAgent,
		Error:      apiErr,
	}

	// Marshal arguments if logging is enabled
	if l.config == nil || l.config.LogArguments {
		if arguments != nil {
			data, err := json.Marshal(arguments)
			if err == nil {
				entry.Arguments = data
			}
		}
	}

	// Marshal result if logging is enabled
	if l.config == nil || l.config.LogResults {
		if result != nil {
			data, err := json.Marshal(result)
			if err == nil {
				entry.Result = data
			}
		}
	}

	return l.Log(entry)
}

// Query retrieves audit entries matching the given criteria.
func (l *Logger) Query(query types.AuditQuery) ([]types.AuditEntry, int, bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var results []types.AuditEntry
	limit := query.Limit
	if limit <= 0 {
		limit = 100
	}

	// Filter in-memory entries (reverse order for most recent first)
	for i := len(l.entries) - 1; i >= 0; i-- {
		entry := l.entries[i]

		// Apply filters
		if query.Tool != "" && entry.Tool != query.Tool {
			continue
		}
		if query.Action != "" && entry.Action != query.Action {
			continue
		}
		if query.RequestID != "" && entry.RequestID != query.RequestID {
			continue
		}
		if query.Since != nil && entry.Timestamp.Before(*query.Since) {
			continue
		}
		if query.Until != nil && entry.Timestamp.After(*query.Until) {
			continue
		}

		results = append(results, entry)

		if len(results) >= limit+1 {
			break
		}
	}

	hasMore := len(results) > limit
	if hasMore {
		results = results[:limit]
	}

	return results, len(results), hasMore, nil
}

// GetEntry retrieves a specific audit entry by request ID.
func (l *Logger) GetEntry(requestID string) (*types.AuditEntry, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i := len(l.entries) - 1; i >= 0; i-- {
		if l.entries[i].RequestID == requestID {
			entry := l.entries[i]
			return &entry, nil
		}
	}

	return nil, nil
}

// sanitizeEntry redacts sensitive fields from the audit entry.
func (l *Logger) sanitizeEntry(entry types.AuditEntry) types.AuditEntry {
	if len(entry.Arguments) == 0 {
		return entry
	}

	var args map[string]interface{}
	if err := json.Unmarshal(entry.Arguments, &args); err != nil {
		return entry
	}

	sanitized := l.sanitizeMap(args)
	data, err := json.Marshal(sanitized)
	if err != nil {
		return entry
	}

	entry.Arguments = data
	return entry
}

// sanitizeMap recursively sanitizes a map, redacting sensitive fields.
func (l *Logger) sanitizeMap(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range m {
		if l.isSensitiveField(k) {
			result[k] = "[REDACTED]"
			continue
		}

		switch val := v.(type) {
		case map[string]interface{}:
			result[k] = l.sanitizeMap(val)
		case []interface{}:
			result[k] = l.sanitizeSlice(val)
		default:
			result[k] = v
		}
	}

	return result
}

// sanitizeSlice recursively sanitizes a slice.
func (l *Logger) sanitizeSlice(s []interface{}) []interface{} {
	result := make([]interface{}, len(s))

	for i, v := range s {
		switch val := v.(type) {
		case map[string]interface{}:
			result[i] = l.sanitizeMap(val)
		case []interface{}:
			result[i] = l.sanitizeSlice(val)
		default:
			result[i] = v
		}
	}

	return result
}

// isSensitiveField checks if a field name is in the sensitive fields list.
func (l *Logger) isSensitiveField(field string) bool {
	for _, sensitive := range l.config.SensitiveFields {
		if field == sensitive {
			return true
		}
	}
	return false
}

// Rotate rotates the log file if it exceeds the configured size.
func (l *Logger) Rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil || l.config.RotateSize <= 0 {
		return nil
	}

	info, err := l.file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat audit log: %w", err)
	}

	if info.Size() < l.config.RotateSize {
		return nil
	}

	// Close current file
	if err := l.file.Close(); err != nil {
		return fmt.Errorf("failed to close audit log: %w", err)
	}

	// Rotate existing files
	for i := l.config.RotateCount - 1; i >= 0; i-- {
		oldPath := l.rotatedPath(i)
		newPath := l.rotatedPath(i + 1)

		if i == l.config.RotateCount-1 {
			// Remove oldest file
			os.Remove(oldPath)
		} else {
			// Rename to next number
			if _, err := os.Stat(oldPath); err == nil {
				os.Rename(oldPath, newPath)
			}
		}
	}

	// Rename current file
	if err := os.Rename(l.config.LogFile, l.rotatedPath(0)); err != nil {
		return fmt.Errorf("failed to rotate audit log: %w", err)
	}

	// Open new file
	file, err := os.OpenFile(l.config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open new audit log: %w", err)
	}

	l.file = file
	l.encoder = json.NewEncoder(file)

	return nil
}

// rotatedPath returns the path for a rotated log file.
func (l *Logger) rotatedPath(index int) string {
	return fmt.Sprintf("%s.%d", l.config.LogFile, index)
}

// Close closes the audit logger.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		if err := l.file.Sync(); err != nil {
			return fmt.Errorf("failed to sync audit log: %w", err)
		}
		if err := l.file.Close(); err != nil {
			return fmt.Errorf("failed to close audit log: %w", err)
		}
		l.file = nil
	}

	return nil
}

// Stats returns statistics about the audit log.
func (l *Logger) Stats() AuditStats {
	l.mu.Lock()
	defer l.mu.Unlock()

	stats := AuditStats{
		EntriesInMemory: len(l.entries),
		Enabled:         l.config.Enabled,
	}

	if l.file != nil {
		if info, err := l.file.Stat(); err == nil {
			stats.FileSize = info.Size()
			stats.FilePath = l.config.LogFile
		}
	}

	// Calculate entries by tool
	stats.EntriesByTool = make(map[string]int)
	for _, entry := range l.entries {
		stats.EntriesByTool[entry.Tool]++
	}

	// Find time range
	if len(l.entries) > 0 {
		stats.OldestEntry = l.entries[0].Timestamp
		stats.NewestEntry = l.entries[len(l.entries)-1].Timestamp
	}

	return stats
}

// AuditStats contains statistics about the audit log.
type AuditStats struct {
	Enabled         bool           `json:"enabled"`
	EntriesInMemory int            `json:"entries_in_memory"`
	FileSize        int64          `json:"file_size"`
	FilePath        string         `json:"file_path"`
	OldestEntry     time.Time      `json:"oldest_entry"`
	NewestEntry     time.Time      `json:"newest_entry"`
	EntriesByTool   map[string]int `json:"entries_by_tool"`
}

// Export exports all in-memory entries to a writer.
func (l *Logger) Export(w io.Writer) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	return encoder.Encode(l.entries)
}

// Replay reads and parses all entries from the log file.
func (l *Logger) Replay() ([]types.AuditEntry, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.config.LogFile == "" {
		return l.entries, nil
	}

	file, err := os.Open(l.config.LogFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []types.AuditEntry{}, nil
		}
		return nil, fmt.Errorf("failed to open audit log for replay: %w", err)
	}
	defer file.Close()

	var entries []types.AuditEntry
	decoder := json.NewDecoder(file)

	for {
		var entry types.AuditEntry
		if err := decoder.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			// Skip malformed entries
			continue
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// NullLogger is a no-op logger for testing or when auditing is disabled.
type NullLogger struct{}

// NewNullLogger creates a new null logger.
func NewNullLogger() *NullLogger {
	return &NullLogger{}
}

// Log is a no-op.
func (l *NullLogger) Log(_ types.AuditEntry) error {
	return nil
}

// LogToolCall is a no-op.
func (l *NullLogger) LogToolCall(_, _, _ string, _, _ interface{}, _ *types.APIError, _ time.Duration, _, _ string) error {
	return nil
}

// Close is a no-op.
func (l *NullLogger) Close() error {
	return nil
}

// AuditLogger interface for dependency injection.
type AuditLogger interface {
	Log(entry types.AuditEntry) error
	LogToolCall(requestID, tool, action string, arguments, result interface{}, apiErr *types.APIError, duration time.Duration, clientIP, userAgent string) error
	Close() error
}

// Ensure Logger implements AuditLogger
var _ AuditLogger = (*Logger)(nil)
var _ AuditLogger = (*NullLogger)(nil)
