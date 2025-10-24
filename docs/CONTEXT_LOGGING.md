# Context Logging Documentation

## Overview

Context logging system tự động inject các thông tin như request ID, trace ID, user ID, IP address, user agent, và duration vào mọi log message mà không cần truyền thủ công.

## Features

- **Automatic Context Injection**: Tự động thêm request ID, trace ID, user ID, IP, user agent vào log
- **Duration Tracking**: Tự động tính toán và log thời gian xử lý request
- **Easy Usage**: Chỉ cần gọi `contextLogger.Info()` thay vì truyền nhiều fields
- **Flexible**: Có thể thêm context fields tùy chỉnh
- **Performance**: Pre-allocate fields để tối ưu performance

## Usage Examples

### 1. Trong Handler Layer

```go
func (h *AuthHandler) Register(c *gin.Context) {
    // Get context logger từ middleware
    contextLogger := middleware.GetContextLogger(c)
    if contextLogger == nil {
        contextLogger = logger.NewContextLogger(h.logger, c.Request.Context())
    }

    var req service.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // Log error với context tự động
        contextLogger.Error("Failed to bind register request", logger.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Log success với context tự động
    contextLogger.Info("User registered successfully", logger.String("email", req.Email))
    c.JSON(http.StatusCreated, response)
}
```

### 2. Trong Service Layer

```go
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
    // Tạo context logger từ context
    contextLogger := logger.NewContextLogger(s.logger, ctx)
    
    // Log với context tự động
    contextLogger.Info("Starting user registration", logger.String("email", req.Email))
    
    // Check if user exists
    existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err == nil && existingUser != nil {
        contextLogger.Error("User already exists", logger.String("email", req.Email))
        return nil, fmt.Errorf("user with email %s already exists", req.Email)
    }

    // Log success
    contextLogger.Info("User registered successfully", 
        logger.String("email", req.Email),
        logger.String("userID", user.ID.Hex()))
    
    return response, nil
}
```

### 3. Với Additional Context Fields

```go
func (s *FoodService) CreateFood(ctx context.Context, food *domain.FoodItem) error {
    contextLogger := logger.NewContextLogger(s.logger, ctx)
    
    // Thêm additional context fields
    contextLogger = contextLogger.With(
        logger.String("service", "FoodService"),
        logger.String("operation", "create_food"),
        logger.String("food_category", food.Category),
    )
    
    // Tất cả logs sau này sẽ include additional context
    contextLogger.Info("Creating food item", logger.String("food_name", food.Name["en"]))
    
    if err := s.foodRepo.Create(ctx, food); err != nil {
        contextLogger.Error("Failed to create food item", logger.Error(err))
        return err
    }
    
    contextLogger.Info("Food item created successfully")
    return nil
}
```

## Context Fields

### Automatic Fields

Các fields này được tự động thêm vào mọi log:

- `request_id`: Unique request identifier
- `trace_id`: Distributed tracing identifier  
- `user_id`: User ID (nếu authenticated)
- `ip_address`: Client IP address
- `user_agent`: Client user agent
- `start_time`: Request start time
- `duration`: Request processing duration

### Custom Fields

Có thể thêm custom fields:

```go
contextLogger = contextLogger.With(
    logger.String("module", "payment"),
    logger.String("transaction_id", "txn_123"),
    logger.Int("amount", 1000),
)
```

## Middleware Setup

### 1. Context Middleware

```go
// Trong routes setup
r.Use(middleware.ContextMiddleware(logger))
```

### 2. Context Helper

```go
contextHelper := logger.NewContextHelper()

// Add context values
ctx = contextHelper.WithRequestID(ctx, "req_123")
ctx = contextHelper.WithTraceID(ctx, "trace_456")
ctx = contextHelper.WithUserID(ctx, "user_789")
ctx = contextHelper.WithIPAddress(ctx, "192.168.1.1")
ctx = contextHelper.WithUserAgent(ctx, "Mozilla/5.0...")
ctx = contextHelper.WithStartTime(ctx, time.Now())
```

