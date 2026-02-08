package github

import (
	"encoding/json"
	"fmt"
	"time"
)

type ReviewComment struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	User      struct {
		Login string `json:"login"`
	} `json:"user"`
}

func (c *Client) FetchReviewComments(owner, repo string) ([]ReviewComment, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/comments?per_page=%d", baseURL, owner, repo, maxPerPage)

	rawItems, err := c.getPaginated(url)
	if err != nil {
		return nil, fmt.Errorf("fetching review comments for %s/%s: %w", owner, repo, err)
	}

	var comments []ReviewComment
	for _, raw := range rawItems {
		var comment ReviewComment
		if err := json.Unmarshal(raw, &comment); err != nil {
			continue
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
