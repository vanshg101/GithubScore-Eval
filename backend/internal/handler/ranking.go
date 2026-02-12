package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Madhur/GithubScoreEval/backend/internal/repository"
	"github.com/Madhur/GithubScoreEval/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type RankingHandler struct {
	rankingService *service.RankingService
	userRepo       repository.UserRepository
}

func NewRankingHandler(rankingService *service.RankingService, userRepo repository.UserRepository) *RankingHandler {
	return &RankingHandler{
		rankingService: rankingService,
		userRepo:       userRepo,
	}
}

// getAccessToken extracts the user's GitHub access token from auth context.
func (h *RankingHandler) getAccessToken(c *gin.Context) string {
	userID, _ := c.Get("user_id")
	if uid, ok := userID.(string); ok {
		user, err := h.userRepo.GetByID(c.Request.Context(), uid)
		if err == nil && user.AccessToken != "" {
			return user.AccessToken
		}
		if err != nil {
			log.Printf("could not retrieve user access token: %v", err)
		}
	}
	return ""
}

// CompareDevelopers handles POST /api/compare
// Body: { "usernames": ["user1", "user2", ...] }
func (h *RankingHandler) CompareDevelopers(c *gin.Context) {
	var req struct {
		Usernames []string `json:"usernames" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide a list of usernames"})
		return
	}

	accessToken := h.getAccessToken(c)

	ranking, err := h.rankingService.CompareDevelopers(c.Request.Context(), req.Usernames, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ranking)
}

// EvaluateOrg handles GET /api/orgs/:org/evaluate
func (h *RankingHandler) EvaluateOrg(c *gin.Context) {
	org := c.Param("org")
	if org == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "org name is required"})
		return
	}

	accessToken := h.getAccessToken(c)

	ranking, err := h.rankingService.EvaluateOrg(c.Request.Context(), org, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ranking)
}

// GetLeaderboard handles GET /api/rankings
func (h *RankingHandler) GetLeaderboard(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	ranking, err := h.rankingService.GetLeaderboard(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ranking)
}
