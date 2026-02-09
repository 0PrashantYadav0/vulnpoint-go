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

// IsConfigured returns true if GitHub OAuth is properly configured
func (s *AuthService) IsConfigured() bool {
	return s.oauth.ClientID != "" && s.oauth.ClientSecret != ""
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

	// Fetch email if not returned (GitHub often hides email)
	if githubUser.Email == "" {
		githubUser.Email, _ = s.getGitHubEmail(ctx, token.AccessToken)
	}

	// Find or create user (FirstOrCreate avoids "record not found" log for new signups)
	githubID := fmt.Sprintf("%d", githubUser.ID)
	var user models.User
	err = s.db.Where("github_id = ?", githubID).
		Assign(models.User{
			AccessToken: token.AccessToken,
			AvatarURL:   githubUser.AvatarURL,
		}).
		Attrs(models.User{
			GitHubID:  githubID,
			Username:  githubUser.Login,
			Email:     githubUser.Email,
			AvatarURL: githubUser.AvatarURL,
		}).
		FirstOrCreate(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// Update email for existing users who may not have had it set
	if user.Email == "" && githubUser.Email != "" {
		user.Email = githubUser.Email
		s.db.Model(&user).Update("email", user.Email)
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

// getGitHubEmail fetches primary email from GitHub /user/emails
func (s *AuthService) getGitHubEmail(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	for _, e := range emails {
		if e.Verified {
			return e.Email, nil
		}
	}
	if len(emails) > 0 {
		return emails[0].Email, nil
	}
	return "", nil
}
