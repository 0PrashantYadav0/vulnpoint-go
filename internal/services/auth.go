package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/datmedevil17/go-vuln/internal/config"
	"github.com/datmedevil17/go-vuln/internal/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	config *config.Config
	oauth  *oauth2.Config
}

type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	oauth := &oauth2.Config{
		ClientID:     cfg.GitHub.ClientID,
		ClientSecret: cfg.GitHub.ClientSecret,
		RedirectURL:  cfg.GitHub.CallbackURL,
		Scopes:       []string{"user:email", "repo"},
		Endpoint:     github.Endpoint,
	}

	return &AuthService{
		db:     db,
		config: cfg,
		oauth:  oauth,
	}
}

// GetAuthURL returns GitHub OAuth URL
func (s *AuthService) GetAuthURL(state string) string {
	return s.oauth.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// HandleCallback processes GitHub OAuth callback
func (s *AuthService) HandleCallback(ctx context.Context, code string) (*models.User, error) {
	// Exchange code for token
	token, err := s.oauth.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from GitHub
	githubUser, err := s.getGitHubUser(ctx, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub user: %w", err)
	}

	// Find or create user
	var user models.User
	result := s.db.Where("github_id = ?", fmt.Sprintf("%d", githubUser.ID)).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		// Create new user
		user = models.User{
			GitHubID:    fmt.Sprintf("%d", githubUser.ID),
			Username:    githubUser.Login,
			Email:       githubUser.Email,
			AvatarURL:   githubUser.AvatarURL,
			AccessToken: token.AccessToken,
		}

		if err := s.db.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else if result.Error != nil {
		return nil, result.Error
	} else {
		// Update existing user
		user.AccessToken = token.AccessToken
		user.AvatarURL = githubUser.AvatarURL
		if err := s.db.Save(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	return &user, nil
}

// GetUserByID retrieves user by ID
func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// getGitHubUser fetches user info from GitHub API
func (s *AuthService) getGitHubUser(ctx context.Context, accessToken string) (*GitHubUser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
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

	var githubUser GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, err
	}

	return &githubUser, nil
}
