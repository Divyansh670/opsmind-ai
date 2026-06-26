import { useState, useRef, useEffect } from 'react';
import { Send, Bot, User, X, AlertTriangle, Building2, DollarSign, Maximize2, Minimize2 } from 'lucide-react';
import type { RAGSource } from '../api/client';

interface Message {
  id: number;
  role: 'user' | 'assistant';
  content: string;
  sources?: RAGSource[];
  streaming?: boolean;
}

interface ChatPanelProps {
  isOpen: boolean;
  onClose: () => void;
}

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';

type ChatSize = 'normal' | 'minimized' | 'maximized';

export default function ChatPanel({ isOpen, onClose }: ChatPanelProps) {
  const [messages, setMessages] = useState<Message[]>([
    {
      id: 0,
      role: 'assistant',
      content: "Hi! I'm OpsMind AI Assistant. Ask me anything about your security findings, cost drift, or architecture violations across your pull requests.",
    },
  ]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [size, setSize] = useState<ChatSize>('normal');
  const bottomRef = useRef<HTMLDivElement>(null);
  const msgId = useRef(1);

  useEffect(() => {
    if (size !== 'minimized') {
      bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages, size]);

  async function handleSend() {
    if (!input.trim() || loading) return;
    const question = input.trim();
    setInput('');
    setLoading(true);

    const userMsg: Message = { id: msgId.current++, role: 'user', content: question };
    const assistantId = msgId.current++;
    const assistantMsg: Message = { id: assistantId, role: 'assistant', content: '', sources: [], streaming: true };

    setMessages(prev => [...prev, userMsg, assistantMsg]);

    try {
      const response = await fetch(`${API_BASE}/api/chat/stream`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ question }),
      });

      const reader = response.body?.getReader();
      const decoder = new TextDecoder();
      if (!reader) throw new Error('No response body');

      let buffer = '';
      while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split('\n');
        buffer = lines.pop() ?? '';

        let eventType = '';
        for (const line of lines) {
          if (line.startsWith('event: ')) {
            eventType = line.slice(7).trim();
          } else if (line.startsWith('data: ')) {
            const data = line.slice(6).trim();
            if (eventType === 'sources') {
              try {
                const sources: RAGSource[] = JSON.parse(data);
                setMessages(prev => prev.map(m => m.id === assistantId ? { ...m, sources } : m));
              } catch {}
            } else if (eventType === 'token') {
              try {
                const token: string = JSON.parse(data);
                setMessages(prev => prev.map(m => m.id === assistantId ? { ...m, content: m.content + token } : m));
              } catch {}
            } else if (eventType === 'done') {
              setMessages(prev => prev.map(m => m.id === assistantId ? { ...m, streaming: false } : m));
            }
            eventType = '';
          }
        }
      }
    } catch {
      setMessages(prev => prev.map(m =>
        m.id === assistantId ? { ...m, content: 'Sorry, something went wrong. Please try again.', streaming: false } : m
      ));
    } finally {
      setLoading(false);
    }
  }

  if (!isOpen) return null;

  const getOuterStyle = (): React.CSSProperties => {
    if (size === 'maximized') {
      return {
        position: 'fixed',
        top: 16,
        left: 16,
        right: 16,
        bottom: 16,
        zIndex: 1000,
        display: 'flex',
        flexDirection: 'column',
        backgroundColor: '#13151f',
        border: '1px solid #1f2330',
        borderRadius: 12,
        boxShadow: '0 20px 60px rgba(0,0,0,0.6)',
      };
    }
    if (size === 'minimized') {
      return {
        position: 'fixed',
        bottom: 24,
        right: 24,
        width: 300,
        zIndex: 1000,
        display: 'flex',
        flexDirection: 'column',
        backgroundColor: '#13151f',
        border: '1px solid #1f2330',
        borderRadius: 12,
        boxShadow: '0 8px 24px rgba(0,0,0,0.4)',
      };
    }
    return {
      position: 'fixed',
      bottom: 24,
      right: 24,
      width: 420,
      height: 560,
      zIndex: 1000,
      display: 'flex',
      flexDirection: 'column',
      backgroundColor: '#13151f',
      border: '1px solid #1f2330',
      borderRadius: 12,
      boxShadow: '0 20px 60px rgba(0,0,0,0.5)',
    };
  };

  return (
    <div style={getOuterStyle()}>

      {/* Header — always visible */}
      <div style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        padding: '12px 16px',
        borderBottom: size !== 'minimized' ? '1px solid #1f2330' : 'none',
        backgroundColor: '#0f1117',
        borderRadius: size === 'minimized' ? 12 : '12px 12px 0 0',
        flexShrink: 0,
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <Bot size={18} color="#60a5fa" />
          <span style={{ fontSize: 14, fontWeight: 700, color: '#e5e7eb' }}>Ask OpsMind AI</span>
          {size === 'normal' && (
            <span style={{ fontSize: 10, color: '#4b5563' }}>RAG • pgvector</span>
          )}
        </div>
        <div style={{ display: 'flex', gap: 2, alignItems: 'center' }}>
          {/* Minimize / Restore */}
          <button
            onClick={() => setSize(s => s === 'minimized' ? 'normal' : 'minimized')}
            title={size === 'minimized' ? 'Restore' : 'Minimize'}
            style={btnStyle}
          >
            <Minimize2 size={14} />
          </button>
          {/* Maximize / Restore — only when not minimized */}
          {size !== 'minimized' && (
            <button
              onClick={() => setSize(s => s === 'maximized' ? 'normal' : 'maximized')}
              title={size === 'maximized' ? 'Restore' : 'Maximize'}
              style={btnStyle}
            >
              <Maximize2 size={14} />
            </button>
          )}
          {/* Close */}
          <button
            onClick={onClose}
            title="Close"
            style={{ ...btnStyle, color: '#f87171' }}
          >
            <X size={14} />
          </button>
        </div>
      </div>

      {/* Body — hidden when minimized */}
      {size !== 'minimized' && (
        <>
          {/* Messages */}
          <div style={{
            flex: 1,
            overflowY: 'auto',
            padding: 16,
            display: 'flex',
            flexDirection: 'column',
            gap: 16,
          }}>
            {messages.map(msg => (
              <div key={msg.id} style={msg.role === 'user' ? styles.userRow : styles.assistantRow}>
                <div style={msg.role === 'user' ? styles.userAvatar : styles.assistantAvatar}>
                  {msg.role === 'user' ? <User size={14} /> : <Bot size={14} />}
                </div>
                <div style={{ maxWidth: '80%' }}>
                  <div style={msg.role === 'user' ? styles.userBubble : styles.assistantBubble}>
                    <pre style={styles.messageText}>{msg.content}</pre>
                    {msg.streaming && <span style={{ color: '#60a5fa' }}>▋</span>}
                  </div>
                  {msg.sources && msg.sources.length > 0 && !msg.streaming && (
                    <div style={{ marginTop: 8 }}>
                      <p style={{ fontSize: 10, color: '#4b5563', marginBottom: 4 }}>Sources used:</p>
                      {msg.sources.map((s, i) => (
                        <div key={i} style={styles.sourceChip}>
                          {getSourceIcon(s.type, s.severity)}
                          <span style={styles.sourceText}>
                            {s.type === 'finding'
                              ? `PR #${s.pr_number} • ${s.repo_name}${s.file_path ? ` • ${s.file_path}` : ''}`
                              : `Rule: ${s.snippet}`}
                          </span>
                          {s.severity && (
                            <span style={{ fontSize: 10, fontWeight: 700, color: getSeverityColor(s.severity) }}>
                              {s.severity}
                            </span>
                          )}
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            ))}
            <div ref={bottomRef} />
          </div>

          {/* Input */}
          <div style={{
            display: 'flex',
            gap: 8,
            padding: '12px 16px',
            borderTop: '1px solid #1f2330',
            backgroundColor: '#0f1117',
            borderRadius: '0 0 12px 12px',
            flexShrink: 0,
          }}>
            <input
              style={styles.input}
              value={input}
              onChange={e => setInput(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && !e.shiftKey && handleSend()}
              placeholder="Ask about findings, cost drift, architecture violations..."
              disabled={loading}
            />
            <button
              style={{
                ...styles.sendBtn,
                opacity: loading || !input.trim() ? 0.5 : 1,
                cursor: loading || !input.trim() ? 'not-allowed' : 'pointer',
              }}
              onClick={handleSend}
              disabled={loading || !input.trim()}
            >
              <Send size={15} />
            </button>
          </div>
        </>
      )}
    </div>
  );
}

const btnStyle: React.CSSProperties = {
  background: 'none',
  border: 'none',
  color: '#6b7280',
  cursor: 'pointer',
  padding: 6,
  borderRadius: 4,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
};

function getSourceIcon(type: string, severity: string) {
  if (type === 'rule') return <Building2 size={12} color="#60a5fa" />;
  if (severity === 'CRITICAL' || severity === 'HIGH') return <AlertTriangle size={12} color="#f87171" />;
  return <DollarSign size={12} color="#fbbf24" />;
}

function getSeverityColor(severity: string) {
  switch (severity) {
    case 'CRITICAL': return '#f87171';
    case 'HIGH': return '#fb923c';
    case 'MEDIUM': return '#fbbf24';
    case 'LOW': return '#34d399';
    default: return '#6b7280';
  }
}

const styles: { [key: string]: React.CSSProperties } = {
  userRow: { display: 'flex', flexDirection: 'row-reverse', gap: 8, alignItems: 'flex-start' },
  assistantRow: { display: 'flex', gap: 8, alignItems: 'flex-start' },
  userAvatar: {
    width: 28, height: 28, borderRadius: '50%', backgroundColor: '#1d4ed8',
    display: 'flex', alignItems: 'center', justifyContent: 'center', flexShrink: 0, color: '#fff',
  },
  assistantAvatar: {
    width: 28, height: 28, borderRadius: '50%', backgroundColor: '#1f2330',
    display: 'flex', alignItems: 'center', justifyContent: 'center', flexShrink: 0, color: '#60a5fa',
  },
  userBubble: { backgroundColor: '#1d4ed8', borderRadius: '12px 12px 2px 12px', padding: '8px 12px' },
  assistantBubble: {
    backgroundColor: '#1a1d29', border: '1px solid #1f2330',
    borderRadius: '12px 12px 12px 2px', padding: '8px 12px',
  },
  messageText: {
    fontSize: 13, color: '#e5e7eb', lineHeight: 1.6,
    whiteSpace: 'pre-wrap', fontFamily: 'inherit', margin: 0,
  },
  sourceChip: {
    display: 'flex', alignItems: 'center', gap: 6, backgroundColor: '#0f1117',
    border: '1px solid #1f2330', borderRadius: 6, padding: '4px 8px', marginBottom: 4,
  },
  sourceText: { fontSize: 11, color: '#6b7280', flex: 1 },
  input: {
    flex: 1, backgroundColor: '#1a1d29', border: '1px solid #1f2330',
    borderRadius: 8, padding: '8px 12px', color: '#e5e7eb',
    fontSize: 13, outline: 'none', fontFamily: 'inherit',
  },
  sendBtn: {
    backgroundColor: '#1d4ed8', border: 'none', borderRadius: 8,
    color: '#fff', padding: '8px 12px', display: 'flex',
    alignItems: 'center', justifyContent: 'center',
  },
};