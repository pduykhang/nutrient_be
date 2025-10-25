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
	// Log with context information automatically included
	s.logger.Info(ctx, "Processing data started",
		logger.String("data", data),
		logger.String("operation", "process_data"))

	// Simulate some processing
	if data == "" {
		s.logger.Error(ctx, "Empty data provided", logger.String("field", "data"))
		return fmt.Errorf("data cannot be empty")
	}

	// Log success with context
	s.logger.Info(ctx, "Data processed successfully",
		logger.String("result", "success"),
		logger.Int("data_length", len(data)))

	return nil
}

// ProcessWithAdditionalContext demonstrates adding more context
func (s *ExampleService) ProcessWithAdditionalContext(ctx context.Context, userID string, operation string) error {
	// Add additional context fields
	enhancedLogger := s.logger.With(
		logger.String("operation", operation),
		logger.String("service", "ExampleService"),
	)

	// All logs will now include the additional context
	enhancedLogger.Info(ctx, "Starting operation", logger.String("user_id", userID))

	// Simulate processing
	if userID == "" {
		enhancedLogger.Error(ctx, "User ID is required")
		return fmt.Errorf("user ID is required")
	}

	enhancedLogger.Info(ctx, "Operation completed successfully")
	return nil
}
