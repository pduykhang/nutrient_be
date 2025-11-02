package response

import "time"

// MealPlanResponse represents a meal plan in API responses
type MealPlanResponse struct {
	ID             string                   `json:"id"`
	UserID         string                   `json:"userId"`
	Name           string                   `json:"name"`
	Description    string                   `json:"description,omitempty"`
	StartDate      time.Time                `json:"startDate"`
	EndDate        time.Time                `json:"endDate"`
	PlanType       string                   `json:"planType"`
	Goal           string                   `json:"goal"`
	TargetCalories float64                  `json:"targetCalories"`
	TargetMacros   MacroNutrientsResponse   `json:"targetMacros"`
	DailyMeals     []DailyMealResponse      `json:"dailyMeals"`
	TotalCalories  float64                  `json:"totalCalories"`
	Status         string                   `json:"status"`
	CreatedAt      time.Time                `json:"createdAt"`
	UpdatedAt      time.Time                `json:"updatedAt"`
}

// DailyMealResponse represents daily meals in API responses
type DailyMealResponse struct {
	Date          time.Time            `json:"date"`
	DayOfWeek     string               `json:"dayOfWeek"`
	Meals         []MealResponse      `json:"meals"`
	TotalCalories float64              `json:"totalCalories"`
	TotalMacros   MacroNutrientsResponse `json:"totalMacros"`
	Notes         string               `json:"notes,omitempty"`
	IsCompleted   bool                `json:"isCompleted"`
}

// MealResponse represents a meal in API responses
type MealResponse struct {
	ID          string                  `json:"id"`
	MealType    string                  `json:"mealType"`
	Time        string                  `json:"time,omitempty"`
	TemplateID  string                  `json:"templateId,omitempty"`
	FoodItems   []MealFoodItemResponse  `json:"foodItems"`
	Calories    float64                 `json:"calories"`
	Macros      MacroNutrientsResponse  `json:"macros"`
	Notes       string                  `json:"notes,omitempty"`
	IsCompleted bool                    `json:"isCompleted"`
}

// MealFoodItemResponse represents a food item in a meal response
type MealFoodItemResponse struct {
	FoodItemID   string                  `json:"foodItemId"`
	FoodName     string                  `json:"foodName"`
	FoodCategory string                  `json:"foodCategory,omitempty"`
	ServingUnit  string                  `json:"servingUnit"`
	Amount       float64                 `json:"amount"`
	Calories     float64                 `json:"calories"`
	Macros       MacroNutrientsResponse  `json:"macros"`
}

