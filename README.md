# VulnPilot - AI-Powered Vulnerability Scanner

![VulnPilot Logo](https://img.shields.io/badge/VulnPilot-Security%20Scanner-blue)
![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

VulnPilot is a comprehensive, AI-powered security vulnerability scanner with GitHub integration, workflow automation, and intelligent code analysis capabilities.

## 🚀 Features

- **GitHub OAuth Integration**: Seamless authentication and repository access
- **AI-Powered Analysis**: Leverages Google Gemini and Groq APIs for intelligent code analysis
-  **Multiple Scan Types**:
  - Nmap (Network port scanning)
  - Nikto (Web server vulnerability scanning)
  - Gobuster (Directory/file bruteforcing)
  - SAST (Static Application Security Testing)
- **Workflow Automation**: Create and schedule custom security workflows
- **AI Chatbot**: Get security guidance and vulnerability explanations
- **Code Analysis**: Deep code analysis with vulnerability pattern detection
- **Notifications**: Email and Slack notifications for scan results
- **Rate Limiting**: Redis-backed distributed rate limiting
- **Secure**: JWT authentication, bcrypt password hashing, AES-256 encryption

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Variables](#environment-variables)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Makefile Commands](#makefile-commands)
- [Architecture](#architecture)
- [Development](#development)

## 🔧 Prerequisites

- **Go**: 1.24 or higher
- **PostgreSQL**: 14 or higher
- **Redis**: 6 or higher
- **Docker & Docker Compose**: (optional, for containerized setup)
- **Security Tools**: nmap, nikto, gobuster (for scanning features)

## 🔐 Environment Variables

Create a `.env` file in the project root. Use `.env.example` as a template:

### Required Variables

```bash
# JWT Configuration (REQUIRED)
JWT_SECRET=your_jwt_secret_key_minimum_32_characters_long

# GitHub OAuth (REQUIRED)
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GITHUB_CALLBACK_URL=http://localhost:8080/api/auth/github/callback

# Database (REQUIRED)
DB_HOST=postgres
DB_PORT=5432
DB_USER=vulnpilot
DB_PASSWORD=your_secure_password_here
DB_NAME=vulnpilot_db
DB_SSLMODE=disable
```

### Optional Variables

```bash
# AI Services (at least one recommended)
GEMINI_API_KEY=your_gemini_api_key_here
GROQ_API_KEY=your_groq_api_key_here

# Email Notifications
EMAIL_ENABLED=true
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password
SMTP_FROM=noreply@vulnpilot.com

# Slack Notifications
SLACK_ENABLED=false
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Frontend
FRONTEND_URL=http://localhost:3000
CORS_ORIGINS=http://localhost:3000,http://localhost:5173

# Server
SERVER_PORT=8080
SERVER_HOST=localhost
SERVER_MODE=development
```

### Getting API Keys

1. **GitHub OAuth**:
   - Go to GitHub Settings > Developer settings > OAuth Apps
   - Create a new OAuth App
   - Set callback URL to `http://localhost:8080/api/auth/github/callback`

2. **Gemini API**:
   - Visit [Google AI Studio](https://makersuite.google.com/app/apikey)
   - Create an API key

3. **Groq API**:
   - Visit [Groq Console](https://console.groq.com/)
   - Create an API key

## 📦 Installation

### Option 1: Using Make (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd go-vuln

# Copy environment file
cp .env.example .env
# Edit .env with your configuration

# Complete setup (Docker + deps + migrations)
make setup

# Start development server
make dev
```

### Option 2: Manual Setup

```bash
# Start PostgreSQL and Redis
make docker-up

# Install dependencies
go mod download
go mod tidy

# Run database migrations
make migrate-up

# Build application
make build

# Run application
./bin/go-vuln
```

## 🎯 Usage

### Starting the Server

```bash
# Development mode with hot reload
make dev

# Production mode
make build
./bin/go-vuln
```

### Running Scans

```bash
# Using the API (after authentication)
curl -X POST http://localhost:8080/api/scan/nmap \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"target": "example.com", "ports": "1-1000"}'
```

## 🛣️ API Endpoints

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/auth/github` | Get GitHub OAuth URL |
| GET | `/api/auth/github/callback` | GitHub OAuth callback |
| GET | `/api/user` | Get current user info |
| POST | `/api/auth/logout` | Logout user |

### Scanning

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/scan/nmap` | Run Nmap scan |
| POST | `/api/scan/nikto` | Run Nikto scan |
| POST | `/api/scan/gobuster` | Run Gobuster scan |
| GET | `/api/scan/results` | List scan results |
| GET | `/api/scan/results/:id` | Get scan result |

### Code Analysis

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/code/analyze` | Analyze code with AI |
| POST | `/api/code/quick-scan` | Quick vulnerability scan |
| POST | `/api/code/compare` | Compare two code snippets |

### Chatbot

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/chatbot/chat` | Chat with AI security expert |
| POST | `/api/chatbot/explain` | Explain vulnerability |
| POST | `/api/chatbot/remediate` | Get remediation steps |
| POST | `/api/chatbot/ask` | Ask security question |

### Workflows

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/workflows` | Create workflow |
| GET | `/api/workflows` | List workflows |
| GET | `/api/workflows/:id` | Get workflow |
| PUT | `/api/workflows/:id` | Update workflow |
| DELETE | `/api/workflows/:id` | Delete workflow |

### GitHub

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/github/repositories` | List user repositories |
| GET | `/api/github/repositories/:owner/:repo/files` | Get repository files |
| GET | `/api/github/repositories/:owner/:repo/content` | Get file content |

## 🛠️ Makefile Commands

```bash
make help              # Show all available commands
make build             # Build the application
make run               # Run the application
make dev               # Run with hot reload (requires air)
make test              # Run tests
make test-coverage     # Run tests with coverage
make clean             # Clean build artifacts

# Database
make migrate-up        # Run migrations
make migrate-down      # Rollback migrations
make db-shell          # Connect to PostgreSQL

# Docker
make docker-up         # Start Docker containers
make docker-down       # Stop Docker containers
make docker-logs       # View container logs
make docker-clean      # Remove containers and volumes

# Code Quality
make lint              # Run linter
make fmt               # Format code
make vet               # Run go vet
make check             # Run all checks (fmt, vet, lint, test)

# Development
make deps              # Update dependencies
make setup             # Complete setup
make install-tools     # Install dev tools (air, golangci-lint)
```

## 🏗️ Architecture

```
go-vuln/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   ├── postgres.go          # PostgreSQL connection
│   │   └── redis.go             # Redis connection
│   ├── handlers/
│   │   ├── auth.go              # Authentication handlers
│   │   ├── chatbot.go           # AI chatbot handlers
│   │   ├── code.go              # Code analysis handlers
│   │   ├── github.go            # GitHub integration handlers
│   │   ├── scan.go              # Security scan handlers
│   │   └── workflow.go          # Workflow management handlers
│   ├── middleware/
│   │   ├── auth.go              # JWT authentication middleware
│   │   ├── cors.go              # CORS middleware
│   │   ├── logger.go            # Request logging middleware
│   │   └── ratelimit.go         # Rate limiting middleware
│   ├── models/
│   │   ├── user.go              # User model
│   │   ├── repository.go        # Repository model
│   │   ├── workflow.go          # Workflow model
│   │   └── embedding.go         # Scan result model
│   ├── routes/
│   │   └── routes.go            # API route definitions
│   ├── services/
│   │   ├── ai.go                # AI service (Gemini/Groq)
│   │   ├── auth.go              # Authentication service
│   │   ├── embedding.go         # Code embedding service
│   │   ├── github.go            # GitHub API service
│   │   ├── notification.go      # Notification service
│   │   ├── scanner.go           # Security scanner service
│   │   └── workflow.go          # Workflow service
│   └── utils/
│       ├── crypto.go            # Encryption utilities
│       ├── jwt.go               # JWT utilities
│       ├── response.go          # HTTP response utilities
│       └── validator.go         # Input validation utilities
├── migrations/
│   └── 000001_init.up.sql       # Database schema
├── docker/
│   └── docker-compose.yml       # Docker services configuration
├── .env.example                 # Environment variables template
├── .air.toml                    # Hot reload configuration
├── Makefile                     # Build automation
└── README.md                    # This file
```

## 👨‍💻 Development

### Hot Reload

Install air for hot reload:

```bash
make install-tools
make dev
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run all checks
make check
```

### Database Management

```bash
# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Connect to database
make db-shell

# Connect to Redis
make redis-cli
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage
```

## 📝 License

MIT License - see LICENSE file for details

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📧 Contact

For questions or support, please open an issue on GitHub.

---

Built with ❤️ using Go, PostgreSQL, Redis, and AI
# vulnpoint-go
