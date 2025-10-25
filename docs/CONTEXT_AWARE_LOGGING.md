# Context-Aware Logging System

## Overview

The context-aware logging system provides automatic injection of request-specific information into log messages without requiring manual passing of context fields. This system makes it easy to trace requests across the entire application stack.

## Features

### Automatic Context Injection
- **Request ID**: Unique identifier for each request
- **Trace ID**: For distributed tracing across services
- **User ID**: Current authenticated user
- **IP Address**: Client IP address
- **User Agent**: Client browser/client information
- **Start Time**: Request start timestamp

### Dual Interface Support
- **Context-aware methods**: `Info(ctx, msg, fields...)` - automatically includes context
- **Legacy methods**: `InfoLegacy(msg, fields...)` - for backward compatibility

## Usage

### Basic Context Logging

```go
// In handlers
func (h *AuthHandler) Login(c *gin.Context) {
    ctx := middleware.GetContext(c)
    
    // Context information automatically included
    h.logger.Info(ctx, "User login attempt", logger.String("email", email))
    
    // Error logging with context
    h.logger.Error(ctx, "Login failed", logger.Error(err))
}
```

### Service Layer Logging

```go
// In services
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) error {
    // Context automatically includes request_id, trace_id, user_id, etc.
    s.logger.Info(ctx, "Processing registration", logger.String("email", req.Email))
    
    if err := s.validateUser(req); err != nil {
        s.logger.Error(ctx, "User validation failed", logger.Error(err))
        return err
    }
    
    s.logger.Info(ctx, "User registered successfully")
    return nil
}
```

### Enhanced Logger with Additional Fields

```go
// Create logger with persistent fields
enhancedLogger := logger.With(
    logger.String("module", "auth"),
    logger.String("version", "v1.0"),
)

// All logs will include both context and persistent fields
enhancedLogger.Info(ctx, "Authentication started")
enhancedLogger.Info(ctx, "Authentication completed")
```

### Formatted Logging

```go
// Formatted messages with context
logger.Infof(ctx, "Processing user %s with %d items", username, itemCount)
logger.Errorf(ctx, "Failed to process request: %v", err)
```

## Context Keys

The system uses the following context keys:

```go
const (
    RequestIDKey ContextKey = "request_id"
    TraceIDKey   ContextKey = "trace_id"
    UserIDKey    ContextKey = "user_id"
    IPAddressKey ContextKey = "ip_address"
    UserAgentKey ContextKey = "user_agent"
    StartTimeKey ContextKey = "start_time"
)
```

## Middleware Integration

### Context Middleware
The `ContextMiddleware` automatically injects context information:

```go
// In routes setup
r.Use(middleware.ContextMiddleware(logger))
```

This middleware:
- Generates request ID and trace ID if not present
- Extracts user information from JWT tokens
- Sets response headers for tracing
- Stores context in Gin context

### Usage in Handlers

```go
func (h *Handler) SomeEndpoint(c *gin.Context) {
    // Get context with all request information
    ctx := middleware.GetContext(c)
    
    // Log with automatic context injection
    h.logger.Info(ctx, "Processing request")
}
```

## Log Output Example

### Context-Aware Logging
```json
{
  "level": "info",
  "timestamp": "2025-10-25T08:54:32.629+0700",
  "message": "User login successful",
  "request_id": "req_12345",
  "trace_id": "trace_67890", 
  "user_id": "user_abc123",
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
  "start_time": "2025-10-25T08:54:30+07:00",
  "email": "user@example.com"
}
```

### Legacy Logging (for comparison)
```json
{
  "level": "info",
  "timestamp": "2025-10-25T08:54:32.630+0700",
  "message": "User login successful",
  "email": "user@example.com"
}
```

## Benefits

### 1. **Automatic Context Injection**
- No need to manually pass request_id, trace_id, user_id in every log call
- Reduces boilerplate code
- Ensures consistent context information across all logs

### 2. **Easy Request Tracing**
- Track requests across multiple services and functions
- Debug issues by filtering logs by request_id or trace_id
- Monitor user behavior by filtering by user_id

### 3. **Backward Compatibility**
- Legacy methods available for gradual migration
- Existing code continues to work without changes
- Can migrate incrementally

### 4. **Performance Optimized**
- Context fields are extracted once per request
- Efficient field building with pre-allocated capacity
- Minimal overhead compared to manual field passing

## Migration Guide

### From Old Context Logger

**Before:**
```go
contextLogger := logger.NewContextLogger(baseLogger, ctx)
contextLogger.Info("User login", logger.String("email", email))
```

**After:**
```go
logger.Info(ctx, "User login", logger.String("email", email))
```

### From Manual Field Passing

**Before:**
```go
logger.Info("User login", 
    logger.String("request_id", requestID),
    logger.String("user_id", userID),
    logger.String("email", email))
```

**After:**
```go
logger.Info(ctx, "User login", logger.String("email", email))
```

## Implementation Details

### Logger Interface
```go
type Logger interface {
    // Context-aware methods
    Debug(ctx context.Context, msg string, fields ...Field)
    Info(ctx context.Context, msg string, fields ...Field)
    Warn(ctx context.Context, msg string, fields ...Field)
    Error(ctx context.Context, msg string, fields ...Field)
    Fatal(ctx context.Context, msg string, fields ...Field)
    Panic(ctx context.Context, msg string, fields ...Field)
    
    // Formatted context-aware methods
    Debugf(ctx context.Context, format string, args ...interface{})
    Infof(ctx context.Context, format string, args ...interface{})
    Warnf(ctx context.Context, format string, args ...interface{})
    Errorf(ctx context.Context, format string, args ...interface{})
    Fatalf(ctx context.Context, format string, args ...interface{})
    Panicf(ctx context.Context, format string, args ...interface{})
    
    // Legacy methods for backward compatibility
    DebugLegacy(msg string, fields ...Field)
    InfoLegacy(msg string, fields ...Field)
    // ... other legacy methods
    
    // With method for persistent fields
    With(fields ...Field) Logger
}
```

### Context Field Building
The system automatically extracts context values and combines them with additional fields:

```go
func (l *zapLogger) buildContextFields(ctx context.Context, additionalFields []Field) []zap.Field {
    fields := make([]zap.Field, 0, len(additionalFields)+10)
    
    // Add context fields
    if requestID := l.getContextValue(ctx, RequestIDKey); requestID != "" {
        fields = append(fields, zap.String("request_id", requestID))
    }
    // ... other context fields
    
    // Add additional fields
    for _, field := range additionalFields {
        fields = append(fields, zap.Any(field.Key, field.Value))
    }
    
    return fields
}
```

## Best Practices

1. **Always use context-aware methods** in handlers and services
2. **Pass context through the call chain** to maintain tracing
3. **Use legacy methods only** for initialization and shutdown logging
4. **Add meaningful additional fields** for business context
5. **Use structured logging** with consistent field names
6. **Filter logs by request_id** for debugging specific requests

## Demo

Run the demo to see context logging in action:

```bash
go run examples/context_logging_demo.go
```

This demonstrates:
- Basic context logging
- Error logging with context
- Formatted logging
- Enhanced logger with persistent fields
- Comparison with legacy logging
