package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"nutrient_be/internal/domain"
	"nutrient_be/internal/pkg/logger"
)

// MealPlanRepository defines the interface for meal plan data operations used by MealPlanService
type MealPlanRepository interface {
	Create(ctx context.Context, plan *domain.MealPlan) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.MealPlan, error)
	GetByUser(ctx context.Context, userID primitive.ObjectID, planType string, limit, offset int) ([]*domain.MealPlan, error)
	GetByUserAndDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate string) ([]*domain.MealPlan, error)
	Update(ctx context.Context, plan *domain.MealPlan) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	UpdateMealCompletion(ctx context.Context, planID primitive.ObjectID, mealID string, isCompleted bool) error
}

// MealPlanTemplateRepository defines the interface for meal template data operations used by MealPlanService
type MealPlanTemplateRepository interface {
	Create(ctx context.Context, template *domain.MealTemplate) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.MealTemplate, error)
	GetByUser(ctx context.Context, userID primitive.ObjectID, mealType string, limit, offset int) ([]*domain.MealTemplate, error)
	GetPublicTemplates(ctx context.Context, mealType string, limit, offset int) ([]*domain.MealTemplate, error)
	Update(ctx context.Context, template *domain.MealTemplate) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// MealPlanService handles meal plan business logic
type MealPlanService struct {
	mealPlanRepo     MealPlanRepository
	mealTemplateRepo MealPlanTemplateRepository
	logger           logger.Logger
}

// NewMealPlanService creates a new meal plan service
func NewMealPlanService(mealPlanRepo MealPlanRepository, mealTemplateRepo MealPlanTemplateRepository, log logger.Logger) *MealPlanService {
	return &MealPlanService{
		mealPlanRepo:     mealPlanRepo,
		mealTemplateRepo: mealTemplateRepo,
		logger:           log,
	}
}
