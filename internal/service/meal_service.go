package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"nutrient_be/internal/domain"
	"nutrient_be/internal/dto/request"
	"nutrient_be/internal/pkg/calculator"
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

// CreateTemplate creates a new meal template with food items and calculates totals
func (s *MealService) CreateTemplate(ctx context.Context, userID string, req *request.CreateMealTemplateRequest) (*domain.MealTemplate, error) {
	s.logger.Info(ctx, "Creating meal template", logger.String("name", req.Name), logger.String("meal_type", req.MealType))

	// Convert userID to ObjectID
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.logger.Error(ctx, "Invalid user ID", logger.Error(err))
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Process food items and calculate nutrients
	foodItems, totalCalories, totalMacros, totalMicros, err := s.processFoodItems(ctx, req.FoodItems)
	if err != nil {
		s.logger.Error(ctx, "Failed to process food items", logger.Error(err))
		return nil, fmt.Errorf("failed to process food items: %w", err)
	}

	// Create template
	template := &domain.MealTemplate{
		ID:            primitive.NewObjectID(),
		UserID:        userIDObj,
		Name:          req.Name,
		Description:   req.Description,
		MealType:      req.MealType,
		FoodItems:     foodItems,
		TotalCalories: totalCalories,
		TotalMacros:   totalMacros,
		TotalMicros:   totalMicros,
		Tags:          req.Tags,
		IsPublic:      req.IsPublic,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save to database
	if err := s.mealTemplateRepo.Create(ctx, template); err != nil {
		s.logger.Error(ctx, "Failed to create meal template", logger.Error(err))
		return nil, fmt.Errorf("failed to create meal template: %w", err)
	}

	s.logger.Info(ctx, "Meal template created successfully", logger.String("template_id", template.ID.Hex()))
	return template, nil
}

// AddFoodToTemplate adds food items to an existing template and recalculates totals
func (s *MealService) AddFoodToTemplate(ctx context.Context, userID string, templateID string, req *request.AddFoodToTemplateRequest) (*domain.MealTemplate, error) {
	s.logger.Info(ctx, "Adding food items to template", logger.String("template_id", templateID))

	// Convert IDs to ObjectID
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.logger.Error(ctx, "Invalid user ID", logger.Error(err))
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	templateIDObj, err := primitive.ObjectIDFromHex(templateID)
	if err != nil {
		s.logger.Error(ctx, "Invalid template ID", logger.Error(err))
		return nil, fmt.Errorf("invalid template ID: %w", err)
	}

	// Get existing template
	template, err := s.mealTemplateRepo.GetByID(ctx, templateIDObj)
	if err != nil {
		s.logger.Error(ctx, "Failed to get template", logger.Error(err))
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// Verify ownership
	if template.UserID != userIDObj {
		s.logger.Error(ctx, "User does not own template")
		return nil, fmt.Errorf("template not found or access denied")
	}

	// Process new food items
	newFoodItems, newCalories, newMacros, newMicros, err := s.processFoodItems(ctx, req.FoodItems)
	if err != nil {
		s.logger.Error(ctx, "Failed to process food items", logger.Error(err))
		return nil, fmt.Errorf("failed to process food items: %w", err)
	}

	// Add new food items to existing ones
	template.FoodItems = append(template.FoodItems, newFoodItems...)

	// Recalculate totals
	template.TotalCalories += newCalories
	template.TotalMacros = calculator.SumMacros(template.TotalMacros, newMacros)
	template.TotalMicros = calculator.SumMicros(template.TotalMicros, newMicros)
	template.UpdatedAt = time.Now()

	// Update in database
	if err := s.mealTemplateRepo.Update(ctx, template); err != nil {
		s.logger.Error(ctx, "Failed to update template", logger.Error(err))
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	s.logger.Info(ctx, "Food items added to template successfully")
	return template, nil
}

// GetTemplate retrieves a meal template with detailed macro and micro information
func (s *MealService) GetTemplate(ctx context.Context, userID string, templateID string) (*domain.MealTemplate, error) {
	s.logger.Info(ctx, "Getting meal template", logger.String("template_id", templateID))

	// Convert IDs to ObjectID
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.logger.Error(ctx, "Invalid user ID", logger.Error(err))
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	templateIDObj, err := primitive.ObjectIDFromHex(templateID)
	if err != nil {
		s.logger.Error(ctx, "Invalid template ID", logger.Error(err))
		return nil, fmt.Errorf("invalid template ID: %w", err)
	}

	// Get template
	template, err := s.mealTemplateRepo.GetByID(ctx, templateIDObj)
	if err != nil {
		s.logger.Error(ctx, "Failed to get template", logger.Error(err))
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// Verify access: user owns it or it's public
	if template.UserID != userIDObj && !template.IsPublic {
		s.logger.Error(ctx, "User does not have access to template")
		return nil, fmt.Errorf("template not found or access denied")
	}

	s.logger.Info(ctx, "Meal template retrieved successfully")
	return template, nil
}

// ListTemplates lists meal templates for a user
func (s *MealService) ListTemplates(ctx context.Context, userID string, mealType string, limit, offset int) ([]*domain.MealTemplate, error) {
	s.logger.Info(ctx, "Listing meal templates", logger.String("meal_type", mealType))

	// Convert userID to ObjectID
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.logger.Error(ctx, "Invalid user ID", logger.Error(err))
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Get templates
	templates, err := s.mealTemplateRepo.GetByUser(ctx, userIDObj, mealType, limit, offset)
	if err != nil {
		s.logger.Error(ctx, "Failed to list templates", logger.Error(err))
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	s.logger.Info(ctx, "Meal templates listed successfully", logger.Int("count", len(templates)))
	return templates, nil
}

// UpdateTemplate updates a meal template
func (s *MealService) UpdateTemplate(ctx context.Context, userID string, templateID string, req *request.UpdateMealTemplateRequest) (*domain.MealTemplate, error) {
	s.logger.Info(ctx, "Updating meal template", logger.String("template_id", templateID))

	// Convert IDs to ObjectID
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.logger.Error(ctx, "Invalid user ID", logger.Error(err))
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	templateIDObj, err := primitive.ObjectIDFromHex(templateID)
	if err != nil {
		s.logger.Error(ctx, "Invalid template ID", logger.Error(err))
		return nil, fmt.Errorf("invalid template ID: %w", err)
	}

	// Get existing template
	template, err := s.mealTemplateRepo.GetByID(ctx, templateIDObj)
	if err != nil {
		s.logger.Error(ctx, "Failed to get template", logger.Error(err))
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// Verify ownership
	if template.UserID != userIDObj {
		s.logger.Error(ctx, "User does not own template")
		return nil, fmt.Errorf("template not found or access denied")
	}

	// Update fields if provided
	if req.Name != "" {
		template.Name = req.Name
	}
	if req.Description != "" {
		template.Description = req.Description
	}
	if req.MealType != "" {
		template.MealType = req.MealType
	}
	if req.Tags != nil {
		template.Tags = req.Tags
	}
	if req.IsPublic != nil {
		template.IsPublic = *req.IsPublic
	}

	// Update food items if provided
	if req.FoodItems != nil && len(req.FoodItems) > 0 {
		foodItems, totalCalories, totalMacros, totalMicros, err := s.processFoodItems(ctx, req.FoodItems)
		if err != nil {
			s.logger.Error(ctx, "Failed to process food items", logger.Error(err))
			return nil, fmt.Errorf("failed to process food items: %w", err)
		}
		template.FoodItems = foodItems
		template.TotalCalories = totalCalories
		template.TotalMacros = totalMacros
		template.TotalMicros = totalMicros
	}

	template.UpdatedAt = time.Now()

	// Update in database
	if err := s.mealTemplateRepo.Update(ctx, template); err != nil {
		s.logger.Error(ctx, "Failed to update template", logger.Error(err))
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	s.logger.Info(ctx, "Meal template updated successfully")
	return template, nil
}

// DeleteTemplate deletes a meal template
func (s *MealService) DeleteTemplate(ctx context.Context, userID string, templateID string) error {
	s.logger.Info(ctx, "Deleting meal template", logger.String("template_id", templateID))

	// Convert IDs to ObjectID
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.logger.Error(ctx, "Invalid user ID", logger.Error(err))
		return fmt.Errorf("invalid user ID: %w", err)
	}

	templateIDObj, err := primitive.ObjectIDFromHex(templateID)
	if err != nil {
		s.logger.Error(ctx, "Invalid template ID", logger.Error(err))
		return fmt.Errorf("invalid template ID: %w", err)
	}

	// Get template to verify ownership
	template, err := s.mealTemplateRepo.GetByID(ctx, templateIDObj)
	if err != nil {
		s.logger.Error(ctx, "Failed to get template", logger.Error(err))
		return fmt.Errorf("failed to get template: %w", err)
	}

	// Verify ownership
	if template.UserID != userIDObj {
		s.logger.Error(ctx, "User does not own template")
		return fmt.Errorf("template not found or access denied")
	}

	// Delete template
	if err := s.mealTemplateRepo.Delete(ctx, templateIDObj); err != nil {
		s.logger.Error(ctx, "Failed to delete template", logger.Error(err))
		return fmt.Errorf("failed to delete template: %w", err)
	}

	s.logger.Info(ctx, "Meal template deleted successfully")
	return nil
}

// processFoodItems processes food items from request, calculates nutrients, and returns totals
func (s *MealService) processFoodItems(
	ctx context.Context,
	foodItemReqs []request.MealTemplateFoodItemRequest,
) ([]domain.MealTemplateFoodItem, float64, domain.MacroNutrients, domain.MicroNutrients, error) {
	foodItems := make([]domain.MealTemplateFoodItem, 0, len(foodItemReqs))
	var totalCalories float64
	var allMacros []domain.MacroNutrients
	var allMicros []domain.MicroNutrients

	for _, foodItemReq := range foodItemReqs {
		// Convert food item ID to ObjectID
		foodIDObj, err := primitive.ObjectIDFromHex(foodItemReq.FoodItemID)
		if err != nil {
			return nil, 0, domain.MacroNutrients{}, domain.MicroNutrients{},
				fmt.Errorf("invalid food item ID '%s': %w", foodItemReq.FoodItemID, err)
		}

		// Get food item from repository
		food, err := s.foodRepo.GetByID(ctx, foodIDObj)
		if err != nil {
			return nil, 0, domain.MacroNutrients{}, domain.MicroNutrients{},
				fmt.Errorf("food item not found: %w", err)
		}

		// Calculate nutrients for the specified serving
		calories, macros, micros, err := calculator.CalculateNutrientsForServing(
			food,
			foodItemReq.ServingUnit,
			foodItemReq.Amount,
		)
		if err != nil {
			return nil, 0, domain.MacroNutrients{}, domain.MicroNutrients{},
				fmt.Errorf("failed to calculate nutrients for food '%s': %w", foodItemReq.FoodItemID, err)
		}

		// Get food name (prefer English, fallback to first available)
		foodName := food.Name["en"]
		if foodName == "" {
			for _, name := range food.Name {
				foodName = name
				break
			}
		}

		// Create meal template food item
		mealFoodItem := domain.MealTemplateFoodItem{
			FoodItemID:  foodIDObj,
			FoodName:    foodName,
			ServingUnit: foodItemReq.ServingUnit,
			Amount:      foodItemReq.Amount,
			Calories:    calories,
			Macros:      macros,
			Micros:      micros,
		}

		foodItems = append(foodItems, mealFoodItem)
		totalCalories += calories
		allMacros = append(allMacros, macros)
		allMicros = append(allMicros, micros)
	}

	// Calculate totals
	totalMacros := calculator.SumMacros(allMacros...)
	totalMicros := calculator.SumMicros(allMicros...)

	return foodItems, totalCalories, totalMacros, totalMicros, nil
}
