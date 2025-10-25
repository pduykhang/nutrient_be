package logger

import (
	"context"
)

// Logger interface defines the contract for logging operations
// This abstraction allows easy switching between different logging libraries
type Logger interface {
	// Context-aware logging methods
	Debug(ctx context.Context, msg string, fields ...Field)
	Debugf(ctx context.Context, format string, args ...interface{})
	Info(ctx context.Context, msg string, fields ...Field)
	Infof(ctx context.Context, format string, args ...interface{})
	Warn(ctx context.Context, msg string, fields ...Field)
	Warnf(ctx context.Context, format string, args ...interface{})
	Error(ctx context.Context, msg string, fields ...Field)
	Errorf(ctx context.Context, format string, args ...interface{})
	Fatal(ctx context.Context, msg string, fields ...Field)
	Fatalf(ctx context.Context, format string, args ...interface{})
	Panic(ctx context.Context, msg string, fields ...Field)
	Panicf(ctx context.Context, format string, args ...interface{})

	// Legacy methods without context (for backward compatibility)
	DebugLegacy(msg string, fields ...Field)
	DebugfLegacy(format string, args ...interface{})
	InfoLegacy(msg string, fields ...Field)
	InfofLegacy(format string, args ...interface{})
	WarnLegacy(msg string, fields ...Field)
	WarnfLegacy(format string, args ...interface{})
	ErrorLegacy(msg string, fields ...Field)
	ErrorfLegacy(format string, args ...interface{})
	FatalLegacy(msg string, fields ...Field)
	FatalfLegacy(format string, args ...interface{})
	PanicLegacy(msg string, fields ...Field)
	PanicfLegacy(format string, args ...interface{})

	// With method for adding persistent fields
	With(fields ...Field) Logger
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// Helper functions to create fields
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}
