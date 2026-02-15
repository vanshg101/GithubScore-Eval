package cron

import (
	"context"
	"log"
	"strconv"
	"time"

	gh "github.com/Madhur/GithubScoreEval/backend/internal/github"
	"github.com/Madhur/GithubScoreEval/backend/internal/mlclient"
	"github.com/Madhur/GithubScoreEval/backend/internal/model"
	"github.com/Madhur/GithubScoreEval/backend/internal/repository"
	"github.com/Madhur/GithubScoreEval/backend/internal/scoring"
)

// RefreshJob re-fetches and re-scores all tracked developers,
// skipping those updated within the threshold window.
type RefreshJob struct {
	devRepo         repository.DeveloperRepository
	scoreRepo       repository.ScoreRepository
	rankingRepo     repository.RankingRepository
	ghClient        *gh.Client
	mlClient        *mlclient.Client
	engine          *scoring.Engine
	refreshHoursStr string
}

// NewRefreshJob creates a new RefreshJob.
func NewRefreshJob(
	devRepo repository.DeveloperRepository,
	scoreRepo repository.ScoreRepository,
	rankingRepo repository.RankingRepository,
	ghClient *gh.Client,
	mlClient *mlclient.Client,
	refreshHoursStr string,
) *RefreshJob {
	return &RefreshJob{
		devRepo:         devRepo,
		scoreRepo:       scoreRepo,
		rankingRepo:     rankingRepo,
		ghClient:        ghClient,
		mlClient:        mlClient,
		engine:          scoring.NewEngine(),
		refreshHoursStr: refreshHoursStr,
	}
}

// Run executes the refresh job: re-fetch, re-score, update rankings.
func (j *RefreshJob) Run() {
	start := time.Now()
	log.Println("[CRON] Refresh job started")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	threshold := j.getThreshold()

	// 1. Get all tracked developers
	developers, err := j.devRepo.GetAll(ctx)
	if err != nil {
		log.Printf("[CRON] Failed to list developers: %v", err)
		return
	}

	if len(developers) == 0 {
		log.Println("[CRON] No developers to refresh")
		return
	}

	log.Printf("[CRON] Found %d tracked developers", len(developers))

	refreshed := 0
	skipped := 0
	failed := 0
	var scores []*model.Score

	for _, dev := range developers {
		// 2. Skip recently updated developers
		if j.isRecent(dev, threshold) {
			log.Printf("[CRON] Skipping %s (updated %s ago)", dev.Username, time.Since(dev.FetchedAt).Round(time.Minute))
			skipped++

			// Still include existing score in ranking update
			existingScore, err := j.scoreRepo.GetByUsername(ctx, dev.Username)
			if err == nil {
				scores = append(scores, existingScore)
			}
			continue
		}

		// 3. Re-fetch developer data from GitHub
		log.Printf("[CRON] Refreshing %s...", dev.Username)
		freshDev, err := j.ghClient.FetchDeveloperData(dev.Username)
		if err != nil {
			log.Printf("[CRON] Failed to fetch %s: %v", dev.Username, err)
			failed++
			continue
		}

		if err := j.devRepo.Save(ctx, freshDev); err != nil {
			log.Printf("[CRON] Failed to store %s: %v", dev.Username, err)
			failed++
			continue
		}

		// 4. Re-score
		score := j.engine.Compute(dev.Username, &freshDev.Metrics)

		// 5. ML prediction (with fallback)
		if j.mlClient != nil {
			req := mlclient.MapMetrics(&freshDev.Metrics)
			mlScore, mlErr := j.mlClient.Predict(ctx, req)
			if mlErr != nil {
				log.Printf("[CRON] ML prediction failed for %s: %v", dev.Username, mlErr)
			} else {
				score.MLImpactScore = mlScore
			}
		}

		if err := j.scoreRepo.Save(ctx, score); err != nil {
			log.Printf("[CRON] Failed to save score for %s: %v", dev.Username, err)
			failed++
			continue
		}

		scores = append(scores, score)
		refreshed++
	}

	// 6. Update global rankings if we have scores
	if len(scores) > 0 {
		ranking := BuildRanking(scores)
		if err := j.rankingRepo.Save(ctx, ranking); err != nil {
			log.Printf("[CRON] Failed to save ranking snapshot: %v", err)
		} else {
			log.Printf("[CRON] Rankings updated with %d developers", len(scores))
		}
	}

	elapsed := time.Since(start).Round(time.Second)
	log.Printf("[CRON] Refresh job completed in %s — refreshed: %d, skipped: %d, failed: %d",
		elapsed, refreshed, skipped, failed)
}

// isRecent returns true if the developer was updated within the threshold.
func (j *RefreshJob) isRecent(dev *model.Developer, threshold time.Duration) bool {
	if dev.FetchedAt.IsZero() {
		return false
	}
	return time.Since(dev.FetchedAt) < threshold
}

// getThreshold parses the refresh hours string to a duration.
func (j *RefreshJob) getThreshold() time.Duration {
	hours, err := strconv.Atoi(j.refreshHoursStr)
	if err != nil || hours <= 0 {
		hours = 6 // default to 6 hours
	}
	return time.Duration(hours) * time.Hour
}
