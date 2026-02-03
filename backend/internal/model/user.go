package model

import "time"

type User struct {
	ID            string    `json:"id"`
	GitHubID      int64     `json:"github_id"`
	Username      string    `json:"username"`
	DisplayName   string    `json:"display_name"`
	AvatarURL     string    `json:"avatar_url"`
	Email         string    `json:"email"`
	AccessToken   string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	LastLoginAt   time.Time `json:"last_login_at"`
}
