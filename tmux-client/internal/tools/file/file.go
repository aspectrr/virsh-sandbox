// Package file provides secure filesystem operations for the agent API.
package file

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"tmux-client/internal/config"
	"tmux-client/internal/types"
)

// Tool provides file operations with safety constraints.
type Tool struct {
	config *config.FileConfig
}

// NewTool creates a new file tool.
func NewTool(cfg *config.FileConfig) (*Tool, error) {
	// Ensure root directory exists if specified
	if cfg.RootDirectory != "" {
		absRoot, err := filepath.Abs(cfg.RootDirectory)
		if err != nil {
			return nil, fmt.Errorf("invalid root directory: %w", err)
		}

		info, err := os.Stat(absRoot)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("root directory does not exist: %s", absRoot)
			}
			return nil, fmt.Errorf("failed to stat root directory: %w", err)
		}

		if !info.IsDir() {
			return nil, fmt.Errorf("root path is not a directory: %s", absRoot)
		}
	}

	// Ensure backup directory exists if specified and backup is enabled
	if cfg.BackupOnEdit && cfg.BackupDirectory != "" {
		if err := os.MkdirAll(cfg.BackupDirectory, 0750); err != nil {
			return nil, fmt.Errorf("failed to create backup directory: %w", err)
		}
	}

	return &Tool{config: cfg}, nil
}

// validatePath checks if a path is allowed and returns the absolute path.
func (t *Tool) validatePath(path string) (string, error) {
	// Clean the path
	path = filepath.Clean(path)

	// Make path absolute
	var absPath string
	var err error

	if t.config.RootDirectory != "" {
		// If root directory is set, paths are relative to it
		if filepath.IsAbs(path) {
			absPath = path
		} else {
			absPath = filepath.Join(t.config.RootDirectory, path)
		}
		absPath, err = filepath.Abs(absPath)
		if err != nil {
			return "", fmt.Errorf("invalid path: %w", err)
		}

		// Ensure path is under root directory
		rootAbs, _ := filepath.Abs(t.config.RootDirectory)
		if !strings.HasPrefix(absPath, rootAbs+string(filepath.Separator)) && absPath != rootAbs {
			return "", fmt.Errorf("path escapes root directory")
		}
	} else {
		absPath, err = filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("invalid path: %w", err)
		}
	}

	// Check denied paths
	for _, denied := range t.config.DeniedPaths {
		deniedAbs, err := filepath.Abs(expandPath(denied))
		if err != nil {
			continue
		}
		if strings.HasPrefix(absPath, deniedAbs) || absPath == deniedAbs {
			return "", fmt.Errorf("access denied: path is in denied list")
		}
	}

	// Check denied extensions
	ext := strings.ToLower(filepath.Ext(absPath))
	for _, deniedExt := range t.config.DeniedExtensions {
		if ext == strings.ToLower(deniedExt) {
			return "", fmt.Errorf("access denied: file extension is not allowed")
		}
	}

	// If allowed paths are specified, check that path is within one of them
	if len(t.config.AllowedPaths) > 0 {
		allowed := false
		for _, allowedPath := range t.config.AllowedPaths {
			allowedAbs, err := filepath.Abs(expandPath(allowedPath))
			if err != nil {
				continue
			}
			if strings.HasPrefix(absPath, allowedAbs) || absPath == allowedAbs {
				allowed = true
				break
			}
		}
		if !allowed {
			return "", fmt.Errorf("access denied: path is not in allowed list")
		}
	}

	return absPath, nil
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	} else if path == "~" {
		home, err := os.UserHomeDir()
		if err == nil {
			path = home
		}
	}
	return path
}

