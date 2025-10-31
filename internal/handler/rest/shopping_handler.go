package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nutrient_be/internal/pkg/logger"
	"nutrient_be/internal/service"
)

// ShoppingHandler handles shopping list endpoints
type ShoppingHandler struct {
	shoppingService *service.ShoppingService
	logger          logger.Logger
}

// NewShoppingHandler creates a new shopping handler
func NewShoppingHandler(shoppingService *service.ShoppingService, log logger.Logger) *ShoppingHandler {
	return &ShoppingHandler{
		shoppingService: shoppingService,
		logger:          log,
	}
}

// Generate handles shopping list generation
func (h *ShoppingHandler) Generate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Shopping list generation not implemented yet"})
}

// List handles listing shopping lists
func (h *ShoppingHandler) List(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Shopping list listing not implemented yet"})
}

// ToggleItem handles toggling shopping list item
func (h *ShoppingHandler) ToggleItem(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Shopping list item toggle not implemented yet"})
}
