package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ContextKey represents a key for context values
type ContextKey string

const (
	RequestIDKey ContextKey = "request_id"
	TraceIDKey   ContextKey = "trace_id"
	UserIDKey    ContextKey = "user_id"
	IPAddressKey ContextKey = "ip_address"
	UserAgentKey ContextKey = "user_agent"
	StartTimeKey ContextKey = "start_time"
)

// zapLogger implements the Logger interface using Zap
type zapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates a new Zap logger instance
func NewZapLogger(isDevelopment bool) (Logger, error) {
	var zapLog *zap.Logger
	var err error

	if isDevelopment {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.CallerKey = "caller"
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		// Only show stack trace for Panic and Fatal levels, not for Error
		// This prevents verbose stack traces in normal error logging
		// Note: AddStacktrace with PanicLevel means stack trace only appears for Panic/Fatal
		zapLog, err = config.Build(
			zap.AddStacktrace(zapcore.PanicLevel), // Only add stack trace for Panic and above
			zap.AddCallerSkip(1),                  // Include caller information (file:line)
		)
	} else {
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		// Only show stack trace for Panic and Fatal in production
		zapLog, err = config.Build(zap.AddStacktrace(zapcore.PanicLevel))
	}

	if err != nil {
		return nil, err
	}

	return &zapLogger{logger: zapLog}, nil
}

// Context-aware logging methods

// Debug logs a debug message with context
func (l *zapLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	allFields := l.buildContextFields(ctx, fields)
	l.logger.Debug(msg, allFields...)
}

// Debugf logs a formatted debug message with context
func (l *zapLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	allFields := l.buildContextFields(ctx, nil)
	l.logger.Debug(msg, allFields...)
}

// Info logs an info message with context
func (l *zapLogger) Info(ctx context.Context, msg string, fields ...Field) {
	allFields := l.buildContextFields(ctx, fields)
	l.logger.Info(msg, allFields...)
}

// Infof logs a formatted info message with context
func (l *zapLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	allFields := l.buildContextFields(ctx, nil)
	l.logger.Info(msg, allFields...)
}

// Warn logs a warning message with context
func (l *zapLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	allFields := l.buildContextFields(ctx, fields)
	l.logger.Warn(msg, allFields...)
}

// Warnf logs a formatted warning message with context
func (l *zapLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	allFields := l.buildContextFields(ctx, nil)
	l.logger.Warn(msg, allFields...)
}

// Error logs an error message with context
func (l *zapLogger) Error(ctx context.Context, msg string, fields ...Field) {
	allFields := l.buildContextFields(ctx, fields)
	l.logger.Error(msg, allFields...)
}

// Errorf logs a formatted error message with context
func (l *zapLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	allFields := l.buildContextFields(ctx, nil)
	l.logger.Error(msg, allFields...)
}

// Fatal logs a fatal message with context and exits
func (l *zapLogger) Fatal(ctx context.Context, msg string, fields ...Field) {
	allFields := l.buildContextFields(ctx, fields)
	l.logger.Fatal(msg, allFields...)
}

// Fatalf logs a formatted fatal message with context and exits
func (l *zapLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	allFields := l.buildContextFields(ctx, nil)
	l.logger.Fatal(msg, allFields...)
}

// Panic logs a panic message with context
func (l *zapLogger) Panic(ctx context.Context, msg string, fields ...Field) {
	allFields := l.buildContextFields(ctx, fields)
	l.logger.Panic(msg, allFields...)
}

// Panicf logs a formatted panic message with context
func (l *zapLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	allFields := l.buildContextFields(ctx, nil)
	l.logger.Panic(msg, allFields...)
}

// Legacy methods without context (for backward compatibility)

// DebugLegacy logs a debug message without context
func (l *zapLogger) DebugLegacy(msg string, fields ...Field) {
	l.logger.Debug(msg, l.toZapFields(fields)...)
}

// DebugfLegacy logs a formatted debug message without context
func (l *zapLogger) DebugfLegacy(format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

// InfoLegacy logs an info message without context
func (l *zapLogger) InfoLegacy(msg string, fields ...Field) {
	l.logger.Info(msg, l.toZapFields(fields)...)
}

// InfofLegacy logs a formatted info message without context
func (l *zapLogger) InfofLegacy(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

// WarnLegacy logs a warning message without context
func (l *zapLogger) WarnLegacy(msg string, fields ...Field) {
	l.logger.Warn(msg, l.toZapFields(fields)...)
}

// WarnfLegacy logs a formatted warning message without context
func (l *zapLogger) WarnfLegacy(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

// ErrorLegacy logs an error message without context
func (l *zapLogger) ErrorLegacy(msg string, fields ...Field) {
	l.logger.Error(msg, l.toZapFields(fields)...)
}

// ErrorfLegacy logs a formatted error message without context
func (l *zapLogger) ErrorfLegacy(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

// FatalLegacy logs a fatal message without context and exits
func (l *zapLogger) FatalLegacy(msg string, fields ...Field) {
	l.logger.Fatal(msg, l.toZapFields(fields)...)
}

// FatalfLegacy logs a formatted fatal message without context and exits
func (l *zapLogger) FatalfLegacy(format string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, args...))
}

// PanicLegacy logs a panic message without context
func (l *zapLogger) PanicLegacy(msg string, fields ...Field) {
	l.logger.Panic(msg, l.toZapFields(fields)...)
}

// PanicfLegacy logs a formatted panic message without context
func (l *zapLogger) PanicfLegacy(format string, args ...interface{}) {
	l.logger.Panic(fmt.Sprintf(format, args...))
}

// With creates a new logger with additional fields
func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{
		logger: l.logger.With(l.toZapFields(fields)...),
	}
}

// buildContextFields builds fields from context and additional fields
func (l *zapLogger) buildContextFields(ctx context.Context, additionalFields []Field) []zap.Field {
	fields := make([]zap.Field, 0, len(additionalFields)+10) // Pre-allocate with extra capacity

	// Add context fields
	if requestID := l.getContextValue(ctx, RequestIDKey); requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}
	if traceID := l.getContextValue(ctx, TraceIDKey); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}
	if userID := l.getContextValue(ctx, UserIDKey); userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}
	if ipAddress := l.getContextValue(ctx, IPAddressKey); ipAddress != "" {
		fields = append(fields, zap.String("ip_address", ipAddress))
	}
	if userAgent := l.getContextValue(ctx, UserAgentKey); userAgent != "" {
		fields = append(fields, zap.String("user_agent", userAgent))
	}
	if startTime := l.getContextValue(ctx, StartTimeKey); startTime != "" {
		fields = append(fields, zap.String("start_time", startTime))
	}

	// Add additional fields
	for _, field := range additionalFields {
		// Special handling for error fields to avoid stack traces
		if field.Key == "error" {
			if err, ok := field.Value.(error); ok {
				// Only log error message, not full stack trace
				fields = append(fields, zap.String("error", err.Error()))
				continue
			}
		}
		fields = append(fields, zap.Any(field.Key, field.Value))
	}

	return fields
}

// getContextValue gets a value from context
func (l *zapLogger) getContextValue(ctx context.Context, key ContextKey) string {
	if value := ctx.Value(key); value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// toZapFields converts our Field structs to Zap fields
func (l *zapLogger) toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}
