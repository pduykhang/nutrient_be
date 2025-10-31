package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"nutrient_be/internal/domain"
	"nutrient_be/internal/pkg/logger"
)

// ReportMealPlanRepository defines the interface for meal plan data operations used by ReportService
type ReportMealPlanRepository interface {
	Create(ctx context.Context, plan *domain.MealPlan) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.MealPlan, error)
	GetByUser(ctx context.Context, userID primitive.ObjectID, planType string, limit, offset int) ([]*domain.MealPlan, error)
	GetByUserAndDateRange(ctx context.Context, userID primitive.ObjectID, startDate, endDate string) ([]*domain.MealPlan, error)
	Update(ctx context.Context, plan *domain.MealPlan) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	UpdateMealCompletion(ctx context.Context, planID primitive.ObjectID, mealID string, isCompleted bool) error
}

// ReportService handles report business logic
type ReportService struct {
	mealPlanRepo ReportMealPlanRepository
	logger       logger.Logger
}

// NewReportService creates a new report service
func NewReportService(mealPlanRepo ReportMealPlanRepository, log logger.Logger) *ReportService {
	return &ReportService{
		mealPlanRepo: mealPlanRepo,
		logger:       log,
	}
}
