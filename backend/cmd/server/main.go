package main

import (
	"log"

	"github.com/Madhur/GithubScoreEval/backend/internal/config"
	cronpkg "github.com/Madhur/GithubScoreEval/backend/internal/cron"
	fs "github.com/Madhur/GithubScoreEval/backend/internal/firestore"
	gh "github.com/Madhur/GithubScoreEval/backend/internal/github"
	"github.com/Madhur/GithubScoreEval/backend/internal/mlclient"
	"github.com/Madhur/GithubScoreEval/backend/internal/repository"
	"github.com/Madhur/GithubScoreEval/backend/internal/router"
	"github.com/Madhur/GithubScoreEval/backend/internal/service"
)

func main() {
	cfg := config.Load()
	log.Printf("Starting server in %s mode on port %s", cfg.Environment, cfg.Port)

	fsClient := fs.NewClient(cfg.GCPProjectID, cfg.FirestoreCredentials)
	defer fsClient.Close()

	userRepo := repository.NewFirestoreUserRepo(fsClient)
	devRepo := repository.NewFirestoreDeveloperRepo(fsClient)

	scoreRepo := repository.NewFirestoreScoreRepo(fsClient)
	rankingRepo := repository.NewFirestoreRankingRepo(fsClient)

	ghClient := gh.NewClient("")
	mlClient := mlclient.NewClient(cfg.MLServiceURL)
	log.Printf("ML service configured at %s", cfg.MLServiceURL)

	devService := service.NewDeveloperService(ghClient, devRepo)
	scoringService := service.NewScoringService(devRepo, scoreRepo, mlClient)
	rankingService := service.NewRankingService(devRepo, scoreRepo, rankingRepo, userRepo, ghClient, mlClient)

	// Cron scheduler for periodic data refresh
	refreshJob := cronpkg.NewRefreshJob(devRepo, scoreRepo, rankingRepo, ghClient, mlClient, cfg.CronRefreshHours)
	scheduler := cronpkg.NewScheduler()
	if err := scheduler.AddJob(cfg.CronSchedule, refreshJob.Run); err != nil {
		log.Printf("Failed to register cron job: %v", err)
	} else {
		scheduler.Start()
		defer scheduler.Stop()
		log.Printf("Cron refresh job scheduled: %s (skip if updated within %sh)", cfg.CronSchedule, cfg.CronRefreshHours)
	}

	r := router.Setup(cfg, userRepo, devService, scoringService, rankingService)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
