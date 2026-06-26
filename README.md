# OpsMind AI — DevOps AI Gatekeeper

An autonomous, agentic code review system that intercepts GitHub Pull Requests and analyzes them for security vulnerabilities, cloud cost drift, and architectural violations — before code reaches production.

## 🚀 Live Demo
- **Dashboard**: https://opsmind-frontend.onrender.com
- **Backend API**: https://opsmind-backend-xqmc.onrender.com

> Hosted on Render's free tier — the backend may take 30-60 seconds to wake up after periods of inactivity. The dashboard shows real findings from actual GitHub pull requests analyzed by live AI agents.

---

## The Problem

Engineering teams merge code every day without knowing the full blast radius of what they're shipping. A single pull request can:

- **Introduce a hardcoded API key** that gets scraped by bots within minutes of going live
- **Resize a database instance** and silently add $15,000/month to the AWS bill
- **Bypass the repository pattern** and create a direct database connection in an HTTP handler — a ticking technical debt time bomb

Traditional code review is manual, inconsistent, and doesn't scale. Senior engineers get review fatigue. Junior engineers don't know what they don't know. Security and FinOps concerns fall through the cracks.

---

## What OpsMind AI Solves

OpsMind AI acts as an always-on, automated senior reviewer that intercepts every pull request before it merges and runs three specialized AI agents concurrently:

- **Security Sentinel** — Scans the actual code diff for hardcoded secrets, SQL injection vectors, insecure authentication patterns, and known vulnerability classes. Tags each finding with CWE IDs and provides exact file/line references with remediation code snippets.

- **Cost Predictor** — Reads infrastructure-as-code changes (Terraform, CloudFormation) and estimates the real-dollar monthly cost impact using AWS pricing reasoning. Flags anything that would meaningfully change the cloud bill.

- **Architecture Supervisor** — Validates code against both generic best practices (repository pattern, separation of concerns, error handling) and your own custom company rules stored as vector embeddings, retrieved via semantic similarity search so fuzzy matches still get caught.

Within seconds of a PR being opened, the engineer gets a detailed Markdown comment posted directly to the pull request, and the engineering dashboard updates automatically.

---

## What Makes It Different

Most code review tools are either static analysis (rule-based, brittle, high false-positive rate) or expensive SaaS products (Snyk, Datadog, Veracode) that cost thousands per month.

OpsMind AI is different in four ways:

**1. Three concurrent AI agents, not one.** Security, cost, and architecture are completely separate concerns requiring different mental models. Running them in parallel via goroutines means each agent gets full context and reasons independently — no compromise between competing priorities.

**2. Self-improving via MLOps feedback.** When an engineer dismisses a finding as a false positive or approves an exception, that decision is logged to a feedback table. Every override teaches the system what matters to your team specifically. The architecture rules system takes this further — engineers define rules in plain English, they get embedded into pgvector, and future PRs are checked against them via semantic similarity search, not keyword matching.

**3. Interactive RAG chatbot.** Instead of reading static reports, engineers can ask natural language questions — "what security issues exist in admin.go?" or "show me CWE-798 findings" — and get streamed, grounded answers with clickable source citations backed by hybrid search (pgvector + PostgreSQL full-text).

**4. Built entirely free.** The entire stack (Go backend, PostgreSQL + pgvector, React dashboard, GitHub Actions CI, cloud deployment on Render) costs $0/month. This makes it accessible to individual engineers, small teams, and open source projects that can't afford enterprise security tooling.

---

## How People Can Use It

### Option 1 — Use the live demo
Visit https://opsmind-frontend.onrender.com to see the dashboard with real PR findings. Click the blue chat button (💬) in the bottom right to ask the RAG chatbot anything about the findings.

