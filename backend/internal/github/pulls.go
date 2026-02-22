package github

import (
	"encoding/json"
	"fmt"
	"time"
)

type PullRequest struct {
	Number    int        `json:"number"`
	Title     string     `json:"title"`
	State     string     `json:"state"`
	CreatedAt time.Time  `json:"created_at"`
	MergedAt  *time.Time `json:"merged_at"`
	Additions int        `json:"additions"`
	Deletions int        `json:"deletions"`
	User      struct {
		Login string `json:"login"`
	} `json:"user"`
}

func (c *Client) FetchPullRequests(owner, repo string) ([]PullRequest, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls?state=all&per_page=%d", baseURL, owner, repo, maxPerPage)

	rawItems, err := c.getPaginated(url)
	if err != nil {
		return nil, fmt.Errorf("fetching PRs for %s/%s: %w", owner, repo, err)
	}

	var prs []PullRequest
	for _, raw := range rawItems {
		var pr PullRequest
		if err := json.Unmarshal(raw, &pr); err != nil {
			continue
		}
		prs = append(prs, pr)
	}

	return prs, nil
}
