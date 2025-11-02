package validator

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"nutrient_be/internal/dto/request"
	"nutrient_be/internal/pkg/logger"
)

// FoodValidator handles food data validation
type FoodValidator struct {
	logger logger.Logger
	// Configuration for validation rules
	maxNameLength        int
	maxDescriptionLength int
	maxCalories          float64
	maxMacroValue        float64
	caloriesTolerance    float64
}

// NewFoodValidator creates a new food validator with default rules
func NewFoodValidator(logger logger.Logger) *FoodValidator {
	return &FoodValidator{
		logger:               logger,
		maxNameLength:        200,
		maxDescriptionLength: 1000,
		maxCalories:          1000,
		maxMacroValue:        100,
		caloriesTolerance:    10, // Allow ±10 calories difference
	}
}

// ValidateCreateRequest validates a CreateFoodRequest
func (v *FoodValidator) ValidateCreateRequest(ctx context.Context, req *request.CreateFoodRequest) error {
	// 1. Validate Name
	if err := v.validateName(req.Name); err != nil {
		return fmt.Errorf("name validation failed: %w", err)
	}

	// 2. Validate Description (optional field)
	if req.Description != nil {
		descRaw := req.Description.GetRaw()
		if len(descRaw) > 0 {
			if err := v.validateDescription(req.Description); err != nil {
				return fmt.Errorf("description validation failed: %w", err)
			}
		}
	}

	// 3. Validate Nutrition Values
	if err := v.validateNutrition(req); err != nil {
		return fmt.Errorf("nutrition validation failed: %w", err)
	}

	// 4. Validate Serving Sizes
	if err := v.validateServingSizes(ctx, req.ServingSizes); err != nil {
		return fmt.Errorf("serving sizes validation failed: %w", err)
	}

	// 5. Validate Calories Consistency
	if err := v.validateCaloriesConsistency(req); err != nil {
		return fmt.Errorf("calories consistency validation failed: %w", err)
	}

	// 6. Validate Image URL (optional)
	if req.ImageURL != "" {
		if err := v.validateImageURL(req.ImageURL); err != nil {
			return fmt.Errorf("image URL validation failed: %w", err)
		}
	}

	return nil
}

// validateName validates multi-language name
func (v *FoodValidator) validateName(name request.MultiLanguage) error {
	raw := name.GetRaw()
	if len(raw) == 0 {
		return fmt.Errorf("name must have at least one language")
	}

	// Must have English
	if name.Get("en") == "" {
		return fmt.Errorf("name must have English (en) translation")
	}

	// Validate each language value
	for lang, value := range raw {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return fmt.Errorf("name value for language '%s' cannot be empty", lang)
		}
		if len(trimmed) > v.maxNameLength {
			return fmt.Errorf("name for language '%s' exceeds maximum length (%d chars)", lang, v.maxNameLength)
		}
	}

	return nil
}

// validateDescription validates multi-language description
func (v *FoodValidator) validateDescription(desc request.MultiLanguage) error {
	raw := desc.GetRaw()
	if len(raw) == 0 {
		return nil // Optional field
	}

	// Validate each language value
	for lang, value := range raw {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" && len(trimmed) > v.maxDescriptionLength {
			return fmt.Errorf("description for language '%s' exceeds maximum length (%d chars)", lang, v.maxDescriptionLength)
		}
	}

	return nil
}

// validateNutrition validates nutrition values
func (v *FoodValidator) validateNutrition(req *request.CreateFoodRequest) error {
	macros := req.Macros

	// At least one macro must be > 0
	if macros.Protein == 0 && macros.Carbohydrates == 0 && macros.Fat == 0 {
		return fmt.Errorf("at least one macro nutrient (protein, carbs, or fat) must be greater than 0")
	}

	// Validate individual macro values
	if macros.Protein < 0 || macros.Protein > v.maxMacroValue {
		return fmt.Errorf("protein must be between 0 and %.2fg per 100g", v.maxMacroValue)
	}
	if macros.Carbohydrates < 0 || macros.Carbohydrates > v.maxMacroValue {
		return fmt.Errorf("carbohydrates must be between 0 and %.2fg per 100g", v.maxMacroValue)
	}
	if macros.Fat < 0 || macros.Fat > v.maxMacroValue {
		return fmt.Errorf("fat must be between 0 and %.2fg per 100g", v.maxMacroValue)
	}
	if macros.Fiber < 0 || macros.Fiber > v.maxMacroValue {
		return fmt.Errorf("fiber must be between 0 and %.2fg per 100g", v.maxMacroValue)
	}
	if macros.Sugar < 0 {
		return fmt.Errorf("sugar cannot be negative")
	}

	// Validate total macros
	totalMacros := macros.Protein + macros.Carbohydrates + macros.Fat + macros.Fiber
	if totalMacros > 999 {
		return fmt.Errorf("total macros exceed maximum (999g per 100g)")
	}

	// Validate calories
	if req.Calories < 0 {
		return fmt.Errorf("calories cannot be negative")
	}
	if req.Calories > v.maxCalories {
		return fmt.Errorf("calories exceed maximum (%.2f per 100g)", v.maxCalories)
	}

	// Validate micros if provided (optional fields, but if set must be >= 0)
	// Note: Micros can be 0 or omitted, but cannot be negative
	if req.Micros.VitaminA < 0 {
		return fmt.Errorf("vitaminA cannot be negative")
	}
	if req.Micros.VitaminC < 0 {
		return fmt.Errorf("vitaminC cannot be negative")
	}
	if req.Micros.Calcium < 0 {
		return fmt.Errorf("calcium cannot be negative")
	}
	if req.Micros.Iron < 0 {
		return fmt.Errorf("iron cannot be negative")
	}
	if req.Micros.Sodium < 0 {
		return fmt.Errorf("sodium cannot be negative")
	}
	if req.Micros.Potassium < 0 {
		return fmt.Errorf("potassium cannot be negative")
	}

	return nil
}

