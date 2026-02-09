package scoring

import (
	"math"

	"github.com/Madhur/GithubScoreEval/backend/internal/model"
)

type Indicator struct {
	Name   string
	Weight float64
	Max    float64
	Extract func(m *model.DeveloperMetrics) float64
}

func DefaultIndicators() []Indicator {
	return []Indicator{
		{
			Name:   "total_commits",
			Weight: 0.15,
			Max:    500,
			Extract: func(m *model.DeveloperMetrics) float64 {
				return float64(m.TotalCommits)
			},
		},
		{
			Name:   "pr_merge_rate",
			Weight: 0.12,
			Max:    1.0,
			Extract: func(m *model.DeveloperMetrics) float64 {
				if m.TotalPRs == 0 {
					return 0
				}
				return float64(m.MergedPRs) / float64(m.TotalPRs)
			},
		},
		{
			Name:   "issue_resolution_ratio",
			Weight: 0.08,
			Max:    1.0,
			Extract: func(m *model.DeveloperMetrics) float64 {
				if m.TotalIssuesOpened == 0 {
					return 0
				}
				return float64(m.TotalIssuesClosed) / float64(m.TotalIssuesOpened)
			},
		},
		{
			Name:   "code_review_participation",
			Weight: 0.10,
			Max:    100,
			Extract: func(m *model.DeveloperMetrics) float64 {
				return float64(m.ReviewComments)
			},
		},
		{
			Name:   "contribution_consistency",
			Weight: 0.12,
			Max:    1.0,
			Extract: func(m *model.DeveloperMetrics) float64 {
				return float64(m.ActiveWeeks) / 52.0
			},
		},
		{
			Name:   "repo_diversity",
			Weight: 0.08,
			Max:    20,
			Extract: func(m *model.DeveloperMetrics) float64 {
				return float64(m.ReposContributed)
			},
		},
		{
			Name:   "stars_earned",
			Weight: 0.08,
			Max:    500,
			Extract: func(m *model.DeveloperMetrics) float64 {
				return float64(m.TotalStars)
			},
		},
		{
			Name:   "fork_impact",
			Weight: 0.05,
			Max:    100,
			Extract: func(m *model.DeveloperMetrics) float64 {
				return float64(m.TotalForks)
			},
		},
		{
			Name:   "avg_pr_size",
			Weight: 0.07,
			Max:    1.0,
			Extract: func(m *model.DeveloperMetrics) float64 {
				if m.AvgPRLinesChanged == 0 {
					return 0
				}
				ideal := 200.0
				diff := math.Abs(m.AvgPRLinesChanged - ideal)
				score := 1.0 - (diff / ideal)
				if score < 0 {
					score = 0
				}
				return score
			},
		},
		{
			Name:   "issue_response_time",
			Weight: 0.05,
			Max:    1.0,
			Extract: func(m *model.DeveloperMetrics) float64 {
				if m.AvgIssueResponseHours == 0 {
					return 0
				}
				maxHours := 168.0
				if m.AvgIssueResponseHours >= maxHours {
					return 0
				}
				return 1.0 - (m.AvgIssueResponseHours / maxHours)
			},
		},
		{
			Name:   "language_diversity",
			Weight: 0.05,
			Max:    10,
			Extract: func(m *model.DeveloperMetrics) float64 {
				return float64(len(m.Languages))
			},
		},
		{
			Name:   "commit_trend",
			Weight: 0.05,
			Max:    1.0,
			Extract: func(m *model.DeveloperMetrics) float64 {
				switch m.CommitTrend {
				case "increasing":
					return 1.0
				case "stable":
					return 0.5
				case "decreasing":
					return 0.2
				default:
					return 0.5
				}
			},
		},
	}
}
