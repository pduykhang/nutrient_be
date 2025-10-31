package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nutrient_be/internal/pkg/logger"
	"nutrient_be/internal/service"
)

// FoodHandler handles food-related endpoints
type FoodHandler struct {
	foodService *service.FoodService
	logger      logger.Logger
}

// NewFoodHandler creates a new food handler
func NewFoodHandler(foodService *service.FoodService, log logger.Logger) *FoodHandler {
	return &FoodHandler{
		foodService: foodService,
		logger:      log,
	}
}

// Create handles food creation
func (h *FoodHandler) Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Food creation not implemented yet"})
}

// Search handles food search
func (h *FoodHandler) Search(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Food search not implemented yet"})
}

// Get handles getting a food item
func (h *FoodHandler) Get(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Food get not implemented yet"})
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
