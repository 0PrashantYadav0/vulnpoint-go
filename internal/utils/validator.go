package utils

import (
	"net/url"
	"regexp"
	"strings"
)

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateURL validates URL format
func ValidateURL(urlString string) bool {
	if urlString == "" {
		return false
	}

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return false
	}

	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

// ValidateUsername validates username format
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 32 {
		return false
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return usernameRegex.MatchString(username)
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Trim whitespace
	input = strings.TrimSpace(input)

	return input
}

// ValidatePort validates port number
func ValidatePort(port string) bool {
	portRegex := regexp.MustCompile(`^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`)
	return portRegex.MatchString(port)
}

// ValidateIPAddress validates IPv4 address
func ValidateIPAddress(ip string) bool {
	ipRegex := regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`)
	return ipRegex.MatchString(ip)
}

// ValidateLanguage validates programming language identifier
func ValidateLanguage(lang string) bool {
	validLanguages := map[string]bool{
		"go":         true,
		"golang":     true,
		"javascript": true,
		"js":         true,
		"typescript": true,
		"ts":         true,
		"python":     true,
		"py":         true,
		"java":       true,
		"c":          true,
		"cpp":        true,
		"c++":        true,
		"csharp":     true,
		"cs":         true,
		"ruby":       true,
		"rb":         true,
		"php":        true,
		"rust":       true,
		"rs":         true,
		"kotlin":     true,
		"swift":      true,
		"shell":      true,
		"bash":       true,
		"sh":         true,
	}

	return validLanguages[strings.ToLower(lang)]
}
