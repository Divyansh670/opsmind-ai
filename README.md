# OpsMind AI — DevOps AI Gatekeeper

An autonomous, agentic code review system that intercepts GitHub Pull Requests 
and analyzes them for security vulnerabilities, cloud cost drift, and 
architectural violations — before code reaches production.

## What It Does
- 🔐 **Security Sentinel** — Scans diffs for hardcoded secrets, injection vectors, and CVEs
- 💰 **Cost Predictor** — Detects cloud infrastructure cost drift from IaC changes
- 🏗️ **Architecture Supervisor** — Validates code against architectural best practices
- 💬 **Automated PR Comments** — Posts findings directly to GitHub pull requests
- 📊 **Engineering Dashboard** — Live React UI showing posture metrics, PR audits, and detailed findings
- 🔁 **MLOps Feedback Loop** — Engineers can dismiss false positives or approve exceptions, logged for future retraining
- 🐳 **Fully Containerized** — Entire stack (DB + backend + frontend) runs via a single `docker compose up`

## Stack
| Layer | Technology |
|-------|-----------|
| Backend | Go (Golang) |
| Database | PostgreSQL 16 + pgvector |
| AI/LLM | Groq API — Llama 3.3 70B |
| Frontend | React + TypeScript (Vite), served via Nginx in production |
| Infra | Docker, Docker Compose, ngrok (dev), GitHub Actions (planned) |

## Current Build Status
✅ Project structure and Git repository  
✅ PostgreSQL + pgvector running in Docker  
✅ Database schema — 5 tables (repositories, pull_requests, agent_findings, feedback_logs, architecture_rules)  
✅ Go backend with concurrent HTTP server and graceful shutdown  
✅ Environment-based config management (local + Docker networking)  
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
✅ Dashboard REST API — metrics, PR list, per-PR findings (with multi-origin CORS)  
✅ React dashboard — live MetricsGrid (critical flaws, cost drift, pass rate)  
✅ React dashboard — PullRequestTable with status/score badges, row selection  
✅ React dashboard — FindingDetails panel with severity, file/line, remediation  
✅ MLOps feedback loop — dismiss/approve exception buttons, logged to `feedback_logs`  
✅ Auto-refresh polling — dashboard updates every 30s with manual refresh option  
✅ Multi-stage Dockerfiles for backend (Go/Alpine) and frontend (Node build → Nginx)  
✅ Full stack containerized via Docker Compose (Postgres + backend + frontend)  
✅ **Verified end-to-end on a real GitHub repository — both local dev and fully Dockerized**  
⏳ pgvector-based architecture rule embeddings  
⏳ CI/CD pipeline and free-tier cloud deployment  

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
                    ┌───────────────────┴───────────────────┐
                    ▼                                       ▼
          GitHub PR Comment (Markdown)          React Dashboard (REST API)
                                                              ↓
                                        Engineer dismisses/approves findings
                                                              ↓
                                              feedback_logs (MLOps loop)
```

## Running with Docker (recommended)

### Prerequisites
- Docker Desktop
- ngrok (for local webhook testing)

### Setup
```bash
cd infra

# Fill in real secrets in .env.docker (gitignored)
# GITHUB_WEBHOOK_SECRET, GITHUB_TOKEN, GROQ_API_KEY

docker compose --env-file .env.docker up -d

# Run migrations (first time only)
Get-Content migrations/001_init_schema.sql | docker exec -i opsmind-postgres psql -U opsmind_user -d opsmind_db

# Expose backend for GitHub webhooks (separate terminal)
ngrok http 8080
```

Dashboard: `http://localhost:3000`  
Backend API: `http://localhost:8080`

## Running Locally (without Docker)

### Prerequisites
- Go 1.24+
- Docker Desktop (for Postgres only)
- Node.js 22+ (LTS)
- ngrok (for local webhook testing)

### Setup
```bash
# Start PostgreSQL only
cd infra && docker compose up -d postgres

# Run migrations
Get-Content migrations/001_init_schema.sql | docker exec -i opsmind-postgres psql -U opsmind_user -d opsmind_db

# Start backend
cd backend && go run cmd/api/main.go

# Start frontend (separate terminal)
cd frontend && npm run dev

# Expose local server for GitHub webhooks (separate terminal)
ngrok http 8080
```

Dashboard: `http://localhost:5173`  
Backend API: `http://localhost:8080`

### Environment Variables
Create `backend/.env` for local dev:
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

Create `infra/.env.docker` for the containerized stack (gitignored):
```
GITHUB_WEBHOOK_SECRET=your_webhook_secret_here
GITHUB_TOKEN=your_github_pat_here
GROQ_API_KEY=your_groq_key_here
```

### Endpoints
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Server + database health check |
| `/webhook/github` | POST | GitHub webhook receiver (pull_request events) |
| `/test/trigger` | POST | Manually injects a fake PR job for local agent testing |
| `/api/metrics` | GET | Dashboard top-level metrics (critical flaws, cost drift, pass rate) |
| `/api/pull-requests` | GET | List of all analyzed PRs |
| `/api/pull-requests/{id}/findings` | GET | All agent findings for a specific PR |
| `/api/findings/{id}/dismiss` | POST | Dismiss a finding as false positive or approved exception |

## Project Structure
```
opsmind-ai/
├── backend/
│   ├── cmd/api/
│   │   └── main.go                       # Entry point, router, graceful shutdown
│   ├── internal/
│   │   ├── api/
│   │   │   ├── dashboard_handler.go      # Metrics, PR list, findings, dismiss endpoints
│   │   │   └── cors.go                   # Multi-origin CORS middleware
│   │   ├── config/
│   │   │   └── config.go                 # Environment config management
│   │   ├── db/
│   │   │   ├── postgres.go               # PostgreSQL connection pool (env-driven DSN)
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
│   ├── Dockerfile                        # Multi-stage Go build → Alpine runtime
│   ├── go.mod
│   └── .env                              # Local secrets (gitignored)
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   │   └── client.ts                 # Axios client for backend API
│   │   ├── components/
│   │   │   ├── Layout.tsx                # Header/nav shell
│   │   │   ├── MetricsGrid.tsx           # 3 top-level metric cards
│   │   │   ├── PullRequestTable.tsx      # PR list with badges
│   │   │   └── FindingDetails.tsx        # Per-PR findings panel with dismiss actions
│   │   ├── hooks/
│   │   │   └── useAuditStream.ts         # Auto-refresh polling hook
│   │   ├── types/
│   │   │   └── api.ts                    # TypeScript types mirroring Go models
│   │   └── App.tsx                       # Dashboard page composition
│   ├── Dockerfile                        # Multi-stage Node build → Nginx runtime
│   ├── nginx.conf                        # SPA routing + API reverse proxy
│   └── package.json
├── infra/
│   ├── docker-compose.yml                # Full stack: Postgres + backend + frontend
│   ├── .env.docker                       # Docker stack secrets (gitignored)
│   └── migrations/
│       └── 001_init_schema.sql           # Full database schema
└── docs/                                 # Architecture docs (coming soon)
```

## Status
🚧 Core AI review engine, live dashboard, MLOps feedback loop, and full containerization 
are all complete and verified end-to-end on a real GitHub repository — both in local dev 
and as a fully Dockerized stack. Next up: CI/CD pipeline and free-tier cloud deployment.