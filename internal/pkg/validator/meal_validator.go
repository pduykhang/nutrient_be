package validator

import (
	"context"
	"fmt"
	"strings"

	"nutrient_be/internal/dto/request"
	"nutrient_be/internal/pkg/logger"
)

// MealValidator handles meal template data validation
type MealValidator struct {
	logger          logger.Logger
	maxNameLength   int
	maxDescriptionLength int
	maxTags         int
	maxTagLength    int
}

// NewMealValidator creates a new meal validator with default rules
func NewMealValidator(logger logger.Logger) *MealValidator {
	return &MealValidator{
		logger:               logger,
		maxNameLength:         200,
		maxDescriptionLength: 1000,
		maxTags:               20,
		maxTagLength:          50,
	}
}

// ValidateCreateRequest validates a CreateMealTemplateRequest
func (v *MealValidator) ValidateCreateRequest(ctx context.Context, req *request.CreateMealTemplateRequest) error {
	// Validate name
	if err := v.validateName(req.Name); err != nil {
		return fmt.Errorf("name validation failed: %w", err)
	}

	// Validate description (optional)
	if req.Description != "" {
		if err := v.validateDescription(req.Description); err != nil {
			return fmt.Errorf("description validation failed: %w", err)
		}
	}

	// Validate meal type
	if err := v.validateMealType(req.MealType); err != nil {
		return fmt.Errorf("meal type validation failed: %w", err)
	}

	// Validate food items
	if err := v.validateFoodItems(req.FoodItems); err != nil {
		return fmt.Errorf("food items validation failed: %w", err)
	}

	// Validate tags (optional)
	if len(req.Tags) > 0 {
		if err := v.validateTags(req.Tags); err != nil {
			return fmt.Errorf("tags validation failed: %w", err)
		}
	}

	return nil
}

// ValidateUpdateRequest validates an UpdateMealTemplateRequest
func (v *MealValidator) ValidateUpdateRequest(ctx context.Context, req *request.UpdateMealTemplateRequest) error {
	// Validate name if provided
	if req.Name != "" {
		if err := v.validateName(req.Name); err != nil {
			return fmt.Errorf("name validation failed: %w", err)
		}
	}

	// Validate description if provided
	if req.Description != "" {
		if err := v.validateDescription(req.Description); err != nil {
			return fmt.Errorf("description validation failed: %w", err)
		}
	}

	// Validate meal type if provided
	if req.MealType != "" {
		if err := v.validateMealType(req.MealType); err != nil {
			return fmt.Errorf("meal type validation failed: %w", err)
		}
	}

	// Validate food items if provided
	if req.FoodItems != nil && len(req.FoodItems) > 0 {
		if err := v.validateFoodItems(req.FoodItems); err != nil {
			return fmt.Errorf("food items validation failed: %w", err)
		}
	}

	// Validate tags if provided
	if req.Tags != nil && len(req.Tags) > 0 {
		if err := v.validateTags(req.Tags); err != nil {
			return fmt.Errorf("tags validation failed: %w", err)
		}
	}

	return nil
}

// ValidateAddFoodRequest validates an AddFoodToTemplateRequest
func (v *MealValidator) ValidateAddFoodRequest(ctx context.Context, req *request.AddFoodToTemplateRequest) error {
	// Validate food items
	if err := v.validateFoodItems(req.FoodItems); err != nil {
		return fmt.Errorf("food items validation failed: %w", err)
	}

	return nil
}

// validateName validates template name
func (v *MealValidator) validateName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(trimmed) > v.maxNameLength {
		return fmt.Errorf("name exceeds maximum length (%d chars)", v.maxNameLength)
	}
	return nil
}

// validateDescription validates template description
func (v *MealValidator) validateDescription(description string) error {
	if len(description) > v.maxDescriptionLength {
		return fmt.Errorf("description exceeds maximum length (%d chars)", v.maxDescriptionLength)
	}
	return nil
}

// validateMealType validates meal type
func (v *MealValidator) validateMealType(mealType string) error {
	validTypes := map[string]bool{
		"breakfast": true,
		"lunch":     true,
		"dinner":    true,
		"snack":     true,
	}

	if !validTypes[mealType] {
		return fmt.Errorf("invalid meal type '%s', must be one of: breakfast, lunch, dinner, snack", mealType)
	}

	return nil
}

// validateFoodItems validates food items array
func (v *MealValidator) validateFoodItems(foodItems []request.MealTemplateFoodItemRequest) error {
	if len(foodItems) == 0 {
		return fmt.Errorf("at least one food item is required")
	}

	// Check for duplicate food items (same food ID and serving unit)
	foodItemMap := make(map[string]bool)
	for i, item := range foodItems {
		// Validate food item ID
		if strings.TrimSpace(item.FoodItemID) == "" {
			return fmt.Errorf("food item %d: foodItemId is required", i+1)
		}

		// Validate serving unit
		if strings.TrimSpace(item.ServingUnit) == "" {
			return fmt.Errorf("food item %d: servingUnit is required", i+1)
		}

		// Validate amount
		if item.Amount <= 0 {
			return fmt.Errorf("food item %d: amount must be greater than 0", i+1)
		}

		// Check for duplicates
		key := fmt.Sprintf("%s:%s", item.FoodItemID, item.ServingUnit)
		if foodItemMap[key] {
			return fmt.Errorf("food item %d: duplicate food item with same foodItemId and servingUnit", i+1)
		}
		foodItemMap[key] = true
	}

	return nil
}

// validateTags validates tags array
func (v *MealValidator) validateTags(tags []string) error {
	if len(tags) > v.maxTags {
		return fmt.Errorf("maximum number of tags is %d", v.maxTags)
	}

	for i, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			return fmt.Errorf("tag %d: cannot be empty", i+1)
		}
		if len(trimmed) > v.maxTagLength {
			return fmt.Errorf("tag %d: exceeds maximum length (%d chars)", i+1, v.maxTagLength)
		}
	}

	return nil
}
