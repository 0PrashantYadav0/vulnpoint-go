package handlers

import (
	"github.com/datmedevil17/go-vuln/internal/middleware"
	"github.com/datmedevil17/go-vuln/internal/services"
	"github.com/datmedevil17/go-vuln/internal/utils"
	"github.com/gin-gonic/gin"
)

type GitHubHandler struct {
	githubService *services.GitHubService
	authService   *services.AuthService
}

func NewGitHubHandler(githubService *services.GitHubService, authService *services.AuthService) *GitHubHandler {
	return &GitHubHandler{
		githubService: githubService,
		authService:   authService,
	}
}

// ListRepositories fetches user's GitHub repositories
func (h *GitHubHandler) ListRepositories(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Get user to retrieve access token
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	repositories, err := h.githubService.ListRepositories(c.Request.Context(), user.AccessToken, userID)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to fetch repositories: "+err.Error())
		return
	}

	utils.SuccessResponse(c, repositories)
}

// GetRepositoryFiles fetches files in a repository
func (h *GitHubHandler) GetRepositoryFiles(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	owner := c.Param("owner")
	repo := c.Param("repo")
	path := c.DefaultQuery("path", "")

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	files, err := h.githubService.GetRepositoryFiles(c.Request.Context(), user.AccessToken, owner, repo, path)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to fetch files: "+err.Error())
		return
	}

	utils.SuccessResponse(c, files)
}

// GetFileContent fetches content of a specific file
func (h *GitHubHandler) GetFileContent(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	owner := c.Param("owner")
	repo := c.Param("repo")
	path := c.Query("path")

	if path == "" {
		utils.BadRequestResponse(c, "File path required")
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	content, err := h.githubService.GetFileContent(c.Request.Context(), user.AccessToken, owner, repo, path)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to fetch file content: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"path":    path,
		"content": content,
	})
}
