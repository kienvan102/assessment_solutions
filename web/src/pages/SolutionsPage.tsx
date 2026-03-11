import { useEffect, useState } from 'react';
import SolutionCard from '../components/SolutionCard';
import CodeReviewCard from '../components/CodeReview/CodeReviewCard';
import { getHealth, getSolutions } from '../services/api';
import type { Solution } from '../types/solution';

export default function SolutionsPage() {
  const [health, setHealth] = useState<string>('checking...');
  const [solutions, setSolutions] = useState<Solution[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    getHealth()
      .then(setHealth)
      .catch(err => setHealth('error: ' + err.message));

    getSolutions()
      .then(data => {
        setSolutions(data);
        setLoading(false);
      })
      .catch(err => {
        console.error('Failed to fetch solutions:', err);
        setLoading(false);
      });
  }, []);

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'flex-end', marginBottom: '1.5rem' }}>
        <div
          style={{
            padding: '0.5rem 1rem',
            backgroundColor: health === 'ok' ? '#e6ffe6' : '#ffe6e6',
            borderRadius: '20px',
            fontSize: '0.9rem'
          }}
        >
          API Status: <strong style={{ color: health === 'ok' ? '#008000' : '#cc0000' }}>{health}</strong>
        </div>
      </div>

      <div>
        {loading ? (
          <p>Loading solutions...</p>
        ) : solutions.length === 0 ? (
          <p>No solutions found.</p>
        ) : (
          solutions.map(solution => {
            if (solution.id === 'q3') {
              return <CodeReviewCard key={solution.id} solution={solution} />;
            }
            return <SolutionCard key={solution.id} solution={solution} />;
          })
        )}
      </div>
    </div>
  );
}
