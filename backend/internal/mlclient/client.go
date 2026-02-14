package mlclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Madhur/GithubScoreEval/backend/internal/model"
)

// PredictRequest is the JSON body sent to the ML service.
type PredictRequest struct {
	TotalCommits          int     `json:"total_commits"`
	TotalPRs              int     `json:"total_prs"`
	MergedPRs             int     `json:"merged_prs"`
	TotalIssuesOpened     int     `json:"total_issues_opened"`
	TotalIssuesClosed     int     `json:"total_issues_closed"`
	ReviewComments        int     `json:"review_comments"`
	ActiveWeeks           int     `json:"active_weeks"`
	ReposContributed      int     `json:"repos_contributed"`
	TotalStars            int     `json:"total_stars"`
	TotalForks            int     `json:"total_forks"`
	AvgPRLinesChanged     float64 `json:"avg_pr_lines_changed"`
	AvgIssueResponseHours float64 `json:"avg_issue_response_hours"`
	CommitTrendScore      float64 `json:"commit_trend_score"`
	LanguageCount         int     `json:"language_count"`
}

// PredictResponse is the JSON response from the ML service.
type PredictResponse struct {
	ImpactScore float64 `json:"impact_score"`
}

// Client communicates with the FastAPI ML microservice.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new ML service client.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// MapMetrics converts Go DeveloperMetrics to a PredictRequest suitable for the ML API.
func MapMetrics(m *model.DeveloperMetrics) *PredictRequest {
	// Clamp active_weeks to 52 (ML model constraint)
	activeWeeks := m.ActiveWeeks
	if activeWeeks > 52 {
		activeWeeks = 52
	}

	return &PredictRequest{
		TotalCommits:          m.TotalCommits,
		TotalPRs:              m.TotalPRs,
		MergedPRs:             m.MergedPRs,
		TotalIssuesOpened:     m.TotalIssuesOpened,
		TotalIssuesClosed:     m.TotalIssuesClosed,
		ReviewComments:        m.ReviewComments,
		ActiveWeeks:           activeWeeks,
		ReposContributed:      m.ReposContributed,
		TotalStars:            m.TotalStars,
		TotalForks:            m.TotalForks,
		AvgPRLinesChanged:     m.AvgPRLinesChanged,
		AvgIssueResponseHours: m.AvgIssueResponseHours,
		CommitTrendScore:      commitTrendToScore(m.CommitTrend),
		LanguageCount:         len(m.Languages),
	}
}

// commitTrendToScore converts the string trend label to a numeric score for the ML model.
func commitTrendToScore(trend string) float64 {
	switch trend {
	case "increasing":
		return 1.0
	case "stable":
		return 0.5
	case "decreasing":
		return 0.2
	default:
		return 0.5
	}
}

// Predict sends developer metrics to the ML service and returns the impact score.
// Returns (score, error). On failure the caller can use fallback logic.
func (c *Client) Predict(ctx context.Context, req *PredictRequest) (float64, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("ml client: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/predict", bytes.NewReader(body))
	if err != nil {
		return 0, fmt.Errorf("ml client: create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("ml client: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("ml client: non-200 status %d: %s", resp.StatusCode, string(respBody))
	}

	var result PredictResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("ml client: decode response: %w", err)
	}

	return result.ImpactScore, nil
}

// IsHealthy checks if the ML service is reachable and the model is loaded.
func (c *Client) IsHealthy(ctx context.Context) bool {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return false
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Printf("ML service health check failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
