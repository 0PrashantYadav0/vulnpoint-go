package handlers

import (
	"github.com/datmedevil17/go-vuln/internal/middleware"
	"github.com/datmedevil17/go-vuln/internal/services"
	"github.com/datmedevil17/go-vuln/internal/utils"
	"github.com/gin-gonic/gin"
)

type CodeHandler struct {
	aiService        *services.AIService
	embeddingService *services.EmbeddingService
}

type AnalyzeCodeRequest struct {
	Code     string `json:"code" binding:"required"`
	Language string `json:"language" binding:"required"`
	Filename string `json:"filename,omitempty"`
}

type AnalyzeCodeResponse struct {
	Analysis           string   `json:"analysis"`
	Vulnerabilities    []string `json:"vulnerabilities"`
	SecurityScore      int      `json:"security_score"`
	Recommendations    string   `json:"recommendations"`
	VulnerabilityCount int      `json:"vulnerability_count"`
}

func NewCodeHandler(aiService *services.AIService, embeddingService *services.EmbeddingService) *CodeHandler {
	return &CodeHandler{
		aiService:        aiService,
		embeddingService: embeddingService,
	}
}

// AnalyzeCode performs AI-powered code analysis
func (h *CodeHandler) AnalyzeCode(c *gin.Context) {
	_, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req AnalyzeCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error())
		return
	}

	// Generate embedding and find patterns
	_, err := h.embeddingService.GenerateCodeEmbedding(req.Code, req.Language)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to generate code embedding")
		return
	}

	// Find vulnerability patterns
	vulnerabilities := h.embeddingService.FindVulnerabilityPatterns(req.Code, req.Language)

	// Get AI analysis
	analysis, err := h.aiService.AnalyzeCode(c.Request.Context(), req.Code, req.Language)
	if err != nil {
		// If AI fails, still return pattern-based results
		analysis = "AI analysis unavailable. Showing pattern-based analysis only."
	}

	// Calculate security score (0-100)
	securityScore := calculateSecurityScore(vulnerabilities)

	// Generate recommendations if vulnerabilities found
	recommendations := ""
	if len(vulnerabilities) > 0 {
		vulnDescription := "Found: " + joinStrings(vulnerabilities, ", ")
		recs, err := h.aiService.GenerateSecurityRecommendations(c.Request.Context(), vulnDescription)
		if err == nil {
			recommendations = recs
		}
	}

	response := AnalyzeCodeResponse{
		Analysis:           analysis,
		Vulnerabilities:    vulnerabilities,
		SecurityScore:      securityScore,
		Recommendations:    recommendations,
		VulnerabilityCount: len(vulnerabilities),
	}

	utils.SuccessResponse(c, response)
}

// QuickScan performs a quick vulnerability scan without AI
func (h *CodeHandler) QuickScan(c *gin.Context) {
	_, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req AnalyzeCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error())
		return
	}

	// Quick pattern-based scan
	vulnerabilities := h.embeddingService.FindVulnerabilityPatterns(req.Code, req.Language)
	securityScore := calculateSecurityScore(vulnerabilities)

	response := gin.H{
		"vulnerabilities":     vulnerabilities,
		"vulnerability_count": len(vulnerabilities),
		"security_score":      securityScore,
		"scan_type":           "quick",
	}

	utils.SuccessResponse(c, response)
}

// CompareCode compares two code snippets for similarity
func (h *CodeHandler) CompareCode(c *gin.Context) {
	_, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req struct {
		Code1     string `json:"code1" binding:"required"`
		Code2     string `json:"code2" binding:"required"`
		Language1 string `json:"language1" binding:"required"`
		Language2 string `json:"language2" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request: "+err.Error())
		return
	}

	// Generate embeddings
	embedding1, err := h.embeddingService.GenerateCodeEmbedding(req.Code1, req.Language1)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to analyze first code snippet")
		return
	}

	embedding2, err := h.embeddingService.GenerateCodeEmbedding(req.Code2, req.Language2)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to analyze second code snippet")
		return
	}

	// Calculate similarity
	similarity := h.embeddingService.CalculateSimilarity(embedding1, embedding2)

	response := gin.H{
		"similarity":         similarity,
		"similarity_percent": similarity * 100,
		"is_duplicate":       similarity > 0.8,
		"common_keywords":    findCommonKeywords(embedding1.Keywords, embedding2.Keywords),
	}

	utils.SuccessResponse(c, response)
}

// Helper functions

func calculateSecurityScore(vulnerabilities []string) int {
	// Start with perfect score
	score := 100

	// Deduct points for each vulnerability
	// Critical patterns deduct more points
	for _, vuln := range vulnerabilities {
		switch {
		case contains(vuln, "Injection"):
			score -= 25
		case contains(vuln, "Credentials"):
			score -= 20
		case contains(vuln, "XSS"):
			score -= 15
		case contains(vuln, "Traversal"):
			score -= 15
		default:
			score -= 10
		}
	}

	if score < 0 {
		score = 0
	}

	return score
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

func findCommonKeywords(keywords1, keywords2 []string) []string {
	keywordSet := make(map[string]bool)
	for _, k := range keywords1 {
		keywordSet[k] = true
	}

	common := []string{}
	for _, k := range keywords2 {
		if keywordSet[k] {
			common = append(common, k)
		}
	}

	return common
}
