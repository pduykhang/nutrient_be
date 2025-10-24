# Response Middleware Documentation

## Overview

Response middleware chuẩn hóa format của tất cả API responses để đảm bảo consistency và dễ dàng xử lý ở frontend.

## Response Format

### Standard Response Structure

```json
{
  "code": 200,
  "message": "Success",
  "data": {
    // Response data here
  },
  "meta": {
    "request_id": "req_12345",
    "trace_id": "trace_67890",
    "timestamp": 1704067200,
    "version": "v1.0"
  }
}
```

### Error Response Structure

```json
{
  "code": 400,
  "message": "Bad request",
  "error": {
    // Error details here
  },
  "meta": {
    "request_id": "req_12345",
    "trace_id": "trace_67890",
    "timestamp": 1704067200,
    "version": "v1.0"
  }
}
```

## Usage Examples

### 1. Success Response

```go
func (h *UserHandler) GetUser(c *gin.Context) {
    user, err := h.userService.GetUser(c.Request.Context(), userID)
    if err != nil {
        h.responseHelper.NotFound(c, gin.H{"user_id": userID}, "User not found")
        return
    }
    
    h.responseHelper.Success(c, user, "User retrieved successfully")
}
```

**Response:**
```json
{
  "code": 200,
  "message": "User retrieved successfully",
  "data": {
    "id": "123",
    "name": "John Doe",
    "email": "john@example.com"
  },
  "meta": {
    "request_id": "req_12345",
    "trace_id": "trace_67890",
    "timestamp": 1704067200,
    "version": "v1.0"
  }
}
```

### 2. Created Response

```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    user, err := h.userService.CreateUser(c.Request.Context(), req)
    if err != nil {
        h.responseHelper.Conflict(c, gin.H{"details": err.Error()}, "User already exists")
        return
    }
    
    h.responseHelper.Created(c, user, "User created successfully")
}
```

**Response:**
```json
{
  "code": 201,
  "message": "User created successfully",
  "data": {
    "id": "456",
    "name": "Jane Doe",
    "email": "jane@example.com"
  },
  "meta": {
    "request_id": "req_12345",
    "trace_id": "trace_67890",
    "timestamp": 1704067200,
    "version": "v1.0"
  }
}
```

### 3. Validation Error Response

```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.responseHelper.BadRequest(c, gin.H{"details": err.Error()}, "Invalid request body")
        return
    }
    
    if err := h.validator.Struct(&req); err != nil {
        validationErrors := gin.H{
            "email":    []string{"Email is required", "Invalid email format"},
            "password": []string{"Password must be at least 6 characters"},
        }
        h.responseHelper.ValidationError(c, validationErrors, "Validation failed")
        return
    }
    
    // ... rest of the logic
}
```

**Response:**
```json
{
  "code": 422,
  "message": "Validation failed",
  "error": {
    "email": ["Email is required", "Invalid email format"],
    "password": ["Password must be at least 6 characters"]
  },
  "meta": {
    "request_id": "req_12345",
    "trace_id": "trace_67890",
    "timestamp": 1704067200,
    "version": "v1.0"
  }
}
```

### 4. Not Found Response

```go
func (h *UserHandler) GetUser(c *gin.Context) {
    userID := c.Param("id")
    user, err := h.userService.GetUser(c.Request.Context(), userID)
    if err != nil {
        h.responseHelper.NotFound(c, gin.H{"user_id": userID}, "User not found")
        return
    }
    
    h.responseHelper.Success(c, user, "User retrieved successfully")
}
```

**Response:**
```json
{
  "code": 404,
  "message": "User not found",
  "error": {
    "user_id": "999"
  },
  "meta": {
    "request_id": "req_12345",
    "trace_id": "trace_67890",
    "timestamp": 1704067200,
    "version": "v1.0"
  }
}
```

## Response Helper Methods

### Success Responses

```go
// 200 OK
responseHelper.Success(c, data, "Operation successful")

// 201 Created
responseHelper.Created(c, data, "Resource created successfully")
```

### Error Responses

```go
// 400 Bad Request
responseHelper.BadRequest(c, error, "Invalid request")

// 401 Unauthorized
responseHelper.Unauthorized(c, error, "Authentication required")

// 403 Forbidden
responseHelper.Forbidden(c, error, "Access denied")

// 404 Not Found
responseHelper.NotFound(c, error, "Resource not found")

// 409 Conflict
responseHelper.Conflict(c, error, "Resource already exists")

// 422 Validation Error
responseHelper.ValidationError(c, errors, "Validation failed")

// 500 Internal Server Error
responseHelper.InternalError(c, error, "Internal server error")
```

## Setup

### 1. Add Middleware to Routes

```go
func SetupRoutes(r *gin.Engine, handlers *Handlers) {
    // Add response middleware
    r.Use(middleware.ResponseMiddleware(logger))
    
    // ... other middleware and routes
}
```

### 2. Use Response Helper in Handlers

```go
type UserHandler struct {
    userService    *service.UserService
    responseHelper *middleware.ResponseHelper
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{
        userService:    userService,
        responseHelper: middleware.NewResponseHelper(),
    }
}
```

