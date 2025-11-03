package response

import (
	"time"
)

// AuthResponse represents an authentication response
type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
	ExpiresAt    time.Time     `json:"expiresAt"`
}
