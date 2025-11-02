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

		// Create context with all the information
		ctx := context.Background()
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

func GetContextData(c *gin.Context) *ContextData {
	if ctxValue, exists := c.Get("context"); exists {
		if ctx, ok := ctxValue.(context.Context); ok {
			return &ContextData{
				RequestID: ctx.Value(logger.RequestIDKey).(string),
				TraceID:   ctx.Value(logger.TraceIDKey).(string),
				UserID:    ctx.Value(logger.UserIDKey).(string),
				IPAddress: ctx.Value(logger.IPAddressKey).(string),
				UserAgent: ctx.Value(logger.UserAgentKey).(string),
				StartTime: ctx.Value(logger.StartTimeKey).(time.Time),
			}
		}
	}
	return nil
}
