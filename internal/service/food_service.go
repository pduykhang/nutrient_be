package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"nutrient_be/internal/domain"
	"nutrient_be/internal/pkg/logger"
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
	foodRepo FoodRepository
	logger   logger.Logger
}

// NewFoodService creates a new food service
func NewFoodService(foodRepo FoodRepository, log logger.Logger) *FoodService {
	return &FoodService{
		foodRepo: foodRepo,
		logger:   log,
	}
}
