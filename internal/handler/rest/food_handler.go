package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"nutrient_be/internal/domain"
	"nutrient_be/internal/dto/request"
	"nutrient_be/internal/dto/response"
	"nutrient_be/internal/handler/middleware"
	"nutrient_be/internal/pkg/logger"
	"nutrient_be/internal/service"
)

// FoodHandler handles food-related endpoints
type FoodHandler struct {
	foodService     *service.FoodService
	structValidator *validator.Validate
	logger          logger.Logger
	responseHelper  *middleware.ResponseHelper
}

// NewFoodHandler creates a new food handler
func NewFoodHandler(foodService *service.FoodService, log logger.Logger) *FoodHandler {
	return &FoodHandler{
		foodService:     foodService,
		structValidator: validator.New(),
		logger:          log,
		responseHelper:  middleware.NewResponseHelper(),
	}
}

// Create handles food creation
func (h *FoodHandler) Create(c *gin.Context) {
	ctx := middleware.GetContext(c)

	// Get user ID from context (set by auth middleware)
	userIDStr, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		h.logger.Error(ctx, "User ID not found in context")
		h.responseHelper.Unauthorized(c, gin.H{"error": "User not authenticated"}, "Authentication required")
		return
	}

	// Bind request
	var req request.CreateFoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(ctx, "Failed to bind create food request", logger.Error(err))
		h.responseHelper.BadRequest(c, gin.H{"details": err.Error()}, "Invalid request body")
		return
	}

	// Format/structure validation (struct tags)
	if err := h.structValidator.Struct(&req); err != nil {
		h.logger.Error(ctx, "Food request validation failed", logger.Error(err))
		h.responseHelper.ValidationError(c, gin.H{"validation_errors": err.Error()}, "Validation failed")
		return
	}

	// Call service - service will do business logic validation
	if err := h.foodService.CreateFood(c.Request.Context(), userIDStr, &req); err != nil {
		h.logger.Error(ctx, "Failed to create food", logger.Error(err))
		h.responseHelper.InternalError(c, gin.H{"details": err.Error()}, "Failed to create food")
		return
	}

	h.logger.Info(ctx, "Food created successfully")
	h.responseHelper.Created(c, gin.H{"message": "Food created successfully"}, "Food created successfully")
}

// Search handles food search
func (h *FoodHandler) Search(c *gin.Context) {
	var req request.SearchFoodRequest
	ctx := middleware.GetContext(c)
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error(ctx, "Failed to bind search food request", logger.Error(err))
		h.responseHelper.BadRequest(c, gin.H{"details": err.Error()}, "Invalid request body")
		return
	}

	// Call service - service will do business logic validation
	foods, err := h.foodService.SearchFood(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error(ctx, "Failed to search food", logger.Error(err))
		h.responseHelper.InternalError(c, gin.H{"details": err.Error()}, "Failed to search food")
		return
	}

	// Convert domain entities to response DTOs
	foodResponses := make([]response.FoodItemResponse, len(foods))
	for i, food := range foods {
		foodResponses[i] = foodItemToResponse(food)
	}

	h.logger.Info(ctx, "Food search successful")
	h.responseHelper.Success(c, foodResponses, "Food search successful")
}

// Get handles getting a food item
func (h *FoodHandler) Get(c *gin.Context) {
	ctx := middleware.GetContext(c)
	foodID := c.Param("id")

	// Validate food ID format
	if foodID == "" {
		h.logger.Error(ctx, "Food ID is required")
		h.responseHelper.BadRequest(c, gin.H{"error": "Food ID is required"}, "Food ID is required")
		return
	}

	// Call service to get food by ID
	food, err := h.foodService.GetFoodByID(c.Request.Context(), foodID)
	if err != nil {
		h.logger.Error(ctx, "Failed to get food by ID", logger.Error(err))
		// Check if it's a not found error
		if err.Error() == "food item not found" {
			h.responseHelper.NotFound(c, gin.H{"error": "Food item not found"}, "Food item not found")
			return
		}
		h.responseHelper.InternalError(c, gin.H{"details": err.Error()}, "Failed to get food")
		return
	}

	// Convert domain entity to response DTO
	foodResponse := foodItemToResponse(food)

	h.logger.Info(ctx, "Food retrieved successfully")
	h.responseHelper.Success(c, foodResponse, "Food retrieved successfully")
}

// Update handles food update
func (h *FoodHandler) Update(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Food update not implemented yet"})
}

// Delete handles food deletion
func (h *FoodHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Food deletion not implemented yet"})
}

// ImportExcel handles Excel import
func (h *FoodHandler) ImportExcel(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Excel import not implemented yet"})
}

// foodItemToResponse converts a domain FoodItem to a response FoodItemResponse
func foodItemToResponse(food *domain.FoodItem) response.FoodItemResponse {
	// Convert serving sizes
	servingSizes := make([]response.ServingSizeResponse, len(food.ServingSizes))
	for i, servingSize := range food.ServingSizes {
		servingSizes[i] = response.ServingSizeResponse{
			Unit:           servingSize.Unit,
			Amount:         servingSize.Amount,
			Description:    servingSize.Description,
			GramEquivalent: servingSize.GramEquivalent,
		}
	}

	// Build response
	foodResponse := response.FoodItemResponse{
		ID:          food.ID.Hex(),
		Name:        food.Name,
		SearchTerms: food.SearchTerms,
		Description: food.Description,
		Category:    food.Category,
		Macros: response.MacroNutrientsResponse{
			Protein:       food.Macros.Protein,
			Carbohydrates: food.Macros.Carbohydrates,
			Fat:           food.Macros.Fat,
			Fiber:         food.Macros.Fiber,
			Sugar:         food.Macros.Sugar,
		},
		Micros: response.MicroNutrientsResponse{
			VitaminA:  food.Micros.VitaminA,
			VitaminC:  food.Micros.VitaminC,
			Calcium:   food.Micros.Calcium,
			Iron:      food.Micros.Iron,
			Sodium:    food.Micros.Sodium,
			Potassium: food.Micros.Potassium,
		},
		ServingSizes: servingSizes,
		Calories:     food.Calories,
		CreatedBy:    food.CreatedBy.Hex(),
		Visibility:   food.Visibility,
		Source:       food.Source,
		ImageURL:     food.ImageURL,
		CreatedAt:    food.CreatedAt,
		UpdatedAt:    food.UpdatedAt,
	}

	return foodResponse
}
