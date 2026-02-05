package repository

import (
	"context"

	"github.com/Madhur/GithubScoreEval/backend/internal/model"
)

type UserRepository interface {
	Save(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByGitHubID(ctx context.Context, githubID int64) (*model.User, error)
	Delete(ctx context.Context, id string) error
}

type DeveloperRepository interface {
	Save(ctx context.Context, dev *model.Developer) error
	GetByUsername(ctx context.Context, username string) (*model.Developer, error)
	GetAll(ctx context.Context) ([]*model.Developer, error)
	Delete(ctx context.Context, username string) error
}

type ScoreRepository interface {
	Save(ctx context.Context, score *model.Score) error
	GetByUsername(ctx context.Context, username string) (*model.Score, error)
	GetAll(ctx context.Context) ([]*model.Score, error)
	Delete(ctx context.Context, username string) error
}

type RankingRepository interface {
	Save(ctx context.Context, ranking *model.Ranking) error
	GetLatest(ctx context.Context) (*model.Ranking, error)
	GetByDate(ctx context.Context, date string) (*model.Ranking, error)
}
