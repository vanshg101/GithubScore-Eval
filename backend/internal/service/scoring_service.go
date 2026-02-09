package service

import (
	"context"
	"fmt"

	"github.com/Madhur/GithubScoreEval/backend/internal/model"
	"github.com/Madhur/GithubScoreEval/backend/internal/repository"
	"github.com/Madhur/GithubScoreEval/backend/internal/scoring"
)

// ScoringService orchestrates developer score computation and persistence.
type ScoringService struct {
	devRepo   repository.DeveloperRepository
	scoreRepo repository.ScoreRepository
	engine    *scoring.Engine
}

// NewScoringService creates a new ScoringService.
func NewScoringService(
	devRepo repository.DeveloperRepository,
	scoreRepo repository.ScoreRepository,
) *ScoringService {
	return &ScoringService{
		devRepo:   devRepo,
		scoreRepo: scoreRepo,
		engine:    scoring.NewEngine(),
	}
}

// ComputeAndStore fetches developer metrics, scores them, and stores the result.
func (s *ScoringService) ComputeAndStore(ctx context.Context, username string) (*model.Score, error) {
	dev, err := s.devRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("developer not found: %w", err)
	}

	score := s.engine.Compute(username, &dev.Metrics)

	if err := s.scoreRepo.Save(ctx, score); err != nil {
		return nil, fmt.Errorf("failed to save score: %w", err)
	}

	return score, nil
}

// GetScore retrieves a previously computed score.
func (s *ScoringService) GetScore(ctx context.Context, username string) (*model.Score, error) {
	return s.scoreRepo.GetByUsername(ctx, username)
}

// GetAllScores retrieves all scores.
func (s *ScoringService) GetAllScores(ctx context.Context) ([]*model.Score, error) {
	return s.scoreRepo.GetAll(ctx)
}