// ReadFile reads a file and returns its content.
func (t *Tool) ReadFile(req types.ReadFileRequest) (*types.ReadFileResponse, error) {
	absPath, err := t.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	// Get file info
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", req.Path)
		}
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("path is a directory, not a file")
	}

	// Check file size
	if info.Size() > t.config.MaxFileSize {
		return nil, fmt.Errorf("file too large: %d bytes (max: %d)", info.Size(), t.config.MaxFileSize)
	}

	// Read file
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024) // 10MB max line size

	lineNum := 0
	totalLines := 0
	fromLine := req.FromLine
	toLine := req.ToLine

	if fromLine <= 0 {
		fromLine = 1
	}

	for scanner.Scan() {
		totalLines++
		lineNum++

		// Apply line range filter
		if lineNum < fromLine {
			continue
		}
		if toLine > 0 && lineNum > toLine {
			break
		}

		// Check max lines limit
		if req.MaxLines > 0 && (lineNum-fromLine+1) > req.MaxLines {
			break
		}

		// Check config max lines
		if t.config.MaxReadLines > 0 && (lineNum-fromLine+1) > t.config.MaxReadLines {
			break
		}

		if content.Len() > 0 {
			content.WriteString("\n")
		}
		content.WriteString(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Count total lines if we didn't read to the end
	if toLine > 0 || (req.MaxLines > 0 && lineNum-fromLine+1 > req.MaxLines) {
		for scanner.Scan() {
			totalLines++
		}
	}

	truncated := false
	if toLine > 0 && totalLines > toLine {
		truncated = true
	}
	if req.MaxLines > 0 && (lineNum-fromLine+1) > req.MaxLines {
		truncated = true
	}

	return &types.ReadFileResponse{
		Path:       req.Path,
		Content:    content.String(),
		TotalLines: totalLines,
		FromLine:   fromLine,
		ToLine:     lineNum,
		Truncated:  truncated,
		Size:       info.Size(),
		Mode:       info.Mode().String(),
		ModTime:    info.ModTime().Format(time.RFC3339),
	}, nil
}

// WriteFile writes content to a file.
func (t *Tool) WriteFile(req types.WriteFileRequest) (*types.WriteFileResponse, error) {
	absPath, err := t.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	// Check if file exists
	info, err := os.Stat(absPath)
	exists := err == nil

	if exists && info.IsDir() {
		return nil, fmt.Errorf("path is a directory")
	}

	if exists && !req.Overwrite {
		if !t.config.AllowOverwrite {
			return nil, fmt.Errorf("file exists and overwrite is disabled in config")
		}
		return nil, fmt.Errorf("file exists and overwrite not requested")
	}

	// Determine file mode
	var mode os.FileMode = 0644
	if req.Mode != "" {
		parsed, err := parseFileMode(req.Mode)
		if err != nil {
			return nil, fmt.Errorf("invalid file mode: %w", err)
		}
		mode = parsed
	}

	// Create parent directory if needed
	if req.CreateDir {
		dir := filepath.Dir(absPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create parent directory: %w", err)
		}
	}

	// Check content size
	if int64(len(req.Content)) > t.config.MaxFileSize {
		return nil, fmt.Errorf("content too large: %d bytes (max: %d)", len(req.Content), t.config.MaxFileSize)
	}

	// Backup if file exists and backup is enabled
	if exists && t.config.BackupOnEdit {
		if err := t.createBackup(absPath); err != nil {
			return nil, fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Write file
	if err := os.WriteFile(absPath, []byte(req.Content), mode); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &types.WriteFileResponse{
		Path:         req.Path,
		Written:      true,
		BytesWritten: int64(len(req.Content)),
		Created:      !exists,
	}, nil
}

// EditFile edits a file using find/replace semantics.
func (t *Tool) EditFile(req types.EditFileRequest) (*types.EditFileResponse, error) {
	absPath, err := t.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	// Check if file exists
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", req.Path)
		}
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("path is a directory, not a file")
	}

	// Check file size
	if info.Size() > t.config.MaxFileSize {
		return nil, fmt.Errorf("file too large: %d bytes (max: %d)", info.Size(), t.config.MaxFileSize)
	}

	// Read current content
	contentBytes, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	content := string(contentBytes)
	contentBefore := content

	// Check if old text exists
	if !strings.Contains(content, req.OldText) {
		return nil, fmt.Errorf("old text not found in file")
	}

	// Perform replacement
	var newContent string
	var replacements int

	if req.All {
		replacements = strings.Count(content, req.OldText)
		newContent = strings.ReplaceAll(content, req.OldText, req.NewText)
	} else {
		replacements = 1
		newContent = strings.Replace(content, req.OldText, req.NewText, 1)
	}

	if newContent == content {
		return &types.EditFileResponse{
			Path:          req.Path,
			Edited:        false,
			Replacements:  0,
			Diff:          "",
			ContentBefore: contentBefore,
			ContentAfter:  content,
		}, nil
	}

	// Generate diff
	diff := generateUnifiedDiff(req.Path, contentBefore, newContent)

	// Backup before editing
	if t.config.BackupOnEdit {
		if err := t.createBackup(absPath); err != nil {
			return nil, fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Write new content
	if err := os.WriteFile(absPath, []byte(newContent), info.Mode()); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &types.EditFileResponse{
		Path:          req.Path,
		Edited:        true,
		Replacements:  replacements,
		Diff:          diff,
		ContentBefore: contentBefore,
		ContentAfter:  newContent,
	}, nil
}

// CopyFile copies a file from source to destination.
func (t *Tool) CopyFile(req types.CopyFileRequest) (*types.CopyFileResponse, error) {
	srcPath, err := t.validatePath(req.Source)
	if err != nil {
		return nil, fmt.Errorf("source path invalid: %w", err)
	}

	dstPath, err := t.validatePath(req.Destination)
	if err != nil {
		return nil, fmt.Errorf("destination path invalid: %w", err)
	}

	// Check source exists
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("source file not found: %s", req.Source)
		}
		return nil, fmt.Errorf("failed to stat source: %w", err)
	}

	if srcInfo.IsDir() {
		return nil, fmt.Errorf("source is a directory; use recursive copy for directories")
	}

	// Check destination
	dstInfo, err := os.Stat(dstPath)
	if err == nil {
		if dstInfo.IsDir() {
			// Copy to directory with same filename
			dstPath = filepath.Join(dstPath, filepath.Base(srcPath))
			// Check if file exists in directory
			_, err = os.Stat(dstPath)
		}

		if err == nil && !req.Overwrite {
			return nil, fmt.Errorf("destination file exists and overwrite not requested")
		}
	}

	// Check file size
	if srcInfo.Size() > t.config.MaxFileSize {
		return nil, fmt.Errorf("source file too large: %d bytes (max: %d)", srcInfo.Size(), t.config.MaxFileSize)
	}

	// Perform copy
	src, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source: %w", err)
	}
	defer src.Close()

	dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return nil, fmt.Errorf("failed to create destination: %w", err)
	}
	defer dst.Close()

	bytesCopied, err := io.Copy(dst, src)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	return &types.CopyFileResponse{
		Source:      req.Source,
		Destination: req.Destination,
		Copied:      true,
		BytesCopied: bytesCopied,
	}, nil
}

