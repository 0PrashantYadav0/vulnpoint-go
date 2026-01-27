package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type EmbeddingService struct{}

type CodeEmbedding struct {
	Hash        string
	Language    string
	Content     string
	Keywords    []string
	Fingerprint string
}

func NewEmbeddingService() *EmbeddingService {
	return &EmbeddingService{}
}

// GenerateCodeEmbedding creates an embedding representation of code
func (s *EmbeddingService) GenerateCodeEmbedding(code string, language string) (*CodeEmbedding, error) {
	// Generate hash for deduplication
	hash := sha256.Sum256([]byte(code))
	hashStr := hex.EncodeToString(hash[:])

	// Extract keywords and patterns
	keywords := s.extractKeywords(code, language)

	// Generate fingerprint for similarity matching
	fingerprint := s.generateFingerprint(code, keywords)

	return &CodeEmbedding{
		Hash:        hashStr,
		Language:    language,
		Content:     code,
		Keywords:    keywords,
		Fingerprint: fingerprint,
	}, nil
}

// extractKeywords extracts important keywords based on language
func (s *EmbeddingService) extractKeywords(code string, language string) []string {
	keywords := make(map[string]bool)

	// Common security-related keywords
	securityKeywords := []string{
		"password", "token", "secret", "key", "api", "auth",
		"sql", "query", "execute", "eval", "exec",
		"file", "path", "upload", "download",
		"crypto", "hash", "encrypt", "decrypt",
		"session", "cookie", "header",
		"input", "output", "sanitize", "validate",
	}

	codeLower := strings.ToLower(code)
	for _, keyword := range securityKeywords {
		if strings.Contains(codeLower, keyword) {
			keywords[keyword] = true
		}
	}

	// Language-specific patterns
	switch language {
	case "go", "golang":
		goPatterns := []string{"import", "func", "struct", "interface", "defer", "goroutine", "channel"}
		for _, pattern := range goPatterns {
			if strings.Contains(codeLower, pattern) {
				keywords[pattern] = true
			}
		}
	case "javascript", "js", "typescript", "ts":
		jsPatterns := []string{"require", "import", "export", "async", "await", "promise", "fetch"}
		for _, pattern := range jsPatterns {
			if strings.Contains(codeLower, pattern) {
				keywords[pattern] = true
			}
		}
	case "python", "py":
		pyPatterns := []string{"import", "from", "def", "class", "async", "await"}
		for _, pattern := range pyPatterns {
			if strings.Contains(codeLower, pattern) {
				keywords[pattern] = true
			}
		}
	}

	// Convert map to slice
	result := make([]string, 0, len(keywords))
	for keyword := range keywords {
		result = append(result, keyword)
	}

	return result
}

// generateFingerprint creates a fingerprint for similarity matching
func (s *EmbeddingService) generateFingerprint(code string, keywords []string) string {
	// Simple fingerprint based on code structure and keywords
	fingerprint := fmt.Sprintf("len:%d|keywords:%s", len(code), strings.Join(keywords, ","))
	hash := sha256.Sum256([]byte(fingerprint))
	return hex.EncodeToString(hash[:16]) // Use first 16 bytes
}

// CalculateSimilarity calculates similarity between two code embeddings
func (s *EmbeddingService) CalculateSimilarity(e1, e2 *CodeEmbedding) float64 {
	// If exact match
	if e1.Hash == e2.Hash {
		return 1.0
	}

	// Calculate keyword overlap
	keywordSet1 := make(map[string]bool)
	for _, k := range e1.Keywords {
		keywordSet1[k] = true
	}

	commonKeywords := 0
	for _, k := range e2.Keywords {
		if keywordSet1[k] {
			commonKeywords++
		}
	}

	totalKeywords := len(e1.Keywords) + len(e2.Keywords) - commonKeywords
	if totalKeywords == 0 {
		return 0.0
	}

	return float64(commonKeywords) / float64(totalKeywords)
}

// FindVulnerabilityPatterns identifies common vulnerability patterns
func (s *EmbeddingService) FindVulnerabilityPatterns(code string, language string) []string {
	patterns := []string{}
	codeLower := strings.ToLower(code)

	// SQL Injection patterns
	if strings.Contains(codeLower, "sql") || strings.Contains(codeLower, "query") {
		if strings.Contains(codeLower, "execute") || strings.Contains(codeLower, "+") {
			patterns = append(patterns, "Potential SQL Injection")
		}
	}

	// Command Injection patterns
	if strings.Contains(codeLower, "exec") || strings.Contains(codeLower, "system") {
		patterns = append(patterns, "Potential Command Injection")
	}

	// XSS patterns
	if strings.Contains(codeLower, "innerhtml") || strings.Contains(codeLower, "document.write") {
		patterns = append(patterns, "Potential XSS Vulnerability")
	}

	// Path Traversal
	if strings.Contains(codeLower, "../") || strings.Contains(codeLower, "..\\") {
		patterns = append(patterns, "Potential Path Traversal")
	}

	// Hardcoded credentials
	if strings.Contains(codeLower, "password") && (strings.Contains(codeLower, "=") || strings.Contains(codeLower, ":")) {
		if !strings.Contains(codeLower, "input") && !strings.Contains(codeLower, "prompt") {
			patterns = append(patterns, "Potential Hardcoded Credentials")
		}
	}

	// Insecure deserialization
	if strings.Contains(codeLower, "deserialize") || strings.Contains(codeLower, "unmarshal") {
		patterns = append(patterns, "Potential Insecure Deserialization")
	}

	return patterns
}
