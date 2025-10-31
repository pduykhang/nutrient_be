package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"nutrient_be/internal/domain"
	"nutrient_be/internal/pkg/logger"
)

// ShoppingListRepository defines the interface for shopping list data operations used by ShoppingService
type ShoppingListRepository interface {
	Create(ctx context.Context, list *domain.ShoppingList) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.ShoppingList, error)
	GetByUser(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*domain.ShoppingList, error)
	GetByMealPlan(ctx context.Context, mealPlanID primitive.ObjectID) (*domain.ShoppingList, error)
	Update(ctx context.Context, list *domain.ShoppingList) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	ToggleItemChecked(ctx context.Context, listID primitive.ObjectID, itemID primitive.ObjectID, checked bool) error
}

// ShoppingMealPlanRepository defines the interface for meal plan data operations used by ShoppingService
type ShoppingMealPlanRepository interface {
	Create(ctx context.Context, plan *domain.MealPlan) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.MealPlan, error)
	GetByUser(ctx context.Context, userID primitive.ObjectID, planType string, limit, offset int) ([]*domain.MealPlan, error)
	GetByUserAndDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate string) ([]*domain.MealPlan, error)
	Update(ctx context.Context, plan *domain.MealPlan) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	UpdateMealCompletion(ctx context.Context, planID primitive.ObjectID, mealID string, isCompleted bool) error
}

// ShoppingService handles shopping list business logic
type ShoppingService struct {
	shoppingRepo ShoppingListRepository
	mealPlanRepo ShoppingMealPlanRepository
	logger       logger.Logger
}

// NewShoppingService creates a new shopping service
func NewShoppingService(shoppingRepo ShoppingListRepository, mealPlanRepo ShoppingMealPlanRepository, log logger.Logger) *ShoppingService {
	return &ShoppingService{
		shoppingRepo: shoppingRepo,
		mealPlanRepo: mealPlanRepo,
		logger:       log,
	}
}