// DeleteFile deletes a file or directory.
func (t *Tool) DeleteFile(req types.DeleteFileRequest) (*types.DeleteFileResponse, error) {
	if !t.config.AllowDelete {
		return nil, fmt.Errorf("delete operations are disabled in config")
	}

	absPath, err := t.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	// Check if path exists
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("path not found: %s", req.Path)
		}
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}

	wasDir := info.IsDir()

	// For directories, require recursive flag
	if wasDir && !req.Recursive {
		return nil, fmt.Errorf("path is a directory; set recursive=true to delete")
	}

	// Backup before deleting (only for files)
	if !wasDir && t.config.BackupOnEdit {
		if err := t.createBackup(absPath); err != nil {
			return nil, fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Delete
	var deleteErr error
	if wasDir {
		deleteErr = os.RemoveAll(absPath)
	} else {
		deleteErr = os.Remove(absPath)
	}

	if deleteErr != nil {
		return nil, fmt.Errorf("failed to delete: %w", deleteErr)
	}

	return &types.DeleteFileResponse{
		Path:    req.Path,
		Deleted: true,
		WasDir:  wasDir,
	}, nil
}

// ListDir lists the contents of a directory.
func (t *Tool) ListDir(req types.ListDirRequest) (*types.ListDirResponse, error) {
	absPath, err := t.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("directory not found: %s", req.Path)
		}
		return nil, fmt.Errorf("failed to stat path: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory")
	}

	var files []types.FileInfo

	if req.Recursive {
		maxDepth := req.MaxDepth
		if maxDepth <= 0 {
			maxDepth = 10 // Default max depth
		}

		err = t.walkDir(absPath, absPath, 0, maxDepth, &files)
		if err != nil {
			return nil, fmt.Errorf("failed to walk directory: %w", err)
		}
	} else {
		entries, err := os.ReadDir(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory: %w", err)
		}

		for _, entry := range entries {
			fileInfo, err := entry.Info()
			if err != nil {
				continue
			}

			files = append(files, types.FileInfo{
				Name:    entry.Name(),
				Path:    filepath.Join(req.Path, entry.Name()),
				Size:    fileInfo.Size(),
				Mode:    fileInfo.Mode().String(),
				ModTime: fileInfo.ModTime(),
				IsDir:   entry.IsDir(),
			})
		}
	}

	return &types.ListDirResponse{
		Path:  req.Path,
		Files: files,
	}, nil
}

// walkDir recursively walks a directory up to maxDepth.
func (t *Tool) walkDir(root, current string, depth, maxDepth int, files *[]types.FileInfo) error {
	if depth > maxDepth {
		return nil
	}

	entries, err := os.ReadDir(current)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(current, entry.Name())
		relPath, _ := filepath.Rel(root, fullPath)

		fileInfo, err := entry.Info()
		if err != nil {
			continue
		}

		*files = append(*files, types.FileInfo{
			Name:    entry.Name(),
			Path:    relPath,
			Size:    fileInfo.Size(),
			Mode:    fileInfo.Mode().String(),
			ModTime: fileInfo.ModTime(),
			IsDir:   entry.IsDir(),
		})

		if entry.IsDir() {
			if err := t.walkDir(root, fullPath, depth+1, maxDepth, files); err != nil {
				// Continue on error
				continue
			}
		}
	}

	return nil
}

