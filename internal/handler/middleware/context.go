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
	contextHelper := logger.NewContextHelper()

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
		ctx = contextHelper.WithRequestID(ctx, requestID)
		ctx = contextHelper.WithTraceID(ctx, traceID)
		ctx = contextHelper.WithUserID(ctx, userID)
		ctx = contextHelper.WithIPAddress(ctx, c.ClientIP())
		ctx = contextHelper.WithUserAgent(ctx, c.GetHeader("User-Agent"))
		ctx = contextHelper.WithStartTime(ctx, time.Now())

		// Store context in Gin context
		c.Set("context", ctx)

		// Set response headers
		c.Header("X-Request-ID", requestID)
		c.Header("X-Trace-ID", traceID)

		// Create context logger
		contextLogger := logger.NewContextLogger(log, ctx)
		c.Set("logger", contextLogger)

		// Continue to next handler
		c.Next()
	}
}

// GetContextLogger extracts context logger from Gin context
func GetContextLogger(c *gin.Context) *logger.ContextLogger {
	if loggerValue, exists := c.Get("logger"); exists {
		if contextLogger, ok := loggerValue.(*logger.ContextLogger); ok {
			return contextLogger
		}
	}
	return nil
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
