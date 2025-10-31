package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"nutrient_be/internal/domain"
	"nutrient_be/internal/pkg/logger"
)

// MealTemplateRepository defines the interface for meal template data operations used by MealService
type MealTemplateRepository interface {
	Create(ctx context.Context, template *domain.MealTemplate) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.MealTemplate, error)
	GetByUser(ctx context.Context, userID primitive.ObjectID, mealType string, limit, offset int) ([]*domain.MealTemplate, error)
	GetPublicTemplates(ctx context.Context, mealType string, limit, offset int) ([]*domain.MealTemplate, error)
	Update(ctx context.Context, template *domain.MealTemplate) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// MealFoodRepository defines the interface for food data operations used by MealService
type MealFoodRepository interface {
	Create(ctx context.Context, food *domain.FoodItem) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.FoodItem, error)
	Search(ctx context.Context, query string, userID primitive.ObjectID, limit, offset int) ([]*domain.FoodItem, error)
	GetByCategory(ctx context.Context, category string, userID primitive.ObjectID, limit, offset int) ([]*domain.FoodItem, error)
	GetByUser(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*domain.FoodItem, error)
	Update(ctx context.Context, food *domain.FoodItem) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetPublicFoods(ctx context.Context, limit, offset int) ([]*domain.FoodItem, error)
}

// MealService handles meal template business logic
type MealService struct {
	mealTemplateRepo MealTemplateRepository
	foodRepo         MealFoodRepository
	logger           logger.Logger
}

// NewMealService creates a new meal service
func NewMealService(mealTemplateRepo MealTemplateRepository, foodRepo MealFoodRepository, log logger.Logger) *MealService {
	return &MealService{
		mealTemplateRepo: mealTemplateRepo,
		foodRepo:         foodRepo,
		logger:           log,
	}
}
