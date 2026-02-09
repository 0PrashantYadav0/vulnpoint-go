package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	JWT       JWTConfig
	GitHub    GitHubConfig
	AI        AIConfig
	Email     EmailConfig
	Slack     SlackConfig
	RateLimit RateLimitConfig
	Logging   LoggingConfig
	Scanning  ScanningConfig
	Frontend  FrontendConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
	Host string
	Mode string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
	DSN      string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Address  string
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret            string
	Expiration        time.Duration
	RefreshExpiration time.Duration
}

// GitHubConfig holds GitHub OAuth configuration
type GitHubConfig struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
}

// AIConfig holds AI service configuration
type AIConfig struct {
	GeminiAPIKey string
	GroqAPIKey   string
}

// EmailConfig holds email service configuration
type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	User     string
	Password string
	From     string
	Enabled  bool
}

// SlackConfig holds Slack notification configuration
type SlackConfig struct {
	WebhookURL string
	Enabled    bool
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled  bool
	Requests int
	Window   time.Duration
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string // json or text
	File   string
}

// ScanningConfig holds security tool paths
type ScanningConfig struct {
	NmapPath     string
	NiktoPath    string
	GobusterPath string
	SQLMapPath   string
	WPScanPath   string
}

// FrontendConfig holds frontend-related configuration
type FrontendConfig struct {
	URL         string
	CORSOrigins []string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
			Mode: getEnv("SERVER_MODE", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "vulnpilot"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "vulnpilot_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			TimeZone: getEnv("DB_TIMEZONE", "UTC"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiration:        getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
			RefreshExpiration: getEnvAsDuration("JWT_REFRESH_EXPIRATION", 7*24*time.Hour),
		},
		GitHub: GitHubConfig{
			ClientID:     getEnv("GITHUB_CLIENT_ID", ""),
			ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
			CallbackURL:  getEnv("GITHUB_CALLBACK_URL", "http://localhost:8080/api/auth/github/callback"),
		},
		AI: AIConfig{
			GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
			GroqAPIKey:   getEnv("GROQ_API_KEY", ""),
		},
		Email: EmailConfig{
			SMTPHost: getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort: getEnvAsInt("SMTP_PORT", 587),
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@vulnpilot.com"),
			Enabled:  getEnvAsBool("EMAIL_ENABLED", false),
		},
		Slack: SlackConfig{
			WebhookURL: getEnv("SLACK_WEBHOOK_URL", ""),
			Enabled:    getEnvAsBool("SLACK_ENABLED", false),
		},
		RateLimit: RateLimitConfig{
			Enabled:  getEnvAsBool("RATE_LIMIT_ENABLED", true),
			Requests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
			Window:   getEnvAsDuration("RATE_LIMIT_WINDOW", 15*time.Minute),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
			File:   getEnv("LOG_FILE", "logs/vulnpilot.log"),
		},
		Scanning: ScanningConfig{
			NmapPath:     getEnv("NMAP_PATH", "/usr/bin/nmap"),
			NiktoPath:    getEnv("NIKTO_PATH", "/usr/bin/nikto"),
			GobusterPath: getEnv("GOBUSTER_PATH", "/usr/local/bin/gobuster"),
			SQLMapPath:   getEnv("SQLMAP_PATH", "/usr/bin/sqlmap"),
			WPScanPath:   getEnv("WPSCAN_PATH", "/usr/bin/wpscan"),
		},
		Frontend: FrontendConfig{
			URL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		},
	}

	// Build database DSN
	config.Database.DSN = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.SSLMode,
		config.Database.TimeZone,
	)

	// Build Redis address
	config.Redis.Address = fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port)

	// Parse CORS origins
	corsOrigins := getEnv("CORS_ORIGINS", "http://localhost:3000")
	config.Frontend.CORSOrigins = strings.Split(corsOrigins, ",")

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.JWT.Secret == "your-secret-key-change-in-production" && c.Server.Mode == "production" {
		return fmt.Errorf("JWT_SECRET must be changed from default value in production")
	}

	if len(c.JWT.Secret) < 32 && c.Server.Mode == "production" {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long in production")
	}

	if (c.GitHub.ClientID == "" || c.GitHub.ClientSecret == "") && c.Server.Mode != "development" {
		return fmt.Errorf("GitHub OAuth credentials are required")
	}

	if c.GitHub.ClientID == "" || c.GitHub.ClientSecret == "" {
		log.Println("WARNING: GitHub OAuth not configured - auth will not work. Set GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET")
	}

	if c.Database.Password == "" {
		log.Println("WARNING: Database password is empty")
	}

	return nil
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
