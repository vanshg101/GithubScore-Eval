// Package middleware provides Gin middleware for authentication, CORS, and logging.
package middleware

import (
	"net/http"
	"strings"

	"github.com/Madhur/GithubScoreEval/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

// AuthRequired returns a Gin middleware that validates JWT tokens.
// It extracts the token from the Authorization header (Bearer scheme)
// or the "token" cookie, validates it, and injects user_id and username
// into the Gin context for downstream handlers.
func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
			return
		}

		claims, err := auth.ValidateToken(tokenStr, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// extractToken attempts to read a JWT from two sources in order:
//  1. Authorization: Bearer <token> header (API clients)
//  2. "token" cookie (browser sessions after OAuth callback)
func extractToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	if header != "" {
		parts := strings.SplitN(header, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	token, err := c.Cookie("token")
	if err == nil && token != "" {
		return token
	}

	return ""
}
