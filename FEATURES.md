# VulnPilot - Features & Workflow Documentation

## Table of Contents

1. [Overview](#overview)
2. [Core Features](#core-features)
3. [User Workflows](#user-workflows)
4. [Technical Architecture](#technical-architecture)
5. [Feature Deep Dive](#feature-deep-dive)
6. [Integration Workflows](#integration-workflows)
7. [Use Cases](#use-cases)

---

## Overview

VulnPilot is an AI-powered security vulnerability scanner designed to help developers identify, understand, and fix security issues in their code. It combines traditional security scanning tools with cutting-edge AI technology to provide comprehensive security analysis.

**Target Users:**
- Security Engineers
- DevOps Teams
- Software Developers
- Code Reviewers
- Security Auditors

**Core Mission:** Democratize security testing by making advanced vulnerability detection accessible through AI-powered automation.

---

## Core Features

### 1. 🔐 Authentication & Authorization

#### GitHub OAuth Integration
- **Single Sign-On (SSO)** via GitHub
- Automatic repository access
- Secure token management
- No password storage required

**Workflow:**
```
User → GitHub Login → OAuth Consent → Token Exchange → VulnPilot Access
```

**Benefits:**
- Seamless onboarding
- Repository permissions inherited from GitHub
- Automatic synchronization with GitHub account

#### JWT Token Management
- **Stateless Authentication** using JSON Web Tokens
- Token expiration: 24 hours (configurable)
- Refresh tokens: 7 days validity
- Secure token storage with bcrypt hashing

**Security Features:**
- HMAC-SHA256 signing
- Automatic token refresh
- Revocation support
- Rate-limited token generation

---

### 2. 🤖 AI-Powered Code Analysis

#### Dual AI Provider Support

**Google Gemini Integration**
- **Model:** Gemini Pro
- **Strengths:** Detailed, comprehensive analysis
- **Use Cases:** Deep code review, complex vulnerability detection
- **Response Time:** 2-5 seconds

**Groq Integration**
- **Model:** Mixtral-8x7B
- **Strengths:** Ultra-fast responses (sub-second)
- **Use Cases:** Quick scans, chatbot interactions
- **Response Time:** < 1 second

**Automatic Failover:**
```
Request → Try Gemini → If fails → Fallback to Groq → Response
```

#### Code Analysis Capabilities

**1. Vulnerability Detection**
- SQL Injection patterns
- Cross-Site Scripting (XSS)
- Command Injection
- Path Traversal
- Hardcoded Credentials
- Insecure Deserialization
- Cryptographic weaknesses

**2. Security Scoring**
- 0-100 security score
- Weighted by vulnerability severity
- Color-coded risk levels (Critical/High/Medium/Low)
- Actionable recommendations

**3. Code Similarity Analysis**
- Duplicate code detection
- Plagiarism checking
- Code reuse identification
- Fingerprint-based matching

---

### 3. 🛡️ Security Scanning Suite

#### Network Scanning (Nmap)

**Purpose:** Discover open ports and network services

**Features:**
- Port range scanning (customizable)
- Service version detection
- OS fingerprinting
- Vulnerability mapping

**Workflow:**
```
User Request → Background Job → Nmap Execution → Parse Results → Store in DB → Notify User
```

**Output:**
- Open ports list
- Service versions
- Potential vulnerabilities
- Risk assessment

#### Web Vulnerability Scanning (Nikto)

**Purpose:** Identify web server vulnerabilities

**Features:**
- 6700+ vulnerability checks
- SSL/TLS configuration analysis
- Outdated software detection
- Misconfigurations

**Use Cases:**
- Web application auditing
- Server hardening validation
- Compliance checking

#### Directory Brute-forcing (Gobuster)

**Purpose:** Discover hidden files and directories

**Features:**
- Wordlist-based enumeration
- Custom wordlist support
- Concurrent scanning
- Pattern matching

**Common Findings:**
- Backup files
- Admin panels
- Configuration files
- API endpoints

---

### 4. 🔄 Workflow Automation

#### Visual Workflow Builder

**Components:**
- **Nodes:** Individual scan/analysis tasks
- **Edges:** Connections defining execution order
- **Triggers:** Manual or scheduled execution

**Node Types:**
1. **Scan Nodes:** Nmap, Nikto, Gobuster
2. **Analysis Nodes:** Code review, AI analysis
3. **Notification Nodes:** Email, Slack alerts
4. **Conditional Nodes:** Branch based on results

**Workflow Example:**
```
GitHub Repo → Clone → Code Analysis → If vulnerabilities found → Notify Team
                                    → If clean → Archive Report
```

#### Scheduling

**Cron-like Scheduling:**
- Hourly, Daily, Weekly, Monthly
- Custom cron expressions
- Timezone support
- Next run prediction

**Use Cases:**
- Continuous security monitoring
- Regular compliance scans
- Automated security gates in CI/CD

---

### 5. 💬 AI Security Chatbot

#### Conversational AI Assistant

**Capabilities:**
1. **Vulnerability Explanations**
   - Plain English descriptions
   - Real-world exploit examples
   - Impact assessment
   - CVSS scoring context

2. **Remediation Guidance**
   - Step-by-step fixes
   - Code examples
   - Best practices
   - Testing recommendations

3. **Security Q&A**
   - General security questions
   - Technology-specific guidance
   - Threat modeling assistance
   - Compliance queries

**Conversation Flow:**
```
User Question → Context Analysis → AI Processing → Structured Response → Follow-up Suggestions
```

**Example Interactions:**

**Q:** "What is SQL Injection?"
**A:** Detailed explanation with:
- Definition
- How it works
- Attack examples
- Prevention methods
- Code samples

**Q:** "How do I fix this vulnerability in my Python code?"
**A:** Contextual response with:
- Fixed code snippet
- Explanation of changes
- Additional security measures
- Testing approach

---

### 6. 📧 Multi-Channel Notifications

#### Email Notifications (SMTP)

**Triggers:**
- Scan completion
- Critical vulnerabilities found
- Workflow failures
- Schedule reminders

**Templates:**
- Scan Summary
- Vulnerability Alert
- Weekly Digest
- Custom messages

**Configuration:**
- Gmail compatible
- Custom SMTP servers
- TLS/SSL support
- HTML/Plain text

#### Slack Integration

**Webhook-based Notifications:**
- Real-time alerts
- Rich formatting with attachments
- Color-coded by severity
- Interactive buttons (future)

**Message Types:**
- Success/Failure indicators
- Vulnerability counts
- Quick stats
- Deep link to results

**Notification Workflow:**
```
Event Triggered → Check User Preferences → Format Message → Send via Email + Slack → Log Delivery
```

---

### 7. 🔒 Security & Rate Limiting

#### Distributed Rate Limiting

**Technology:** Redis-backed token bucket algorithm

**Features:**
- Per-user rate limits
- Per-IP fallback for anonymous requests
- Configurable thresholds
- Automatic quota reset

**Default Limits:**
- 100 requests per 15 minutes per user
- Burst allowance: 10 requests
- Rate limit headers included in responses

**Protection Against:**
- API abuse
- DDoS attacks
- Credential stuffing
- Resource exhaustion

#### Encryption & Hashing

**Password Storage:**
- Bcrypt hashing (cost factor: 10)
- Salted automatically
- No plaintext storage

**Data Encryption:**
- AES-256-GCM for sensitive data
- Secure key derivation
- Nonce-based encryption
- Authenticated encryption

**Token Generation:**
- Cryptographically secure random
- 32-byte minimum entropy
- URL-safe encoding

---

### 8. 🗄️ Data Management

#### Database Schema

**Users Table:**
- UUID primary keys
- GitHub OAuth data
- Access tokens (encrypted)
- Activity tracking

**Repositories Table:**
- GitHub metadata
- Language detection
- Last scan timestamp
- Vulnerability history

**Workflows Table:**
- JSONB node/edge storage
- Schedule configuration
- Execution history
- Success/failure tracking

**Scan Results Table:**
- Scan type and target
- JSONB results storage
- Status tracking
- Error logging

**Performance Optimizations:**
- Indexed foreign keys
- Composite indexes for common queries
- JSONB GIN indexes for search
- Automatic timestamp updates

---

## User Workflows

### Workflow 1: First-Time User Onboarding

```
1. Landing Page → Click "Sign in with GitHub"
2. GitHub OAuth → Authorize VulnPilot
3. Redirect → Dashboard with repository list
4. Select Repository → View files
5. Choose File → Run Analysis
6. View Results → Get AI recommendations
7. Fix Code → Re-scan to verify
```

**Time to First Scan:** < 2 minutes

---

### Workflow 2: Automated Repository Scanning

```
1. User creates workflow:
   - Add "GitHub Clone" node
   - Add "Code Analysis" node
   - Add "Notification" node
   - Connect nodes

2. Configure schedule:
   - Set frequency (e.g., daily at 2 AM)
   - Enable notifications

3. Save workflow → Automatic execution begins

4. Results:
   - Daily security reports
   - Email alerts for new vulnerabilities
   - Historical tracking in dashboard
```

**Benefits:**
- Zero manual intervention
- Continuous monitoring
- Trend analysis over time

---

### Workflow 3: AI-Assisted Vulnerability Remediation

```
1. User receives vulnerability alert
2. Opens chatbot: "Explain SQL Injection in my code"
3. AI provides detailed explanation
4. User asks: "How do I fix it?"
5. AI generates fixed code snippet
6. User applies fix
7. User asks: "How do I test this?"
8. AI provides testing approach
9. User runs code analysis to verify fix
```

**Learning Loop:** Users understand vulnerabilities while fixing them

---

### Workflow 4: CI/CD Integration (Future)

```
Pull Request → Webhook → VulnPilot Scan → Report to GitHub → Block if critical → Notify Developer
```

**Prevents:**
- Vulnerable code reaching production
- Security debt accumulation
- Compliance violations

---

## Technical Architecture

### Request Flow

```
Client Request
    ↓
CORS Middleware → Validate origin
    ↓
Logger Middleware → Log request
    ↓
Rate Limit Middleware → Check quota (Redis)
    ↓
Auth Middleware → Validate JWT
    ↓
Handler → Process request
    ↓
Service Layer → Business logic
    ↓
Database/External APIs → Data operations
    ↓
Response → JSON formatted
```

### Scan Execution Architecture

```
API Request (Scan endpoint)
    ↓
Create Scan Record (Status: pending)
    ↓
Launch Background Goroutine
    ↓
Update Status (Status: running)
    ↓
Execute Tool (nmap/nikto/gobuster)
    ↓
Parse Output → JSON
    ↓
Update Record (Status: completed, Results: JSON)
    ↓
Trigger Notification Service
    ↓
Send Email + Slack notification
```

**Concurrency:** Non-blocking scans using Go routines

---

### AI Analysis Pipeline

```
Code Input
    ↓
Embedding Generation → SHA256 hash + Keywords
    ↓
Pattern Matching → Vulnerability detection
    ↓
AI Provider Selection (Gemini/Groq)
    ↓
Prompt Engineering → Structured query
    ↓
API Call → AI provider
    ↓
Response Parsing → Structured data
    ↓
Security Scoring → 0-100 calculation
    ↓
Recommendation Generation → Action items
    ↓
Response to User
```

---

## Feature Deep Dive

### Code Embedding System

**Purpose:** Enable similarity detection and pattern matching without AI

**How It Works:**

1. **Fingerprinting:**
   ```go
   Input Code → Hash (SHA256) → Keywords Extraction → Fingerprint Generation
   ```

2. **Keyword Extraction:**
   - Security keywords (password, token, SQL, exec)
   - Language-specific keywords (import, async, func)
   - Pattern matching (regex-based)

3. **Similarity Calculation:**
   ```
   Jaccard Similarity = |Keywords1 ∩ Keywords2| / |Keywords1 ∪ Keywords2|
   ```

**Use Cases:**
- Duplicate vulnerability detection
- Code reuse analysis
- Similar pattern matching
- License compliance

---

### Notification System Design

**Event-Driven Architecture:**

```
Event Trigger → Notification Service
    ↓
Check User Preferences
    ↓
Format Message (Template)
    ↓
Parallel Execution:
    - Email Thread → SMTP Send
    - Slack Thread → Webhook POST
    ↓
Log Delivery Status
    ↓
Update User Notification History
```

**Reliability:**
- Async execution (non-blocking)
- Retry logic (3 attempts)
- Error logging
- Graceful degradation

---

## Integration Workflows

### GitHub Integration

**OAuth Flow:**
```
User → VulnPilot → GitHub Authorization → User Approves → 
GitHub → Authorization Code → VulnPilot → Token Exchange → 
GitHub → Access Token → VulnPilot → Store Encrypted Token
```

**Repository Access:**
```
List Repos → GET /user/repos (GitHub API)
Get Files → GET /repos/:owner/:repo/contents/:path
Read Content → GET /repos/:owner/:repo/contents/:file (Accept: application/vnd.github.v3.raw)
```

**Permissions Required:**
- `user:email` - User email access
- `repo` - Repository access (public and private)

---

### AI Provider Integration

**Gemini API:**
```
POST https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent
Headers:
  - Content-Type: application/json
Body:
  {
    "contents": [{"parts": [{"text": "prompt"}]}]
  }
```

**Groq API:**
```
POST https://api.groq.com/openai/v1/chat/completions
Headers:
  - Authorization: Bearer <API_KEY>
  - Content-Type: application/json
Body:
  {
    "model": "mixtral-8x7b-32768",
    "messages": [{"role": "user", "content": "prompt"}]
  }
```

---

## Use Cases

### Use Case 1: Startup Security Audit

**Scenario:** Early-stage startup needs security assessment before funding round

**Workflow:**
1. Connect GitHub repositories
2. Run comprehensive scan on all repos
3. AI analyzes code for vulnerabilities
4. Generate executive summary report
5. Remediate critical issues
6. Re-scan to verify fixes

**Outcome:** Clean security posture for investor due diligence

---

### Use Case 2: Enterprise Continuous Monitoring

**Scenario:** Large enterprise with 100+ repositories

**Workflow:**
1. Create workflow for each team
2. Schedule daily scans at off-peak hours
3. Route alerts to team Slack channels
4. Weekly digest reports to security team
5. Quarterly compliance reports

**Benefits:**
- Proactive vulnerability detection
- Reduced mean time to remediation
- Compliance documentation

---

### Use Case 3: Security Training

**Scenario:** Developer training on secure coding

**Workflow:**
1. Instructor uploads vulnerable code samples
2. Students analyze code using VulnPilot
3. AI chatbot explains vulnerabilities
4. Students fix code based on recommendations
5. Verify fixes with re-scan
6. Compare before/after security scores

**Learning Outcomes:**
- Hands-on vulnerability identification
- AI-assisted learning
- Immediate feedback loop

---

### Use Case 4: Open Source Security

**Scenario:** Maintaining security in public repositories

**Workflow:**
1. Scan public repositories
2. Identify common vulnerability patterns
3. Create automated workflows for new PRs
4. Notify contributors of security issues
5. Track remediation progress

**Community Benefits:**
- Improved OSS security
- Educational resource
- Security best practices enforcement

---

## Summary

VulnPilot combines **traditional security tools** with **modern AI technology** to provide:

✅ **Comprehensive Coverage** - Multiple scan types (network, web, code)  
✅ **AI-Powered Insights** - Intelligent analysis and recommendations  
✅ **Automation** - Workflow builder for continuous monitoring  
✅ **Educational** - Learn security through AI chatbot  
✅ **Scalable** - From individual developers to enterprises  
✅ **Integrated** - GitHub OAuth, Email, Slack  
✅ **Secure** - Enterprise-grade encryption and authentication  

**Vision:** Make enterprise-grade security accessible to everyone through AI automation.
