// Package scoring implements the weighted scoring algorithm used to evaluate
// GitHub developer contributions across 12 configurable indicators.
package scoring

import (
	"math"

	"github.com/Madhur/GithubScoreEval/backend/internal/model"
)

// Indicator defines a single scoring dimension.
//   - Name:    unique identifier used as a map key in the final score.
//   - Weight:  fractional weight (all weights should sum to 1.0).
//   - Max:     upper bound for normalization (raw values above Max are clamped to 1.0).
//   - Extract: function that extracts the raw metric value from DeveloperMetrics.
type Indicator struct {
	Name    string
	Weight  float64
	Max     float64
	Extract func(m *model.DeveloperMetrics) float64
}

// DefaultIndicators returns the production set of 12 scoring indicators.
// The weights sum to 1.0 and are tuned to balance activity volume,
// code quality, collaboration, community impact, and consistency.
func DefaultIndicators() []Indicator {
	return []Indicator{
		// ── Activity Volume ──────────────────────────────────────────
		{
			Name:   "total_commits",
			Weight: 0.15, // highest weight — primary productivity signal
			Max:    500,  // 500 commits/year = 1.0 normalized
			Extract: func(m *model.DeveloperMetrics) float64 {
				return float64(m.TotalCommits)
			},
		},
		// ── Code Quality ────────────────────────────────────────────
		{
			Name:   "pr_merge_rate",
			Weight: 0.12, // ratio of merged PRs to total PRs
			Max:    1.0,  // already a ratio, 100% = 1.0
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
		// ── Consistency & Collaboration ─────────────────────────────
		{
			Name:   "contribution_consistency",
			Weight: 0.12, // active weeks out of 52 — rewards sustained effort
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
		// ── Community Impact ────────────────────────────────────────
		{
			Name:   "stars_earned",
			Weight: 0.08,
			Max:    500, // 500 stars = 1.0 normalized
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
			// Penalises both trivially small and excessively large PRs.
			// The ideal PR size is 200 lines changed; deviation from that
			// linearly reduces the score toward 0.
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
			// Inverse linear scale: 0 hours → 1.0, 168 hours (1 week) → 0.
			// Faster response times earn higher scores.
			Extract: func(m *model.DeveloperMetrics) float64 {
				if m.AvgIssueResponseHours == 0 {
					return 0
				}
				maxHours := 168.0 // one week in hours
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
		// ── Trend Signal ────────────────────────────────────────────
		{
			Name:   "commit_trend",
			Weight: 0.05,
			Max:    1.0,
			// Maps qualitative trend label to a numeric score.
			// "increasing" is strongly rewarded; "decreasing" still
			// earns partial credit (0.2) since the developer is still active.
			Extract: func(m *model.DeveloperMetrics) float64 {
				switch m.CommitTrend {
				case "increasing":
					return 1.0
				case "stable":
					return 0.5
				case "decreasing":
					return 0.2
				default:
					return 0.5 // treat unknown as "stable"
				}
			},
		},
	}
}
