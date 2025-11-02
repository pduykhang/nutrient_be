package request

import (
	"encoding/json"
	"fmt"
)

// CreateFoodRequest represents a request to create a food item
type CreateFoodRequest struct {
	Name         MultiLanguage         `json:"name" validate:"required"`
	SearchTerms  []string              `json:"searchTerms"`
	Description  MultiLanguage         `json:"description,omitempty"`
	Category     string                `json:"category" validate:"required,oneof=protein vegetable fruit dairy grain"`
	Macros       MacroNutrientsRequest `json:"macros" validate:"required"`
	Micros       MicroNutrientsRequest `json:"micros,omitempty"`
	ServingSizes []ServingSizeRequest  `json:"servingSizes" validate:"required,min=1"`
	Calories     float64               `json:"calories" validate:"required,min=0"`
	Visibility   string                `json:"visibility" validate:"required,oneof=public private"`
	ImageURL     string                `json:"imageUrl,omitempty"`
}

// UpdateFoodRequest represents a request to update a food item
type UpdateFoodRequest struct {
	Name         MultiLanguage          `json:"name,omitempty"`
	SearchTerms  []string               `json:"searchTerms,omitempty"`
	Description  MultiLanguage          `json:"description,omitempty"`
	Category     string                 `json:"category,omitempty"`
	Macros       *MacroNutrientsRequest `json:"macros,omitempty"`
	Micros       *MicroNutrientsRequest `json:"micros,omitempty"`
	ServingSizes []ServingSizeRequest   `json:"servingSizes,omitempty"`
	Calories     *float64               `json:"calories,omitempty"`
	Visibility   string                 `json:"visibility,omitempty"`
	ImageURL     string                 `json:"imageUrl,omitempty"`
}

// SearchFoodRequest represents a request to search food items
type SearchFoodRequest struct {
	Query  string `form:"query" validate:"required"`
	Limit  int    `form:"limit,default=20"`
	Offset int    `form:"offset,default=0"`
}

// MacroNutrientsRequest represents macronutrient values in requests
type MacroNutrientsRequest struct {
	Protein       float64 `json:"protein" validate:"min=0"`
	Carbohydrates float64 `json:"carbohydrates" validate:"min=0"`
	Fat           float64 `json:"fat" validate:"min=0"`
	Fiber         float64 `json:"fiber" validate:"min=0"`
	Sugar         float64 `json:"sugar,omitempty" validate:"min=0"`
}

// MicroNutrientsRequest represents micronutrient values in requests
type MicroNutrientsRequest struct {
	VitaminA  float64 `json:"vitaminA,omitempty" validate:"min=0"`
	VitaminC  float64 `json:"vitaminC,omitempty" validate:"min=0"`
	Calcium   float64 `json:"calcium,omitempty" validate:"min=0"`
	Iron      float64 `json:"iron,omitempty" validate:"min=0"`
	Sodium    float64 `json:"sodium,omitempty" validate:"min=0"`
	Potassium float64 `json:"potassium,omitempty" validate:"min=0"`
}

// ServingSizeRequest represents a serving size in requests
type ServingSizeRequest struct {
	Unit           string  `json:"unit" validate:"required"`
	Amount         float64 `json:"amount" validate:"required,min=0"`
	Description    string  `json:"description,omitempty"`
	GramEquivalent float64 `json:"gramEquivalent" validate:"required,min=0"`
}

var (
	supportedLanguages = map[string]bool{
		"en": true,
		"vi": true,
	}
)

type MultiLanguage map[string]string

func (m MultiLanguage) Get(lang string) string {
	if value, ok := m[lang]; ok {
		return value
	}
	return ""
}
func (m MultiLanguage) GetRaw() map[string]string {
	return m
}

func (m MultiLanguage) Validate() error {
	for lang, value := range m {
		if value == "" {
			return fmt.Errorf("value for language %s is required", lang)
		}
	}

	return nil
}

func (m *MultiLanguage) UnmarshalJSON(data []byte) error {
	if *m == nil {
		*m = make(map[string]string)
	}
	tmp := make(map[string]string)
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	for lang, value := range tmp {
		if !supportedLanguages[lang] {
			return fmt.Errorf("language %s is not supported", lang)
		}
		(*m)[lang] = value
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (m MultiLanguage) MarshalJSON() ([]byte, error) {
	result := make(map[string]string)
	for lang, value := range m {
		if !supportedLanguages[lang] {
			continue
		}
		result[lang] = value
	}
	return json.Marshal(result)
}
