package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nutrient_be/internal/pkg/logger"
	"nutrient_be/internal/service"
)

// MealHandler handles meal template endpoints
type MealHandler struct {
	mealService *service.MealService
	logger      logger.Logger
}

// NewMealHandler creates a new meal handler
func NewMealHandler(mealService *service.MealService, log logger.Logger) *MealHandler {
	return &MealHandler{
		mealService: mealService,
		logger:      log,
	}
}

// CreateTemplate handles meal template creation
func (h *MealHandler) CreateTemplate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Meal template creation not implemented yet"})
}

// ListTemplates handles listing meal templates
func (h *MealHandler) ListTemplates(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Meal template listing not implemented yet"})
}

// GetTemplate handles getting a meal template
func (h *MealHandler) GetTemplate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Meal template get not implemented yet"})
}

// UpdateTemplate handles meal template update
func (h *MealHandler) UpdateTemplate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Meal template update not implemented yet"})
}

// DeleteTemplate handles meal template deletion
func (h *MealHandler) DeleteTemplate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Meal template deletion not implemented yet"})
}
