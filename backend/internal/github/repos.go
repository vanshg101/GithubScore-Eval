package github

import (
	"encoding/json"
	"fmt"
)

type Repository struct {
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	Fork            bool   `json:"fork"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
	Language        string `json:"language"`
	DefaultBranch   string `json:"default_branch"`
	Owner           struct {
		Login string `json:"login"`
	} `json:"owner"`
}

func (c *Client) FetchRepos(username string) ([]Repository, error) {
	url := fmt.Sprintf("%s/users/%s/repos?per_page=%d&type=owner&sort=updated", baseURL, username, maxPerPage)

	rawItems, err := c.getPaginated(url)
	if err != nil {
		return nil, fmt.Errorf("fetching repos: %w", err)
	}

	var repos []Repository
	for _, raw := range rawItems {
		var repo Repository
		if err := json.Unmarshal(raw, &repo); err != nil {
			continue
		}
		if !repo.Fork {
			repos = append(repos, repo)
		}
	}

	return repos, nil
}
