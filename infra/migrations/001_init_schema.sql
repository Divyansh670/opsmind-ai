-- =========================================
-- OpsMind AI - Initial Schema
-- =========================================

-- Enable pgvector (in case this runs on a fresh DB)
CREATE EXTENSION IF NOT EXISTS vector;

-- =========================================
-- Table: repositories
-- =========================================
CREATE TABLE IF NOT EXISTS repositories (
    id          SERIAL PRIMARY KEY,
    repo_name   VARCHAR(255) NOT NULL UNIQUE,
    git_url     VARCHAR(512) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =========================================
-- Table: pull_requests
-- =========================================
CREATE TABLE IF NOT EXISTS pull_requests (
    id            SERIAL PRIMARY KEY,
    repo_id       INTEGER NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    pr_number     INTEGER NOT NULL,
    head_commit   VARCHAR(64) NOT NULL,
    author        VARCHAR(255),
    status        VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- PENDING, APPROVED, FLAGGED
    security_score INTEGER DEFAULT 0,
    cost_drift_usd NUMERIC(12,2) DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (repo_id, pr_number)
);

-- =========================================
-- Table: agent_findings
-- =========================================
CREATE TABLE IF NOT EXISTS agent_findings (
    id            SERIAL PRIMARY KEY,
    pr_id         INTEGER NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    agent_name    VARCHAR(100) NOT NULL, -- SecuritySentinel, CostPredictor, ArchitectureSupervisor
    severity      VARCHAR(20) NOT NULL,  -- CRITICAL, HIGH, MEDIUM, LOW
    cwe_id        VARCHAR(50),
    file_path     VARCHAR(512),
    line_number   INTEGER,
    description   TEXT,
    remediation   TEXT,
    dismissed     BOOLEAN NOT NULL DEFAULT FALSE,
    dismiss_reason TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =========================================
-- Table: feedback_logs (for MLOps loop, used later)
-- =========================================
CREATE TABLE IF NOT EXISTS feedback_logs (
    id              SERIAL PRIMARY KEY,
    finding_id      INTEGER NOT NULL REFERENCES agent_findings(id) ON DELETE CASCADE,
    diff_snippet    TEXT,
    agent_explanation TEXT,
    engineer_reason TEXT,
    action          VARCHAR(50) NOT NULL, -- DISMISSED, APPROVED_EXCEPTION
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =========================================
-- Table: architecture_rules (for vector search, used later)
-- =========================================
CREATE TABLE IF NOT EXISTS architecture_rules (
    id          SERIAL PRIMARY KEY,
    rule_text   TEXT NOT NULL,
    embedding   vector(768), -- dimension depends on embedding model chosen later
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =========================================
-- Indexes for performance
-- =========================================
CREATE INDEX IF NOT EXISTS idx_pr_repo_id ON pull_requests(repo_id);
CREATE INDEX IF NOT EXISTS idx_findings_pr_id ON agent_findings(pr_id);
CREATE INDEX IF NOT EXISTS idx_findings_severity ON agent_findings(severity);