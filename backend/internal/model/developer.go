package model

import "time"

type DeveloperProfile struct {
	Name        string `json:"name" firestore:"name"`
	Bio         string `json:"bio" firestore:"bio"`
	PublicRepos int    `json:"public_repos" firestore:"public_repos"`
	Followers   int    `json:"followers" firestore:"followers"`
	AvatarURL   string `json:"avatar_url" firestore:"avatar_url"`
}

type DeveloperMetrics struct {
	TotalCommits          int      `json:"total_commits" firestore:"total_commits"`
	TotalPRs              int      `json:"total_prs" firestore:"total_prs"`
	MergedPRs             int      `json:"merged_prs" firestore:"merged_prs"`
	TotalIssuesOpened     int      `json:"total_issues_opened" firestore:"total_issues_opened"`
	TotalIssuesClosed     int      `json:"total_issues_closed" firestore:"total_issues_closed"`
	ReviewComments        int      `json:"review_comments" firestore:"review_comments"`
	ActiveWeeks           int      `json:"active_weeks" firestore:"active_weeks"`
	ReposContributed      int      `json:"repos_contributed" firestore:"repos_contributed"`
	TotalStars            int      `json:"total_stars" firestore:"total_stars"`
	TotalForks            int      `json:"total_forks" firestore:"total_forks"`
	AvgPRLinesChanged     float64  `json:"avg_pr_lines_changed" firestore:"avg_pr_lines_changed"`
	AvgIssueResponseHours float64  `json:"avg_issue_response_hours" firestore:"avg_issue_response_hours"`
	CommitTrend           string   `json:"commit_trend" firestore:"commit_trend"`
	Languages             []string `json:"languages" firestore:"languages"`
}

type Developer struct {
	Username  string           `json:"username" firestore:"username"`
	Profile   DeveloperProfile `json:"profile" firestore:"profile"`
	Metrics   DeveloperMetrics `json:"metrics" firestore:"metrics"`
	FetchedAt time.Time        `json:"fetched_at" firestore:"fetched_at"`
}
