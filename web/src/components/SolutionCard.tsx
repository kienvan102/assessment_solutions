import { useState } from 'react';
import { requestSolutionEndpoint } from '../services/api';
import type { RunRecord, Solution } from '../types/solution';

export default function SolutionCard({ solution }: { solution: Solution }) {
  const [advancedMode, setAdvancedMode] = useState(false);
  const [payload, setPayload] = useState(solution.sample_payload || '');
  const [formValues, setFormValues] = useState<Record<string, string>>(() => {
    const init: Record<string, string> = {};
    (solution.input_fields || []).forEach(f => {
      init[f.name] = f.default ?? '';
    });
    return init;
  });

  const [history, setHistory] = useState<RunRecord[]>([]);
  const [selectedRunIndex, setSelectedRunIndex] = useState<number | null>(null);

  const selectedRun = selectedRunIndex === null ? null : history[selectedRunIndex];
  const [loading, setLoading] = useState(false);

  const buildRequestBody = (): string | undefined => {
    if ((solution.method || 'GET') === 'GET') return undefined;

    if (advancedMode) {
      return payload;
    }

    const fields = solution.input_fields || [];
    if (fields.length === 0) {
      return payload;
    }

    const obj: Record<string, any> = {};
    for (const f of fields) {
      const raw = formValues[f.name] ?? '';
      if (f.required && raw === '') {
        throw new Error(`Field '${f.label || f.name}' is required`);
      }
      if (raw === '') continue;

      if (f.type === 'number') {
        const n = Number(raw);
        if (Number.isNaN(n)) {
          throw new Error(`Field '${f.label || f.name}' must be a valid number`);
        }
        obj[f.name] = n;
      } else if (f.type === 'boolean') {
        obj[f.name] = raw === 'true';
      } else {
        obj[f.name] = raw;
      }
    }
    return JSON.stringify(obj, null, 2);
  };

  const handleTest = async () => {
    if (!solution.endpoint) return;

    setLoading(true);
    setSelectedRunIndex(null);
    try {
      const body = buildRequestBody();
      const method = solution.method || 'GET';
      const res = await requestSolutionEndpoint({
        endpoint: solution.endpoint,
        method,
        body,
      });

      const record: RunRecord = {
        at: new Date().toISOString(),
        request: {
          endpoint: solution.endpoint,
          method,
          body: method !== 'GET' && body ? String(body) : undefined,
        },
        response: { status: res.status, data: res.data },
      };
      setHistory(prev => [record, ...prev]);
      setSelectedRunIndex(0);
    } catch (err: any) {
      const record: RunRecord = {
        at: new Date().toISOString(),
        request: {
          endpoint: solution.endpoint,
          method: solution.method || 'GET',
          body: undefined,
        },
        response: { status: 0, data: err.message },
      };
      setHistory(prev => [record, ...prev]);
      setSelectedRunIndex(0);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      style={{
        border: '1px solid #e0e0e0',
        borderRadius: '8px',
        padding: '1.5rem',
        marginBottom: '2rem',
        backgroundColor: '#ffffff',
        boxShadow: '0 2px 4px rgba(0,0,0,0.05)'
      }}
    >
      <h3 style={{ marginTop: 0, borderBottom: '1px solid #eee', paddingBottom: '0.5rem' }}>
        {solution.title}
      </h3>

      <div style={{ marginBottom: '1.5rem' }}>
        <strong>Description:</strong>
        <p style={{ margin: '0.5rem 0 0 0' }}>{solution.description}</p>
      </div>

      {solution.expected_behaviors && solution.expected_behaviors.length > 0 && (
        <div style={{ marginBottom: '1.5rem', padding: '1rem', backgroundColor: '#f0f7ff', borderRadius: '4px' }}>
          <strong>Expected Behaviors:</strong>
          <ul style={{ margin: '0.5rem 0 0 0', paddingLeft: '1.5rem' }}>
            {solution.expected_behaviors.map((behavior, idx) => (
              <li key={idx} style={{ marginBottom: '0.25rem' }}>{behavior}</li>
            ))}
          </ul>
        </div>
      )}

      <div style={{ marginBottom: '1.5rem' }}>
        <strong>Core Implementation:</strong>
        <pre
          style={{
            backgroundColor: '#2d2d2d',
            color: '#fff',
            padding: '1rem',
            borderRadius: '4px',
            overflowX: 'auto',
            fontSize: '0.9rem',
            marginTop: '0.5rem'
          }}
        >
          <code>{solution.code}</code>
        </pre>
      </div>

      {solution.endpoint && (
        <div style={{ borderTop: '1px solid #eee', paddingTop: '1.5rem' }}>
          <strong style={{ display: 'block', marginBottom: '0.5rem' }}>
            Interactive Test ({solution.method} {solution.endpoint})
          </strong>

          {solution.method !== 'GET' && (
            <div style={{ marginBottom: '1rem' }}>
              <label style={{ display: 'inline-flex', alignItems: 'center', gap: '0.5rem', marginBottom: '0.75rem' }}>
                <input
                  type="checkbox"
                  checked={advancedMode}
                  onChange={(e) => setAdvancedMode(e.target.checked)}
                />
                Advanced JSON mode
              </label>

              {!advancedMode && (solution.input_fields || []).length > 0 && (
                <div
                  style={{
                    display: 'grid',
                    gridTemplateColumns: 'repeat(auto-fit, minmax(220px, 1fr))',
                    gap: '0.75rem',
                    marginBottom: '0.75rem'
                  }}
                >
                  {(solution.input_fields || []).map((field) => (
                    <div key={field.name}>
                      <label style={{ display: 'block', fontWeight: 600, marginBottom: '0.25rem' }}>
                        {field.label}{field.required ? ' *' : ''}
                      </label>
                      <input
                        type={field.type === 'number' ? 'number' : 'text'}
                        value={formValues[field.name] ?? ''}
                        placeholder={field.placeholder || ''}
                        onChange={(e) => setFormValues(prev => ({ ...prev, [field.name]: e.target.value }))}
                        style={{
                          width: '100%',
                          padding: '0.5rem',
                          borderRadius: '4px',
                          border: '1px solid #ccc'
                        }}
                      />
                    </div>
                  ))}
                </div>
              )}

              {(advancedMode || (solution.input_fields || []).length === 0) && (
                <textarea
                  value={payload}
                  onChange={(e) => setPayload(e.target.value)}
                  style={{
                    width: '100%',
                    height: '120px',
                    fontFamily: 'monospace',
                    padding: '0.5rem',
                    borderRadius: '4px',
                    border: '1px solid #ccc'
                  }}
                />
              )}
            </div>
          )}

          <button
            onClick={handleTest}
            disabled={loading}
            style={{
              padding: '0.5rem 1rem',
              backgroundColor: '#0066cc',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: loading ? 'not-allowed' : 'pointer'
            }}
          >
            {loading ? 'Testing...' : 'Run Test'}
          </button>

          {history.length > 0 && (
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 2fr', gap: '1rem', marginTop: '1rem' }}>
              <div>
                <div style={{ fontWeight: 700, marginBottom: '0.5rem' }}>Run History</div>
                <div style={{ border: '1px solid #ddd', borderRadius: '4px', maxHeight: '220px', overflowY: 'auto' }}>
                  {history.map((run, idx) => (
                    <button
                      key={run.at + idx}
                      onClick={() => setSelectedRunIndex(idx)}
                      style={{
                        display: 'block',
                        width: '100%',
                        textAlign: 'left',
                        padding: '0.5rem',
                        border: 'none',
                        borderBottom: '1px solid #eee',
                        backgroundColor: selectedRunIndex === idx ? '#f0f7ff' : '#fff',
                        cursor: 'pointer'
                      }}
                    >
                      <div style={{ fontSize: '0.85rem' }}>{new Date(run.at).toLocaleString()}</div>
                      <div style={{ fontSize: '0.8rem', color: '#666' }}>{run.request.method} {run.request.endpoint}</div>
                      <div style={{ fontSize: '0.8rem', fontWeight: 700 }}>Status: {run.response.status}</div>
                    </button>
                  ))}
                </div>
              </div>

              {selectedRun && (
                <div>
                  <div style={{ fontWeight: 700, marginBottom: '0.5rem' }}>Selected Run Output</div>
                  <div
                    style={{
                      padding: '1rem',
                      backgroundColor: selectedRun.response.status >= 200 && selectedRun.response.status < 300 ? '#e6ffe6' : '#ffe6e6',
                      borderLeft: `4px solid ${selectedRun.response.status >= 200 && selectedRun.response.status < 300 ? '#00cc00' : '#cc0000'}`,
                      borderRadius: '4px'
                    }}
                  >
                    <div style={{ marginBottom: '0.75rem', fontWeight: 'bold' }}>
                      Status: {selectedRun.response.status}
                    </div>

                    {selectedRun.request.body && (
                      <div style={{ marginBottom: '0.75rem' }}>
                        <div style={{ fontWeight: 700, marginBottom: '0.25rem' }}>Request Body</div>
                        <pre style={{ margin: 0, whiteSpace: 'pre-wrap', wordWrap: 'break-word', background: '#fff', padding: '0.5rem', borderRadius: '4px' }}>
                          {selectedRun.request.body}
                        </pre>
                      </div>
                    )}

                    <div style={{ fontWeight: 700, marginBottom: '0.25rem' }}>Response</div>
                    <pre style={{ margin: 0, whiteSpace: 'pre-wrap', wordWrap: 'break-word', background: '#fff', padding: '0.5rem', borderRadius: '4px' }}>
                      {typeof selectedRun.response.data === 'object' ? JSON.stringify(selectedRun.response.data, null, 2) : selectedRun.response.data}
                    </pre>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
