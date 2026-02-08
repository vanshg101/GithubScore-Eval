package github

import (
	"encoding/json"
	"fmt"
	"time"
)

type Commit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
}

func (c *Client) FetchCommits(owner, repo, username string) ([]Commit, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/commits?author=%s&per_page=%d", baseURL, owner, repo, username, maxPerPage)

	rawItems, err := c.getPaginated(url)
	if err != nil {
		return nil, fmt.Errorf("fetching commits for %s/%s: %w", owner, repo, err)
	}

	var commits []Commit
	for _, raw := range rawItems {
		var commit Commit
		if err := json.Unmarshal(raw, &commit); err != nil {
			continue
		}
		commits = append(commits, commit)
	}

	return commits, nil
}
