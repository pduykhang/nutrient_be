package request

// CreateMealTemplateRequest represents a request to create a meal template
type CreateMealTemplateRequest struct {
	Name        string                     `json:"name" validate:"required"`
	Description string                     `json:"description,omitempty"`
	MealType    string                     `json:"mealType" validate:"required,oneof=breakfast lunch dinner snack"`
	FoodItems   []MealTemplateFoodItemRequest `json:"foodItems" validate:"required,min=1"`
	Tags        []string                   `json:"tags,omitempty"`
	IsPublic    bool                       `json:"isPublic"`
}

// UpdateMealTemplateRequest represents a request to update a meal template
type UpdateMealTemplateRequest struct {
	Name        string                      `json:"name,omitempty"`
	Description string                      `json:"description,omitempty"`
	MealType    string                      `json:"mealType,omitempty"`
	FoodItems   []MealTemplateFoodItemRequest `json:"foodItems,omitempty"`
	Tags        []string                    `json:"tags,omitempty"`
	IsPublic    *bool                       `json:"isPublic,omitempty"`
}

// MealTemplateFoodItemRequest represents a food item in a meal template request
type MealTemplateFoodItemRequest struct {
	FoodItemID  string                `json:"foodItemId" validate:"required"`
	ServingUnit string                `json:"servingUnit" validate:"required"`
	Amount      float64               `json:"amount" validate:"required,min=0"`
}

