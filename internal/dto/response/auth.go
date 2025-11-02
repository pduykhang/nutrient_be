package response

import (
	"time"
)

// UserResponse represents a user in API responses
type UserResponse struct {
	ID          string                  `json:"id"`
	Email       string                  `json:"email"`
	Profile     UserProfileResponse     `json:"profile"`
	Preferences UserPreferencesResponse `json:"preferences"`
	CreatedAt   time.Time               `json:"createdAt"`
	UpdatedAt   time.Time               `json:"updatedAt"`
}

// UserProfileResponse represents user profile in API responses
type UserProfileResponse struct {
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Weight float64 `json:"weight"`
	Height float64 `json:"height"`
	Gender string  `json:"gender"`
	Goal   string  `json:"goal"`
}

// UserPreferencesResponse represents user preferences in API responses
type UserPreferencesResponse struct {
	Language      string                 `json:"language"`
	CalorieTarget float64                `json:"calorieTarget"`
	MacroTargets  MacroNutrientsResponse `json:"macroTargets"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
	ExpiresAt    time.Time     `json:"expiresAt"`
}
