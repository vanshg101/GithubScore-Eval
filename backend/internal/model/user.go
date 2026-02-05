package model

import "time"

type User struct {
	ID            string    `json:"id" firestore:"id"`
	GitHubID      int64     `json:"github_id" firestore:"github_id"`
	Username      string    `json:"username" firestore:"username"`
	DisplayName   string    `json:"display_name" firestore:"display_name"`
	AvatarURL     string    `json:"avatar_url" firestore:"avatar_url"`
	Email         string    `json:"email" firestore:"email"`
	AccessToken   string    `json:"-" firestore:"access_token"`
	CreatedAt     time.Time `json:"created_at" firestore:"created_at"`
	LastLoginAt   time.Time `json:"last_login_at" firestore:"last_login_at"`
}
