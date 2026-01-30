package handlers

import (
	"net/http"

	"github.com/datmedevil17/go-vuln/internal/services"
	"github.com/gin-gonic/gin"
)

type AIWorkflowHandler struct {
	aiService *services.AIService
}

func NewAIWorkflowHandler(aiService *services.AIService) *AIWorkflowHandler {
	return &AIWorkflowHandler{aiService: aiService}
}

type GenerateWorkflowRequest struct {
	Prompt string `json:"prompt"`
}

func (h *AIWorkflowHandler) GenerateWorkflow(c *gin.Context) {
	var req GenerateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Prompt is required"})
		return
	}

	jsonConfig, err := h.aiService.GenerateWorkflowJSON(c.Request.Context(), req.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate workflow: " + err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", []byte(jsonConfig))
}
