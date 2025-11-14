package rest

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"nutrient_be/internal/domain"
	"nutrient_be/internal/dto/request"
	"nutrient_be/internal/dto/response"
	"nutrient_be/internal/handler/middleware"
	"nutrient_be/internal/pkg/logger"
	mealValidator "nutrient_be/internal/pkg/validator"
	"nutrient_be/internal/service"
)

// MealHandler handles meal template endpoints
type MealHandler struct {
	mealService     *service.MealService
	structValidator *validator.Validate
	mealValidator   *mealValidator.MealValidator
	logger          logger.Logger
	responseHelper  *middleware.ResponseHelper
}

// NewMealHandler creates a new meal handler
func NewMealHandler(mealService *service.MealService, log logger.Logger) *MealHandler {
	return &MealHandler{
		mealService:     mealService,
		structValidator: validator.New(),
		mealValidator:   mealValidator.NewMealValidator(log),
		logger:          log,
		responseHelper:  middleware.NewResponseHelper(),
	}
}

// getUserIDFromContext extracts user ID from context or returns error response
// Returns userID and true if successful, false if error response was sent
func (h *MealHandler) getUserIDFromContext(c *gin.Context, ctx context.Context) (string, bool) {
	userIDStr, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		h.logger.Error(ctx, "User ID not found in context")
		h.responseHelper.Unauthorized(c, gin.H{"error": "User not authenticated"}, "Authentication required")
		return "", false
	}
	return userIDStr, true
}

// getTemplateIDFromParams extracts template ID from URL params or returns error response
// Returns templateID and true if successful, false if error response was sent
func (h *MealHandler) getTemplateIDFromParams(c *gin.Context, ctx context.Context) (string, bool) {
	templateID := c.Param("id")
	if templateID == "" {
		h.logger.Error(ctx, "Template ID is required")
		h.responseHelper.BadRequest(c, gin.H{"error": "Template ID is required"}, "Template ID is required")
		return "", false
	}
	return templateID, true
}

// bindRequest binds JSON request
// Returns true if successful, false if error response was sent
func (h *MealHandler) bindRequest(c *gin.Context, ctx context.Context, req interface{}, requestType string) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		h.logger.Error(ctx, "Failed to bind request", logger.String("type", requestType), logger.Error(err))
		h.responseHelper.BadRequest(c, gin.H{"details": err.Error()}, "Invalid request body")
		return false
	}
	return true
}

// validateRequest validates request using struct validator and business validator
// Returns true if successful, false if error response was sent
func (h *MealHandler) validateRequest(c *gin.Context, ctx context.Context, req interface{}, requestType string) bool {
	// Format/structure validation (struct tags)
	if err := h.structValidator.Struct(req); err != nil {
		h.logger.Error(ctx, "Request validation failed", logger.String("type", requestType), logger.Error(err))
		h.responseHelper.ValidationError(c, gin.H{"validation_errors": err.Error()}, "Validation failed")
		return false
	}
	return true
}

// validateBusinessLogic validates request using business logic validator
// Returns true if successful, false if error response was sent
func (h *MealHandler) validateBusinessLogic(c *gin.Context, ctx context.Context, err error, requestType string) bool {
	if err != nil {
		h.logger.Error(ctx, "Business validation failed", logger.String("type", requestType), logger.Error(err))
		h.responseHelper.ValidationError(c, gin.H{"validation_errors": err.Error()}, "Validation failed")
		return false
	}
	return true
}

// handleServiceError handles service errors and sends appropriate response
// Returns true if error was handled, false if no error
func (h *MealHandler) handleServiceError(
	c *gin.Context,
	ctx context.Context,
	err error,
	operation string,
) bool {
	if err == nil {
		return false
	}

	h.logger.Error(ctx, "Service operation failed", logger.String("operation", operation), logger.Error(err))

	// Check for specific error types
	errMsg := err.Error()
	if errMsg == "template not found or access denied" {
		h.responseHelper.NotFound(c, gin.H{"error": "Meal template not found"}, "Meal template not found")
		return true
	}

	// Default to internal error
	h.responseHelper.InternalError(c, gin.H{"details": errMsg}, "Operation failed")
	return true
}

