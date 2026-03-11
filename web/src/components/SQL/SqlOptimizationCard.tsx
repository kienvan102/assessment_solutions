import { useState } from 'react';
import type { Solution } from '../../types/solution';
import { requestSolutionEndpoint } from '../../services/api';

export default function SqlOptimizationCard({ solution }: { solution: Solution }) {
  const [activeTab, setActiveTab] = useState<'query' | 'explain' | 'index'>('query');
  
  // States for query execution
  const [query, setQuery] = useState("SELECT * FROM transactions \nWHERE user_id = 123 \nORDER BY created_at DESC \nLIMIT 10;");
  const [loadingQuery, setLoadingQuery] = useState(false);
  const [queryResult, setQueryResult] = useState<{ columns?: string[], rows?: any[][], error?: string } | null>(null);

  // States for Explain Plan
  const [explainResult, setExplainResult] = useState<{ explain?: string, error?: string } | null>(null);
  const [loadingExplain, setLoadingExplain] = useState(false);

  // States for Index Creation
  const [indexQuery, setIndexQuery] = useState("CREATE INDEX idx_user_created ON transactions(user_id, created_at);");
  const [loadingIndex, setLoadingIndex] = useState(false);
  const [indexResult, setIndexResult] = useState<{ message?: string, error?: string } | null>(null);

  const [loadingReset, setLoadingReset] = useState(false);

  const handleRunQuery = async () => {
    setLoadingQuery(true);
    setQueryResult(null);
    try {
      const res = await requestSolutionEndpoint({
        endpoint: solution.endpoint!,
        method: 'POST',
        body: JSON.stringify({ action: 'query', query })
      });
      setQueryResult(res.data);
    } catch (err: any) {
      setQueryResult({ error: err.response?.data?.error || err.message });
    } finally {
      setLoadingQuery(false);
    }
  };

  const handleExplainPlan = async () => {
    setLoadingExplain(true);
    setExplainResult(null);
    try {
      const res = await requestSolutionEndpoint({
        endpoint: solution.endpoint!,
        method: 'POST',
        body: JSON.stringify({ action: 'explain', query })
      });
      setExplainResult(res.data);
    } catch (err: any) {
      setExplainResult({ error: err.response?.data?.error || err.message });
    } finally {
      setLoadingExplain(false);
    }
  };

  const handleCreateIndex = async () => {
    setLoadingIndex(true);
    setIndexResult(null);
    try {
      const res = await requestSolutionEndpoint({
        endpoint: solution.endpoint!,
        method: 'POST',
        body: JSON.stringify({ action: 'exec', query: indexQuery })
      });
      setIndexResult({ message: res.data.message });
    } catch (err: any) {
      setIndexResult({ error: err.response?.data?.error || err.message });
    } finally {
      setLoadingIndex(false);
    }
  };

  const handleResetDatabase = async () => {
    setLoadingReset(true);
    try {
      await requestSolutionEndpoint({
        endpoint: solution.endpoint!,
        method: 'POST',
        body: JSON.stringify({ action: 'reset', query: '' })
      });
      setIndexResult(null);
      setExplainResult(null);
      setQueryResult(null);
      alert("Database reset successfully! All indexes removed.");
    } catch (err: any) {
      alert("Failed to reset database: " + (err.response?.data?.error || err.message));
    } finally {
      setLoadingReset(false);
    }
  };

  return (
    <div style={{
      border: '1px solid #e0e0e0', borderRadius: '8px', padding: '1.5rem',
      marginBottom: '2rem', backgroundColor: '#ffffff', boxShadow: '0 2px 4px rgba(0,0,0,0.05)'
    }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '1px solid #eee', paddingBottom: '0.5rem', marginBottom: '1.5rem' }}>
        <h3 style={{ margin: 0 }}>{solution.title}</h3>
        <button 
          onClick={handleResetDatabase}
          disabled={loadingReset}
          style={{ padding: '0.4rem 0.8rem', backgroundColor: '#fef2f2', color: '#b91c1c', border: '1px solid #fca5a5', borderRadius: '4px', cursor: 'pointer', fontSize: '0.85rem' }}
        >
          {loadingReset ? 'Resetting...' : 'Reset Database (Drop Indexes)'}
        </button>
      </div>

      <div style={{ marginBottom: '1.5rem' }}>
        <strong>Description:</strong>
        <p style={{ margin: '0.5rem 0 0 0' }}>{solution.description}</p>
      </div>

      {/* Testing Guide / Instructions */}
      <div style={{ marginBottom: '1.5rem', padding: '1rem', backgroundColor: '#e0f2fe', borderRadius: '4px', borderLeft: '4px solid #0ea5e9' }}>
        <h4 style={{ margin: '0 0 0.5rem 0', color: '#0369a1' }}>Testing Guide</h4>
        <ol style={{ margin: 0, paddingLeft: '1.5rem', color: '#0c4a6e', fontSize: '0.95rem', lineHeight: '1.6' }}>
          <li>Go to <strong>Tab 2 (Explain Query Plan)</strong> first and click "Analyze". You will see it performs a slow <code style={{ backgroundColor: '#bae6fd', padding: '2px 4px', borderRadius: '3px' }}>SCAN TABLE</code> (Full Table Scan).</li>
          <li>Go to <strong>Tab 3 (Apply Optimization Index)</strong> and execute the proposed <code>CREATE INDEX</code> statement.</li>
          <li>Go back to <strong>Tab 2 (Explain Query Plan)</strong> and analyze it again. You will see the index successfully upgraded the plan to a highly efficient <code style={{ backgroundColor: '#bae6fd', padding: '2px 4px', borderRadius: '3px' }}>SEARCH TABLE ... USING INDEX</code>.</li>
        </ol>
      </div>

      {/* Database Schema Visualizer */}
      <div style={{ marginBottom: '1.5rem', padding: '1rem', backgroundColor: '#f8f9fa', borderRadius: '4px', border: '1px solid #e2e8f0' }}>
        <h4 style={{ margin: '0 0 0.5rem 0' }}>Database Schema</h4>
        <p style={{ margin: '0 0 1rem 0', fontSize: '0.9rem', color: '#64748b' }}>
          Table: <strong>transactions</strong> (Seeded with 10,000 randomized rows to simulate volume)
        </p>
        <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.9rem', maxWidth: '400px' }}>
          <thead>
            <tr style={{ backgroundColor: '#edf2f7' }}>
              <th style={{ border: '1px solid #cbd5e1', padding: '0.5rem', textAlign: 'left' }}>Column</th>
              <th style={{ border: '1px solid #cbd5e1', padding: '0.5rem', textAlign: 'left' }}>Type</th>
            </tr>
          </thead>
          <tbody>
            <tr><td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>id</td><td style={{ border: '1px solid #cbd5e1', padding: '0.5rem', color: '#64748b' }}>INTEGER PRIMARY KEY AUTOINCREMENT</td></tr>
            <tr><td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>user_id</td><td style={{ border: '1px solid #cbd5e1', padding: '0.5rem', color: '#64748b' }}>INTEGER</td></tr>
            <tr><td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>amount</td><td style={{ border: '1px solid #cbd5e1', padding: '0.5rem', color: '#64748b' }}>DECIMAL(10, 2)</td></tr>
            <tr><td style={{ border: '1px solid #cbd5e1', padding: '0.5rem' }}>created_at</td><td style={{ border: '1px solid #cbd5e1', padding: '0.5rem', color: '#64748b' }}>TIMESTAMP</td></tr>
          </tbody>
        </table>
      </div>

      {/* Tabs Navigation */}
      <div style={{ display: 'flex', borderBottom: '1px solid #ddd', marginBottom: '1.5rem', gap: '1rem' }}>
        <button 
          onClick={() => setActiveTab('query')}
          style={{
            padding: '0.75rem 1rem', border: 'none', background: 'none',
            borderBottom: activeTab === 'query' ? '3px solid #0066cc' : '3px solid transparent',
            fontWeight: activeTab === 'query' ? 'bold' : 'normal',
            color: activeTab === 'query' ? '#0066cc' : '#666', cursor: 'pointer'
          }}
        >
          1. Query Execution
        </button>
        <button 
          onClick={() => setActiveTab('explain')}
          style={{
            padding: '0.75rem 1rem', border: 'none', background: 'none',
            borderBottom: activeTab === 'explain' ? '3px solid #0066cc' : '3px solid transparent',
            fontWeight: activeTab === 'explain' ? 'bold' : 'normal',
            color: activeTab === 'explain' ? '#0066cc' : '#666', cursor: 'pointer'
          }}
        >
          2. Explain Query Plan
        </button>
        <button 
          onClick={() => setActiveTab('index')}
          style={{
            padding: '0.75rem 1rem', border: 'none', background: 'none',
            borderBottom: activeTab === 'index' ? '3px solid #0066cc' : '3px solid transparent',
            fontWeight: activeTab === 'index' ? 'bold' : 'normal',
            color: activeTab === 'index' ? '#0066cc' : '#666', cursor: 'pointer'
          }}
        >
          3. Apply Optimization Index
        </button>
      </div>

      {/* Tab: Query Execution */}
      <div style={{ display: activeTab === 'query' ? 'block' : 'none' }}>
        <label style={{ display: 'block', fontWeight: 'bold', marginBottom: '0.5rem' }}>Write Data Query:</label>
        <textarea
          value={query} onChange={(e) => setQuery(e.target.value)}
          style={{ width: '100%', height: '120px', fontFamily: 'monospace', padding: '1rem', backgroundColor: '#1e1e1e', color: '#d4d4d4', borderRadius: '6px', border: 'none', boxSizing: 'border-box' }}
          spellCheck="false"
        />
        <button
          onClick={handleRunQuery} disabled={loadingQuery}
          style={{ marginTop: '1rem', padding: '0.5rem 1.5rem', backgroundColor: '#0066cc', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer', fontWeight: 'bold' }}
        >
          {loadingQuery ? 'Running...' : 'Run Query'}
        </button>

        {queryResult && (
          <div style={{ marginTop: '1.5rem' }}>
            {queryResult.error ? (
              <div style={{ padding: '1rem', backgroundColor: '#fee2e2', color: '#991b1b', borderRadius: '4px', borderLeft: '4px solid #ef4444' }}>
                <strong>Error: </strong> {queryResult.error}
              </div>
            ) : queryResult.columns && queryResult.rows && (
              <div style={{ overflowX: 'auto', border: '1px solid #e2e8f0', borderRadius: '4px' }}>
                <table style={{ minWidth: '100%', borderCollapse: 'collapse' }}>
                  <thead>
                    <tr style={{ backgroundColor: '#f8fafc' }}>
                      {queryResult.columns.map((col, i) => (
                        <th key={i} style={{ borderBottom: '2px solid #cbd5e1', padding: '0.75rem', textAlign: 'left', fontWeight: 'bold' }}>{col}</th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {queryResult.rows.map((row, i) => (
                      <tr key={i} style={{ borderBottom: '1px solid #e2e8f0', backgroundColor: i % 2 === 0 ? '#ffffff' : '#f8fafc' }}>
                        {row.map((cell, j) => (
                          <td key={j} style={{ padding: '0.75rem' }}>{cell === null ? 'NULL' : String(cell)}</td>
                        ))}
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}
      </div>

      {/* Tab: Explain Plan */}
      <div style={{ display: activeTab === 'explain' ? 'block' : 'none' }}>
        <div style={{ padding: '1rem', backgroundColor: '#e0f2fe', borderRadius: '4px', borderLeft: '4px solid #0ea5e9', marginBottom: '1.5rem' }}>
          <strong>What is an Explain Plan?</strong>
          <p style={{ margin: '0.5rem 0 0 0', fontSize: '0.9rem' }}>
            It shows how the database engine executes your query. Before adding an index, you will likely see <code>SCAN TABLE transactions</code> (Full Table Scan). After adding the correct index, it should change to <code>SEARCH TABLE transactions USING INDEX</code>.
          </p>
        </div>

        <button
          onClick={handleExplainPlan} disabled={loadingExplain}
          style={{ padding: '0.5rem 1.5rem', backgroundColor: '#0ea5e9', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer', fontWeight: 'bold' }}
        >
          {loadingExplain ? 'Analyzing...' : 'Explain Current Query'}
        </button>

        {explainResult && (
          <div style={{ marginTop: '1.5rem' }}>
            {explainResult.error ? (
              <div style={{ padding: '1rem', backgroundColor: '#fee2e2', color: '#991b1b', borderRadius: '4px', borderLeft: '4px solid #ef4444' }}>
                <strong>Error: </strong> {explainResult.error}
              </div>
            ) : (
              <pre style={{ backgroundColor: '#2d2d2d', color: explainResult.explain?.includes('USING INDEX') ? '#4ade80' : '#f87171', padding: '1.5rem', borderRadius: '6px', whiteSpace: 'pre-wrap', fontWeight: 'bold' }}>
                {explainResult.explain || 'No execution plan generated.'}
              </pre>
            )}
          </div>
        )}
      </div>

      {/* Tab: Apply Index */}
      <div style={{ display: activeTab === 'index' ? 'block' : 'none' }}>
        <label style={{ display: 'block', fontWeight: 'bold', marginBottom: '0.5rem' }}>Write CREATE INDEX Statement:</label>
        <textarea
          value={indexQuery} onChange={(e) => setIndexQuery(e.target.value)}
          style={{ width: '100%', height: '80px', fontFamily: 'monospace', padding: '1rem', backgroundColor: '#1e1e1e', color: '#d4d4d4', borderRadius: '6px', border: 'none', boxSizing: 'border-box' }}
          spellCheck="false"
        />
        <button
          onClick={handleCreateIndex} disabled={loadingIndex}
          style={{ marginTop: '1rem', padding: '0.5rem 1.5rem', backgroundColor: '#10b981', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer', fontWeight: 'bold' }}
        >
          {loadingIndex ? 'Creating...' : 'Execute Index Creation'}
        </button>

        {indexResult && (
          <div style={{ marginTop: '1.5rem' }}>
            {indexResult.error ? (
              <div style={{ padding: '1rem', backgroundColor: '#fee2e2', color: '#991b1b', borderRadius: '4px', borderLeft: '4px solid #ef4444' }}>
                <strong>Error: </strong> {indexResult.error}
              </div>
            ) : (
              <div style={{ padding: '1rem', backgroundColor: '#dcfce3', color: '#166534', borderRadius: '4px', borderLeft: '4px solid #22c55e' }}>
                <strong>Success: </strong> {indexResult.message}
                <p style={{ margin: '0.5rem 0 0 0', fontSize: '0.9rem' }}>
                  Now go back to the <strong>Explain Query Plan</strong> tab to see if your query is using the new index!
                </p>
              </div>
            )}
          </div>
        )}
      </div>

    </div>
  );
}
