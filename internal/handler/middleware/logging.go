package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"nutrient_be/internal/pkg/logger"
)

const (
	RequestIDHeader = "X-Request-ID"
)

// LoggingMiddleware logs HTTP requests
// Note: This middleware should be placed after ContextMiddleware to use enriched context
func LoggingMiddleware(log logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Try to get enriched context from Gin context (set by ContextMiddleware)
		// If not available, create a basic context with request ID
		var ctx context.Context
		if ctxValue, exists := param.Keys["context"]; exists {
			if enrichedCtx, ok := ctxValue.(context.Context); ok {
				ctx = enrichedCtx
			} else {
				// Fallback: create basic context
				ctx = context.Background()
				ctx = context.WithValue(ctx, logger.RequestIDKey, param.Request.Header.Get(RequestIDHeader))
				ctx = context.WithValue(ctx, logger.IPAddressKey, param.ClientIP)
				ctx = context.WithValue(ctx, logger.UserAgentKey, param.Request.UserAgent())
			}
		} else {
			// Fallback: create basic context
			ctx = context.Background()
			ctx = context.WithValue(ctx, logger.RequestIDKey, param.Request.Header.Get(RequestIDHeader))
			ctx = context.WithValue(ctx, logger.IPAddressKey, param.ClientIP)
			ctx = context.WithValue(ctx, logger.UserAgentKey, param.Request.UserAgent())
		}

		log.Info(ctx, "HTTP Request",
			logger.String("method", param.Method),
			logger.String("path", param.Path),
			logger.Int("status", param.StatusCode),
			logger.String("latency", param.Latency.String()),
			logger.String("clientIP", param.ClientIP),
			logger.String("userAgent", param.Request.UserAgent()),
		)
		return ""
	})
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(log logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Get context from Gin context
		ctx := GetContext(c)

		log.Error(ctx, "Panic recovered",
			logger.Any("error", recovered),
			logger.String("path", c.Request.URL.Path),
			logger.String("method", c.Request.Method),
		)
		c.AbortWithStatus(500)
	})
}

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestTimeoutMiddleware sets request timeout
func RequestTimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set timeout for the request
		c.Request = c.Request.WithContext(c.Request.Context())
		c.Next()
	}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(RequestIDHeader, requestID)
		c.Header(RequestIDHeader, requestID)
		c.Next()
	}
}
