package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"nutrient_be/internal/domain"
	"nutrient_be/internal/dto/request"
	"nutrient_be/internal/pkg/logger"
	"nutrient_be/internal/pkg/validator"
)

// FoodRepository defines the interface for food data operations used by FoodService
type FoodRepository interface {
	Create(ctx context.Context, food *domain.FoodItem) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.FoodItem, error)
	Search(ctx context.Context, query string, userID primitive.ObjectID, limit, offset int) ([]*domain.FoodItem, error)
	GetByCategory(ctx context.Context, category string, userID primitive.ObjectID, limit, offset int) ([]*domain.FoodItem, error)
	GetByUser(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*domain.FoodItem, error)
	Update(ctx context.Context, food *domain.FoodItem) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetPublicFoods(ctx context.Context, limit, offset int) ([]*domain.FoodItem, error)
}

// FoodService handles food-related business logic
type FoodService struct {
	foodRepo  FoodRepository
	validator *validator.FoodValidator
	logger    logger.Logger
}

// NewFoodService creates a new food service
func NewFoodService(foodRepo FoodRepository, log logger.Logger) *FoodService {
	return &FoodService{
		foodRepo:  foodRepo,
		validator: validator.NewFoodValidator(log),
		logger:    log,
	}
}

// CreateFood creates a new food item with validation
func (s *FoodService) CreateFood(ctx context.Context, userID string, req *request.CreateFoodRequest) error {
	s.logger.Info(ctx, "Creating food", logger.String("food_name", req.Name.Get("en")))

	// Validate request using centralized validator
	if err := s.validator.ValidateCreateRequest(ctx, req); err != nil {
		s.logger.Error(ctx, "Food validation failed", logger.Error(err))
		return fmt.Errorf("validation failed: %w", err)
	}

	// Convert request to domain entity
	foodDB := domain.FoodItemFromRequest(ctx, req, userID)

	// Save to database
	if err := s.foodRepo.Create(ctx, foodDB); err != nil {
		s.logger.Error(ctx, "Failed to create food", logger.Error(err))
		return fmt.Errorf("failed to create food: %w", err)
	}

	s.logger.Info(ctx, "Food created successfully", logger.String("food_id", foodDB.ID.Hex()))
	return nil
}

// SearchFood searches for food items based on query
// It extracts userID from context to filter results (public foods + user's own foods)
func (s *FoodService) SearchFood(ctx context.Context, req *request.SearchFoodRequest) ([]*domain.FoodItem, error) {
	s.logger.Info(ctx, "Searching food", logger.String("query", req.Query))

	// Extract userID from context (set by auth middleware and context middleware)
	userIDStr := ""
	if userIDValue := ctx.Value(logger.UserIDKey); userIDValue != nil {
		if userID, ok := userIDValue.(string); ok {
			userIDStr = userID
		}
	}

	// Convert userID string to ObjectID for repository call
	// If userID is empty, use empty ObjectID (will only return public foods)
	var userIDObj primitive.ObjectID
	if userIDStr != "" {
		var err error
		userIDObj, err = primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			s.logger.Warn(ctx, "Invalid userID in context, using empty ObjectID", logger.Error(err))
			userIDObj = primitive.NilObjectID
		}
	}

	foods, err := s.foodRepo.Search(ctx, req.Query, userIDObj, req.Limit, req.Offset)
	if err != nil {
		s.logger.Error(ctx, "Failed to search food", logger.Error(err))
		return nil, fmt.Errorf("failed to search food: %w", err)
	}
	s.logger.Info(ctx, "Food search successful", logger.Int("total_foods", len(foods)))
	return foods, nil
}

func (s *FoodService) GetFoodByID(ctx context.Context, id string) (*domain.FoodItem, error) {
	s.logger.Info(ctx, "Getting food by ID", logger.String("food_id", id))
	foodID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error(ctx, "Failed to convert food ID to object ID", logger.Error(err))
		return nil, fmt.Errorf("failed to convert food ID to object ID: %w", err)
	}
	food, err := s.foodRepo.GetByID(ctx, foodID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get food by ID", logger.Error(err))
		return nil, fmt.Errorf("failed to get food by ID: %w", err)
	}
	s.logger.Info(ctx, "Food retrieved successfully", logger.String("food_id", food.ID.Hex()))
	return food, nil
}