// Exists checks if a path exists.
func (t *Tool) Exists(path string) (bool, bool, error) {
	absPath, err := t.validatePath(path)
	if err != nil {
		return false, false, err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false, nil
		}
		return false, false, err
	}

	return true, info.IsDir(), nil
}

// Hash returns the SHA256 hash of a file.
func (t *Tool) Hash(path string) (string, error) {
	absPath, err := t.validatePath(path)
	if err != nil {
		return "", err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// createBackup creates a backup of a file.
func (t *Tool) createBackup(absPath string) error {
	// Determine backup directory
	backupDir := t.config.BackupDirectory
	if backupDir == "" {
		backupDir = filepath.Dir(absPath)
	}

	// Ensure backup directory exists
	if err := os.MkdirAll(backupDir, 0750); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	filename := filepath.Base(absPath)
	backupName := fmt.Sprintf("%s.%s.bak", filename, timestamp)
	backupPath := filepath.Join(backupDir, backupName)

	// Copy file to backup
	src, err := os.Open(absPath)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy to backup: %w", err)
	}

	return nil
}

// parseFileMode parses a file mode string like "0644".
func parseFileMode(s string) (os.FileMode, error) {
	s = strings.TrimPrefix(s, "0")
	mode, err := strconv.ParseUint(s, 8, 32)
	if err != nil {
		return 0, err
	}
	return os.FileMode(mode), nil
}

// generateUnifiedDiff generates a simple unified diff between two strings.
func generateUnifiedDiff(filename, before, after string) string {
	beforeLines := strings.Split(before, "\n")
	afterLines := strings.Split(after, "\n")

	var diff bytes.Buffer

	diff.WriteString(fmt.Sprintf("--- a/%s\n", filename))
	diff.WriteString(fmt.Sprintf("+++ b/%s\n", filename))

	// Simple diff algorithm - find differences
	maxLines := len(beforeLines)
	if len(afterLines) > maxLines {
		maxLines = len(afterLines)
	}

	// Group changes into hunks
	type hunk struct {
		startBefore int
		startAfter  int
		linesBefore []string
		linesAfter  []string
	}

	var hunks []hunk
	var currentHunk *hunk
	contextLines := 3

	for i := 0; i < maxLines; i++ {
		var beforeLine, afterLine string
		hasBefore := i < len(beforeLines)
		hasAfter := i < len(afterLines)

		if hasBefore {
			beforeLine = beforeLines[i]
		}
		if hasAfter {
			afterLine = afterLines[i]
		}

		isDifferent := !hasBefore || !hasAfter || beforeLine != afterLine

		if isDifferent {
			if currentHunk == nil {
				// Start new hunk with context
				start := i - contextLines
				if start < 0 {
					start = 0
				}
				currentHunk = &hunk{
					startBefore: start + 1,
					startAfter:  start + 1,
				}
				// Add context lines
				for j := start; j < i; j++ {
					if j < len(beforeLines) {
						currentHunk.linesBefore = append(currentHunk.linesBefore, " "+beforeLines[j])
						currentHunk.linesAfter = append(currentHunk.linesAfter, " "+afterLines[j])
					}
				}
			}

			if hasBefore {
				currentHunk.linesBefore = append(currentHunk.linesBefore, "-"+beforeLine)
			}
			if hasAfter {
				currentHunk.linesAfter = append(currentHunk.linesAfter, "+"+afterLine)
			}
		} else if currentHunk != nil {
			// Add context after changes
			currentHunk.linesBefore = append(currentHunk.linesBefore, " "+beforeLine)
			currentHunk.linesAfter = append(currentHunk.linesAfter, " "+afterLine)

			// Check if we should close the hunk
			if i >= len(beforeLines)-1 || (i+contextLines < len(beforeLines) && beforeLines[i+contextLines] == afterLines[i+contextLines]) {
				hunks = append(hunks, *currentHunk)
				currentHunk = nil
			}
		}
	}

	if currentHunk != nil {
		hunks = append(hunks, *currentHunk)
	}

	// Output hunks
	for _, h := range hunks {
		// Combine lines, putting removals before additions
		var combined []string
		for _, l := range h.linesBefore {
			if strings.HasPrefix(l, "-") || strings.HasPrefix(l, " ") {
				combined = append(combined, l)
			}
		}
		for _, l := range h.linesAfter {
			if strings.HasPrefix(l, "+") {
				combined = append(combined, l)
			}
		}

		diff.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@\n", h.startBefore, len(h.linesBefore), h.startAfter, len(h.linesAfter)))
		for _, l := range combined {
			diff.WriteString(l)
			diff.WriteString("\n")
		}
	}

	return diff.String()
}
