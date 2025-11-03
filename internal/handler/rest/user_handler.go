package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"nutrient_be/internal/dto/request"
	"nutrient_be/internal/handler/middleware"
	"nutrient_be/internal/pkg/logger"
	"nutrient_be/internal/service"
)

// UserHandler handles user profile and preferences endpoints
type UserHandler struct {
	userService    *service.UserService
	structValidator *validator.Validate
	logger          logger.Logger
	responseHelper  *middleware.ResponseHelper
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService, log logger.Logger) *UserHandler {
	return &UserHandler{
		userService:     userService,
		structValidator: validator.New(),
		logger:          log,
		responseHelper:  middleware.NewResponseHelper(),
	}
}

// GetProfile handles getting user profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	ctx := middleware.GetContext(c)

	// Get user ID from context (set by auth middleware)
	userIDStr, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		h.logger.Error(ctx, "User ID not found in context")
		h.responseHelper.Unauthorized(c, gin.H{"error": "User not authenticated"}, "Authentication required")
		return
	}

	// Get user profile
	profile, err := h.userService.GetProfile(c.Request.Context(), userIDStr)
	if err != nil {
		h.logger.Error(ctx, "Failed to get user profile", logger.Error(err))
		h.responseHelper.InternalError(c, gin.H{"details": err.Error()}, "Failed to get user profile")
		return
	}

	h.logger.Info(ctx, "User profile retrieved successfully")
	h.responseHelper.Success(c, profile, "User profile retrieved successfully")
}

// UpdateProfile handles updating user profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	ctx := middleware.GetContext(c)

	// Get user ID from context
	userIDStr, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		h.logger.Error(ctx, "User ID not found in context")
		h.responseHelper.Unauthorized(c, gin.H{"error": "User not authenticated"}, "Authentication required")
		return
	}

	// Bind request
	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx, "Failed to bind update profile request", logger.Error(err))
		h.responseHelper.BadRequest(c, gin.H{"details": err.Error()}, "Invalid request body")
		return
	}

	// Validate request
	if err := h.structValidator.Struct(&req); err != nil {
		h.logger.Error(ctx, "Update profile validation failed", logger.Error(err))
		h.responseHelper.ValidationError(c, gin.H{"validation_errors": err.Error()}, "Validation failed")
		return
	}

	// Update profile
	updatedProfile, err := h.userService.UpdateProfile(c.Request.Context(), userIDStr, &req)
	if err != nil {
		h.logger.Error(ctx, "Failed to update user profile", logger.Error(err))
		h.responseHelper.InternalError(c, gin.H{"details": err.Error()}, "Failed to update user profile")
		return
	}

	h.logger.Info(ctx, "User profile updated successfully")
	h.responseHelper.Success(c, updatedProfile, "User profile updated successfully")
}

// UpdatePreferences handles updating user preferences
func (h *UserHandler) UpdatePreferences(c *gin.Context) {
	ctx := middleware.GetContext(c)

	// Get user ID from context
	userIDStr, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		h.logger.Error(ctx, "User ID not found in context")
		h.responseHelper.Unauthorized(c, gin.H{"error": "User not authenticated"}, "Authentication required")
		return
	}

	// Bind request
	var req request.UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx, "Failed to bind update preferences request", logger.Error(err))
		h.responseHelper.BadRequest(c, gin.H{"details": err.Error()}, "Invalid request body")
		return
	}

	// Validate request
	if err := h.structValidator.Struct(&req); err != nil {
		h.logger.Error(ctx, "Update preferences validation failed", logger.Error(err))
		h.responseHelper.ValidationError(c, gin.H{"validation_errors": err.Error()}, "Validation failed")
		return
	}

	// Update preferences
	updatedUser, err := h.userService.UpdatePreferences(c.Request.Context(), userIDStr, &req)
	if err != nil {
		h.logger.Error(ctx, "Failed to update user preferences", logger.Error(err))
		h.responseHelper.InternalError(c, gin.H{"details": err.Error()}, "Failed to update user preferences")
		return
	}

	h.logger.Info(ctx, "User preferences updated successfully")
	h.responseHelper.Success(c, updatedUser, "User preferences updated successfully")
}

// ChangePassword handles changing user password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	ctx := middleware.GetContext(c)

	// Get user ID from context
	userIDStr, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		h.logger.Error(ctx, "User ID not found in context")
		h.responseHelper.Unauthorized(c, gin.H{"error": "User not authenticated"}, "Authentication required")
		return
	}

	// Bind request
	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx, "Failed to bind change password request", logger.Error(err))
		h.responseHelper.BadRequest(c, gin.H{"details": err.Error()}, "Invalid request body")
		return
	}

	// Validate request
	if err := h.structValidator.Struct(&req); err != nil {
		h.logger.Error(ctx, "Change password validation failed", logger.Error(err))
		h.responseHelper.ValidationError(c, gin.H{"validation_errors": err.Error()}, "Validation failed")
		return
	}

	// Change password
	if err := h.userService.ChangePassword(c.Request.Context(), userIDStr, &req); err != nil {
		h.logger.Error(ctx, "Failed to change password", logger.Error(err))
		h.responseHelper.InternalError(c, gin.H{"details": err.Error()}, "Failed to change password")
		return
	}

	h.logger.Info(ctx, "Password changed successfully")
	h.responseHelper.Success(c, gin.H{"message": "Password changed successfully"}, "Password changed successfully")
}