// CreateTemplate handles meal template creation
func (h *MealHandler) CreateTemplate(c *gin.Context) {
	ctx := middleware.GetContext(c)

	// Get user ID from context
	userIDStr, ok := h.getUserIDFromContext(c, ctx)
	if !ok {
		return
	}

	// Bind request
	var req request.CreateMealTemplateRequest
	if !h.bindRequest(c, ctx, &req, "CreateMealTemplateRequest") {
		return
	}

	// Validate request
	if !h.validateRequest(c, ctx, &req, "CreateMealTemplateRequest") {
		return
	}

	// Business logic validation
	if !h.validateBusinessLogic(c, ctx, h.mealValidator.ValidateCreateRequest(ctx, &req), "CreateMealTemplateRequest") {
		return
	}

	// Call service
	template, err := h.mealService.CreateTemplate(ctx, userIDStr, &req)
	if h.handleServiceError(c, ctx, err, "create meal template") {
		return
	}

	// Convert to response and send success
	templateResponse := mealTemplateToResponse(template)
	h.logger.Info(ctx, "Meal template created successfully")
	h.responseHelper.Created(c, templateResponse, "Meal template created successfully")
}

// ListTemplates handles listing meal templates
func (h *MealHandler) ListTemplates(c *gin.Context) {
	ctx := middleware.GetContext(c)

	// Get user ID from context
	userIDStr, ok := h.getUserIDFromContext(c, ctx)
	if !ok {
		return
	}

	// Get and validate query parameters
	mealType := c.Query("mealType")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Call service
	templates, err := h.mealService.ListTemplates(ctx, userIDStr, mealType, limit, offset)
	if h.handleServiceError(c, ctx, err, "list meal templates") {
		return
	}

	// Convert to response
	templateResponses := make([]response.MealTemplateResponse, len(templates))
	for i, template := range templates {
		templateResponses[i] = mealTemplateToResponse(template)
	}

	h.logger.Info(ctx, "Meal templates listed successfully")
	h.responseHelper.Success(c, templateResponses, "Meal templates listed successfully")
}

// GetTemplate handles getting a meal template
func (h *MealHandler) GetTemplate(c *gin.Context) {
	ctx := middleware.GetContext(c)

	// Get user ID from context
	userIDStr, ok := h.getUserIDFromContext(c, ctx)
	if !ok {
		return
	}

	// Get template ID from params
	templateID, ok := h.getTemplateIDFromParams(c, ctx)
	if !ok {
		return
	}

	// Call service
	template, err := h.mealService.GetTemplate(ctx, userIDStr, templateID)
	if h.handleServiceError(c, ctx, err, "get meal template") {
		return
	}

	// Convert to response and send success
	templateResponse := mealTemplateToResponse(template)
	h.logger.Info(ctx, "Meal template retrieved successfully")
	h.responseHelper.Success(c, templateResponse, "Meal template retrieved successfully")
}

// AddFoodToTemplate handles adding food items to a meal template
func (h *MealHandler) AddFoodToTemplate(c *gin.Context) {
	ctx := middleware.GetContext(c)

	// Get user ID from context
	userIDStr, ok := h.getUserIDFromContext(c, ctx)
	if !ok {
		return
	}

	// Get template ID from params
	templateID, ok := h.getTemplateIDFromParams(c, ctx)
	if !ok {
		return
	}

	// Bind request
	var req request.AddFoodToTemplateRequest
	if !h.bindRequest(c, ctx, &req, "AddFoodToTemplateRequest") {
		return
	}

	// Validate request
	if !h.validateRequest(c, ctx, &req, "AddFoodToTemplateRequest") {
		return
	}

	// Business logic validation
	if !h.validateBusinessLogic(c, ctx, h.mealValidator.ValidateAddFoodRequest(ctx, &req), "AddFoodToTemplateRequest") {
		return
	}

	// Call service
	template, err := h.mealService.AddFoodToTemplate(ctx, userIDStr, templateID, &req)
	if h.handleServiceError(c, ctx, err, "add food to template") {
		return
	}

	// Convert to response and send success
	templateResponse := mealTemplateToResponse(template)
	h.logger.Info(ctx, "Food items added to template successfully")
	h.responseHelper.Success(c, templateResponse, "Food items added to template successfully")
}

