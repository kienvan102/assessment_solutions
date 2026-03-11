import { useState, useEffect } from 'react';
import type { Solution } from '../../types/solution';
import { requestSolutionEndpoint } from '../../services/api';

export default function CodeReviewCard({ solution }: { solution: Solution }) {
  const [activeTab, setActiveTab] = useState<'analysis' | 'playground'>('playground');
  const [analysis, setAnalysis] = useState<{ problems: any[], improvements: any[] } | null>(null);
  const [loadingAnalysis, setLoadingAnalysis] = useState(false);

  // Playground state
  const [action, setAction] = useState<'bad' | 'good' | 'simulate'>('bad');
  const [payload, setPayload] = useState('<script>alert(1)</script>');
  const [testResult, setTestResult] = useState<any>(null);
  const [testStatus, setTestStatus] = useState<number | null>(null);
  const [loadingTest, setLoadingTest] = useState(false);

  useEffect(() => {
    if (activeTab === 'analysis' && !analysis) {
      fetchAnalysis();
    }
  }, [activeTab]);

  const fetchAnalysis = async () => {
    setLoadingAnalysis(true);
    try {
      const res = await requestSolutionEndpoint({
        endpoint: solution.endpoint || '/api/codereview1',
        method: 'POST',
        body: JSON.stringify({ action: 'analysis', payload: '' })
      });
      setAnalysis(res.data);
    } catch (err) {
      console.error("Failed to load analysis", err);
    } finally {
      setLoadingAnalysis(false);
    }
  };

  const handleTest = async () => {
    setLoadingTest(true);
    setTestResult(null);
    setTestStatus(null);
    
    try {
      const res = await requestSolutionEndpoint({
        endpoint: solution.endpoint || '/api/codereview1',
        method: 'POST',
        body: JSON.stringify({ action, payload })
      });
      
      setTestStatus(res.status);
      setTestResult(res.data);
    } catch (err: any) {
      setTestStatus(err.response?.status || 500);
      setTestResult(err.message || 'Unknown error occurred');
    } finally {
      setLoadingTest(false);
    }
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

      <div style={{ marginBottom: '1.5rem' }}>
        <strong>Code to Review:</strong>
        <pre style={{
          backgroundColor: '#2d2d2d', color: '#fff', padding: '1rem',
          borderRadius: '4px', overflowX: 'auto', fontSize: '0.9rem', marginTop: '0.5rem'
        }}>
          <code>{solution.code}</code>
        </pre>
      </div>

      {/* Tabs Navigation */}
      <div style={{ display: 'flex', borderBottom: '1px solid #ddd', marginBottom: '1.5rem' }}>
        <button 
          onClick={() => setActiveTab('playground')}
          style={{
            padding: '0.75rem 1.5rem', border: 'none', background: 'none',
            borderBottom: activeTab === 'playground' ? '3px solid #0066cc' : '3px solid transparent',
            fontWeight: activeTab === 'playground' ? 'bold' : 'normal',
            color: activeTab === 'playground' ? '#0066cc' : '#666', cursor: 'pointer'
          }}
        >
          Interactive Playground
        </button>
        <button 
          onClick={() => setActiveTab('analysis')}
          style={{
            padding: '0.75rem 1.5rem', border: 'none', background: 'none',
            borderBottom: activeTab === 'analysis' ? '3px solid #0066cc' : '3px solid transparent',
            fontWeight: activeTab === 'analysis' ? 'bold' : 'normal',
            color: activeTab === 'analysis' ? '#0066cc' : '#666', cursor: 'pointer'
          }}
        >
          Written Analysis
        </button>
      </div>

      {/* Tab Content: Playground */}
      {activeTab === 'playground' && (
        <div>
          <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '1rem', marginBottom: '1rem' }}>
            <div>
              <label style={{ display: 'block', fontWeight: 'bold', marginBottom: '0.5rem' }}>
                Test Scenario
              </label>
              <select 
                value={action} 
                onChange={(e) => setAction(e.target.value as any)}
                style={{ width: '100%', padding: '0.5rem', borderRadius: '4px', border: '1px solid #ccc' }}
              >
                <option value="bad">1. Test Original Flawed Implementation</option>
                <option value="good">2. Test My Fixed Implementation</option>
                <option value="simulate">3. Simulate Race Condition</option>
              </select>
            </div>
          </div>

          {(action === 'bad' || action === 'good') && (
            <div style={{ marginBottom: '1rem' }}>
              <label style={{ display: 'block', fontWeight: 'bold', marginBottom: '0.5rem' }}>
                Payload Data
              </label>
              <textarea 
                value={payload}
                onChange={(e) => setPayload(e.target.value)}
                style={{ width: '100%', height: '80px', padding: '0.5rem', borderRadius: '4px', border: '1px solid #ccc', fontFamily: 'monospace' }}
                placeholder="Enter text to post to the handler..."
              />
              <small style={{ color: '#666', display: 'block', marginTop: '0.25rem' }}>
                {action === 'bad' ? '💡 Try an XSS payload like <script>alert(1)</script> to see the vulnerability.' : '💡 My fixed handler safely encodes this into JSON and sets proper headers.'}
              </small>
            </div>
          )}

          {action === 'simulate' && (
            <div style={{ padding: '1rem', backgroundColor: '#fff3cd', borderRadius: '4px', borderLeft: '4px solid #ffc107', marginBottom: '1rem' }}>
              <strong>What does this do?</strong>
              <p style={{ margin: '0.5rem 0 0 0', fontSize: '0.9rem' }}>
                This sends 100 concurrent HTTP requests to the original flawed handler. 
                Because the handler writes to a global variable without a mutex, requests overwrite each other's data before writing to the response stream.
              </p>
            </div>
          )}

          <button
            onClick={handleTest}
            disabled={loadingTest}
            style={{
              padding: '0.5rem 1.5rem', backgroundColor: '#0066cc', color: 'white',
              border: 'none', borderRadius: '4px', cursor: loadingTest ? 'not-allowed' : 'pointer',
              fontWeight: 'bold'
            }}
          >
            {loadingTest ? 'Running...' : 'Run Scenario'}
          </button>

          {testResult && (
            <div style={{ marginTop: '1.5rem' }}>
              <h4 style={{ margin: '0 0 0.5rem 0' }}>Result (Status: {testStatus})</h4>
              <pre style={{
                backgroundColor: testStatus === 200 ? '#e6ffe6' : '#ffe6e6',
                borderLeft: `4px solid ${testStatus === 200 ? '#00cc00' : '#cc0000'}`,
                padding: '1rem', borderRadius: '4px', overflowX: 'auto', margin: 0,
                whiteSpace: 'pre-wrap'
              }}>
                {typeof testResult === 'object' ? JSON.stringify(testResult, null, 2) : testResult}
              </pre>
            </div>
          )}
        </div>
      )}

      {/* Tab Content: Analysis */}
      {activeTab === 'analysis' && (
        <div>
          {loadingAnalysis ? (
            <p>Loading analysis...</p>
          ) : analysis ? (
            <div>
              <h4 style={{ color: '#cc0000', borderBottom: '1px solid #eee', paddingBottom: '0.5rem' }}>Identified Problems</h4>
              <div style={{ display: 'grid', gap: '1rem', marginBottom: '2rem' }}>
                {analysis.problems.map((p, i) => (
                  <div key={i} style={{ padding: '1rem', backgroundColor: '#fff', border: '1px solid #ddd', borderRadius: '4px', borderLeft: `4px solid ${p.severity === 'Critical' ? '#cc0000' : p.severity === 'High' ? '#ff6600' : '#ffcc00'}` }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.5rem' }}>
                      <strong style={{ fontSize: '1.1rem' }}>{p.title}</strong>
                      <span style={{ fontSize: '0.8rem', padding: '0.2rem 0.5rem', backgroundColor: '#eee', borderRadius: '12px', fontWeight: 'bold' }}>{p.severity}</span>
                    </div>
                    <p style={{ margin: 0, color: '#444' }}>{p.description}</p>
                  </div>
                ))}
              </div>

              <h4 style={{ color: '#008000', borderBottom: '1px solid #eee', paddingBottom: '0.5rem' }}>Suggested Improvements</h4>
              <div style={{ display: 'grid', gap: '1rem' }}>
                {analysis.improvements.map((imp, i) => (
                  <div key={i} style={{ padding: '1rem', backgroundColor: '#f0fdf4', border: '1px solid #bbf7d0', borderRadius: '4px' }}>
                    <strong style={{ display: 'block', marginBottom: '0.5rem', color: '#166534' }}>✓ {imp.title}</strong>
                    <p style={{ margin: 0, color: '#15803d' }}>{imp.description}</p>
                  </div>
                ))}
              </div>
            </div>
          ) : (
            <p>No analysis data available.</p>
          )}
        </div>
      )}
    </div>
  );
}
