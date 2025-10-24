package logger

import (
	"context"
	"time"
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

// ContextLogger wraps the logger with context information
type ContextLogger struct {
	logger Logger
	ctx    context.Context
}

// NewContextLogger creates a new context logger
func NewContextLogger(log Logger, ctx context.Context) *ContextLogger {
	return &ContextLogger{
		logger: log,
		ctx:    ctx,
	}
}

// Debug logs a debug message with context information
func (cl *ContextLogger) Debug(msg string, fields ...Field) {
	cl.logger.Debug(msg, cl.buildFields(fields)...)
}

// Info logs an info message with context information
func (cl *ContextLogger) Info(msg string, fields ...Field) {
	cl.logger.Info(msg, cl.buildFields(fields)...)
}

// Warn logs a warning message with context information
func (cl *ContextLogger) Warn(msg string, fields ...Field) {
	cl.logger.Warn(msg, cl.buildFields(fields)...)
}

// Error logs an error message with context information
func (cl *ContextLogger) Error(msg string, fields ...Field) {
	cl.logger.Error(msg, cl.buildFields(fields)...)
}

// Fatal logs a fatal message with context information
func (cl *ContextLogger) Fatal(msg string, fields ...Field) {
	cl.logger.Fatal(msg, cl.buildFields(fields)...)
}

// With creates a new context logger with additional fields
func (cl *ContextLogger) With(fields ...Field) *ContextLogger {
	return &ContextLogger{
		logger: cl.logger.With(fields...),
		ctx:    cl.ctx,
	}
}

// buildFields builds fields from context and additional fields
func (cl *ContextLogger) buildFields(additionalFields []Field) []Field {
	fields := make([]Field, 0, len(additionalFields)+10) // Pre-allocate with extra capacity

	// Add context fields
	if requestID := cl.getContextValue(RequestIDKey); requestID != "" {
		fields = append(fields, String("request_id", requestID))
	}
	if traceID := cl.getContextValue(TraceIDKey); traceID != "" {
		fields = append(fields, String("trace_id", traceID))
	}
	if userID := cl.getContextValue(UserIDKey); userID != "" {
		fields = append(fields, String("user_id", userID))
	}
	if ipAddress := cl.getContextValue(IPAddressKey); ipAddress != "" {
		fields = append(fields, String("ip_address", ipAddress))
	}
	if userAgent := cl.getContextValue(UserAgentKey); userAgent != "" {
		fields = append(fields, String("user_agent", userAgent))
	}
	if startTime := cl.getContextValue(StartTimeKey); startTime != "" {
		fields = append(fields, String("start_time", startTime))
	}

	// Add duration if start_time is available
	if startTimeStr := cl.getContextValue(StartTimeKey); startTimeStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			duration := time.Since(startTime)
			fields = append(fields, String("duration", duration.String()))
		}
	}

	// Add additional fields
	fields = append(fields, additionalFields...)

	return fields
}

// getContextValue gets a value from context
func (cl *ContextLogger) getContextValue(key ContextKey) string {
	if value := cl.ctx.Value(key); value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// ContextHelper provides helper functions for context management
type ContextHelper struct{}

// WithRequestID adds request ID to context
func (ch *ContextHelper) WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithTraceID adds trace ID to context
func (ch *ContextHelper) WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithUserID adds user ID to context
func (ch *ContextHelper) WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// WithIPAddress adds IP address to context
func (ch *ContextHelper) WithIPAddress(ctx context.Context, ipAddress string) context.Context {
	return context.WithValue(ctx, IPAddressKey, ipAddress)
}

// WithUserAgent adds user agent to context
func (ch *ContextHelper) WithUserAgent(ctx context.Context, userAgent string) context.Context {
	return context.WithValue(ctx, UserAgentKey, userAgent)
}

// WithStartTime adds start time to context
func (ch *ContextHelper) WithStartTime(ctx context.Context, startTime time.Time) context.Context {
	return context.WithValue(ctx, StartTimeKey, startTime.Format(time.RFC3339))
}

// GetRequestID gets request ID from context
func (ch *ContextHelper) GetRequestID(ctx context.Context) string {
	if value := ctx.Value(RequestIDKey); value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// GetTraceID gets trace ID from context
func (ch *ContextHelper) GetTraceID(ctx context.Context) string {
	if value := ctx.Value(TraceIDKey); value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// GetUserID gets user ID from context
func (ch *ContextHelper) GetUserID(ctx context.Context) string {
	if value := ctx.Value(UserIDKey); value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// NewContextHelper creates a new context helper
func NewContextHelper() *ContextHelper {
	return &ContextHelper{}
}
