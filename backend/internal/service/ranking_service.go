package service

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	gh "github.com/Madhur/GithubScoreEval/backend/internal/github"
	"github.com/Madhur/GithubScoreEval/backend/internal/mlclient"
	"github.com/Madhur/GithubScoreEval/backend/internal/model"
	"github.com/Madhur/GithubScoreEval/backend/internal/repository"
	"github.com/Madhur/GithubScoreEval/backend/internal/scoring"
)

// RankingService handles bulk scoring, ranking, percentile, comparison and org evaluation.
type RankingService struct {
	devRepo     repository.DeveloperRepository
	scoreRepo   repository.ScoreRepository
	rankingRepo repository.RankingRepository
	userRepo    repository.UserRepository
	engine      *scoring.Engine
	ghClient    *gh.Client
	mlClient    *mlclient.Client
}

// NewRankingService creates a new RankingService.
func NewRankingService(
	devRepo repository.DeveloperRepository,
	scoreRepo repository.ScoreRepository,
	rankingRepo repository.RankingRepository,
	userRepo repository.UserRepository,
	ghClient *gh.Client,
	mlClient *mlclient.Client,
) *RankingService {
	return &RankingService{
		devRepo:     devRepo,
		scoreRepo:   scoreRepo,
		rankingRepo: rankingRepo,
		userRepo:    userRepo,
		engine:      scoring.NewEngine(),
		ghClient:    ghClient,
		mlClient:    mlClient,
	}
}

// BulkScore scores multiple users, stores scores, and returns a ranking snapshot.
func (s *RankingService) BulkScore(ctx context.Context, usernames []string, accessToken string) (*model.Ranking, error) {
	var scores []*model.Score

	for _, username := range usernames {
		// Ensure developer data exists — fetch if missing
		dev, err := s.devRepo.GetByUsername(ctx, username)
		if err != nil {
			log.Printf("Developer %s not in store, fetching...", username)
			client := s.ghClient
			if accessToken != "" {
				client = gh.NewClient(accessToken)
			}
			fetchedDev, fetchErr := client.FetchDeveloperData(username)
			if fetchErr != nil {
				log.Printf("Failed to fetch %s: %v", username, fetchErr)
				continue
			}
			if saveErr := s.devRepo.Save(ctx, fetchedDev); saveErr != nil {
				log.Printf("Failed to store %s: %v", username, saveErr)
				continue
			}
			dev = fetchedDev
		}

		score := s.engine.Compute(username, &dev.Metrics)

		// Call ML service for impact prediction (with fallback)
		if s.mlClient != nil {
			req := mlclient.MapMetrics(&dev.Metrics)
			mlScore, mlErr := s.mlClient.Predict(ctx, req)
			if mlErr != nil {
				log.Printf("ML prediction failed for %s (using fallback): %v", username, mlErr)
			} else {
				score.MLImpactScore = mlScore
			}
		}

		if err := s.scoreRepo.Save(ctx, score); err != nil {
			log.Printf("Failed to save score for %s: %v", username, err)
			continue
		}
		scores = append(scores, score)
	}

	if len(scores) == 0 {
		return nil, fmt.Errorf("no developers could be scored")
	}

	ranking := buildRanking(scores)

	if err := s.rankingRepo.Save(ctx, ranking); err != nil {
		return nil, fmt.Errorf("saving ranking: %w", err)
	}

	return ranking, nil
}

// CompareDevelopers scores the given usernames and returns a comparison with ranking.
func (s *RankingService) CompareDevelopers(ctx context.Context, usernames []string, accessToken string) (*model.Ranking, error) {
	if len(usernames) < 2 {
		return nil, fmt.Errorf("need at least 2 usernames to compare")
	}
	if len(usernames) > 10 {
		return nil, fmt.Errorf("maximum 10 usernames per comparison")
	}
	return s.BulkScore(ctx, usernames, accessToken)
}

// EvaluateOrg fetches org members, scores them all, and returns a ranking.
func (s *RankingService) EvaluateOrg(ctx context.Context, org string, accessToken string) (*model.Ranking, error) {
	client := s.ghClient
	if accessToken != "" {
		client = gh.NewClient(accessToken)
	}

	members, err := client.FetchOrgMembers(org)
	if err != nil {
		return nil, fmt.Errorf("fetching org members: %w", err)
	}
	if len(members) == 0 {
		return nil, fmt.Errorf("no public members found for org %s", org)
	}

	usernames := make([]string, len(members))
	for i, m := range members {
		usernames[i] = m.Login
	}

	return s.BulkScore(ctx, usernames, accessToken)
}

// GetLeaderboard returns the latest ranking snapshot.
func (s *RankingService) GetLeaderboard(ctx context.Context, page, pageSize int) (*model.Ranking, error) {
	ranking, err := s.rankingRepo.GetLatest(ctx)
	if err != nil {
		// If no ranking exists, build one from all existing scores
		allScores, scoreErr := s.scoreRepo.GetAll(ctx)
		if scoreErr != nil || len(allScores) == 0 {
			return nil, fmt.Errorf("no rankings available")
		}
		ranking = buildRanking(allScores)
	}

	// Paginate
	total := len(ranking.Rankings)
	start := (page - 1) * pageSize
	if start >= total {
		ranking.Rankings = []model.RankEntry{}
		return ranking, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	ranking.Rankings = ranking.Rankings[start:end]

	return ranking, nil
}

// buildRanking sorts scores descending, assigns ranks and percentiles.
func buildRanking(scores []*model.Score) *model.Ranking {
	// Sort descending by WeightedScore
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].WeightedScore > scores[j].WeightedScore
	})

	total := len(scores)
	entries := make([]model.RankEntry, total)

	for i, sc := range scores {
		rank := i + 1
		// Percentile: percentage of developers scored below this rank
		percentile := float64(total-rank) / float64(total) * 100

		// Update the score's Percentile field in the score object
		sc.Percentile = percentile

		entries[i] = model.RankEntry{
			Rank:     rank,
			Username: sc.Username,
			Score:    sc.WeightedScore,
			MLScore:  sc.MLImpactScore,
		}
	}

	return &model.Ranking{
		SnapshotDate:    time.Now().Format("2006-01-02"),
		Rankings:        entries,
		TotalDevelopers: total,
		CreatedAt:       time.Now(),
	}
}
