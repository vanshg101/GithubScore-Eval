package cron

import (
	"sort"
	"time"

	"github.com/Madhur/GithubScoreEval/backend/internal/model"
)

// BuildRanking sorts scores descending, assigns ranks and percentiles, and returns a Ranking snapshot.
func BuildRanking(scores []*model.Score) *model.Ranking {
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].WeightedScore > scores[j].WeightedScore
	})

	total := len(scores)
	entries := make([]model.RankEntry, total)

	for i, sc := range scores {
		rank := i + 1
		percentile := float64(total-rank) / float64(total) * 100
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
