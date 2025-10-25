package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"nutrient_be/internal/config"
	"nutrient_be/internal/domain"
	"nutrient_be/internal/pkg/logger"
	"nutrient_be/internal/repository"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo repository.UserRepository
	config   config.AuthConfig
	logger   logger.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, cfg config.AuthConfig, log logger.Logger) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
		logger:   log,
	}
}

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

// AuthResponse represents an authentication response
type AuthResponse struct {
	User         *domain.User `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	ExpiresAt    time.Time    `json:"expiresAt"`
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		ID:           primitive.NewObjectID(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Profile: domain.UserProfile{
			Name:   req.Name,
			Age:    req.Age,
			Weight: req.Weight,
			Height: req.Height,
			Gender: req.Gender,
			Goal:   req.Goal,
		},
		Preferences: domain.UserPreferences{
			Language:      "en",
			CalorieTarget: s.calculateCalorieTarget(req.Weight, req.Height, req.Age, req.Gender, req.Goal),
			MacroTargets:  s.calculateMacroTargets(req.Goal),
		},
	}

	// Save user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	accessToken, refreshToken, expiresAt, err := s.generateTokens(user.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	s.logger.Info(ctx, "User registered successfully",
		logger.String("email", req.Email),
		logger.String("userID", user.ID.Hex()))

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate tokens
	accessToken, refreshToken, expiresAt, err := s.generateTokens(user.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	s.logger.Info(ctx, "User logged in successfully",
		logger.String("email", req.Email),
		logger.String("userID", user.ID.Hex()))

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// RefreshToken refreshes an access token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// Parse and validate refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate new tokens
	accessToken, newRefreshToken, expiresAt, err := s.generateTokens(user.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// generateTokens generates access and refresh tokens
func (s *AuthService) generateTokens(userID string) (string, string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.config.JWTExpiration * time.Second)

	// Access token
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
		"iat":     now.Unix(),
		"type":    "access",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Refresh token
	refreshExpiresAt := now.Add(s.config.RefreshExpiration * time.Second)
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     refreshExpiresAt.Unix(),
		"iat":     now.Unix(),
		"type":    "refresh",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessTokenString, refreshTokenString, expiresAt, nil
}

// calculateCalorieTarget calculates daily calorie target based on user profile
func (s *AuthService) calculateCalorieTarget(weight, height float64, age int, gender, goal string) float64 {
	// Basic BMR calculation (Mifflin-St Jeor Equation)
	var bmr float64
	if gender == "male" {
		bmr = 10*weight + 6.25*height - 5*float64(age) + 5
	} else {
		bmr = 10*weight + 6.25*height - 5*float64(age) - 161
	}

	// Activity factor (sedentary)
	activityFactor := 1.2
	maintenanceCalories := bmr * activityFactor

	// Adjust based on goal
	switch goal {
	case "weight_loss":
		return maintenanceCalories - 500 // 500 calorie deficit
	case "muscle_gain":
		return maintenanceCalories + 300 // 300 calorie surplus
	default: // maintenance
		return maintenanceCalories
	}
}

// calculateMacroTargets calculates macro targets based on goal
func (s *AuthService) calculateMacroTargets(goal string) domain.MacroNutrients {
	switch goal {
	case "weight_loss":
		return domain.MacroNutrients{
			Protein:       1.6, // g per kg body weight
			Carbohydrates: 2.0,
			Fat:           0.8,
			Fiber:         0.03,
		}
	case "muscle_gain":
		return domain.MacroNutrients{
			Protein:       2.2,
			Carbohydrates: 4.0,
			Fat:           1.0,
			Fiber:         0.03,
		}
	default: // maintenance
		return domain.MacroNutrients{
			Protein:       1.8,
			Carbohydrates: 3.0,
			Fat:           0.9,
			Fiber:         0.03,
		}
	}
}