## Log Output Examples

### Console Output (Development)

```
2025-01-15T10:30:00.123Z	INFO	User registered successfully	{"request_id": "req_123", "trace_id": "trace_456", "user_id": "", "ip_address": "192.168.1.1", "user_agent": "Mozilla/5.0...", "start_time": "2025-01-15T10:29:59.000Z", "duration": "1.123s", "email": "user@example.com"}
```

### JSON Output (Production)

```json
{
  "level": "info",
  "timestamp": "2025-01-15T10:30:00.123Z",
  "message": "User registered successfully",
  "request_id": "req_123",
  "trace_id": "trace_456", 
  "user_id": "user_789",
  "ip_address": "192.168.1.1",
  "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
  "start_time": "2025-01-15T10:29:59.000Z",
  "duration": "1.123s",
  "email": "user@example.com"
}
```

## Best Practices

### 1. Always Use Context Logger

```go
// ✅ Good
contextLogger := middleware.GetContextLogger(c)
contextLogger.Info("Operation completed")

// ❌ Bad - không có context
h.logger.Info("Operation completed")
```

### 2. Add Meaningful Fields

```go
// ✅ Good
contextLogger.Info("User login successful", 
    logger.String("email", req.Email),
    logger.String("login_method", "password"))

// ❌ Bad - không có context cụ thể
contextLogger.Info("Login successful")
```

### 3. Use With() for Additional Context

```go
// ✅ Good - thêm context cho cả operation
contextLogger = contextLogger.With(
    logger.String("module", "auth"),
    logger.String("operation", "login"),
)
contextLogger.Info("Starting login process")
contextLogger.Info("Validating credentials")
contextLogger.Info("Login completed")

// ❌ Bad - lặp lại fields
contextLogger.Info("Starting login process", logger.String("module", "auth"))
contextLogger.Info("Validating credentials", logger.String("module", "auth"))
contextLogger.Info("Login completed", logger.String("module", "auth"))
```

### 4. Handle Missing Context Logger

```go
// ✅ Good - fallback nếu không có context logger
contextLogger := middleware.GetContextLogger(c)
if contextLogger == nil {
    contextLogger = logger.NewContextLogger(h.logger, c.Request.Context())
}
```

## Performance Considerations

### 1. Pre-allocate Fields

Context logger pre-allocates fields để tránh memory allocation:

```go
// Pre-allocate với capacity
fields := make([]Field, 0, len(additionalFields)+10)
```

### 2. Lazy Context Extraction

Context values chỉ được extract khi cần thiết:

```go
if requestID := cl.getContextValue(RequestIDKey); requestID != "" {
    fields = append(fields, String("request_id", requestID))
}
```

### 3. Reuse Context Logger

Tạo context logger một lần và reuse:

```go
// ✅ Good - reuse context logger
contextLogger := logger.NewContextLogger(s.logger, ctx)
contextLogger.Info("Step 1 completed")
contextLogger.Info("Step 2 completed")
contextLogger.Info("All steps completed")

// ❌ Bad - tạo mới mỗi lần
logger.NewContextLogger(s.logger, ctx).Info("Step 1 completed")
logger.NewContextLogger(s.logger, ctx).Info("Step 2 completed")
```

## Migration Guide

### From Regular Logger

```go
// Before
h.logger.Info("User registered", 
    logger.String("email", req.Email),
    logger.String("userID", user.ID.Hex()),
    logger.String("requestID", requestID),
    logger.String("traceID", traceID),
)

// After
contextLogger := middleware.GetContextLogger(c)
contextLogger.Info("User registered", 
    logger.String("email", req.Email),
    logger.String("userID", user.ID.Hex()),
)
```

### Benefits

1. **Less Code**: Không cần truyền request_id, trace_id, user_id mỗi lần
2. **Consistency**: Tất cả logs đều có cùng context fields
3. **Maintainability**: Dễ maintain và update context fields
4. **Debugging**: Dễ trace requests qua logs
5. **Performance**: Tối ưu memory allocation
