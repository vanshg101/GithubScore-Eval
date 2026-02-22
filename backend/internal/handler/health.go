package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler provides a simple liveness/readiness check endpoint.
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "github-score-eval-backend",
	})
}
