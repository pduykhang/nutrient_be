package validator

import "nutrient_be/internal/pkg/logger"

// Validator provides centralized validation for all entities
type Validator struct {
	Food     *FoodValidator
	Meal     *MealValidator
	MealPlan *MealPlanValidator
}

// New creates a new validator instance with all validators
func New(logger logger.Logger) *Validator {
	return &Validator{
		Food:     NewFoodValidator(logger),
		Meal:     NewMealValidator(logger),
		MealPlan: NewMealPlanValidator(logger),
	}
}
