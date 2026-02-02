package router

import (
	"github.com/Madhur/GithubScoreEval/backend/internal/handler"
	"github.com/Madhur/GithubScoreEval/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// Health check
	healthHandler := handler.NewHealthHandler()
	router.GET("/health", healthHandler.HealthCheck)



	return router
}
