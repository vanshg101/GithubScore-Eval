package handler

import (
	"net/http"

	"github.com/Madhur/GithubScoreEval/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ScoreHandler struct {
	scoringService *service.ScoringService
}

func NewScoreHandler(scoringService *service.ScoringService) *ScoreHandler {
	return &ScoreHandler{scoringService: scoringService}
}

// ComputeScore computes (or recomputes) the score for a developer.
// POST /api/developers/:username/score
func (h *ScoreHandler) ComputeScore(c *gin.Context) {
	username := c.Param("username")

	score, err := h.scoringService.ComputeAndStore(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":         score.Username,
		"weighted_score":   score.WeightedScore,
		"indicator_scores": score.IndicatorScores,
		"computed_at":      score.ComputedAt,
	})
}

// GetScore retrieves the stored score for a developer.
// GET /api/developers/:username/score
func (h *ScoreHandler) GetScore(c *gin.Context) {
	username := c.Param("username")

	score, err := h.scoringService.GetScore(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "score not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":         score.Username,
		"weighted_score":   score.WeightedScore,
		"indicator_scores": score.IndicatorScores,
		"computed_at":      score.ComputedAt,
	})
}
