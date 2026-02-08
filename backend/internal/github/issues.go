package github

import (
	"encoding/json"
	"fmt"
	"time"
)

type Issue struct {
	Number           int        `json:"number"`
	Title            string     `json:"title"`
	State            string     `json:"state"`
	CreatedAt        time.Time  `json:"created_at"`
	ClosedAt         *time.Time `json:"closed_at"`
	PullRequestLinks *struct{}  `json:"pull_request"`
	User             struct {
		Login string `json:"login"`
	} `json:"user"`
}

func (c *Client) FetchIssues(owner, repo string) ([]Issue, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues?state=all&per_page=%d", baseURL, owner, repo, maxPerPage)

	rawItems, err := c.getPaginated(url)
	if err != nil {
		return nil, fmt.Errorf("fetching issues for %s/%s: %w", owner, repo, err)
	}

	var issues []Issue
	for _, raw := range rawItems {
		var issue Issue
		if err := json.Unmarshal(raw, &issue); err != nil {
			continue
		}
		if issue.PullRequestLinks == nil {
			issues = append(issues, issue)
		}
	}

	return issues, nil
}
