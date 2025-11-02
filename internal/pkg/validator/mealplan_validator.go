package validator

import (
	"fmt"
	"strings"
	"time"

	"nutrient_be/internal/dto/request"
	"nutrient_be/internal/pkg/logger"
)

// MealPlanValidator handles meal plan validation
type MealPlanValidator struct {
	minCalories          float64
	maxCalories          float64
	minDateRangeDays     int
	maxDateRangeDays     int
	maxNameLength        int
	maxDescriptionLength int
	logger               logger.Logger
}

// NewMealPlanValidator creates a new meal plan validator
func NewMealPlanValidator(logger logger.Logger) *MealPlanValidator {
	return &MealPlanValidator{
		minCalories:          500,
		maxCalories:          5000,
		minDateRangeDays:     1,
		maxDateRangeDays:     90, // Max 3 months
		maxNameLength:        100,
		maxDescriptionLength: 500,
		logger:               logger,
	}
}

// ValidateCreateRequest validates a CreateMealPlanRequest
func (v *MealPlanValidator) ValidateCreateRequest(req *request.CreateMealPlanRequest) error {
	// 1. Validate Name
	if err := v.validateName(req.Name); err != nil {
		return fmt.Errorf("name validation failed: %w", err)
	}

	// 2. Validate Description (optional)
	if req.Description != "" && len(req.Description) > v.maxDescriptionLength {
		return fmt.Errorf("description exceeds maximum length (%d chars)", v.maxDescriptionLength)
	}

	// 3. Validate Date Range
	if err := v.validateDateRange(req.StartDate, req.EndDate); err != nil {
		return fmt.Errorf("date range validation failed: %w", err)
	}

	// 4. Validate Plan Type
	if err := v.validatePlanType(req.PlanType); err != nil {
		return fmt.Errorf("plan type validation failed: %w", err)
	}

	// 5. Validate Goal
	if err := v.validateGoal(req.Goal); err != nil {
		return fmt.Errorf("goal validation failed: %w", err)
	}

	// 6. Validate Target Calories
	if err := v.validateTargetCalories(req.TargetCalories); err != nil {
		return fmt.Errorf("target calories validation failed: %w", err)
	}

	return nil
}

// validateName validates meal plan name
func (v *MealPlanValidator) validateName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(trimmed) > v.maxNameLength {
		return fmt.Errorf("name exceeds maximum length (%d chars)", v.maxNameLength)
	}
	return nil
}

// validateDateRange validates start and end dates
func (v *MealPlanValidator) validateDateRange(startDate, endDate time.Time) error {
	now := time.Now().Truncate(24 * time.Hour)

	// Start date cannot be in the past (allow today)
	if startDate.Before(now) {
		return fmt.Errorf("start date cannot be in the past")
	}

	// End date must be after start date
	if !endDate.After(startDate) {
		return fmt.Errorf("end date must be after start date")
	}

	// Calculate date range
	daysDiff := int(endDate.Sub(startDate).Hours() / 24)

	// Validate minimum range
	if daysDiff < v.minDateRangeDays {
		return fmt.Errorf("date range must be at least %d day(s)", v.minDateRangeDays)
	}

	// Validate maximum range
	if daysDiff > v.maxDateRangeDays {
		return fmt.Errorf("date range exceeds maximum (%d days)", v.maxDateRangeDays)
	}

	return nil
}

// validatePlanType validates plan type
func (v *MealPlanValidator) validatePlanType(planType string) error {
	validTypes := map[string]bool{
		"weekly":  true,
		"monthly": true,
	}

	if !validTypes[planType] {
		return fmt.Errorf("invalid plan type '%s'. Valid types: weekly, monthly", planType)
	}

	return nil
}

// validateGoal validates goal value
func (v *MealPlanValidator) validateGoal(goal string) error {
	validGoals := map[string]bool{
		"weight_loss": true,
		"muscle_gain": true,
		"maintenance": true,
	}

	if !validGoals[goal] {
		return fmt.Errorf("invalid goal '%s'. Valid goals: weight_loss, muscle_gain, maintenance", goal)
	}

	return nil
}

// validateTargetCalories validates target calories
func (v *MealPlanValidator) validateTargetCalories(calories float64) error {
	if calories < v.minCalories {
		return fmt.Errorf("target calories (%.2f) is below minimum (%.2f)", calories, v.minCalories)
	}
	if calories > v.maxCalories {
		return fmt.Errorf("target calories (%.2f) exceeds maximum (%.2f)", calories, v.maxCalories)
	}
	return nil
}
