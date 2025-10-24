package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShoppingItem represents an item in a shopping list
type ShoppingItem struct {
	FoodItemID  primitive.ObjectID `bson:"foodItemId" json:"foodItemId"`
	FoodName    string             `bson:"foodName" json:"foodName"`
	TotalAmount float64            `bson:"totalAmount" json:"totalAmount"`
	Unit        string             `bson:"unit" json:"unit"`
	Checked     bool               `bson:"checked" json:"checked"`
}

// ShoppingList represents a shopping list generated from a meal plan
type ShoppingList struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"userId" json:"userId"`
	MealPlanID primitive.ObjectID `bson:"mealPlanId" json:"mealPlanId"`
	Items      []ShoppingItem     `bson:"items" json:"items"`
	TotalCost  float64            `bson:"totalCost,omitempty" json:"totalCost,omitempty"` // Optional
	Status     string             `bson:"status" json:"status"`                           // "pending", "completed"
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}