### Option 2 — Connect it to your own GitHub repo (free)
1. Fork or clone this repo
2. Create a free account on [Render](https://render.com), [Groq](https://console.groq.com), and [Google AI Studio](https://aistudio.google.com/apikey)
3. Deploy the backend and frontend following the setup instructions below
4. Add a GitHub webhook pointing at your Render backend URL
5. Every PR you open will now get automatically analyzed and commented on

### Option 3 — Run it locally with Docker
```bash
cd infra
# Fill in your secrets in .env.docker
docker compose --env-file .env.docker up -d
ngrok http 8080  # expose for GitHub webhooks
```
The full stack (Postgres, backend, React dashboard) comes up in one command.

---

## What It Does
- 🔐 **Security Sentinel** — Scans diffs for hardcoded secrets, injection vectors, and CVEs with CWE tagging
- 💰 **Cost Predictor** — Detects cloud infrastructure cost drift from IaC changes with real USD/month estimates
- 🏗️ **Architecture Supervisor** — Validates code against best practices + your custom pgvector-embedded rules
- 💬 **Automated PR Comments** — Posts detailed Markdown findings directly to GitHub pull requests
- 📊 **Engineering Dashboard** — Live React UI with metrics, trend charts, PR audits, and per-finding detail
- 🤖 **RAG Chatbot** — Hybrid search (pgvector + full-text) with SSE streaming and source citations
- 🔁 **MLOps Feedback Loop** — Engineers dismiss false positives or approve exceptions, logged for retraining
- 🐳 **Fully Containerized** — Entire stack runs via a single `docker compose up`
- ☁️ **Live in Production** — Deployed free on Render with auto-deploy on every push
- ⚙️ **CI Pipeline** — GitHub Actions validates backend, frontend, and Docker builds on every commit

---

## Stack
| Layer | Technology |
|-------|-----------|
| Backend | Go (Golang) — concurrent, compiled, zero-dependency binary |
| Database | PostgreSQL 16 + pgvector — semantic + full-text hybrid search |
| AI/LLM | Groq API — Llama 3.3 70B for agents + RAG; Gemini embedding-001 for vectors |
| Frontend | React + TypeScript (Vite), served via Nginx in production |
| Infra | Docker, Docker Compose, Render (hosting), GitHub Actions (CI), ngrok (dev) |

---

## Current Build Status
✅ Project structure and Git repository  
✅ PostgreSQL + pgvector running in Docker  
✅ Database schema — 5 tables (repositories, pull_requests, agent_findings, feedback_logs, architecture_rules)  
✅ Go backend with concurrent HTTP server and graceful shutdown  
✅ Environment-based config management (local, Docker, and Render-compatible)  
✅ PostgreSQL connection pool (pgx)  
✅ GitHub webhook handler with HMAC-SHA256 signature validation  
✅ Database models and structs for all entities  
✅ Job queue for async PR processing  
✅ Concurrent worker pool (goroutines + channels)  
✅ Groq LLM client integration  
✅ Security Sentinel agent — hardcoded secrets, injection vectors, CWE tagging  
✅ Cost Predictor agent — IaC cost drift estimation in USD/month  
✅ Architecture Supervisor agent — pattern violation detection + custom pgvector rules  
✅ All 3 agents run concurrently per PR and findings persist to PostgreSQL  
✅ Real diff fetching from GitHub API  
✅ Automated Markdown comment posting back to GitHub PRs  
✅ Dashboard REST API — metrics, PR list, per-PR findings (with multi-origin CORS)  
✅ React dashboard — live MetricsGrid (critical flaws, cost drift, pass rate)  
✅ React dashboard — PullRequestTable with status/score badges, row selection  
✅ React dashboard — FindingDetails panel with severity, file/line, remediation  
✅ React dashboard — Trend charts (security score + cost drift over time)  
✅ React dashboard — Repositories page with per-repo stats and risk badges  
✅ React dashboard — Settings page to manage custom architecture rules  
✅ MLOps feedback loop — dismiss/approve exception buttons, logged to `feedback_logs`  
✅ Auto-refresh polling — dashboard updates every 30s with manual refresh option  
✅ Gemini embedding-001 integration for architecture rule vector embeddings  
✅ pgvector semantic similarity search wired into Architecture Supervisor agent  
✅ RAG chatbot — hybrid pgvector + PostgreSQL full-text search over findings and rules  
✅ RAG chatbot — SSE streaming with word-by-word answer delivery and source citations  
✅ RAG chatbot — minimize/maximize/close UI with floating chat button  
✅ Multi-stage Dockerfiles for backend (Go/Alpine) and frontend (Node build → Nginx)  
✅ Full stack containerized via Docker Compose (Postgres + backend + frontend)  
✅ GitHub Actions CI pipeline — backend build + go vet + tests, frontend build, Docker validation  
✅ Go test suite — 10 tests covering HMAC validation, comment formatting, and JSON parsing  
✅ **Deployed live on Render** — managed PostgreSQL, Dockerized Go backend, static React frontend  
✅ **Verified end-to-end on a real GitHub repository — local dev, fully Dockerized, and live in production**  

---

## Architecture
```
GitHub PR → Webhook (HMAC verified) → Go Backend → Job Queue
                                                        ↓
                                    Concurrent Worker Pool (goroutines)
                                                        ↓
                        ┌───────────────┬───────────────┬─────────────────────┐
                        ▼               ▼               ▼
              Security Sentinel   Cost Predictor   Architecture Supervisor
                                                          ↑
                                               pgvector similarity search
                                               (custom rules via Gemini embeddings)
                        │               │               │
                        └───────────────┴───────────────┘
                                        ↓
                              PostgreSQL (findings, PR status)
                                        ↓
                    ┌───────────────────┴───────────────────┐
                    ▼                                       ▼
          GitHub PR Comment (Markdown)          React Dashboard (REST API)
                                                              ↓
                                        ┌─────────────────────────────────┐
                                        ▼                                 ▼
                             Engineer dismisses/approves          RAG Chatbot
                               findings (MLOps loop)        (hybrid search + SSE stream)
```

---

## Deployment (Render)

The production stack runs on Render's free tier:
- **PostgreSQL** — managed database with pgvector extension enabled
- **Backend** — Dockerized Go service, auto-deploys from `main` on every push
- **Frontend** — static site built from `frontend/`, auto-deploys from `main` on every push

---

## Running with Docker (recommended for local dev)

### Prerequisites
- Docker Desktop
- ngrok (for local webhook testing)

### Setup
```bash
cd infra

# Fill in real secrets in .env.docker (gitignored)
# GITHUB_WEBHOOK_SECRET, GITHUB_TOKEN, GROQ_API_KEY, GEMINI_API_KEY

docker compose --env-file .env.docker up -d

# Run migrations (first time only)
Get-Content migrations/001_init_schema.sql | docker exec -i opsmind-postgres psql -U opsmind_user -d opsmind_db

# Expose backend for GitHub webhooks (separate terminal)
ngrok http 8080
```

Dashboard: `http://localhost:3000`  
Backend API: `http://localhost:8080`

---

## Running Locally (without Docker)

### Prerequisites
- Go 1.25+
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
GEMINI_API_KEY=your_gemini_key_here
MAX_WORKERS=5
ENABLE_TEST_TRIGGER=true
```

Create `infra/.env.docker` for the containerized stack (gitignored):
```
GITHUB_WEBHOOK_SECRET=your_webhook_secret_here
GITHUB_TOKEN=your_github_pat_here
GROQ_API_KEY=your_groq_key_here
GEMINI_API_KEY=your_gemini_key_here
```

---

## API Endpoints
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Server + database health check |
| `/webhook/github` | POST | GitHub webhook receiver (pull_request events) |
| `/test/trigger` | POST | Injects a synthetic vulnerable PR for demo (rate-limited to 1/30s) |
| `/api/metrics` | GET | Dashboard top-level metrics (critical flaws, cost drift, pass rate) |
| `/api/pull-requests` | GET | List of all analyzed PRs |
| `/api/pull-requests/{id}/findings` | GET | All agent findings for a specific PR |
| `/api/findings/{id}/dismiss` | POST | Dismiss a finding as false positive or approved exception |
| `/api/rules` | GET | List all custom architecture rules |
| `/api/rules` | POST | Create a new rule (auto-embeds via Gemini) |
| `/api/rules/{id}` | DELETE | Delete an architecture rule |
| `/api/trend` | GET | PR trend data for dashboard charts |
| `/api/repos` | GET | Repository stats with aggregated PR metrics |
| `/api/chat` | POST | RAG query — returns answer + sources as JSON |
| `/api/chat/stream` | POST | RAG query — streams answer token-by-token via SSE |

---

## CI Pipeline

Every push to `main` triggers `.github/workflows/ci.yml`:
1. Builds and `go vet`'s the backend (Go 1.25, clean environment)
2. Runs 10 unit tests (HMAC validation, comment formatting, JSON parsing)
3. Type-checks and builds the React frontend with Node 22
4. Validates both Dockerfiles build successfully end-to-end

All on GitHub's free runners — zero cost.

---

## Project Structure
```
opsmind-ai/
├── .github/
│   └── workflows/
│       └── ci.yml                        # GitHub Actions CI pipeline
├── backend/
│   ├── cmd/api/
│   │   └── main.go                       # Entry point, router, graceful shutdown
│   ├── internal/
│   │   ├── api/
│   │   │   ├── dashboard_handler.go      # All REST API endpoints
│   │   │   ├── rag_handler.go            # RAG chatbot endpoints (JSON + SSE)
│   │   │   └── cors.go                   # Multi-origin CORS middleware
│   │   ├── config/
│   │   │   └── config.go                 # Environment config (supports Render's PORT)
│   │   ├── db/
│   │   │   ├── postgres.go               # PostgreSQL connection pool (env-driven DSN)
│   │   │   └── repository.go             # All DB read/write + pgvector + FTS operations
│   │   ├── models/
│   │   │   └── structures.go             # All entity structs and constants
│   │   ├── webhook/
│   │   │   ├── handler.go                # GitHub webhook + HMAC-SHA256 validation
│   │   │   ├── handler_test.go           # Webhook signature unit tests
│   │   │   └── test_trigger.go           # Rate-limited synthetic PR endpoint
│   │   └── agents/
│   │       ├── engine.go                 # Concurrent worker pool / orchestrator
│   │       ├── groq_client.go            # Groq LLM API client (complete + stream)
│   │       ├── groq_client_test.go       # JSON parsing unit tests
│   │       ├── gemini_client.go          # Gemini embeddings API client
│   │       ├── github_client.go          # GitHub API (diff fetch + PR comments)
│   │       ├── security_sentinel.go      # Security agent
│   │       ├── cost_predictor.go         # Cost agent
│   │       ├── architecture_supervisor.go # Architecture agent + pgvector rule injection
│   │       ├── comment_formatter.go      # Markdown comment builder
│   │       └── comment_formatter_test.go # Comment formatting unit tests
│   ├── Dockerfile                        # Multi-stage Go build → Alpine runtime
│   ├── go.mod
│   └── .env                              # Local secrets (gitignored)
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   │   └── client.ts                 # Axios client (60s timeout for cold starts)
│   │   ├── components/
│   │   │   ├── Layout.tsx                # Header/nav shell with page routing
│   │   │   ├── MetricsGrid.tsx           # 3 top-level metric cards with skeletons
│   │   │   ├── PullRequestTable.tsx      # PR list with status/score badges
│   │   │   ├── FindingDetails.tsx        # Per-PR findings with dismiss actions
│   │   │   ├── TrendCharts.tsx           # Security score + cost drift over time
│   │   │   ├── RepositoriesPage.tsx      # Per-repo stats and risk badges
│   │   │   ├── RulesManager.tsx          # Custom architecture rules CRUD UI
│   │   │   └── ChatPanel.tsx             # RAG chatbot with SSE streaming + resize
│   │   ├── hooks/
│   │   │   └── useAuditStream.ts         # Auto-refresh polling hook (30s interval)
│   │   ├── types/
│   │   │   └── api.ts                    # TypeScript types mirroring Go models
│   │   └── App.tsx                       # Page router and dashboard composition
│   ├── Dockerfile                        # Multi-stage Node build → Nginx runtime
│   ├── nginx.conf                        # SPA routing + API reverse proxy
│   └── package.json
├── infra/
│   ├── docker-compose.yml                # Full stack: Postgres + backend + frontend
│   ├── .env.docker                       # Docker stack secrets (gitignored)
│   └── migrations/
│       └── 001_init_schema.sql           # Full database schema with pgvector + FTS index
└── docs/
```

---

## Status
🟢 **Complete and live in production.** The full OpsMind AI system — autonomous AI agents, live dashboard, RAG chatbot with hybrid search and SSE streaming, MLOps feedback loop, pgvector semantic search, trend charts, repositories page, custom rules management, full containerization, CI pipeline with tests, and free-tier cloud deployment — is complete and verified end-to-end on a real GitHub repository.

Built entirely free. No AWS. No paid APIs beyond free tiers. No credit card required anywhere.