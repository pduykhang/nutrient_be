package request

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Name     string  `json:"name" validate:"required"`
	Age      int     `json:"age" validate:"required,min=1,max=120"`
	Weight   float64 `json:"weight" validate:"required,min=1"`
	Height   float64 `json:"height" validate:"required,min=1"`
	Gender   string  `json:"gender" validate:"required,oneof=male female other"`
	Goal     string  `json:"goal" validate:"required,oneof=weight_loss muscle_gain maintenance"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

