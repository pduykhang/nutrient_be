package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"nutrient_be/internal/pkg/logger"
)

// ContextMiddleware creates middleware that injects context information
func ContextMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID if not present
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Generate trace ID if not present
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Extract user information from JWT token if available
		userID := ""
		if userIDValue, exists := c.Get("userID"); exists {
			if userIDStr, ok := userIDValue.(string); ok {
				userID = userIDStr
			}
		}

		// Create context by extending the request context with all the information
		// This preserves cancellation and timeout capabilities from the HTTP request
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, logger.RequestIDKey, requestID)
		ctx = context.WithValue(ctx, logger.TraceIDKey, traceID)
		ctx = context.WithValue(ctx, logger.UserIDKey, userID)
		ctx = context.WithValue(ctx, logger.IPAddressKey, c.ClientIP())
		ctx = context.WithValue(ctx, logger.UserAgentKey, c.GetHeader("User-Agent"))
		ctx = context.WithValue(ctx, logger.StartTimeKey, time.Now().Format(time.RFC3339))

		// Store context in Gin context
		c.Set("context", ctx)

		// Set response headers
		c.Header("X-Request-ID", requestID)
		c.Header("X-Trace-ID", traceID)

		// Continue to next handler
		c.Next()
	}
}

// GetContext extracts context from Gin context
func GetContext(c *gin.Context) context.Context {
	if ctxValue, exists := c.Get("context"); exists {
		if ctx, ok := ctxValue.(context.Context); ok {
			return ctx
		}
	}
	return context.Background()
}

type ContextData struct {
	RequestID string
	TraceID   string
	UserID    string
	IPAddress string
	UserAgent string
	StartTime time.Time
}

// getStringValueFromContext safely extracts a string value from context
func getStringValueFromContext(ctx context.Context, key logger.ContextKey) string {
	if val := ctx.Value(key); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// parseStartTimeFromContext parses StartTime from context (stored as RFC3339 string)
func parseStartTimeFromContext(ctx context.Context) time.Time {
	startTimeStr := getStringValueFromContext(ctx, logger.StartTimeKey)
	if startTimeStr == "" {
		return time.Time{}
	}
	parsed, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func GetContextData(c *gin.Context) *ContextData {
	ctxValue, exists := c.Get("context")
	if !exists {
		return nil
	}

	ctx, ok := ctxValue.(context.Context)
	if !ok {
		return nil
	}

	return &ContextData{
		RequestID: getStringValueFromContext(ctx, logger.RequestIDKey),
		TraceID:   getStringValueFromContext(ctx, logger.TraceIDKey),
		UserID:    getStringValueFromContext(ctx, logger.UserIDKey),
		IPAddress: getStringValueFromContext(ctx, logger.IPAddressKey),
		UserAgent: getStringValueFromContext(ctx, logger.UserAgentKey),
		StartTime: parseStartTimeFromContext(ctx),
	}
}
