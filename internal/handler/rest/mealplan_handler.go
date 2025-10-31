package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nutrient_be/internal/pkg/logger"
	"nutrient_be/internal/service"
)

// MealPlanHandler handles meal plan endpoints
type MealPlanHandler struct {
	mealPlanService *service.MealPlanService
	logger          logger.Logger
}

// NewMealPlanHandler creates a new meal plan handler
func NewMealPlanHandler(mealPlanService *service.MealPlanService, log logger.Logger) *MealPlanHandler {
	return &MealPlanHandler{
		mealPlanService: mealPlanService,
		logger:          log,
	}
}

// Create handles meal plan creation
func (h *MealPlanHandler) Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Meal plan creation not implemented yet"})
}

// List handles listing meal plans
func (h *MealPlanHandler) List(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Meal plan listing not implemented yet"})
}

// Get handles getting a meal plan
func (h *MealPlanHandler) Get(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Meal plan get not implemented yet"})
}

// Update handles meal plan update
func (h *MealPlanHandler) Update(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Meal plan update not implemented yet"})
}

// Delete handles meal plan deletion
func (h *MealPlanHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Meal plan deletion not implemented yet"})
}
