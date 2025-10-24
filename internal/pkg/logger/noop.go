package logger

// noopLogger implements the Logger interface with no-op operations
// Useful for testing or when logging should be disabled
type noopLogger struct{}

// NewNoopLogger creates a new no-op logger instance
func NewNoopLogger() Logger {
	return &noopLogger{}
}

// Debug does nothing
func (l *noopLogger) Debug(msg string, fields ...Field) {}

// Info does nothing
func (l *noopLogger) Info(msg string, fields ...Field) {}

// Warn does nothing
func (l *noopLogger) Warn(msg string, fields ...Field) {}

// Error does nothing
func (l *noopLogger) Error(msg string, fields ...Field) {}

// Fatal does nothing
func (l *noopLogger) Fatal(msg string, fields ...Field) {}

// With returns the same no-op logger
func (l *noopLogger) With(fields ...Field) Logger {
	return l
}
