package request

import "time"

// CreateMealPlanRequest represents a request to create a meal plan
type CreateMealPlanRequest struct {
	Name          string    `json:"name" validate:"required"`
	Description   string    `json:"description,omitempty"`
	StartDate     time.Time `json:"startDate" validate:"required"`
	EndDate       time.Time `json:"endDate" validate:"required"`
	PlanType      string    `json:"planType" validate:"required,oneof=weekly monthly"`
	Goal          string    `json:"goal" validate:"required,oneof=weight_loss muscle_gain maintenance"`
	TargetCalories float64  `json:"targetCalories" validate:"required,min=0"`
}

// UpdateMealPlanRequest represents a request to update a meal plan
type UpdateMealPlanRequest struct {
	Name          string    `json:"name,omitempty"`
	Description   string    `json:"description,omitempty"`
	StartDate     *time.Time `json:"startDate,omitempty"`
	EndDate       *time.Time `json:"endDate,omitempty"`
	PlanType      string    `json:"planType,omitempty"`
	Goal          string    `json:"goal,omitempty"`
	TargetCalories *float64 `json:"targetCalories,omitempty"`
	Status        string    `json:"status,omitempty"`
}

