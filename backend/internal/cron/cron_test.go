package cron

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Madhur/GithubScoreEval/backend/internal/model"
)

// ---- Mock repositories ----

type mockDeveloperRepo struct {
	developers []*model.Developer
	saved      []*model.Developer
}

func (m *mockDeveloperRepo) Save(_ context.Context, dev *model.Developer) error {
	m.saved = append(m.saved, dev)
	return nil
}
func (m *mockDeveloperRepo) GetByUsername(_ context.Context, username string) (*model.Developer, error) {
	for _, d := range m.developers {
		if d.Username == username {
			return d, nil
		}
	}
	return nil, fmt.Errorf("not found")
}
func (m *mockDeveloperRepo) GetAll(_ context.Context) ([]*model.Developer, error) {
	return m.developers, nil
}
func (m *mockDeveloperRepo) Delete(_ context.Context, _ string) error { return nil }

type mockScoreRepo struct {
	scores []*model.Score
	saved  []*model.Score
}

func (m *mockScoreRepo) Save(_ context.Context, score *model.Score) error {
	m.saved = append(m.saved, score)
	// Also add to main scores list for GetByUsername lookups
	m.scores = append(m.scores, score)
	return nil
}
func (m *mockScoreRepo) GetByUsername(_ context.Context, username string) (*model.Score, error) {
	for _, s := range m.scores {
		if s.Username == username {
			return s, nil
		}
	}
	return nil, fmt.Errorf("not found")
}
func (m *mockScoreRepo) GetAll(_ context.Context) ([]*model.Score, error) {
	return m.scores, nil
}
func (m *mockScoreRepo) Delete(_ context.Context, _ string) error { return nil }

type mockRankingRepo struct {
	saved *model.Ranking
}

func (m *mockRankingRepo) Save(_ context.Context, ranking *model.Ranking) error {
	m.saved = ranking
	return nil
}
func (m *mockRankingRepo) GetLatest(_ context.Context) (*model.Ranking, error) {
	if m.saved == nil {
		return nil, fmt.Errorf("no rankings")
	}
	return m.saved, nil
}
func (m *mockRankingRepo) GetByDate(_ context.Context, _ string) (*model.Ranking, error) {
	return m.saved, nil
}

// ---- Tests ----

func TestRefreshJob_SkipsRecentDevelopers(t *testing.T) {
	devRepo := &mockDeveloperRepo{
		developers: []*model.Developer{
			{Username: "recent-user", FetchedAt: time.Now().Add(-1 * time.Hour)},
		},
	}
	scoreRepo := &mockScoreRepo{
		scores: []*model.Score{
			{Username: "recent-user", WeightedScore: 50.0},
		},
	}
	rankingRepo := &mockRankingRepo{}

	job := &RefreshJob{
		devRepo:         devRepo,
		scoreRepo:       scoreRepo,
		rankingRepo:     rankingRepo,
		ghClient:        nil, // won't be called since we skip
		mlClient:        nil,
		engine:          nil,
		refreshHoursStr: "6",
	}

	// Developer was updated 1 hour ago, threshold is 6 hours → should skip
	if !job.isRecent(devRepo.developers[0], 6*time.Hour) {
		t.Error("Expected recent developer to be skipped")
	}
}

func TestRefreshJob_DoesNotSkipStaleDevelopers(t *testing.T) {
	staleDev := &model.Developer{
		Username:  "stale-user",
		FetchedAt: time.Now().Add(-12 * time.Hour),
	}

	job := &RefreshJob{
		refreshHoursStr: "6",
	}

	// Developer was updated 12 hours ago, threshold is 6 hours → should NOT skip
	if job.isRecent(staleDev, 6*time.Hour) {
		t.Error("Expected stale developer to NOT be skipped")
	}
}

