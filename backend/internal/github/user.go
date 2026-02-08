package github

import "fmt"

type UserProfile struct {
	Login       string `json:"login"`
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatar_url"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
}

func (c *Client) FetchUserProfile(username string) (*UserProfile, error) {
	url := fmt.Sprintf("%s/users/%s", baseURL, username)

	var profile UserProfile
	if err := c.getJSON(url, &profile); err != nil {
		return nil, fmt.Errorf("fetching user profile: %w", err)
	}

	return &profile, nil
}
