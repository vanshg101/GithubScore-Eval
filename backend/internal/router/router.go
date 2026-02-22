// Package router defines the HTTP routing table for the backend API.
package router

import (
	"github.com/Madhur/GithubScoreEval/backend/internal/config"
	"github.com/Madhur/GithubScoreEval/backend/internal/handler"
	"github.com/Madhur/GithubScoreEval/backend/internal/middleware"
	"github.com/Madhur/GithubScoreEval/backend/internal/repository"
	"github.com/Madhur/GithubScoreEval/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// Setup creates and returns a fully configured Gin engine with all routes,
// middleware (Logger, CORS, Recovery, AuthRequired), and handler registrations.
// Route groups:
//   - /health        — public liveness check
//   - /auth/*        — OAuth login/callback (public) + /me (authenticated)
//   - /api/*         — all endpoints require JWT authentication
func Setup(cfg *config.Config, userRepo repository.UserRepository, devService *service.DeveloperService, scoringService *service.ScoringService, rankingService *service.RankingService) *gin.Engine {
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
	scoreHandler := handler.NewScoreHandler(scoringService)
	rankingHandler := handler.NewRankingHandler(rankingService, userRepo)
	api := router.Group("/api")
	api.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		api.GET("/developers", devHandler.ListDevelopers)
		api.GET("/developers/:username", devHandler.GetDeveloper)
		api.POST("/developers/:username/fetch", devHandler.FetchDeveloper)
		api.POST("/developers/:username/score", scoreHandler.ComputeScore)
		api.GET("/developers/:username/score", scoreHandler.GetScore)
		api.POST("/compare", rankingHandler.CompareDevelopers)
		api.GET("/orgs/:org/evaluate", rankingHandler.EvaluateOrg)
		api.GET("/rankings", rankingHandler.GetLeaderboard)
	}

	return router
}