// UpdateTemplate handles meal template update
func (h *MealHandler) UpdateTemplate(c *gin.Context) {
	ctx := middleware.GetContext(c)

	// Get user ID from context
	userIDStr, ok := h.getUserIDFromContext(c, ctx)
	if !ok {
		return
	}

	// Get template ID from params
	templateID, ok := h.getTemplateIDFromParams(c, ctx)
	if !ok {
		return
	}

	// Bind request
	var req request.UpdateMealTemplateRequest
	if !h.bindRequest(c, ctx, &req, "UpdateMealTemplateRequest") {
		return
	}

	// Validate request
	if !h.validateRequest(c, ctx, &req, "UpdateMealTemplateRequest") {
		return
	}

	// Business logic validation
	if !h.validateBusinessLogic(c, ctx, h.mealValidator.ValidateUpdateRequest(ctx, &req), "UpdateMealTemplateRequest") {
		return
	}

	// Call service
	template, err := h.mealService.UpdateTemplate(ctx, userIDStr, templateID, &req)
	if h.handleServiceError(c, ctx, err, "update meal template") {
		return
	}

	// Convert to response and send success
	templateResponse := mealTemplateToResponse(template)
	h.logger.Info(ctx, "Meal template updated successfully")
	h.responseHelper.Success(c, templateResponse, "Meal template updated successfully")
}

// DeleteTemplate handles meal template deletion
func (h *MealHandler) DeleteTemplate(c *gin.Context) {
	ctx := middleware.GetContext(c)

	// Get user ID from context
	userIDStr, ok := h.getUserIDFromContext(c, ctx)
	if !ok {
		return
	}

	// Get template ID from params
	templateID, ok := h.getTemplateIDFromParams(c, ctx)
	if !ok {
		return
	}

	// Call service
	if h.handleServiceError(c, ctx, h.mealService.DeleteTemplate(ctx, userIDStr, templateID), "delete meal template") {
		return
	}

	// Send success response
	h.logger.Info(ctx, "Meal template deleted successfully")
	h.responseHelper.Success(c, gin.H{"message": "Meal template deleted successfully"}, "Meal template deleted successfully")
}

// mealTemplateToResponse converts a domain MealTemplate to a response MealTemplateResponse
func mealTemplateToResponse(template *domain.MealTemplate) response.MealTemplateResponse {
	// Convert food items
	foodItems := make([]response.MealTemplateFoodItemResponse, len(template.FoodItems))
	for i, foodItem := range template.FoodItems {
		foodItems[i] = response.MealTemplateFoodItemResponse{
			FoodItemID:  foodItem.FoodItemID.Hex(),
			FoodName:    foodItem.FoodName,
			ServingUnit: foodItem.ServingUnit,
			Amount:      foodItem.Amount,
			Calories:    foodItem.Calories,
			Macros: response.MacroNutrientsResponse{
				Protein:       foodItem.Macros.Protein,
				Carbohydrates: foodItem.Macros.Carbohydrates,
				Fat:           foodItem.Macros.Fat,
				Fiber:         foodItem.Macros.Fiber,
				Sugar:         foodItem.Macros.Sugar,
			},
			Micros: response.MicroNutrientsResponse{
				VitaminA:  foodItem.Micros.VitaminA,
				VitaminC:  foodItem.Micros.VitaminC,
				Calcium:   foodItem.Micros.Calcium,
				Iron:      foodItem.Micros.Iron,
				Sodium:    foodItem.Micros.Sodium,
				Potassium: foodItem.Micros.Potassium,
			},
		}
	}

	// Build response
	return response.MealTemplateResponse{
		ID:            template.ID.Hex(),
		UserID:        template.UserID.Hex(),
		Name:          template.Name,
		Description:   template.Description,
		MealType:      template.MealType,
		FoodItems:     foodItems,
		TotalCalories: template.TotalCalories,
		TotalMacros: response.MacroNutrientsResponse{
			Protein:       template.TotalMacros.Protein,
			Carbohydrates: template.TotalMacros.Carbohydrates,
			Fat:           template.TotalMacros.Fat,
			Fiber:         template.TotalMacros.Fiber,
			Sugar:         template.TotalMacros.Sugar,
		},
		TotalMicros: response.MicroNutrientsResponse{
			VitaminA:  template.TotalMicros.VitaminA,
			VitaminC:  template.TotalMicros.VitaminC,
			Calcium:   template.TotalMicros.Calcium,
			Iron:      template.TotalMicros.Iron,
			Sodium:    template.TotalMicros.Sodium,
			Potassium: template.TotalMicros.Potassium,
		},
		Tags:      template.Tags,
		IsPublic:  template.IsPublic,
		CreatedAt: template.CreatedAt,
		UpdatedAt: template.UpdatedAt,
	}
}
