package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nutrient_be/internal/pkg/logger"
	"nutrient_be/internal/service"
)

// ReportHandler handles report endpoints
type ReportHandler struct {
	reportService *service.ReportService
	logger        logger.Logger
}

// NewReportHandler creates a new report handler
func NewReportHandler(reportService *service.ReportService, log logger.Logger) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
		logger:        log,
	}
}

// Weekly handles weekly reports
func (h *ReportHandler) Weekly(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Weekly reports not implemented yet"})
}

// Monthly handles monthly reports
func (h *ReportHandler) Monthly(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Monthly reports not implemented yet"})
}
