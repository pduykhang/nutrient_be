package response

import "time"

// FoodItemResponse represents a food item in API responses
type FoodItemResponse struct {
	ID           string                 `json:"id"`
	Name         map[string]string      `json:"name"`
	SearchTerms  []string               `json:"searchTerms"`
	Description  map[string]string      `json:"description,omitempty"`
	Category     string                 `json:"category"`
	Macros       MacroNutrientsResponse `json:"macros"`
	Micros       MicroNutrientsResponse `json:"micros,omitempty"`
	ServingSizes []ServingSizeResponse  `json:"servingSizes"`
	Calories     float64                `json:"calories"`
	CreatedBy    string                 `json:"createdBy"`
	Visibility   string                 `json:"visibility"`
	Source       string                 `json:"source"`
	ImageURL     string                 `json:"imageUrl,omitempty"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

// MacroNutrientsResponse represents macronutrient values in API responses
type MacroNutrientsResponse struct {
	Protein       float64 `json:"protein"`
	Carbohydrates float64 `json:"carbohydrates"`
	Fat           float64 `json:"fat"`
	Fiber         float64 `json:"fiber"`
	Sugar         float64 `json:"sugar,omitempty"`
}

// MicroNutrientsResponse represents micronutrient values in API responses
type MicroNutrientsResponse struct {
	VitaminA  float64 `json:"vitaminA,omitempty"`
	VitaminC  float64 `json:"vitaminC,omitempty"`
	Calcium   float64 `json:"calcium,omitempty"`
	Iron      float64 `json:"iron,omitempty"`
	Sodium    float64 `json:"sodium,omitempty"`
	Potassium float64 `json:"potassium,omitempty"`
}

// ServingSizeResponse represents a serving size in API responses
type ServingSizeResponse struct {
	Unit           string  `json:"unit"`
	Amount         float64 `json:"amount"`
	Description    string  `json:"description,omitempty"`
	GramEquivalent float64 `json:"gramEquivalent"`
}
