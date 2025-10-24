package service

import (
	"nutrient_be/internal/pkg/logger"
	"nutrient_be/internal/repository"
)

// FoodService handles food-related business logic
type FoodService struct {
	foodRepo repository.FoodRepository
	logger   logger.Logger
}

// NewFoodService creates a new food service
func NewFoodService(foodRepo repository.FoodRepository, log logger.Logger) *FoodService {
	return &FoodService{
		foodRepo: foodRepo,
		logger:   log,
	}
}

// MealService handles meal template business logic
type MealService struct {
	mealTemplateRepo repository.MealTemplateRepository
	foodRepo         repository.FoodRepository
	logger           logger.Logger
}

// NewMealService creates a new meal service
func NewMealService(mealTemplateRepo repository.MealTemplateRepository, foodRepo repository.FoodRepository, log logger.Logger) *MealService {
	return &MealService{
		mealTemplateRepo: mealTemplateRepo,
		foodRepo:         foodRepo,
		logger:           log,
	}
}

// MealPlanService handles meal plan business logic
type MealPlanService struct {
	mealPlanRepo     repository.MealPlanRepository
	mealTemplateRepo repository.MealTemplateRepository
	logger           logger.Logger
}

// NewMealPlanService creates a new meal plan service
func NewMealPlanService(mealPlanRepo repository.MealPlanRepository, mealTemplateRepo repository.MealTemplateRepository, log logger.Logger) *MealPlanService {
	return &MealPlanService{
		mealPlanRepo:     mealPlanRepo,
		mealTemplateRepo: mealTemplateRepo,
		logger:           log,
	}
}

// ShoppingService handles shopping list business logic
type ShoppingService struct {
	shoppingRepo repository.ShoppingListRepository
	mealPlanRepo repository.MealPlanRepository
	logger       logger.Logger
}

// NewShoppingService creates a new shopping service
func NewShoppingService(shoppingRepo repository.ShoppingListRepository, mealPlanRepo repository.MealPlanRepository, log logger.Logger) *ShoppingService {
	return &ShoppingService{
		shoppingRepo: shoppingRepo,
		mealPlanRepo: mealPlanRepo,
		logger:       log,
	}
}

// ReportService handles report business logic
type ReportService struct {
	mealPlanRepo repository.MealPlanRepository
	logger       logger.Logger
}

// NewReportService creates a new report service
func NewReportService(mealPlanRepo repository.MealPlanRepository, log logger.Logger) *ReportService {
	return &ReportService{
		mealPlanRepo: mealPlanRepo,
		logger:       log,
	}
}
