# OpsMind AI — DevOps AI Gatekeeper

An autonomous, agentic code review system that intercepts GitHub Pull Requests 
and analyzes them for security vulnerabilities, cloud cost drift, and 
architectural violations — before code reaches production.

## What It Does
- 🔐 **Security Sentinel** — Scans diffs for hardcoded secrets, injection vectors, and CVEs
- 💰 **Cost Predictor** — Detects cloud infrastructure cost drift from IaC changes
- 🏗️ **Architecture Supervisor** — Validates code against architectural best practices
- 💬 **Automated PR Comments** — Posts findings directly to GitHub pull requests
- 📊 **Engineering Dashboard** — Real-time React UI for posture reporting (coming soon)

## Stack
| Layer | Technology |
|-------|-----------|
| Backend | Go (Golang) |
| Database | PostgreSQL 16 + pgvector |
| AI/LLM | Groq API — Llama 3.3 70B |
| Frontend | React + TypeScript |
| Infra | Docker, ngrok (dev), GitHub Actions (planned) |

## Current Build Status
✅ Project structure and Git repository  
✅ PostgreSQL + pgvector running in Docker  
✅ Database schema — 5 tables (repositories, pull_requests, agent_findings, feedback_logs, architecture_rules)  
✅ Go backend with concurrent HTTP server and graceful shutdown  
✅ Environment-based config management  
✅ PostgreSQL connection pool (pgx)  
✅ GitHub webhook handler with HMAC-SHA256 signature validation  
✅ Database models and structs for all entities  
✅ Job queue for async PR processing  
✅ Concurrent worker pool (goroutines + channels)  
✅ Groq LLM client integration  
✅ Security Sentinel agent — hardcoded secrets, injection vectors, CWE tagging  
✅ Cost Predictor agent — IaC cost drift estimation in USD/month  
✅ Architecture Supervisor agent — pattern violation detection  
✅ All 3 agents run concurrently per PR and findings persist to PostgreSQL  
✅ Real diff fetching from GitHub API  
✅ Automated Markdown comment posting back to GitHub PRs  
✅ **Verified end-to-end on a real GitHub repository with a real pull request**  
⏳ React engineering dashboard  
⏳ MLOps feedback loop (dismiss/approve exceptions)  
⏳ pgvector-based architecture rule embeddings  
⏳ CI/CD pipeline and free-tier deployment  

## Architecture
```
GitHub PR → Webhook (HMAC verified) → Go Backend → Job Queue
                                                        ↓
                                    Concurrent Worker Pool (goroutines)
                                                        ↓
                        ┌───────────────┬───────────────┬─────────────────────┐
                        ▼               ▼               ▼
              Security Sentinel   Cost Predictor   Architecture Supervisor
                        │               │               │
                        └───────────────┴───────────────┘
                                        ↓
                              PostgreSQL (findings, PR status)
                                        ↓
                          GitHub PR Comment (Markdown summary)
```

## Local Development

### Prerequisites
- Go 1.22+
- Docker Desktop
- Node.js LTS
- ngrok (for local webhook testing)

### Setup
```bash
# Start PostgreSQL
cd infra && docker compose up -d

# Run migrations
Get-Content infra/migrations/001_init_schema.sql | docker exec -i opsmind-postgres psql -U opsmind_user -d opsmind_db

# Start backend
cd backend && go run cmd/api/main.go

# In a separate terminal — expose local server for GitHub webhooks
ngrok http 8080
```

### Environment Variables
Create `backend/.env`:
```
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5434
DB_USER=opsmind_user
DB_PASSWORD=opsmind_pass_dev
DB_NAME=opsmind_db
DB_SSL_MODE=disable
GITHUB_WEBHOOK_SECRET=your_webhook_secret_here
GITHUB_TOKEN=your_github_pat_here
GROQ_API_KEY=your_groq_key_here
GROQ_MODEL_ID=llama-3.3-70b-versatile
MAX_WORKERS=5
```

### Endpoints
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Server + database health check |
| `/webhook/github` | POST | GitHub webhook receiver (pull_request events) |
| `/test/trigger` | POST | Manually injects a fake PR job for local agent testing |

## Project Structure
```
opsmind-ai/
├── backend/
│   ├── cmd/api/
│   │   └── main.go                       # Entry point, router, graceful shutdown
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go                 # Environment config management
│   │   ├── db/
│   │   │   ├── postgres.go               # PostgreSQL connection pool
│   │   │   └── repository.go             # All DB read/write operations
│   │   ├── models/
│   │   │   └── structures.go             # All entity structs and constants
│   │   ├── webhook/
│   │   │   ├── handler.go                # GitHub webhook + HMAC-SHA256 validation
│   │   │   └── test_trigger.go           # Manual test trigger endpoint
│   │   └── agents/
│   │       ├── engine.go                 # Concurrent worker pool / orchestrator
│   │       ├── groq_client.go            # Groq LLM API client
│   │       ├── github_client.go          # GitHub API client (diff fetch + comments)
│   │       ├── security_sentinel.go      # Security agent
│   │       ├── cost_predictor.go         # Cost agent
│   │       ├── architecture_supervisor.go # Architecture agent
│   │       └── comment_formatter.go      # Markdown comment builder
│   ├── go.mod
│   └── .env                              # Local secrets (gitignored)
├── frontend/                             # React dashboard (coming soon)
├── infra/
│   ├── docker-compose.yml                # PostgreSQL + pgvector container
│   └── migrations/
│       └── 001_init_schema.sql           # Full database schema
└── docs/                                 # Architecture docs (coming soon)
```

## Status
🚧 Core AI review engine is fully functional and verified on real GitHub PRs. 
Now building the dashboard and feedback loop.