package router

import (
	"github.com/Madhur/GithubScoreEval/backend/internal/config"
	"github.com/Madhur/GithubScoreEval/backend/internal/handler"
	"github.com/Madhur/GithubScoreEval/backend/internal/middleware"
	"github.com/Madhur/GithubScoreEval/backend/internal/repository"
	"github.com/Madhur/GithubScoreEval/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func Setup(cfg *config.Config, userRepo repository.UserRepository, devService *service.DeveloperService) *gin.Engine {
	router := gin.New()

	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	healthHandler := handler.NewHealthHandler()
	router.GET("/health", healthHandler.HealthCheck)

	authHandler := handler.NewAuthHandler(cfg, userRepo)
	authGroup := router.Group("/auth")
	{
		authGroup.GET("/github/login", authHandler.GitHubLogin)
		authGroup.GET("/github/callback", authHandler.GitHubCallback)
		authGroup.POST("/logout", authHandler.Logout)
	}

	authGroup.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		authGroup.GET("/me", authHandler.GetCurrentUser)
	}

	devHandler := handler.NewDeveloperHandler(devService, userRepo)
	api := router.Group("/api")
	api.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		api.GET("/developers", devHandler.ListDevelopers)
		api.GET("/developers/:username", devHandler.GetDeveloper)
		api.POST("/developers/:username/fetch", devHandler.FetchDeveloper)
	}

	return router
}
