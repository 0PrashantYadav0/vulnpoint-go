package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/datmedevil17/go-vuln/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GitHubService struct {
	db *gorm.DB
}

type GitHubRepo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
	Private     bool   `json:"private"`
}

type GitHubFile struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

func NewGitHubService(db *gorm.DB) *GitHubService {
	return &GitHubService{db: db}
}

// ListRepositories fetches repositories from GitHub API
func (s *GitHubService) ListRepositories(ctx context.Context, accessToken string, userID uuid.UUID) ([]models.Repository, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/repos?per_page=100", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s", string(body))
	}

	var githubRepos []GitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&githubRepos); err != nil {
		return nil, err
	}

	// Convert and store in database
	var repositories []models.Repository
	for _, gr := range githubRepos {
		repo := models.Repository{
			UserID:      userID,
			GitHubID:    gr.ID,
			FullName:    gr.FullName,
			Name:        gr.Name,
			Description: gr.Description,
			HTMLURL:     gr.HTMLURL,
			Language:    gr.Language,
			IsPrivate:   gr.Private,
		}

		// Upsert repository
		var existingRepo models.Repository
		result := s.db.Where("git_hub_id = ?", gr.ID).First(&existingRepo)
		if result.Error == gorm.ErrRecordNotFound {
			s.db.Create(&repo)
		} else {
			s.db.Model(&existingRepo).Updates(repo)
			repo = existingRepo
		}

		repositories = append(repositories, repo)
	}

	return repositories, nil
}

// GetRepositoryFiles fetches file tree from GitHub
func (s *GitHubService) GetRepositoryFiles(ctx context.Context, accessToken, owner, repo, path string) ([]GitHubFile, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, path)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch files: %s", resp.Status)
	}

	var files []GitHubFile
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}

	return files, nil
}

// GetFileContent fetches content of a specific file
func (s *GitHubService) GetFileContent(ctx context.Context, accessToken, owner, repo, path string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, path)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3.raw")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
