package validator

import (
	"fmt"
	"strings"

	"nutrient_be/internal/dto/request"
	"nutrient_be/internal/pkg/logger"
)

// MealValidator handles meal template validation
type MealValidator struct {
	maxNameLength        int
	maxDescriptionLength int
	maxTagsPerMeal       int
	maxFoodItemsPerMeal  int
	logger               logger.Logger
}

// NewMealValidator creates a new meal validator
func NewMealValidator(logger logger.Logger) *MealValidator {
	return &MealValidator{
		maxNameLength:        100,
		maxDescriptionLength: 500,
		maxTagsPerMeal:       10,
		maxFoodItemsPerMeal:  50,
		logger:               logger,
	}
}

// ValidateCreateTemplateRequest validates a CreateMealTemplateRequest
func (v *MealValidator) ValidateCreateTemplateRequest(req *request.CreateMealTemplateRequest) error {
	// 1. Validate Name
	if err := v.validateName(req.Name); err != nil {
		return fmt.Errorf("name validation failed: %w", err)
	}

	// 2. Validate Description (optional)
	if req.Description != "" && len(req.Description) > v.maxDescriptionLength {
		return fmt.Errorf("description exceeds maximum length (%d chars)", v.maxDescriptionLength)
	}

	// 3. Validate Meal Type
	if err := v.validateMealType(req.MealType); err != nil {
		return fmt.Errorf("meal type validation failed: %w", err)
	}

	// 4. Validate Food Items
	if err := v.validateFoodItems(req.FoodItems); err != nil {
		return fmt.Errorf("food items validation failed: %w", err)
	}

	// 5. Validate Tags
	if err := v.validateTags(req.Tags); err != nil {
		return fmt.Errorf("tags validation failed: %w", err)
	}

	return nil
}

// validateName validates meal template name
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

// validateMealType validates meal type value
func (v *MealValidator) validateMealType(mealType string) error {
	validTypes := map[string]bool{
		"breakfast": true,
		"lunch":     true,
		"dinner":    true,
		"snack":     true,
	}

	if !validTypes[mealType] {
		return fmt.Errorf("invalid meal type '%s'. Valid types: breakfast, lunch, dinner, snack", mealType)
	}

	return nil
}

// validateFoodItems validates food items in meal template
func (v *MealValidator) validateFoodItems(items []request.MealTemplateFoodItemRequest) error {
	if len(items) == 0 {
		return fmt.Errorf("at least one food item is required")
	}

	if len(items) > v.maxFoodItemsPerMeal {
		return fmt.Errorf("too many food items (%d). Maximum allowed: %d", len(items), v.maxFoodItemsPerMeal)
	}

	for i, item := range items {
		// Validate FoodItemID
		if item.FoodItemID == "" {
			return fmt.Errorf("food item %d: foodItemId is required", i+1)
		}

		// Validate ServingUnit
		validUnits := map[string]bool{
			"gram":  true,
			"kg":    true,
			"piece": true,
			"cup":   true,
			"ml":    true,
			"box":   true,
		}
		if !validUnits[item.ServingUnit] {
			return fmt.Errorf("food item %d: invalid serving unit '%s'", i+1, item.ServingUnit)
		}

		// Validate Amount
		if item.Amount <= 0 {
			return fmt.Errorf("food item %d: amount must be greater than 0", i+1)
		}

		// Validate reasonable amount limit
		if item.Amount > 10000 {
			return fmt.Errorf("food item %d: amount (%.2f) is unreasonably large", i+1, item.Amount)
		}
	}

	return nil
}

// validateTags validates meal template tags
func (v *MealValidator) validateTags(tags []string) error {
	if len(tags) > v.maxTagsPerMeal {
		return fmt.Errorf("too many tags (%d). Maximum allowed: %d", len(tags), v.maxTagsPerMeal)
	}

	for i, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			return fmt.Errorf("tag %d cannot be empty", i+1)
		}
		if len(trimmed) > 50 {
			return fmt.Errorf("tag %d exceeds maximum length (50 chars)", i+1)
		}
	}

	return nil
}
