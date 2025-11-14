package calculator

import (
	"fmt"

	"nutrient_be/internal/domain"
)

// CalculateNutrientsForServing calculates macros, micros, and calories for a given amount of food
// based on the serving unit and amount specified.
//
// Parameters:
//   - food: The food item with base nutrients per 100g
//   - servingUnit: The unit requested (e.g., "gram", "cup", "piece")
//   - amount: The amount in the specified unit
//
// Returns:
//   - calories: Calculated calories for the specified amount
//   - macros: Calculated macros for the specified amount
//   - micros: Calculated micros for the specified amount
//   - error: Error if serving unit not found or calculation fails
func CalculateNutrientsForServing(
	food *domain.FoodItem,
	servingUnit string,
	amount float64,
) (float64, domain.MacroNutrients, domain.MicroNutrients, error) {
	// Find the matching serving size
	var servingSize *domain.ServingSize
	for i := range food.ServingSizes {
		if food.ServingSizes[i].Unit == servingUnit {
			servingSize = &food.ServingSizes[i]
			break
		}
	}

	if servingSize == nil {
		return 0, domain.MacroNutrients{}, domain.MicroNutrients{},
			fmt.Errorf("serving unit '%s' not found for food '%s'", servingUnit, food.Name)
	}

	// Calculate total grams: (amount / servingSize.amount) * servingSize.gramEquivalent
	// Example: 2 cups where 1 cup = 250g -> (2 / 1) * 250 = 500g
	totalGrams := (amount / servingSize.Amount) * servingSize.GramEquivalent

	// Calculate multiplier: totalGrams / 100 (since food nutrients are per 100g)
	multiplier := totalGrams / 100.0

	// Calculate calories, macros, and micros
	calories := food.Calories * multiplier
	macros := domain.MacroNutrients{
		Protein:       food.Macros.Protein * multiplier,
		Carbohydrates: food.Macros.Carbohydrates * multiplier,
		Fat:           food.Macros.Fat * multiplier,
		Fiber:         food.Macros.Fiber * multiplier,
		Sugar:         food.Macros.Sugar * multiplier,
	}
	micros := domain.MicroNutrients{
		VitaminA:  food.Micros.VitaminA * multiplier,
		VitaminC:  food.Micros.VitaminC * multiplier,
		Calcium:   food.Micros.Calcium * multiplier,
		Iron:      food.Micros.Iron * multiplier,
		Sodium:    food.Micros.Sodium * multiplier,
		Potassium: food.Micros.Potassium * multiplier,
	}

	return calories, macros, micros, nil
}

// SumMacros sums multiple macro nutrient values
func SumMacros(macrosList ...domain.MacroNutrients) domain.MacroNutrients {
	result := domain.MacroNutrients{}
	for _, macros := range macrosList {
		result.Protein += macros.Protein
		result.Carbohydrates += macros.Carbohydrates
		result.Fat += macros.Fat
		result.Fiber += macros.Fiber
		result.Sugar += macros.Sugar
	}
	return result
}

// SumMicros sums multiple micro nutrient values
func SumMicros(microsList ...domain.MicroNutrients) domain.MicroNutrients {
	result := domain.MicroNutrients{}
	for _, micros := range microsList {
		result.VitaminA += micros.VitaminA
		result.VitaminC += micros.VitaminC
		result.Calcium += micros.Calcium
		result.Iron += micros.Iron
		result.Sodium += micros.Sodium
		result.Potassium += micros.Potassium
	}
	return result
}

