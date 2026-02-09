package routes

import (
	"github.com/datmedevil17/go-vuln/internal/handlers"
	"github.com/datmedevil17/go-vuln/internal/utils"
	"github.com/gin-gonic/gin"
)

// RouterConfig holds all handler and utility dependencies
type RouterConfig struct {
	AuthHandler       *handlers.AuthHandler
	WorkflowHandler   *handlers.WorkflowHandler
	GitHubHandler     *handlers.GitHubHandler
	ScannerHandler    *handlers.ScannerHandler
	CodeHandler       *handlers.CodeHandler
	ChatbotHandler    *handlers.ChatbotHandler
	AIWorkflowHandler *handlers.AIWorkflowHandler
	JWTUtil           *utils.JWTManager
}

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, cfg *RouterConfig) {
	// Health check (no auth)
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "vulnpilot-go",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		RegisterAuthRoutes(api, cfg.AuthHandler)

		// Protected API routes
		RegisterAPIRoutes(api, &APIRoutesConfig{
			AuthHandler:       cfg.AuthHandler,
			WorkflowHandler:   cfg.WorkflowHandler,
			GitHubHandler:     cfg.GitHubHandler,
			ScannerHandler:    cfg.ScannerHandler,
			CodeHandler:       cfg.CodeHandler,
			ChatbotHandler:    cfg.ChatbotHandler,
			AIWorkflowHandler: cfg.AIWorkflowHandler,
			JWTUtil:           cfg.JWTUtil,
		})
	}
}
