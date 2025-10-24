package service

import (
	"context"
	"fmt"

	"nutrient_be/internal/pkg/logger"
)

// ExampleService demonstrates how to use context logger in service layer
type ExampleService struct {
	logger logger.Logger
}

// NewExampleService creates a new example service
func NewExampleService(log logger.Logger) *ExampleService {
	return &ExampleService{
		logger: log,
	}
}

// ProcessData demonstrates context logger usage in service
func (s *ExampleService) ProcessData(ctx context.Context, data string) error {
	// Create context logger from context
	contextLogger := logger.NewContextLogger(s.logger, ctx)

	// Log with context information automatically included
	contextLogger.Info("Processing data started",
		logger.String("data", data),
		logger.String("operation", "process_data"))

	// Simulate some processing
	if data == "" {
		contextLogger.Error("Empty data provided", logger.String("field", "data"))
		return fmt.Errorf("data cannot be empty")
	}

	// Log success with context
	contextLogger.Info("Data processed successfully",
		logger.String("result", "success"),
		logger.Int("data_length", len(data)))

	return nil
}

// ProcessWithAdditionalContext demonstrates adding more context
func (s *ExampleService) ProcessWithAdditionalContext(ctx context.Context, userID string, operation string) error {
	// Create context logger
	contextLogger := logger.NewContextLogger(s.logger, ctx)

	// Add additional context fields
	contextLogger = contextLogger.With(
		logger.String("operation", operation),
		logger.String("service", "ExampleService"),
	)

	// All logs will now include the additional context
	contextLogger.Info("Starting operation", logger.String("user_id", userID))

	// Simulate processing
	if userID == "" {
		contextLogger.Error("User ID is required")
		return fmt.Errorf("user ID is required")
	}

	contextLogger.Info("Operation completed successfully")
	return nil
}
