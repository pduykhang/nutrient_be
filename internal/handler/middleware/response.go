package middleware

import (
	"net/http"
	"time"

	"nutrient_be/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

// ResponseFormat defines the standard response structure
type ResponseFormat struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta contains additional response metadata
type Meta struct {
	RequestID string `json:"request_id,omitempty"`
	TraceID   string `json:"trace_id,omitempty"`
	Timestamp int64  `json:"timestamp"`
	Version   string `json:"version,omitempty"`
}

// ResponseMiddleware creates middleware for standardizing API responses
func ResponseMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Override the default JSON method to use our standardized format
		c.Next()

		// Skip if response already sent (e.g., by error handler)
		if c.Writer.Written() {
			return
		}

		// Get response data from context
		var responseData interface{}
		if data, exists := c.Get("response_data"); exists {
			responseData = data
		}

		// Get response message from context
		var message string
		if msg, exists := c.Get("response_message"); exists {
			if msgStr, ok := msg.(string); ok {
				message = msgStr
			}
		}

		// Determine response code and message based on HTTP status
		statusCode := c.Writer.Status()
		if message == "" {
			message = getDefaultMessage(statusCode)
		}

		// Create standardized response
		response := ResponseFormat{
			Code:    statusCode,
			Message: message,
			Data:    responseData,
		}

		// Add metadata
		response.Meta = &Meta{
			RequestID: c.GetHeader("X-Request-ID"),
			TraceID:   c.GetHeader("X-Trace-ID"),
			Timestamp: getCurrentTimestamp(),
			Version:   "v1.0",
		}

		// Handle error responses
		if statusCode >= 400 {
			response.Error = responseData
			response.Data = nil
		}

		// Send standardized response
		c.JSON(statusCode, response)
	}
}

// ResponseHelper provides helper functions for setting response data
type ResponseHelper struct{}

// Success sets success response data
func (rh *ResponseHelper) Success(c *gin.Context, data interface{}, message ...string) {
	c.Set("response_data", data)
	if len(message) > 0 {
		c.Set("response_message", message[0])
	}
	c.Status(http.StatusOK)
}

// Created sets created response data
func (rh *ResponseHelper) Created(c *gin.Context, data interface{}, message ...string) {
	c.Set("response_data", data)
	if len(message) > 0 {
		c.Set("response_message", message[0])
	}
	c.Status(http.StatusCreated)
}

// BadRequest sets bad request response
func (rh *ResponseHelper) BadRequest(c *gin.Context, error interface{}, message ...string) {
	c.Set("response_data", error)
	if len(message) > 0 {
		c.Set("response_message", message[0])
	}
	c.Status(http.StatusBadRequest)
}

// Unauthorized sets unauthorized response
func (rh *ResponseHelper) Unauthorized(c *gin.Context, error interface{}, message ...string) {
	c.Set("response_data", error)
	if len(message) > 0 {
		c.Set("response_message", message[0])
	}
	c.Status(http.StatusUnauthorized)
}

// Forbidden sets forbidden response
func (rh *ResponseHelper) Forbidden(c *gin.Context, error interface{}, message ...string) {
	c.Set("response_data", error)
	if len(message) > 0 {
		c.Set("response_message", message[0])
	}
	c.Status(http.StatusForbidden)
}

// NotFound sets not found response
func (rh *ResponseHelper) NotFound(c *gin.Context, error interface{}, message ...string) {
	c.Set("response_data", error)
	if len(message) > 0 {
		c.Set("response_message", message[0])
	}
	c.Status(http.StatusNotFound)
}

// Conflict sets conflict response
func (rh *ResponseHelper) Conflict(c *gin.Context, error interface{}, message ...string) {
	c.Set("response_data", error)
	if len(message) > 0 {
		c.Set("response_message", message[0])
	}
	c.Status(http.StatusConflict)
}

// InternalError sets internal server error response
func (rh *ResponseHelper) InternalError(c *gin.Context, error interface{}, message ...string) {
	c.Set("response_data", error)
	if len(message) > 0 {
		c.Set("response_message", message[0])
	}
	c.Status(http.StatusInternalServerError)
}

// ValidationError sets validation error response
func (rh *ResponseHelper) ValidationError(c *gin.Context, errors interface{}, message ...string) {
	c.Set("response_data", errors)
	if len(message) > 0 {
		c.Set("response_message", message[0])
	}
	c.Status(http.StatusUnprocessableEntity)
}

// NewResponseHelper creates a new response helper
func NewResponseHelper() *ResponseHelper {
	return &ResponseHelper{}
}

// getDefaultMessage returns default message for HTTP status codes
func getDefaultMessage(statusCode int) string {
	switch statusCode {
	case http.StatusOK:
		return "Success"
	case http.StatusCreated:
		return "Created successfully"
	case http.StatusBadRequest:
		return "Bad request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "Not found"
	case http.StatusConflict:
		return "Conflict"
	case http.StatusUnprocessableEntity:
		return "Validation failed"
	case http.StatusInternalServerError:
		return "Internal server error"
	case http.StatusServiceUnavailable:
		return "Service unavailable"
	default:
		return "Unknown error"
	}
}

// getCurrentTimestamp returns current timestamp in Unix format
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}
