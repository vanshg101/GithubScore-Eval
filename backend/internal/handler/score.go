package handler

import (
	"net/http"

	"github.com/Madhur/GithubScoreEval/backend/internal/model"
	"github.com/Madhur/GithubScoreEval/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// ScoreHandler handles score computation and retrieval for individual developers.
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

	c.JSON(http.StatusOK, scoreToJSON(score))
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

	c.JSON(http.StatusOK, scoreToJSON(score))
}

// scoreToJSON builds a consistent JSON response for a score.
func scoreToJSON(score *model.Score) gin.H {
	return gin.H{
		"username":         score.Username,
		"weighted_score":   score.WeightedScore,
		"ml_impact_score":  score.MLImpactScore,
		"indicator_scores": score.IndicatorScores,
		"percentile":       score.Percentile,
		"computed_at":      score.ComputedAt,
	}
}
