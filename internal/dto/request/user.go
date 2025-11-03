package request

// UpdateProfileRequest represents a request to update user profile
type UpdateProfileRequest struct {
	Name   *string  `json:"name,omitempty" validate:"omitempty,min=1"`
	Age    *int     `json:"age,omitempty" validate:"omitempty,min=1,max=120"`
	Weight *float64 `json:"weight,omitempty" validate:"omitempty,min=1"`
	Height *float64 `json:"height,omitempty" validate:"omitempty,min=1"`
	Gender *string  `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	Goal   *string  `json:"goal,omitempty" validate:"omitempty,oneof=weight_loss muscle_gain maintenance"`
}

// UpdatePreferencesRequest represents a request to update user preferences
type UpdatePreferencesRequest struct {
	Language      *string                `json:"language,omitempty" validate:"omitempty,oneof=en vi"`
	CalorieTarget *float64               `json:"calorieTarget,omitempty" validate:"omitempty,min=0"`
	MacroTargets  *MacroNutrientsRequest `json:"macroTargets,omitempty"`
}

// ChangePasswordRequest represents a request to change password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=6"`
}
