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
	"nutrient_be/internal/dto/request"
	"nutrient_be/internal/dto/response"
	"nutrient_be/internal/pkg/logger"
)

// UserRepository defines the interface for user data operations used by AuthService
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// AuthService handles authentication operations
type AuthService struct {
	userRepo UserRepository
	config   config.AuthConfig
	logger   logger.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo UserRepository, cfg config.AuthConfig, log logger.Logger) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
		logger:   log,
	}
}

// LoginRequest is now in internal/dto/request/auth.go

// AuthResponse represents an authentication response
type AuthResponse struct {
	User         *response.UserResponse `json:"user"`
	AccessToken  string                 `json:"accessToken"`
	RefreshToken string                 `json:"refreshToken"`
	ExpiresAt    time.Time              `json:"expiresAt"`
}

// Register registers a new user (email and password only)
// Profile should be set separately via user service
func (s *AuthService) Register(ctx context.Context, req *request.RegisterRequest) (*AuthResponse, error) {
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

	// Create user with default profile
	user := &domain.User{
		ID:           primitive.NewObjectID(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Profile:      domain.UserProfile{
			// Default empty profile - user should set it via user service
		},
		Preferences: domain.UserPreferences{
			Language:      "en",
			CalorieTarget: 0, // Will be calculated when profile is set
			MacroTargets:  domain.MacroNutrients{},
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

	// Convert to response
	userResponse := domainUserToResponse(user)

	return &AuthResponse{
		User:         userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req *request.LoginRequest) (*AuthResponse, error) {
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

	// Convert to response
	userResponse := domainUserToResponse(user)

	return &AuthResponse{
		User:         userResponse,
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

	// Convert to response
	userResponse := domainUserToResponse(user)

	return &AuthResponse{
		User:         userResponse,
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

// ValidateToken validates a JWT token and returns user ID
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	// Check token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return "", fmt.Errorf("invalid token type")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid user ID in token")
	}

	return userIDStr, nil
}

// domainUserToResponse converts domain.User to response.UserResponse
func domainUserToResponse(user *domain.User) *response.UserResponse {
	return &response.UserResponse{
		ID:    user.ID.Hex(),
		Email: user.Email,
		Profile: response.UserProfileResponse{
			Name:   user.Profile.Name,
			Age:    user.Profile.Age,
			Weight: user.Profile.Weight,
			Height: user.Profile.Height,
			Gender: user.Profile.Gender,
			Goal:   user.Profile.Goal,
		},
		Preferences: response.UserPreferencesResponse{
			Language:      user.Preferences.Language,
			CalorieTarget: user.Preferences.CalorieTarget,
			MacroTargets: response.MacroNutrientsResponse{
				Protein:       user.Preferences.MacroTargets.Protein,
				Carbohydrates: user.Preferences.MacroTargets.Carbohydrates,
				Fat:           user.Preferences.MacroTargets.Fat,
				Fiber:         user.Preferences.MacroTargets.Fiber,
				Sugar:         user.Preferences.MacroTargets.Sugar,
			},
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
