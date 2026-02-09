package handlers

import (
	"strings"

	"github.com/datmedevil17/go-vuln/internal/middleware"
	"github.com/datmedevil17/go-vuln/internal/models"
	"github.com/datmedevil17/go-vuln/internal/services"
	"github.com/datmedevil17/go-vuln/internal/utils"
	"github.com/gin-gonic/gin"
)

// RepositoryResponse is the API response format for repositories
type RepositoryResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	GitHubID    int64  `json:"github_id"`
	FullName    string `json:"full_name"`
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
	IsPrivate   bool   `json:"is_private"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func toRepositoryResponse(r models.Repository) RepositoryResponse {
	owner := ""
	if idx := strings.Index(r.FullName, "/"); idx >= 0 {
		owner = r.FullName[:idx]
	}
	return RepositoryResponse{
		ID:          r.ID.String(),
		UserID:      r.UserID.String(),
		GitHubID:    r.GitHubID,
		FullName:    r.FullName,
		Owner:       owner,
		Name:        r.Name,
		Description: r.Description,
		HTMLURL:     r.HTMLURL,
		Language:    r.Language,
		IsPrivate:   r.IsPrivate,
		CreatedAt:   r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

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

	responses := make([]RepositoryResponse, len(repositories))
	for i, r := range repositories {
		responses[i] = toRepositoryResponse(r)
	}
	utils.SuccessResponse(c, responses)
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
