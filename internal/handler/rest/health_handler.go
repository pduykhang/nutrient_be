package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"nutrient_be/internal/pkg/logger"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db     *mongo.Client
	logger logger.Logger
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *mongo.Client, log logger.Logger) *HealthHandler {
	return &HealthHandler{
		db:     db,
		logger: log,
	}
}

// Liveness handles liveness probe
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "UP",
		"timestamp": time.Now().Unix(),
		"service":   "nutrient-api",
	})
}

// Readiness handles readiness probe
func (h *HealthHandler) Readiness(c *gin.Context) {
	// Check database connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx, nil); err != nil {
		h.logger.Error(ctx, "Database health check failed", logger.Error(err))
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "DOWN",
			"timestamp": time.Now().Unix(),
			"checks": gin.H{
				"database": "DOWN",
			},
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "UP",
		"timestamp": time.Now().Unix(),
		"checks": gin.H{
			"database": "UP",
		},
		"service": "nutrient-api",
	})
}
