package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MealTemplateFoodItem represents a food item in a meal template
type MealTemplateFoodItem struct {
	FoodItemID  primitive.ObjectID `bson:"foodItemId" json:"foodItemId"`
	FoodName    string             `bson:"foodName" json:"foodName"` // Denormalized for performance
	ServingUnit string             `bson:"servingUnit" json:"servingUnit"`
	Amount      float64            `bson:"amount" json:"amount"`
	Calories    float64            `bson:"calories" json:"calories"` // Calculated calories for this amount
	Macros      MacroNutrients     `bson:"macros" json:"macros"`     // Calculated macros for this amount
	Micros      MicroNutrients     `bson:"micros,omitempty" json:"micros,omitempty"` // Calculated micros for this amount
}

// MealTemplate represents a reusable meal combination
type MealTemplate struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID     `bson:"userId" json:"userId"`
	Name          string                 `bson:"name" json:"name"`
	Description   string                 `bson:"description,omitempty" json:"description,omitempty"`
	MealType      string                 `bson:"mealType" json:"mealType"` // "breakfast", "lunch", "dinner", "snack"
	FoodItems     []MealTemplateFoodItem `bson:"foodItems" json:"foodItems"`
	TotalCalories float64                `bson:"totalCalories" json:"totalCalories"` // Calculated
	TotalMacros   MacroNutrients         `bson:"totalMacros" json:"totalMacros"`
	TotalMicros   MicroNutrients         `bson:"totalMicros,omitempty" json:"totalMicros,omitempty"`
	Tags          []string               `bson:"tags,omitempty" json:"tags,omitempty"`
	IsPublic      bool                   `bson:"isPublic" json:"isPublic"`
	CreatedAt     time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time              `bson:"updatedAt" json:"updatedAt"`
}

// MealFoodItem represents a food item in a meal
type MealFoodItem struct {
	FoodItemID   primitive.ObjectID `bson:"foodItemId" json:"foodItemId"`
	FoodName     string             `bson:"foodName" json:"foodName"`                             // Denormalized
	FoodCategory string             `bson:"foodCategory,omitempty" json:"foodCategory,omitempty"` // Denormalized
	ServingUnit  string             `bson:"servingUnit" json:"servingUnit"`
	Amount       float64            `bson:"amount" json:"amount"`
	Calories     float64            `bson:"calories" json:"calories"` // Calculated
	Macros       MacroNutrients     `bson:"macros" json:"macros"`     // Calculated
}

// Meal represents a single meal within a day
type Meal struct {
	ID          string              `bson:"id" json:"id"` // Unique ID within the meal plan
	MealType    string              `bson:"mealType" json:"mealType"`
	Time        string              `bson:"time,omitempty" json:"time,omitempty"` // "07:00"
	TemplateID  *primitive.ObjectID `bson:"templateId,omitempty" json:"templateId,omitempty"`
	FoodItems   []MealFoodItem      `bson:"foodItems" json:"foodItems"`
	Calories    float64             `bson:"calories" json:"calories"` // Sum for this meal
	Macros      MacroNutrients      `bson:"macros" json:"macros"`     // Sum for this meal
	Notes       string              `bson:"notes,omitempty" json:"notes,omitempty"`
	IsCompleted bool                `bson:"isCompleted" json:"isCompleted"`
}

// DailyMeal represents all meals for a single day
type DailyMeal struct {
	Date          time.Time      `bson:"date" json:"date"`
	DayOfWeek     string         `bson:"dayOfWeek" json:"dayOfWeek"`
	Meals         []Meal         `bson:"meals" json:"meals"`
	TotalCalories float64        `bson:"totalCalories" json:"totalCalories"` // Sum for this day
	TotalMacros   MacroNutrients `bson:"totalMacros" json:"totalMacros"`     // Sum for this day
	Notes         string         `bson:"notes,omitempty" json:"notes,omitempty"`
	IsCompleted   bool           `bson:"isCompleted" json:"isCompleted"`
}

// MealPlan represents a complete eating schedule for a time period
type MealPlan struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	Name           string             `bson:"name" json:"name"`
	Description    string             `bson:"description,omitempty" json:"description,omitempty"`
	StartDate      time.Time          `bson:"startDate" json:"startDate"`
	EndDate        time.Time          `bson:"endDate" json:"endDate"`
	PlanType       string             `bson:"planType" json:"planType"`             // "weekly" or "monthly"
	Goal           string             `bson:"goal" json:"goal"`                     // "weight_loss", "muscle_gain", "maintenance"
	TargetCalories float64            `bson:"targetCalories" json:"targetCalories"` // Daily target
	TargetMacros   MacroNutrients     `bson:"targetMacros" json:"targetMacros"`     // Daily target
	DailyMeals     []DailyMeal        `bson:"dailyMeals" json:"dailyMeals"`
	TotalCalories  float64            `bson:"totalCalories" json:"totalCalories"` // Total for entire period
	Status         string             `bson:"status" json:"status"`               // "draft", "active", "completed"
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}