// validateServingSizes validates serving sizes
func (v *FoodValidator) validateServingSizes(ctx context.Context, sizes []request.ServingSizeRequest) error {
	if len(sizes) == 0 {
		return fmt.Errorf("at least one serving size is required")
	}

	validUnits := map[string]bool{
		"gram":  true,
		"kg":    true,
		"piece": true,
		"cup":   true,
		"ml":    true,
		"box":   true,
	}

	hasGramBase := false

	for i, size := range sizes {
		// Validate unit
		if !validUnits[size.Unit] {
			return fmt.Errorf("serving size %d: invalid unit '%s'. Valid units: gram, kg, piece, cup, ml, box", i+1, size.Unit)
		}

		// Validate amount
		if size.Amount <= 0 {
			return fmt.Errorf("serving size %d: amount must be greater than 0", i+1)
		}

		// Validate gramEquivalent
		if size.GramEquivalent <= 0 {
			return fmt.Errorf("serving size %d: gramEquivalent must be greater than 0", i+1)
		}

		// Check for base gram serving (100g recommended)
		if size.Unit == "gram" && size.Amount == 100 && size.GramEquivalent == 100 {
			hasGramBase = true
		}

		// Validate consistency: for gram unit, amount should equal gramEquivalent
		if size.Unit == "gram" && size.Amount != size.GramEquivalent {
			return fmt.Errorf("serving size %d: for gram unit, amount (%.2f) should equal gramEquivalent (%.2f)", i+1, size.Amount, size.GramEquivalent)
		}

		// Validate gramEquivalent is reasonable (not too large)
		if size.GramEquivalent > 100000 {
			return fmt.Errorf("serving size %d: gramEquivalent (%.2f) is unreasonably large", i+1, size.GramEquivalent)
		}
	}

	// Recommend having 100g base serving (warn but don't fail)
	if !hasGramBase {
		v.logger.Warn(ctx, "No gram base serving size found")
	}

	return nil
}

// validateCaloriesConsistency validates that calories match calculated value from macros
func (v *FoodValidator) validateCaloriesConsistency(req *request.CreateFoodRequest) error {
	// Calculate expected calories from macros
	// Formula: Protein (4 cal/g) + Carbs (4 cal/g) + Fat (9 cal/g) + Fiber (~2 cal/g)
	expectedCalories := (req.Macros.Protein * 4) +
		(req.Macros.Carbohydrates * 4) +
		(req.Macros.Fat * 9) +
		(req.Macros.Fiber * 2) // Fiber has approximately 2 calories per gram

	// Allow tolerance
	diff := req.Calories - expectedCalories
	if diff < -v.caloriesTolerance || diff > v.caloriesTolerance {
		return fmt.Errorf(
			"calories (%.2f) don't match calculated calories from macros (%.2f). Difference: %.2f. Allowed tolerance: ±%.2f",
			req.Calories, expectedCalories, diff, v.caloriesTolerance,
		)
	}

	return nil
}

// validateImageURL validates image URL format
func (v *FoodValidator) validateImageURL(urlStr string) error {
	if strings.TrimSpace(urlStr) == "" {
		return nil // Optional field
	}

	// Validate URL format
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Validate scheme (http/https only)
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme, got: %s", parsedURL.Scheme)
	}

	// Validate host is present
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a valid host")
	}

	// Validate URL length
	if len(urlStr) > 2048 {
		return fmt.Errorf("image URL exceeds maximum length (2048 chars)")
	}

	return nil
}
