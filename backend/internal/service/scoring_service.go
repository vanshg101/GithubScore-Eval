package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Madhur/GithubScoreEval/backend/internal/mlclient"
	"github.com/Madhur/GithubScoreEval/backend/internal/model"
	"github.com/Madhur/GithubScoreEval/backend/internal/repository"
	"github.com/Madhur/GithubScoreEval/backend/internal/scoring"
)

// ScoringService orchestrates developer score computation and persistence.
type ScoringService struct {
	devRepo   repository.DeveloperRepository
	scoreRepo repository.ScoreRepository
	engine    *scoring.Engine
	mlClient  *mlclient.Client
}

// NewScoringService creates a new ScoringService.
func NewScoringService(
	devRepo repository.DeveloperRepository,
	scoreRepo repository.ScoreRepository,
	mlClient *mlclient.Client,
) *ScoringService {
	return &ScoringService{
		devRepo:   devRepo,
		scoreRepo: scoreRepo,
		engine:    scoring.NewEngine(),
		mlClient:  mlClient,
	}
}

// ComputeAndStore fetches developer metrics, scores them, calls the ML service, and stores the result.
func (s *ScoringService) ComputeAndStore(ctx context.Context, username string) (*model.Score, error) {
	dev, err := s.devRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("developer not found: %w", err)
	}

	score := s.engine.Compute(username, &dev.Metrics)

	// Call ML service for impact prediction (with fallback)
	score.MLImpactScore = s.predictMLScore(ctx, &dev.Metrics)

	if err := s.scoreRepo.Save(ctx, score); err != nil {
		return nil, fmt.Errorf("failed to save score: %w", err)
	}

	return score, nil
}

// predictMLScore calls the ML service and returns the impact score.
// Falls back to 0 if the ML service is unavailable.
func (s *ScoringService) predictMLScore(ctx context.Context, metrics *model.DeveloperMetrics) float64 {
	if s.mlClient == nil {
		return 0
	}

	req := mlclient.MapMetrics(metrics)
	score, err := s.mlClient.Predict(ctx, req)
	if err != nil {
		log.Printf("ML prediction failed (using fallback): %v", err)
		return 0
	}

	return score
}

// GetScore retrieves a previously computed score.
func (s *ScoringService) GetScore(ctx context.Context, username string) (*model.Score, error) {
	return s.scoreRepo.GetByUsername(ctx, username)
}

// GetAllScores retrieves all scores.
func (s *ScoringService) GetAllScores(ctx context.Context) ([]*model.Score, error) {
	return s.scoreRepo.GetAll(ctx)
}
