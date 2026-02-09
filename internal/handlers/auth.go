package handlers

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/datmedevil17/go-vuln/internal/config"
	"github.com/datmedevil17/go-vuln/internal/services"
	"github.com/datmedevil17/go-vuln/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *services.AuthService
	jwtUtil     *utils.JWTUtil
	config      *config.Config
}

func NewAuthHandler(authService *services.AuthService, jwtUtil *utils.JWTUtil, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtUtil:     jwtUtil,
		config:      cfg,
	}
}

// GetAuthURL generates GitHub OAuth URL
func (h *AuthHandler) GetAuthURL(c *gin.Context) {
	if h.authService == nil {
		utils.InternalErrorResponse(c, "Authentication service not available")
		return
	}
	if !h.authService.IsConfigured() {
		utils.InternalErrorResponse(c, "GitHub OAuth is not configured. Set GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET.")
		return
	}

	state := generateRandomState()
	authURL := h.authService.GetAuthURL(state)

	utils.SuccessResponse(c, gin.H{
		"url":   authURL,
		"state": state,
	})
}

// HandleCallback processes GitHub OAuth callback
func (h *AuthHandler) HandleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		utils.BadRequestResponse(c, "Authorization code required")
		return
	}

	// Exchange code for user
	user, err := h.authService.HandleCallback(c.Request.Context(), code)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to authenticate: "+err.Error())
		return
	}

	// Generate JWT token
	token, err := h.jwtUtil.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to generate token")
		return
	}

	// Redirect to frontend with token
	frontendURL := h.config.Frontend.URL
	redirectURL := frontendURL + "/auth/callback?token=" + token

	c.Redirect(302, redirectURL)
}

// GetCurrentUser returns current authenticated user info
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		utils.UnauthorizedResponse(c, "Invalid user context")
		return
	}

	user, err := h.authService.GetUserByID(id)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, user)
}

// Logout (client-side token deletion)
func (h *AuthHandler) Logout(c *gin.Context) {
	utils.SuccessMessageResponse(c, "Logged out successfully", nil)
}

func generateRandomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
