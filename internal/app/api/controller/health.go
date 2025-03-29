package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-gin-boilerplate/internal/app/api/binding"
	"go-gin-boilerplate/internal/service"
)

type HealthController struct {
	healthService service.HealthService
}

func NewHealthController(healthService service.HealthService) *HealthController {
	return &HealthController{
		healthService: healthService,
	}
}

// Check godoc
// @Summary Health check endpoint
// @Description Returns the health status of the application
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} binding.HealthResponse
// @Router /health [get]
func (c *HealthController) Check(ctx *gin.Context) {
	status := c.healthService.Check()
	response := binding.HealthResponse{
		Status:    status["status"].(string),
		Uptime:    status["uptime"].(string),
		Timestamp: status["timestamp"].(int64),
	}
	ctx.JSON(http.StatusOK, response)
} 