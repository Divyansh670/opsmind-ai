# OpsMind AI — DevOps AI Gatekeeper

An autonomous, agentic code review system that intercepts GitHub Pull Requests 
and analyzes them for security vulnerabilities, cloud cost drift, and 
architectural violations — before code reaches production.

## What It Does
- 🔐 **Security Sentinel** — Scans diffs for hardcoded secrets, injection vectors, and CVEs
- 💰 **Cost Predictor** — Detects cloud infrastructure cost drift from IaC changes
- 🏗️ **Architecture Supervisor** — Validates code against company-wide architectural rules
- 📊 **Engineering Dashboard** — Real-time React UI for posture reporting and findings

## Stack
| Layer | Technology |
|-------|-----------|
| Backend | Go (Golang) |
| Database | PostgreSQL 16 + pgvector |
| AI/LLM | Groq API — Llama 3.3 70B |
| Frontend | React + TypeScript |
| Infra | Docker, GitHub Actions |

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
⏳ AI agent layer (Security Sentinel, Cost Predictor, Architecture Supervisor)  
⏳ Groq LLM integration and structured output parsing  
⏳ GitHub PR comment posting  
⏳ React dashboard  
⏳ MLOps feedback loop  
⏳ CI/CD pipeline and deployment  

## Architecture
```
GitHub PR → Webhook → Go Backend → Job Queue → AI Agents → PostgreSQL
                                                          → GitHub Comments
                                                          → React Dashboard
```

## Local Development

### Prerequisites
- Go 1.22+
- Docker Desktop
- Node.js LTS

### Setup
```bash
# Start PostgreSQL
cd infra && docker compose up -d

# Run migrations
Get-Content infra/migrations/001_init_schema.sql | docker exec -i opsmind-postgres psql -U opsmind_user -d opsmind_db

# Start backend
cd backend && go run cmd/api/main.go
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
GITHUB_WEBHOOK_SECRET=
GITHUB_TOKEN=
GROQ_API_KEY=your_groq_key_here
GROQ_MODEL_ID=llama-3.3-70b-versatile
MAX_WORKERS=5
```

### Endpoints
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Server + database health check |
| `/webhook/github` | POST | GitHub webhook receiver |

## Project Structure
```
opsmind-ai/
├── backend/
│   ├── cmd/api/
│   │   └── main.go              # Entry point, router, graceful shutdown
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go        # Environment config management
│   │   ├── db/
│   │   │   └── postgres.go      # PostgreSQL connection pool
│   │   ├── models/
│   │   │   └── structures.go    # All entity structs and constants
│   │   ├── webhook/
│   │   │   └── handler.go       # GitHub webhook + HMAC-SHA256 validation
│   │   ├── agents/              # AI agents (coming soon)
│   │   └── middleware/          # HTTP middleware (coming soon)
│   ├── go.mod
│   └── .env                     # Local secrets (gitignored)
├── frontend/                    # React dashboard (coming soon)
├── infra/
│   ├── docker-compose.yml       # PostgreSQL + pgvector container
│   └── migrations/
│       └── 001_init_schema.sql  # Full database schema
└── docs/                        # Architecture docs (coming soon)
```

## Status
🚧 Actively being built — follow along for daily progress