package handlers

import (
	"fmt"

	"github.com/datmedevil17/go-vuln/internal/middleware"
	"github.com/datmedevil17/go-vuln/internal/services"
	"github.com/datmedevil17/go-vuln/internal/utils"
	"github.com/gin-gonic/gin"
)

type ChatbotHandler struct {
	aiService *services.AIService
}

type ChatRequest struct {
	Message             string              `json:"message" binding:"required"`
	ConversationHistory []map[string]string `json:"conversation_history,omitempty"`
}

type ChatResponse struct {
	Response       string `json:"response"`
	ConversationID string `json:"conversation_id,omitempty"`
}

func NewChatbotHandler(aiService *services.AIService) *ChatbotHandler {
	return &ChatbotHandler{
		aiService: aiService,
	}
}

// Chat handles chatbot interactions
func (h *ChatbotHandler) Chat(c *gin.Context) {
	_, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error())
		return
	}

	// Get AI response
	response, err := h.aiService.ChatResponse(c.Request.Context(), req.Message, req.ConversationHistory)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to generate response: "+err.Error())
		return
	}

	chatResponse := ChatResponse{
		Response: response,
	}

	utils.SuccessResponse(c, chatResponse)
}

// ExplainVulnerability explains a specific vulnerability
func (h *ChatbotHandler) ExplainVulnerability(c *gin.Context) {
	_, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req struct {
		VulnerabilityType string `json:"vulnerability_type" binding:"required"`
		Context           string `json:"context,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error())
		return
	}

	// Create specialized prompt for vulnerability explanation
	prompt := "Explain the following security vulnerability in detail: " + req.VulnerabilityType
	if req.Context != "" {
		prompt += "\n\nContext: " + req.Context
	}
	prompt += "\n\nPlease provide:\n1. What it is\n2. Why it's dangerous\n3. How attackers exploit it\n4. How to fix it\n5. Best practices to prevent it"

	response, err := h.aiService.ChatResponse(c.Request.Context(), prompt, nil)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to generate explanation: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"vulnerability_type": req.VulnerabilityType,
		"explanation":        response,
	})
}

// SuggestRemediation suggests remediation steps for vulnerabilities
func (h *ChatbotHandler) SuggestRemediation(c *gin.Context) {
	_, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req struct {
		VulnerabilityType string `json:"vulnerability_type" binding:"required"`
		CodeSnippet       string `json:"code_snippet,omitempty"`
		Language          string `json:"language,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error())
		return
	}

	// Create specialized prompt for remediation
	prompt := "Provide specific remediation steps for this security vulnerability: " + req.VulnerabilityType

	if req.CodeSnippet != "" && req.Language != "" {
		prompt += fmt.Sprintf("\n\nVulnerable Code (%s):\n%s", req.Language, req.CodeSnippet)
		prompt += "\n\nPlease provide:\n1. Fixed version of the code\n2. Explanation of the changes\n3. Additional security measures\n4. Testing recommendations"
	} else {
		prompt += "\n\nPlease provide:\n1. Step-by-step remediation guide\n2. Code examples if applicable\n3. Configuration changes needed\n4. Verification steps"
	}

	response, err := h.aiService.ChatResponse(c.Request.Context(), prompt, nil)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to generate remediation: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"vulnerability_type": req.VulnerabilityType,
		"remediation":        response,
	})
}

// AskSecurityQuestion handles general security questions
func (h *ChatbotHandler) AskSecurityQuestion(c *gin.Context) {
	_, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req struct {
		Question string `json:"question" binding:"required"`
		Category string `json:"category,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error())
		return
	}

	// Add context based on category
	prompt := req.Question
	if req.Category != "" {
		prompt = fmt.Sprintf("Category: %s\n\nQuestion: %s", req.Category, req.Question)
	}

	response, err := h.aiService.ChatResponse(c.Request.Context(), prompt, nil)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to generate answer: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"question": req.Question,
		"answer":   response,
		"category": req.Category,
	})
}
