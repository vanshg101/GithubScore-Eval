package model

import "time"

type RankEntry struct {
	Rank     int     `json:"rank" firestore:"rank"`
	Username string  `json:"username" firestore:"username"`
	Score    float64 `json:"score" firestore:"score"`
	MLScore  float64 `json:"ml_score" firestore:"ml_score"`
}

type Ranking struct {
	SnapshotDate    string      `json:"snapshot_date" firestore:"snapshot_date"`
	Rankings        []RankEntry `json:"rankings" firestore:"rankings"`
	TotalDevelopers int         `json:"total_developers" firestore:"total_developers"`
	CreatedAt       time.Time   `json:"created_at" firestore:"created_at"`
}
