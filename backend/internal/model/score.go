package model

import "time"

type IndicatorScore struct {
	Raw        float64 `json:"raw" firestore:"raw"`
	Normalized float64 `json:"normalized" firestore:"normalized"`
	Weighted   float64 `json:"weighted" firestore:"weighted"`
}

type Score struct {
	Username        string                    `json:"username" firestore:"username"`
	WeightedScore   float64                   `json:"weighted_score" firestore:"weighted_score"`
	MLImpactScore   float64                   `json:"ml_impact_score" firestore:"ml_impact_score"`
	IndicatorScores map[string]IndicatorScore `json:"indicator_scores" firestore:"indicator_scores"`
	Percentile      float64                   `json:"percentile" firestore:"percentile"`
	ComputedAt      time.Time                 `json:"computed_at" firestore:"computed_at"`
}
