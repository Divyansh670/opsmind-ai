import React from "react";
import {
  GitBranch,
  Webhook,
  Bot,
  Shield,
  DollarSign,
  Building2,
  MessageSquare,
  Settings,
  CheckCircle2,
  ExternalLink,
} from "lucide-react";

export default function HowToUsePage() {
  return (
    <div>
      <div style={styles.pageHeader}>
        <h1 style={styles.pageTitle}>HOW TO USE OPSMIND AI</h1>

        <p style={styles.pageSubtitle}>
          Connect OpsMind AI to your GitHub repository in under 5 minutes —
          no code changes required.
        </p>
      </div>

      {/* ---------------- Quick Start ---------------- */}

      <Section title="⚡ Quick Start — 3 Steps">
        <div style={styles.stepsGrid}>
          <StepCard
            number="1"
            title="Connect your GitHub Repository"
            description="Go to GitHub → Repository → Settings → Webhooks → Add Webhook."
            code={`Payload URL
https://opsmind-backend-xqmc.onrender.com/webhook/github

Content Type
application/json

Events
Pull Requests`}
          />

          <StepCard
            number="2"
            title="Configure Webhook Secret"
            description="Create a random secret string. Use exactly the same secret inside your backend."
            code={`Example Secret

my-team-opsmind-secret-2026`}
          />

          <StepCard
            number="3"
            title="Open a Pull Request"
            description="Every new PR is analyzed automatically and an AI review comment is posted."
            code={`OpsMind AI automatically

✓ Security Scan
✓ Cost Drift
✓ Architecture Rules
✓ AI Review`}
          />
        </div>
      </Section>

      {/* ---------------- Dashboard ---------------- */}

      <Section title="📊 Dashboard Features">
        <div style={styles.featuresGrid}>
          <FeatureCard
            icon={<Shield size={20} color="#f87171" />}
            title="Security Findings"
            description="Inspect vulnerabilities, CWE IDs, affected files and AI generated fixes."
            tip="Mark findings as False Positive whenever appropriate."
          />

          <FeatureCard
            icon={<DollarSign size={20} color="#fbbf24" />}
            title="Cost Drift"
            description="Terraform, Kubernetes and CloudFormation changes are analyzed automatically."
            tip="Monitor which PR increases cloud cost the most."
          />

          <FeatureCard
            icon={<Building2 size={20} color="#60a5fa" />}
            title="Architecture Rules"
            description="Create company specific coding rules written in plain English."
            tip="Repository Layer → Service Layer → Controller Layer."
          />

          <FeatureCard
            icon={<MessageSquare size={20} color="#34d399" />}
            title="RAG Chatbot"
            description="Ask questions about your repositories, pull requests and findings."
            tip='Example: "Show all Critical vulnerabilities".'
          />
        </div>
      </Section>
           {/* ---------------- Use Cases ---------------- */}

      <Section title="🎯 Use Cases">
        <div style={styles.useCasesGrid}>
          <UseCaseCard
            emoji="👨‍💻"
            title="Solo Developer"
            description="Automatically review every Pull Request before merging."
            examples={[
              "Catch hardcoded API keys",
              "Estimate Terraform costs",
              "Maintain consistent code quality",
            ]}
          />

          <UseCaseCard
            emoji="👥"
            title="Small Teams"
            description="Enforce your team's coding standards with AI."
            examples={[
              "Create custom architecture rules",
              "Track security debt",
              "Review findings using AI Chat",
            ]}
          />

          <UseCaseCard
            emoji="🎓"
            title="Learning & Portfolio"
            description="Demonstrate a real DevSecOps pipeline."
            examples={[
              "Trigger a fake vulnerable PR",
              "Watch AI agents analyze it",
              "Use it in portfolio demos",
            ]}
          />

          <UseCaseCard
            emoji="🔬"
            title="Security Research"
            description="Analyze security trends across repositories."
            examples={[
              "Connect multiple repositories",
              "Find Critical issues instantly",
              "Export findings through the API",
            ]}
          />
        </div>
      </Section>

      {/* ---------------- Self Host ---------------- */}

      <Section title="🚀 Self-Host for Free">
        <p style={styles.selfHostIntro}>
          Fork the project and deploy your own private instance in under
          30 minutes.
        </p>

        <div style={styles.selfHostSteps}>
          {[
            {
              icon: <GitBranch size={16} color="#60a5fa" />,
              step: "Fork Repository",
              detail: "Fork OpsMind AI on GitHub",
            },
            {
              icon: <Bot size={16} color="#fbbf24" />,
              step: "Create API Keys",
              detail: "Groq + Gemini (Free)",
            },
            {
              icon: <CheckCircle2 size={16} color="#34d399" />,
              step: "Deploy",
              detail: "Deploy Frontend + Backend + PostgreSQL",
            },
            {
              icon: <Webhook size={16} color="#f87171" />,
              step: "Configure Webhooks",
              detail: "Connect GitHub repositories",
            },
            {
              icon: <Settings size={16} color="#9ca3af" />,
              step: "Architecture Rules",
              detail: "Define your own coding standards",
            },
          ].map((item, i) => (
            <div key={i} style={styles.selfHostStep}>
              <div style={styles.selfHostNum}>{i + 1}</div>

              <div style={styles.selfHostIcon}>
                {item.icon}
              </div>

              <div>
                <div style={styles.selfHostStepTitle}>
                  {item.step}
                </div>

                <div style={styles.selfHostStepDetail}>
                  {item.detail}
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* ✅ FIXED: Missing opening <a> tag */}

        <a
          href="https://github.com/Divyansh670/opsmind-ai"
          target="_blank"
          rel="noopener noreferrer"
          style={styles.githubLink}
        >
          <GitBranch size={16} />
          <span>View Full Setup Guide on GitHub</span>
          <ExternalLink size={12} />
        </a>
      </Section>

      {/* ---------------- API ---------------- */}

      <Section title="🔌 REST API">
        <p style={styles.apiIntro}>
          Every feature inside OpsMind AI is exposed through REST APIs.
        </p>

        <div style={styles.apiGrid}>
          {[
            {
              method: "GET",
              path: "/api/metrics",
              desc: "Dashboard metrics",
            },
            {
              method: "GET",
              path: "/api/pull-requests",
              desc: "All Pull Requests",
            },
            {
              method: "GET",
              path: "/api/pull-requests/{id}/findings",
              desc: "Security findings",
            },
            {
              method: "POST",
              path: "/api/findings/{id}/dismiss",
              desc: "Dismiss finding",
            },
            {
              method: "GET",
              path: "/api/repos",
              desc: "Repository statistics",
            },
            {
              method: "POST",
              path: "/api/chat/stream",
              desc: "AI Chatbot",
            },
            {
              method: "GET",
              path: "/api/trend",
              desc: "Security trends",
            },
            {
              method: "POST",
              path: "/test/trigger",
              desc: "Trigger demo PR",
            },
          ].map((ep, i) => (
            <div key={i} style={styles.apiRow}>
              <span
                style={{
                  ...styles.methodBadge,
                  backgroundColor:
                    ep.method === "GET"
                      ? "#17351d"
                      : "#3a1717",
                  color:
                    ep.method === "GET"
                      ? "#34d399"
                      : "#f87171",
                }}
              >
                {ep.method}
              </span>

              <code style={styles.apiPath}>
                {ep.path}
              </code>

              <span style={styles.apiDesc}>
                {ep.desc}
              </span>
            </div>
          ))}
        </div>
      </Section>
    </div>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div style={styles.section}>
      <h2 style={styles.sectionTitle}>{title}</h2>
      {children}
    </div>
  );
}

function StepCard({ number, title, description, code }: {
  number: string; title: string; description: string; code: string;
}) {
  return (
    <div style={styles.stepCard}>
      <div style={styles.stepNumber}>{number}</div>
      <h3 style={styles.stepTitle}>{title}</h3>
      <p style={styles.stepDesc}>{description}</p>
      <pre style={styles.codeBlock}>{code}</pre>
    </div>
  );
}

function FeatureCard({ icon, title, description, tip }: {
  icon: React.ReactNode; title: string; description: string; tip: string;
}) {
  return (
    <div style={styles.featureCard}>
      <div style={styles.featureHeader}>
        {icon}
        <h3 style={styles.featureTitle}>{title}</h3>
      </div>
      <p style={styles.featureDesc}>{description}</p>
      <div style={styles.tipBox}>
        <span style={styles.tipLabel}>💡 Tip: </span>
        <span style={styles.tipText}>{tip}</span>
      </div>
    </div>
  );
}

function UseCaseCard({ emoji, title, description, examples }: {
  emoji: string; title: string; description: string; examples: string[];
}) {
  return (
    <div style={styles.useCaseCard}>
      <div style={styles.useCaseEmoji}>{emoji}</div>
      <h3 style={styles.useCaseTitle}>{title}</h3>
      <p style={styles.useCaseDesc}>{description}</p>
      <ul style={styles.useCaseList}>
        {examples.map((ex, i) => (
          <li key={i} style={styles.useCaseItem}>
            <span style={{ color: '#34d399', marginRight: 6 }}>→</span>
            {ex}
          </li>
        ))}
      </ul>
    </div>
  );
}

const styles: { [key: string]: React.CSSProperties } = {
  pageHeader: { marginBottom: 32 },
  pageTitle: { fontSize: 16, fontWeight: 700, color: '#9ca3af', letterSpacing: 0.5, marginBottom: 8 },
  pageSubtitle: { fontSize: 14, color: '#6b7280', lineHeight: 1.6 },
  section: { marginBottom: 36 },
  sectionTitle: { fontSize: 14, fontWeight: 700, color: '#e5e7eb', marginBottom: 16 },
  stepsGrid: { display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 16 },
  stepCard: { backgroundColor: '#13151f', border: '1px solid #1f2330', borderRadius: 10, padding: 20 },
  stepNumber: {
    width: 28, height: 28, borderRadius: '50%', backgroundColor: '#1d4ed8',
    color: '#fff', fontSize: 13, fontWeight: 700, display: 'flex',
    alignItems: 'center', justifyContent: 'center', marginBottom: 12,
  },
  stepTitle: { fontSize: 13, fontWeight: 700, color: '#e5e7eb', marginBottom: 8 },
  stepDesc: { fontSize: 12, color: '#9ca3af', lineHeight: 1.6, marginBottom: 12 },
  codeBlock: {
    backgroundColor: '#0f1117', border: '1px solid #1f2330', borderRadius: 6,
    padding: '8px 10px', fontSize: 11, color: '#60a5fa',
    whiteSpace: 'pre-wrap' as const, fontFamily: 'monospace', lineHeight: 1.5,
  },
  featuresGrid: { display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 16 },
  featureCard: { backgroundColor: '#13151f', border: '1px solid #1f2330', borderRadius: 10, padding: 20 },
  featureHeader: { display: 'flex', alignItems: 'center', gap: 10, marginBottom: 10 },
  featureTitle: { fontSize: 14, fontWeight: 700, color: '#e5e7eb' },
  featureDesc: { fontSize: 13, color: '#9ca3af', lineHeight: 1.6, marginBottom: 12 },
  tipBox: { backgroundColor: '#0f1117', border: '1px solid #1f3328', borderRadius: 6, padding: '8px 10px', fontSize: 12 },
  tipLabel: { color: '#34d399', fontWeight: 600 },
  tipText: { color: '#6b7280' },
  useCasesGrid: { display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: 16 },
  useCaseCard: { backgroundColor: '#13151f', border: '1px solid #1f2330', borderRadius: 10, padding: 20 },
  useCaseEmoji: { fontSize: 28, marginBottom: 10 },
  useCaseTitle: { fontSize: 14, fontWeight: 700, color: '#e5e7eb', marginBottom: 8 },
  useCaseDesc: { fontSize: 13, color: '#9ca3af', lineHeight: 1.6, marginBottom: 12 },
  useCaseList: { listStyle: 'none', padding: 0, margin: 0 },
  useCaseItem: { fontSize: 12, color: '#6b7280', lineHeight: 1.8, display: 'flex', alignItems: 'flex-start' },
  selfHostIntro: { fontSize: 13, color: '#9ca3af', lineHeight: 1.6, marginBottom: 16 },
  selfHostSteps: { display: 'flex', flexDirection: 'column' as const, gap: 12, marginBottom: 16 },
  selfHostStep: {
    display: 'flex', alignItems: 'center', gap: 12,
    backgroundColor: '#13151f', border: '1px solid #1f2330', borderRadius: 8, padding: '10px 16px',
  },
  selfHostNum: {
    width: 22, height: 22, borderRadius: '50%', backgroundColor: '#1f2330',
    color: '#6b7280', fontSize: 11, fontWeight: 700,
    display: 'flex', alignItems: 'center', justifyContent: 'center', flexShrink: 0,
  },
  selfHostIcon: { flexShrink: 0 },
  selfHostStepTitle: { fontSize: 13, fontWeight: 600, color: '#e5e7eb' },
  selfHostStepDetail: { fontSize: 11, color: '#6b7280', marginTop: 2 },
  githubLink: {
    display: 'inline-flex', alignItems: 'center', gap: 8,
    backgroundColor: '#1d4ed8', color: '#fff', borderRadius: 8,
    padding: '8px 16px', fontSize: 13, fontWeight: 600, textDecoration: 'none',
  },
  apiIntro: { fontSize: 13, color: '#9ca3af', lineHeight: 1.6, marginBottom: 16 },
  apiGrid: { display: 'flex', flexDirection: 'column' as const, gap: 6 },
  apiRow: {
    display: 'flex', alignItems: 'center', gap: 12,
    backgroundColor: '#13151f', border: '1px solid #1f2330', borderRadius: 8, padding: '8px 14px',
  },
  methodBadge: { fontSize: 10, fontWeight: 700, padding: '2px 8px', borderRadius: 4, flexShrink: 0, fontFamily: 'monospace' },
  apiPath: { fontSize: 12, color: '#e5e7eb', fontFamily: 'monospace', minWidth: 260, flexShrink: 0 },
  apiDesc: { fontSize: 12, color: '#6b7280' },
};