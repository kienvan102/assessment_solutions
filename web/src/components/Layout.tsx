import { Link, useLocation } from 'react-router-dom';

export default function Layout({ children }: { children: React.ReactNode }) {
  const location = useLocation();

  return (
    <div style={{ padding: '2rem', fontFamily: 'system-ui, sans-serif', maxWidth: '1000px', margin: '0 auto' }}>
      <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem', borderBottom: '2px solid #eee', paddingBottom: '1rem' }}>
        <h1 style={{ margin: 0, color: '#333' }}>Interactive Assessment</h1>
        <nav style={{ display: 'flex', gap: '1rem' }}>
          <Link
            to="/"
            style={{
              textDecoration: 'none',
              padding: '0.5rem 1rem',
              borderRadius: '4px',
              backgroundColor: location.pathname === '/' ? '#0066cc' : 'transparent',
              color: location.pathname === '/' ? 'white' : '#0066cc',
              fontWeight: 600,
            }}
          >
            Solutions
          </Link>
          <Link
            to="/architecture"
            style={{
              textDecoration: 'none',
              padding: '0.5rem 1rem',
              borderRadius: '4px',
              backgroundColor: location.pathname === '/architecture' ? '#0066cc' : 'transparent',
              color: location.pathname === '/architecture' ? 'white' : '#0066cc',
              fontWeight: 600,
            }}
          >
            Architecture
          </Link>
        </nav>
      </header>
      <main>
        {children}
      </main>
    </div>
  );
}
