package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Madhur/GithubScoreEval/backend/internal/auth"
	"github.com/Madhur/GithubScoreEval/backend/internal/config"
	"github.com/Madhur/GithubScoreEval/backend/internal/model"
	"github.com/Madhur/GithubScoreEval/backend/internal/store"
	"github.com/gin-gonic/gin"
)

const jwtExpiry = 7 * 24 * time.Hour

type AuthHandler struct {
	cfg       *config.Config
	userStore *store.UserStore
}

func NewAuthHandler(cfg *config.Config, userStore *store.UserStore) *AuthHandler {
	return &AuthHandler{
		cfg:       cfg,
		userStore: userStore,
	}
}

func (h *AuthHandler) GitHubLogin(c *gin.Context) {
	state := generateState()
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	url := auth.GetOAuthLoginURL(h.cfg.GithubClientID, h.cfg.GithubRedirectURL, state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	savedState, err := c.Cookie("oauth_state")
	if err != nil || state != savedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid oauth state"})
		return
	}

	tokenResp, err := auth.ExchangeCodeForToken(h.cfg.GithubClientID, h.cfg.GithubClientSecret, code)
	if err != nil {
		log.Printf("token exchange error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange code"})
		return
	}

	ghUser, err := auth.FetchGitHubUser(tokenResp.AccessToken)
	if err != nil {
		log.Printf("fetch user error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	userID := fmt.Sprintf("%d", ghUser.ID)
	now := time.Now()

	existing, _ := h.userStore.GetByID(userID)
	if existing != nil {
		existing.AccessToken = tokenResp.AccessToken
		existing.LastLoginAt = now
		existing.DisplayName = ghUser.Name
		existing.AvatarURL = ghUser.AvatarURL
		existing.Email = ghUser.Email
		h.userStore.Save(existing)
	} else {
		user := &model.User{
			ID:          userID,
			GitHubID:    ghUser.ID,
			Username:    ghUser.Login,
			DisplayName: ghUser.Name,
			AvatarURL:   ghUser.AvatarURL,
			Email:       ghUser.Email,
			AccessToken: tokenResp.AccessToken,
			CreatedAt:   now,
			LastLoginAt: now,
		}
		h.userStore.Save(user)
	}

	jwtToken, err := auth.GenerateToken(userID, ghUser.Login, h.cfg.JWTSecret, jwtExpiry)
	if err != nil {
		log.Printf("jwt generation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	c.SetCookie("token", jwtToken, int(jwtExpiry.Seconds()), "/", "", false, true)

	c.Redirect(http.StatusTemporaryRedirect, h.cfg.FrontendURL+"?token="+jwtToken)
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	user, err := h.userStore.GetByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"display_name": user.DisplayName,
		"avatar_url":   user.AvatarURL,
		"email":        user.Email,
		"created_at":   user.CreatedAt,
		"last_login":   user.LastLoginAt,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
