package logger

import (
	"context"
)

// noopLogger implements the Logger interface with no-op operations
// Useful for testing or when logging should be disabled
type noopLogger struct{}

// NewNoopLogger creates a new no-op logger instance
func NewNoopLogger() Logger {
	return &noopLogger{}
}

// Context-aware logging methods (no-op)

// Debug logs a debug message with context (no-op)
func (l *noopLogger) Debug(ctx context.Context, msg string, fields ...Field) {}

// Debugf logs a formatted debug message with context (no-op)
func (l *noopLogger) Debugf(ctx context.Context, format string, args ...interface{}) {}

// Info logs an info message with context (no-op)
func (l *noopLogger) Info(ctx context.Context, msg string, fields ...Field) {}

// Infof logs a formatted info message with context (no-op)
func (l *noopLogger) Infof(ctx context.Context, format string, args ...interface{}) {}

// Warn logs a warning message with context (no-op)
func (l *noopLogger) Warn(ctx context.Context, msg string, fields ...Field) {}

// Warnf logs a formatted warning message with context (no-op)
func (l *noopLogger) Warnf(ctx context.Context, format string, args ...interface{}) {}

// Error logs an error message with context (no-op)
func (l *noopLogger) Error(ctx context.Context, msg string, fields ...Field) {}

// Errorf logs a formatted error message with context (no-op)
func (l *noopLogger) Errorf(ctx context.Context, format string, args ...interface{}) {}

// Fatal logs a fatal message with context and exits (no-op)
func (l *noopLogger) Fatal(ctx context.Context, msg string, fields ...Field) {}

// Fatalf logs a formatted fatal message with context and exits (no-op)
func (l *noopLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {}

// Panic logs a panic message with context (no-op)
func (l *noopLogger) Panic(ctx context.Context, msg string, fields ...Field) {}

// Panicf logs a formatted panic message with context (no-op)
func (l *noopLogger) Panicf(ctx context.Context, format string, args ...interface{}) {}

// Legacy methods without context (for backward compatibility)

// DebugLegacy logs a debug message without context (no-op)
func (l *noopLogger) DebugLegacy(msg string, fields ...Field) {}

// DebugfLegacy logs a formatted debug message without context (no-op)
func (l *noopLogger) DebugfLegacy(format string, args ...interface{}) {}

// InfoLegacy logs an info message without context (no-op)
func (l *noopLogger) InfoLegacy(msg string, fields ...Field) {}

// InfofLegacy logs a formatted info message without context (no-op)
func (l *noopLogger) InfofLegacy(format string, args ...interface{}) {}

// WarnLegacy logs a warning message without context (no-op)
func (l *noopLogger) WarnLegacy(msg string, fields ...Field) {}

// WarnfLegacy logs a formatted warning message without context (no-op)
func (l *noopLogger) WarnfLegacy(format string, args ...interface{}) {}

// ErrorLegacy logs an error message without context (no-op)
func (l *noopLogger) ErrorLegacy(msg string, fields ...Field) {}

// ErrorfLegacy logs a formatted error message without context (no-op)
func (l *noopLogger) ErrorfLegacy(format string, args ...interface{}) {}

// FatalLegacy logs a fatal message without context and exits (no-op)
func (l *noopLogger) FatalLegacy(msg string, fields ...Field) {}

// FatalfLegacy logs a formatted fatal message without context and exits (no-op)
func (l *noopLogger) FatalfLegacy(format string, args ...interface{}) {}

// PanicLegacy logs a panic message without context (no-op)
func (l *noopLogger) PanicLegacy(msg string, fields ...Field) {}

// PanicfLegacy logs a formatted panic message without context (no-op)
func (l *noopLogger) PanicfLegacy(format string, args ...interface{}) {}

// With creates a new logger with additional fields (no-op)
func (l *noopLogger) With(fields ...Field) Logger {
	return l
}
