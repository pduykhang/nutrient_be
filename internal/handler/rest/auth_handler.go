package rest

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"nutrient_be/internal/dto/request"
	"nutrient_be/internal/handler/middleware"
	"nutrient_be/internal/pkg/logger"
	"nutrient_be/internal/service"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService    *service.AuthService
	validator      *validator.Validate
	logger         logger.Logger
	responseHelper *middleware.ResponseHelper
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		validator:      validator.New(),
		logger:         log,
		responseHelper: middleware.NewResponseHelper(),
	}
}

// Register handles user registration (email and password only)
func (h *AuthHandler) Register(c *gin.Context) {
	// Get context from middleware
	ctx := middleware.GetContext(c)

	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx, "Failed to bind register request", logger.Error(err))
		h.responseHelper.BadRequest(c, gin.H{"details": err.Error()}, "Invalid request body")
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error(ctx, "Register request validation failed", logger.Error(err))
		h.responseHelper.ValidationError(c, gin.H{"validation_errors": err.Error()}, "Validation failed")
		return
	}

	// Register user
	response, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error(ctx, "Failed to register user", logger.Error(err))
		h.responseHelper.Conflict(c, gin.H{"details": err.Error()}, "Registration failed")
		return
	}

	// Log success with context information automatically included
	h.logger.Info(ctx, "User registered successfully", logger.String("email", req.Email))
	h.responseHelper.Created(c, response, "User registered successfully")
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	// Get context from middleware
	ctx := middleware.GetContext(c)

	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx, "Failed to bind login request", logger.Error(err))
		h.responseHelper.BadRequest(c, gin.H{"details": err.Error()}, "Invalid request body")
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error(ctx, "Login request validation failed", logger.Error(err))
		h.responseHelper.ValidationError(c, gin.H{"validation_errors": err.Error()}, "Validation failed")
		return
	}

	// Login user
	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error(ctx, "Failed to login user", logger.Error(err))
		h.responseHelper.Unauthorized(c, gin.H{"details": err.Error()}, "Invalid credentials")
		return
	}

	// Log success with context information automatically included
	h.logger.Info(ctx, "User logged in successfully", logger.String("email", req.Email))
	h.responseHelper.Success(c, response, "Login successful")
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	// Get context from middleware
	ctx := middleware.GetContext(c)

	var req struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx, "Failed to bind refresh request", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error(ctx, "Refresh request validation failed", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	// Refresh token
	response, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Error(ctx, "Failed to refresh token", logger.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info(ctx, "Token refreshed successfully")
	h.responseHelper.Success(c, response, "Token refreshed successfully")
}

// Logout handles user logout (token invalidation would be handled by token blacklist in production)
func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := middleware.GetContext(c)

	// In a production system, you would invalidate the token here
	// For now, we just acknowledge the logout request
	// Token invalidation can be implemented with a token blacklist or by reducing token lifetime

	h.logger.Info(ctx, "User logged out successfully")
	h.responseHelper.Success(c, gin.H{"message": "Logged out successfully"}, "Logged out successfully")
}

// Validate handles token validation
func (h *AuthHandler) Validate(c *gin.Context) {
	ctx := middleware.GetContext(c)

	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		h.responseHelper.Unauthorized(c, gin.H{"error": "Authorization header required"}, "Authorization required")
		return
	}

	// Check if header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		h.responseHelper.Unauthorized(c, gin.H{"error": "Invalid authorization header format"}, "Invalid authorization header")
		return
	}

	// Extract token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Validate token
	userID, err := h.authService.ValidateToken(c.Request.Context(), tokenString)
	if err != nil {
		h.logger.Error(ctx, "Token validation failed", logger.Error(err))
		h.responseHelper.Unauthorized(c, gin.H{"error": "Invalid token"}, "Invalid token")
		return
	}

	h.logger.Info(ctx, "Token validated successfully", logger.String("userID", userID))
	h.responseHelper.Success(c, gin.H{"valid": true, "userID": userID}, "Token is valid")
}
