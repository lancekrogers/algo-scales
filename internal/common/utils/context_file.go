package utils

import (
	"context"
	"fmt"
	"os"
	"time"
)

// ContextFileOperations provides context-aware file operations
type ContextFileOperations struct {
	defaultTimeout time.Duration
}

// NewContextFileOperations creates a new context-aware file operations helper
func NewContextFileOperations(defaultTimeout time.Duration) *ContextFileOperations {
	return &ContextFileOperations{
		defaultTimeout: defaultTimeout,
	}
}

// WriteFileWithContext writes a file with context cancellation support
func (c *ContextFileOperations) WriteFileWithContext(ctx context.Context, filename string, data []byte, perm os.FileMode) error {
	// Create a channel to signal completion
	done := make(chan error, 1)
	
	go func() {
		done <- os.WriteFile(filename, data, perm)
	}()
	
	// Wait for completion or context cancellation
	select {
	case <-ctx.Done():
		return fmt.Errorf("file write cancelled: %w", ctx.Err())
	case err := <-done:
		return err
	}
}

// ReadFileWithContext reads a file with context cancellation support
func (c *ContextFileOperations) ReadFileWithContext(ctx context.Context, filename string) ([]byte, error) {
	// Create a channel to signal completion
	type result struct {
		data []byte
		err  error
	}
	done := make(chan result, 1)
	
	go func() {
		data, err := os.ReadFile(filename)
		done <- result{data: data, err: err}
	}()
	
	// Wait for completion or context cancellation
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("file read cancelled: %w", ctx.Err())
	case res := <-done:
		return res.data, res.err
	}
}

// MkdirTempWithContext creates a temporary directory with context support
func (c *ContextFileOperations) MkdirTempWithContext(ctx context.Context, dir, pattern string) (string, error) {
	// Create a channel to signal completion
	type result struct {
		path string
		err  error
	}
	done := make(chan result, 1)
	
	go func() {
		path, err := os.MkdirTemp(dir, pattern)
		done <- result{path: path, err: err}
	}()
	
	// Wait for completion or context cancellation
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("temp dir creation cancelled: %w", ctx.Err())
	case res := <-done:
		return res.path, res.err
	}
}

// ReadDirWithContext reads a directory with context support
func (c *ContextFileOperations) ReadDirWithContext(ctx context.Context, dirname string) ([]os.DirEntry, error) {
	// Create a channel to signal completion
	type result struct {
		entries []os.DirEntry
		err     error
	}
	done := make(chan result, 1)
	
	go func() {
		entries, err := os.ReadDir(dirname)
		done <- result{entries: entries, err: err}
	}()
	
	// Wait for completion or context cancellation
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("directory read cancelled: %w", ctx.Err())
	case res := <-done:
		return res.entries, res.err
	}
}

// Default package-level instance with 5 second timeout
var DefaultContextFileOps = NewContextFileOperations(5 * time.Second)

// Package-level convenience functions
func WriteFileWithContext(ctx context.Context, filename string, data []byte, perm os.FileMode) error {
	return DefaultContextFileOps.WriteFileWithContext(ctx, filename, data, perm)
}

func ReadFileWithContext(ctx context.Context, filename string) ([]byte, error) {
	return DefaultContextFileOps.ReadFileWithContext(ctx, filename)
}

func MkdirTempWithContext(ctx context.Context, dir, pattern string) (string, error) {
	return DefaultContextFileOps.MkdirTempWithContext(ctx, dir, pattern)
}

func ReadDirWithContext(ctx context.Context, dirname string) ([]os.DirEntry, error) {
	return DefaultContextFileOps.ReadDirWithContext(ctx, dirname)
}