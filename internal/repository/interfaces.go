package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"nutrient_be/internal/domain"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// FoodRepository defines the interface for food data operations
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

// MealTemplateRepository defines the interface for meal template data operations
type MealTemplateRepository interface {
	Create(ctx context.Context, template *domain.MealTemplate) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.MealTemplate, error)
	GetByUser(ctx context.Context, userID primitive.ObjectID, mealType string, limit, offset int) ([]*domain.MealTemplate, error)
	GetPublicTemplates(ctx context.Context, mealType string, limit, offset int) ([]*domain.MealTemplate, error)
	Update(ctx context.Context, template *domain.MealTemplate) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// MealPlanRepository defines the interface for meal plan data operations
type MealPlanRepository interface {
	Create(ctx context.Context, plan *domain.MealPlan) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.MealPlan, error)
	GetByUser(ctx context.Context, userID primitive.ObjectID, planType string, limit, offset int) ([]*domain.MealPlan, error)
	GetByUserAndDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate string) ([]*domain.MealPlan, error)
	Update(ctx context.Context, plan *domain.MealPlan) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	UpdateMealCompletion(ctx context.Context, planID primitive.ObjectID, mealID string, isCompleted bool) error
}

// ShoppingListRepository defines the interface for shopping list data operations
type ShoppingListRepository interface {
	Create(ctx context.Context, list *domain.ShoppingList) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.ShoppingList, error)
	GetByUser(ctx context.Context, userID primitive.ObjectID, limit, offset int) ([]*domain.ShoppingList, error)
	GetByMealPlan(ctx context.Context, mealPlanID primitive.ObjectID) (*domain.ShoppingList, error)
	Update(ctx context.Context, list *domain.ShoppingList) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	ToggleItemChecked(ctx context.Context, listID primitive.ObjectID, itemID primitive.ObjectID, checked bool) error
}
