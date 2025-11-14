package response

import "time"

// MealTemplateResponse represents a meal template in API responses
type MealTemplateResponse struct {
	ID            string                         `json:"id"`
	UserID        string                         `json:"userId"`
	Name          string                         `json:"name"`
	Description   string                         `json:"description,omitempty"`
	MealType      string                         `json:"mealType"`
	FoodItems     []MealTemplateFoodItemResponse `json:"foodItems"`
	TotalCalories float64                        `json:"totalCalories"`
	TotalMacros   MacroNutrientsResponse         `json:"totalMacros"`
	TotalMicros   MicroNutrientsResponse         `json:"totalMicros,omitempty"`
	Tags          []string                       `json:"tags,omitempty"`
	IsPublic      bool                           `json:"isPublic"`
	CreatedAt     time.Time                      `json:"createdAt"`
	UpdatedAt     time.Time                      `json:"updatedAt"`
}

// MealTemplateFoodItemResponse represents a food item in a meal template response
type MealTemplateFoodItemResponse struct {
	FoodItemID  string                 `json:"foodItemId"`
	FoodName    string                 `json:"foodName"`
	ServingUnit string                 `json:"servingUnit"`
	Amount      float64                `json:"amount"`
	Calories    float64                `json:"calories"`
	Macros      MacroNutrientsResponse `json:"macros"`
	Micros      MicroNutrientsResponse `json:"micros,omitempty"`
}

