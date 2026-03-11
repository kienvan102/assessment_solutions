import { useState } from 'react';
import type { Solution } from '../../types/solution';
import { requestSolutionEndpoint } from '../../services/api';

export default function SqlCard({ solution }: { solution: Solution }) {
  const [query, setQuery] = useState(solution.code || 'SELECT * FROM users;');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<{ columns?: string[], rows?: any[][], error?: string } | null>(null);

  const handleTest = async () => {
    setLoading(true);
    setResult(null);
    try {
      const res = await requestSolutionEndpoint({
        endpoint: solution.endpoint!,
        method: 'POST',
        body: JSON.stringify({ query })
      });
      setResult(res.data);
    } catch (err: any) {
      setResult({ error: err.response?.data?.error || err.message || 'Unknown error occurred' });
    } finally {
      setLoading(false);
    }
  };

  const handleLoadSolution = () => {
    setQuery(`SELECT 
    u.name, 
    COALESCE(SUM(o.amount), 0) as total_amount
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id, u.name;`);
  };

  return (
    <div style={{
      border: '1px solid #e0e0e0',
      borderRadius: '8px',
      padding: '1.5rem',
      marginBottom: '2rem',
      backgroundColor: '#ffffff',
      boxShadow: '0 2px 4px rgba(0,0,0,0.05)'
    }}>
      <h3 style={{ marginTop: 0, borderBottom: '1px solid #eee', paddingBottom: '0.5rem' }}>
        {solution.title}
      </h3>

      <div style={{ marginBottom: '1.5rem' }}>
        <strong>Description:</strong>
        <p style={{ margin: '0.5rem 0 0 0' }}>{solution.description}</p>
      </div>

      <div style={{ marginBottom: '1.5rem', padding: '1rem', backgroundColor: '#f8f9fa', borderRadius: '4px', border: '1px solid #e2e8f0' }}>
        <h4 style={{ margin: '0 0 1rem 0' }}>Database Schema (SQLite)</h4>
        
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '2rem' }}>
          {/* Users Table */}
          <div>
            <strong style={{ display: 'block', marginBottom: '0.5rem', color: '#333' }}>users</strong>
            <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.9rem' }}>
              <thead>
                <tr style={{ backgroundColor: '#edf2f7' }}>
                  <th style={{ border: '1px solid #cbd5e1', padding: '0.5rem', textAlign: 'left' }}>id (PK)</th>
                  <th style={{ border: '1px solid #cbd5e1', padding: '0.5rem', textAlign: 'left' }}>name</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>101</td>
                  <td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>Alice</td>
                </tr>
                <tr>
                  <td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>102</td>
                  <td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>Bob</td>
                </tr>
                <tr>
                  <td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>103</td>
                  <td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>Charlie</td>
                </tr>
              </tbody>
            </table>
          </div>

          {/* Orders Table */}
          <div>
            <strong style={{ display: 'block', marginBottom: '0.5rem', color: '#333' }}>orders</strong>
            <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.9rem' }}>
              <thead>
                <tr style={{ backgroundColor: '#edf2f7' }}>
                  <th style={{ border: '1px solid #cbd5e1', padding: '0.5rem', textAlign: 'left' }}>id (PK)</th>
                  <th style={{ border: '1px solid #cbd5e1', padding: '0.5rem', textAlign: 'left' }}>user_id (FK)</th>
                  <th style={{ border: '1px solid #cbd5e1', padding: '0.5rem', textAlign: 'left' }}>amount</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>1</td>
                  <td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>101</td>
                  <td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>50.00</td>
                </tr>
                <tr>
                  <td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>2</td>
                  <td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>101</td>
                  <td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>75.00</td>
                </tr>
                <tr>
                  <td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>3</td>
                  <td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>102</td>
                  <td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>30.00</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <div style={{ marginBottom: '1.5rem' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem' }}>
          <strong style={{ display: 'block' }}>SQL Query Editor</strong>
          <button 
            onClick={handleLoadSolution}
            style={{ padding: '0.25rem 0.75rem', fontSize: '0.85rem', backgroundColor: '#e2e8f0', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
          >
            Load My Solution
          </button>
        </div>
        <textarea
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          style={{
            width: '100%',
            height: '150px',
            fontFamily: 'monospace',
            fontSize: '1rem',
            padding: '1rem',
            backgroundColor: '#1e1e1e',
            color: '#d4d4d4',
            borderRadius: '6px',
            border: 'none',
            boxSizing: 'border-box'
          }}
          spellCheck="false"
        />
        <div style={{ marginTop: '1rem' }}>
          <button
            onClick={handleTest}
            disabled={loading || !query.trim()}
            style={{
              padding: '0.5rem 1.5rem',
              backgroundColor: '#0066cc',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: loading || !query.trim() ? 'not-allowed' : 'pointer',
              fontWeight: 'bold'
            }}
          >
            {loading ? 'Running Query...' : 'Run Query'}
          </button>
        </div>
      </div>

      {result && (
        <div style={{ marginTop: '2rem', borderTop: '1px solid #eee', paddingTop: '1.5rem' }}>
          <h4 style={{ margin: '0 0 1rem 0' }}>Query Results</h4>
          
          {result.error ? (
            <div style={{ padding: '1rem', backgroundColor: '#fee2e2', color: '#991b1b', borderRadius: '4px', borderLeft: '4px solid #ef4444' }}>
              <strong>Error: </strong> {result.error}
            </div>
          ) : result.columns && result.rows ? (
            result.rows.length === 0 ? (
              <div style={{ padding: '1rem', backgroundColor: '#f1f5f9', color: '#475569', borderRadius: '4px' }}>
                No rows returned.
              </div>
            ) : (
              <div style={{ overflowX: 'auto' }}>
                <table style={{ minWidth: '100%', borderCollapse: 'collapse', border: '1px solid #e2e8f0' }}>
                  <thead>
                    <tr style={{ backgroundColor: '#f8fafc' }}>
                      {result.columns.map((col, i) => (
                        <th key={i} style={{ borderBottom: '2px solid #cbd5e1', padding: '0.75rem', textAlign: 'left', fontWeight: 'bold', color: '#334155' }}>
                          {col}
                        </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {result.rows.map((row, i) => (
                      <tr key={i} style={{ borderBottom: '1px solid #e2e8f0', backgroundColor: i % 2 === 0 ? '#ffffff' : '#f8fafc' }}>
                        {row.map((cell, j) => (
                          <td key={j} style={{ padding: '0.75rem', color: '#334155' }}>
                            {cell === null ? <span style={{ color: '#94a3b8', fontStyle: 'italic' }}>NULL</span> : String(cell)}
                          </td>
                        ))}
                      </tr>
                    ))}
                  </tbody>
                </table>
                <div style={{ marginTop: '0.5rem', fontSize: '0.85rem', color: '#64748b', textAlign: 'right' }}>
                  {result.rows.length} row(s) returned
                </div>
              </div>
            )
          ) : null}
        </div>
      )}
    </div>
  );
}
