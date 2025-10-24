package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"passwordHash" json:"-"`
	Profile      UserProfile        `bson:"profile" json:"profile"`
	Preferences  UserPreferences    `bson:"preferences" json:"preferences"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// UserProfile contains user profile information
type UserProfile struct {
	Name   string  `bson:"name" json:"name"`
	Age    int     `bson:"age" json:"age"`
	Weight float64 `bson:"weight" json:"weight"`
	Height float64 `bson:"height" json:"height"`
	Gender string  `bson:"gender" json:"gender"`
	Goal   string  `bson:"goal" json:"goal"` // "weight_loss", "muscle_gain", "maintenance"
}

// UserPreferences contains user preferences
type UserPreferences struct {
	Language      string         `bson:"language" json:"language"`
	CalorieTarget float64        `bson:"calorieTarget" json:"calorieTarget"`
	MacroTargets  MacroNutrients `bson:"macroTargets" json:"macroTargets"`
}

// MacroNutrients represents macronutrient values
type MacroNutrients struct {
	Protein       float64 `bson:"protein" json:"protein"` // grams per 100g
	Carbohydrates float64 `bson:"carbohydrates" json:"carbohydrates"`
	Fat           float64 `bson:"fat" json:"fat"`
	Fiber         float64 `bson:"fiber" json:"fiber"`
	Sugar         float64 `bson:"sugar,omitempty" json:"sugar,omitempty"`
}

// MicroNutrients represents micronutrient values
type MicroNutrients struct {
	VitaminA  float64 `bson:"vitaminA,omitempty" json:"vitaminA,omitempty"`
	VitaminC  float64 `bson:"vitaminC,omitempty" json:"vitaminC,omitempty"`
	Calcium   float64 `bson:"calcium,omitempty" json:"calcium,omitempty"`
	Iron      float64 `bson:"iron,omitempty" json:"iron,omitempty"`
	Sodium    float64 `bson:"sodium,omitempty" json:"sodium,omitempty"`
	Potassium float64 `bson:"potassium,omitempty" json:"potassium,omitempty"`
}

// ServingSize represents a serving size for a food item
type ServingSize struct {
	Unit           string  `bson:"unit" json:"unit"`                                   // "gram", "kg", "box", "cup", "ml", "piece"
	Amount         float64 `bson:"amount" json:"amount"`                               // e.g., 100, 1, 250
	Description    string  `bson:"description,omitempty" json:"description,omitempty"` // "1 medium banana"
	GramEquivalent float64 `bson:"gramEquivalent" json:"gramEquivalent"`               // Convert to grams
}

// FoodItem represents a food item in the database
type FoodItem struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         map[string]string  `bson:"name" json:"name"` // Multi-language support
	SearchTerms  []string           `bson:"searchTerms" json:"searchTerms"`
	Description  map[string]string  `bson:"description,omitempty" json:"description,omitempty"`
	Category     string             `bson:"category" json:"category"` // "protein", "vegetable", "fruit", "dairy", "grain"
	Macros       MacroNutrients     `bson:"macros" json:"macros"`
	Micros       MicroNutrients     `bson:"micros" json:"micros"`
	ServingSizes []ServingSize      `bson:"servingSizes" json:"servingSizes"`
	Calories     float64            `bson:"calories" json:"calories"` // Base calories per 100g
	CreatedBy    primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	Visibility   string             `bson:"visibility" json:"visibility"` // "public" or "private"
	Source       string             `bson:"source" json:"source"`         // "user" or "imported"
	ImageURL     string             `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}
