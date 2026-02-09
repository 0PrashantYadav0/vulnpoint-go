package routes

import (
	"github.com/datmedevil17/go-vuln/internal/handlers"
	"github.com/datmedevil17/go-vuln/internal/middleware"
	"github.com/datmedevil17/go-vuln/internal/utils"
	"github.com/gin-gonic/gin"
)

// RegisterAPIRoutes registers protected API routes
func RegisterAPIRoutes(rg *gin.RouterGroup, cfg *APIRoutesConfig) {
	protected := rg.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTUtil))
	{
		// User
		protected.GET("/user", cfg.AuthHandler.GetCurrentUser)

		// AI Workflow Generation
		protected.POST("/workflow/ai-generate", cfg.AIWorkflowHandler.GenerateWorkflow)

		// Workflows
		workflows := protected.Group("/workflows")
		{
			workflows.POST("", cfg.WorkflowHandler.CreateWorkflow)
			workflows.GET("", cfg.WorkflowHandler.ListWorkflows)
			workflows.GET("/reports", cfg.WorkflowHandler.ListWorkflowExecutions)
			workflows.GET("/:id", cfg.WorkflowHandler.GetWorkflow)
			workflows.PUT("/:id", cfg.WorkflowHandler.UpdateWorkflow)
			workflows.DELETE("/:id", cfg.WorkflowHandler.DeleteWorkflow)
			workflows.POST("/:id/execute", cfg.WorkflowHandler.ExecuteWorkflow)
		}

		// GitHub
		github := protected.Group("/github")
		{
			github.GET("/repositories", cfg.GitHubHandler.ListRepositories)
			github.GET("/repositories/:owner/:repo/files", cfg.GitHubHandler.GetRepositoryFiles)
			github.GET("/repositories/:owner/:repo/content", cfg.GitHubHandler.GetFileContent)
		}

		// Scanner
		scan := protected.Group("/scan")
		{
			scan.POST("/nmap", cfg.ScannerHandler.NmapScan)
			scan.POST("/nikto", cfg.ScannerHandler.NiktoScan)
			scan.POST("/gobuster", cfg.ScannerHandler.GobusterScan)
			scan.GET("/results", cfg.ScannerHandler.ListScanResults)
			scan.GET("/results/:id", cfg.ScannerHandler.GetScanResult)
		}

		// Code analysis
		code := protected.Group("/code")
		{
			code.POST("/analyze", cfg.CodeHandler.AnalyzeCode)
			code.POST("/quick-scan", cfg.CodeHandler.QuickScan)
			code.POST("/compare", cfg.CodeHandler.CompareCode)
		}

		// Chatbot
		chatbot := protected.Group("/chatbot")
		{
			chatbot.POST("/chat", cfg.ChatbotHandler.Chat)
			chatbot.POST("/explain", cfg.ChatbotHandler.ExplainVulnerability)
			chatbot.POST("/remediate", cfg.ChatbotHandler.SuggestRemediation)
			chatbot.POST("/ask", cfg.ChatbotHandler.AskSecurityQuestion)
		}
	}
}

// APIRoutesConfig holds handlers for API routes
type APIRoutesConfig struct {
	AuthHandler       *handlers.AuthHandler
	WorkflowHandler   *handlers.WorkflowHandler
	GitHubHandler     *handlers.GitHubHandler
	ScannerHandler    *handlers.ScannerHandler
	CodeHandler       *handlers.CodeHandler
	ChatbotHandler    *handlers.ChatbotHandler
	AIWorkflowHandler *handlers.AIWorkflowHandler
	JWTUtil           *utils.JWTManager
}