func TestRefreshJob_DoesNotSkipZeroTime(t *testing.T) {
	devWithNoFetchedAt := &model.Developer{
		Username:  "new-user",
		FetchedAt: time.Time{}, // zero value
	}

	job := &RefreshJob{
		refreshHoursStr: "6",
	}

	if job.isRecent(devWithNoFetchedAt, 6*time.Hour) {
		t.Error("Expected developer with zero FetchedAt to NOT be skipped")
	}
}

func TestRefreshJob_GetThreshold(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"6", 6 * time.Hour},
		{"12", 12 * time.Hour},
		{"1", 1 * time.Hour},
		{"0", 6 * time.Hour},   // invalid → default
		{"-1", 6 * time.Hour},  // invalid → default
		{"abc", 6 * time.Hour}, // invalid → default
	}

	for _, tc := range tests {
		job := &RefreshJob{refreshHoursStr: tc.input}
		got := job.getThreshold()
		if got != tc.expected {
			t.Errorf("getThreshold(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}

func TestBuildRanking_SortsAndAssignsRanks(t *testing.T) {
	scores := []*model.Score{
		{Username: "low", WeightedScore: 20.0, MLImpactScore: 15.0},
		{Username: "high", WeightedScore: 80.0, MLImpactScore: 90.0},
		{Username: "mid", WeightedScore: 50.0, MLImpactScore: 45.0},
	}

	ranking := BuildRanking(scores)

	if ranking.TotalDevelopers != 3 {
		t.Errorf("Expected 3 total developers, got %d", ranking.TotalDevelopers)
	}

	if ranking.Rankings[0].Username != "high" {
		t.Errorf("Expected rank 1 to be 'high', got %s", ranking.Rankings[0].Username)
	}
	if ranking.Rankings[0].Rank != 1 {
		t.Errorf("Expected rank 1, got %d", ranking.Rankings[0].Rank)
	}

	if ranking.Rankings[1].Username != "mid" {
		t.Errorf("Expected rank 2 to be 'mid', got %s", ranking.Rankings[1].Username)
	}

	if ranking.Rankings[2].Username != "low" {
		t.Errorf("Expected rank 3 to be 'low', got %s", ranking.Rankings[2].Username)
	}

	// Check MLScore is propagated
	if ranking.Rankings[0].MLScore != 90.0 {
		t.Errorf("Expected MLScore 90.0 for 'high', got %f", ranking.Rankings[0].MLScore)
	}
}

func TestBuildRanking_PercentileCalculation(t *testing.T) {
	scores := []*model.Score{
		{Username: "a", WeightedScore: 100.0},
		{Username: "b", WeightedScore: 75.0},
		{Username: "c", WeightedScore: 50.0},
		{Username: "d", WeightedScore: 25.0},
	}

	ranking := BuildRanking(scores)

	// Rank 1 of 4 → percentile = (4-1)/4 * 100 = 75%
	if scores[0].Percentile != 75.0 {
		t.Errorf("Expected percentile 75.0 for rank 1, got %f", scores[0].Percentile)
	}

	// Rank 4 of 4 → percentile = (4-4)/4 * 100 = 0%
	if scores[3].Percentile != 0.0 {
		t.Errorf("Expected percentile 0.0 for rank 4, got %f", scores[3].Percentile)
	}

	// Snapshot date should be today
	today := time.Now().Format("2006-01-02")
	if ranking.SnapshotDate != today {
		t.Errorf("Expected snapshot date %s, got %s", today, ranking.SnapshotDate)
	}
}

func TestScheduler_AddJobInvalidSchedule(t *testing.T) {
	s := NewScheduler()
	err := s.AddJob("invalid cron expression!!!!", func() {})
	if err == nil {
		t.Error("Expected error for invalid cron expression")
	}
}

func TestScheduler_AddJobValidSchedule(t *testing.T) {
	s := NewScheduler()
	err := s.AddJob("*/5 * * * *", func() {})
	if err != nil {
		t.Errorf("Expected no error for valid cron expression, got %v", err)
	}
}