### 3. Use in Handler Methods

```go
func (h *UserHandler) GetUser(c *gin.Context) {
    // Business logic
    user, err := h.userService.GetUser(c.Request.Context(), userID)
    if err != nil {
        h.responseHelper.NotFound(c, gin.H{"user_id": userID}, "User not found")
        return
    }
    
    h.responseHelper.Success(c, user, "User retrieved successfully")
}
```

## HTTP Status Code Mapping

| Method | HTTP Status | Description |
|--------|-------------|-------------|
| `Success()` | 200 | OK |
| `Created()` | 201 | Created |
| `BadRequest()` | 400 | Bad Request |
| `Unauthorized()` | 401 | Unauthorized |
| `Forbidden()` | 403 | Forbidden |
| `NotFound()` | 404 | Not Found |
| `Conflict()` | 409 | Conflict |
| `ValidationError()` | 422 | Unprocessable Entity |
| `InternalError()` | 500 | Internal Server Error |

## Default Messages

Nếu không truyền message, middleware sẽ sử dụng default messages:

- 200: "Success"
- 201: "Created successfully"
- 400: "Bad request"
- 401: "Unauthorized"
- 403: "Forbidden"
- 404: "Not found"
- 409: "Conflict"
- 422: "Validation failed"
- 500: "Internal server error"

## Meta Information

Mỗi response đều bao gồm meta information:

- `request_id`: Unique request identifier
- `trace_id`: Distributed tracing identifier
- `timestamp`: Unix timestamp
- `version`: API version

## Benefits

### 1. Consistency

Tất cả API responses đều có cùng format:

```json
{
  "code": 200,
  "message": "Success",
  "data": { ... },
  "meta": { ... }
}
```

### 2. Easy Frontend Handling

Frontend có thể xử lý responses một cách consistent:

```javascript
// Frontend code
const response = await fetch('/api/v1/users');
const data = await response.json();

if (data.code === 200) {
    // Success
    console.log(data.data);
} else {
    // Error
    console.error(data.message, data.error);
}
```

### 3. Better Debugging

Meta information giúp debug dễ dàng:

```json
{
  "meta": {
    "request_id": "req_12345",
    "trace_id": "trace_67890",
    "timestamp": 1704067200
  }
}
```

### 4. API Versioning

Dễ dàng thêm version information:

```json
{
  "meta": {
    "version": "v1.0"
  }
}
```

## Best Practices

### 1. Always Use Response Helper

```go
// ✅ Good
h.responseHelper.Success(c, data, "Operation successful")

// ❌ Bad - không consistent
c.JSON(200, data)
```

### 2. Provide Meaningful Messages

```go
// ✅ Good
h.responseHelper.Success(c, user, "User profile retrieved successfully")

// ❌ Bad - generic message
h.responseHelper.Success(c, user, "Success")
```

### 3. Include Relevant Error Details

```go
// ✅ Good
h.responseHelper.ValidationError(c, gin.H{
    "email":    []string{"Email is required"},
    "password": []string{"Password must be at least 6 characters"},
}, "Validation failed")

// ❌ Bad - không có details
h.responseHelper.ValidationError(c, gin.H{"error": "validation failed"}, "Validation failed")
```

### 4. Use Appropriate HTTP Status Codes

```go
// ✅ Good
h.responseHelper.NotFound(c, gin.H{"user_id": userID}, "User not found")

// ❌ Bad - wrong status code
h.responseHelper.BadRequest(c, gin.H{"user_id": userID}, "User not found")
```

## Migration Guide

### From Direct JSON Response

```go
// Before
c.JSON(200, gin.H{
    "user": user,
    "message": "Success",
})

// After
h.responseHelper.Success(c, gin.H{"user": user}, "User retrieved successfully")
```

### From Error Response

```go
// Before
c.JSON(404, gin.H{
    "error": "User not found",
    "user_id": userID,
})

// After
h.responseHelper.NotFound(c, gin.H{"user_id": userID}, "User not found")
```

## Testing

### Unit Test Example

```go
func TestGetUser_Success(t *testing.T) {
    // Setup
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    
    handler := NewUserHandler(mockUserService)
    
    // Execute
    handler.GetUser(c)
    
    // Assert
    assert.Equal(t, 200, w.Code)
    
    var response middleware.ResponseFormat
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, 200, response.Code)
    assert.Equal(t, "User retrieved successfully", response.Message)
    assert.NotNil(t, response.Data)
}
```

## Performance Considerations

### 1. Minimal Overhead

Response middleware có minimal overhead vì chỉ format response cuối cùng.

### 2. Memory Efficient

Sử dụng context để store response data thay vì tạo struct mới.

### 3. Reusable Helper

ResponseHelper được tạo một lần và reuse trong handler.

## Troubleshooting

### Common Issues

1. **Response not formatted**: Đảm bảo middleware được add đúng thứ tự
2. **Missing meta information**: Đảm bảo context middleware được add trước response middleware
3. **Wrong status code**: Sử dụng đúng method của ResponseHelper

### Debug Tips

1. Check middleware order trong routes setup
2. Verify ResponseHelper được inject vào handler
3. Check logs để xem response format
