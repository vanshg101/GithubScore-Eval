package scoring

import (
	"time"

	"github.com/Madhur/GithubScoreEval/backend/internal/model"
)

// Engine computes developer scores using configurable weighted indicators.
type Engine struct {
	Indicators []Indicator
}

// NewEngine creates a scoring engine with default indicators.
func NewEngine() *Engine {
	return &Engine{
		Indicators: DefaultIndicators(),
	}
}

// NewEngineWithIndicators creates a scoring engine with custom indicators.
func NewEngineWithIndicators(indicators []Indicator) *Engine {
	return &Engine{
		Indicators: indicators,
	}
}

// Compute calculates a weighted score (0-100) from developer metrics.
func (e *Engine) Compute(username string, metrics *model.DeveloperMetrics) *model.Score {
	indicatorScores := make(map[string]model.IndicatorScore, len(e.Indicators))
	weightedTotal := 0.0

	for _, ind := range e.Indicators {
		raw := ind.Extract(metrics)
		normalized := normalize(raw, ind.Max)
		weighted := normalized * ind.Weight * 100

		indicatorScores[ind.Name] = model.IndicatorScore{
			Raw:        raw,
			Normalized: normalized,
			Weighted:   weighted,
		}

		weightedTotal += weighted
	}

	return &model.Score{
		Username:        username,
		WeightedScore:   clamp(weightedTotal, 0, 100),
		IndicatorScores: indicatorScores,
		ComputedAt:      time.Now(),
	}
}

// normalize applies min-max normalization clamped to [0, 1].
func normalize(raw, max float64) float64 {
	if max <= 0 {
		return 0
	}
	n := raw / max
	return clamp(n, 0, 1)
}

// clamp restricts a value to [min, max].
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
