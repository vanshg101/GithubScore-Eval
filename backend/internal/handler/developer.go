package handler

import (
	"log"
	"net/http"

	"github.com/Madhur/GithubScoreEval/backend/internal/repository"
	"github.com/Madhur/GithubScoreEval/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type DeveloperHandler struct {
	devService *service.DeveloperService
	userRepo   repository.UserRepository
}

func NewDeveloperHandler(devService *service.DeveloperService, userRepo repository.UserRepository) *DeveloperHandler {
	return &DeveloperHandler{
		devService: devService,
		userRepo:   userRepo,
	}
}

func (h *DeveloperHandler) FetchDeveloper(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	userID, _ := c.Get("user_id")
	var accessToken string
	if uid, ok := userID.(string); ok {
		user, err := h.userRepo.GetByID(c.Request.Context(), uid)
		if err == nil && user.AccessToken != "" {
			accessToken = user.AccessToken
		} else if err != nil {
			log.Printf("could not retrieve user access token: %v", err)
		}
	}

	developer, err := h.devService.FetchAndStore(c.Request.Context(), username, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, developer)
}

func (h *DeveloperHandler) GetDeveloper(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	developer, err := h.devService.GetByUsername(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "developer not found"})
		return
	}

	c.JSON(http.StatusOK, developer)
}

func (h *DeveloperHandler) ListDevelopers(c *gin.Context) {
	developers, err := h.devService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, developers)
}
