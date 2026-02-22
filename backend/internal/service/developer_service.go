package service

import (
	"context"
	"fmt"
	"log"

	gh "github.com/Madhur/GithubScoreEval/backend/internal/github"
	"github.com/Madhur/GithubScoreEval/backend/internal/model"
	"github.com/Madhur/GithubScoreEval/backend/internal/repository"
)

type DeveloperService struct {
	ghClient *gh.Client
	devRepo  repository.DeveloperRepository
}

func NewDeveloperService(ghClient *gh.Client, devRepo repository.DeveloperRepository) *DeveloperService {
	return &DeveloperService{
		ghClient: ghClient,
		devRepo:  devRepo,
	}
}

func (s *DeveloperService) FetchAndStore(ctx context.Context, username string, accessToken string) (*model.Developer, error) {
	log.Printf("Fetching GitHub data for user: %s", username)

	client := s.ghClient
	if accessToken != "" {
		client = gh.NewClient(accessToken)
	}

	developer, err := client.FetchDeveloperData(username)
	if err != nil {
		return nil, fmt.Errorf("fetching developer data: %w", err)
	}

	if err := s.devRepo.Save(ctx, developer); err != nil {
		return nil, fmt.Errorf("storing developer data: %w", err)
	}

	log.Printf("Stored data for %s: %d commits, %d PRs, %d repos",
		username, developer.Metrics.TotalCommits, developer.Metrics.TotalPRs, developer.Metrics.ReposContributed)

	return developer, nil
}

func (s *DeveloperService) GetByUsername(ctx context.Context, username string) (*model.Developer, error) {
	return s.devRepo.GetByUsername(ctx, username)
}

func (s *DeveloperService) GetAll(ctx context.Context) ([]*model.Developer, error) {
	return s.devRepo.GetAll(ctx)
}
