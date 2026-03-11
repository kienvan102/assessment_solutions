export default function ArchitecturePage() {
  return (
    <div style={{ lineHeight: '1.6', color: '#333' }}>
      <h2 style={{ borderBottom: '1px solid #eee', paddingBottom: '0.5rem' }}>Project Architecture & Design</h2>
      
      <p style={{ fontSize: '1.1rem' }}>
        This project is designed as an interactive technical assessment platform. It demonstrates full-stack software engineering principles, specifically focusing on clean architecture, dependency injection, and decoupled systems.
      </p>

      <div style={{ marginTop: '2rem' }}>
        <h3>Backend Structure (Go)</h3>
        <p>The Go backend follows a domain-driven, clean architecture approach, eliminating global state and hardcoded dependencies.</p>
        <ul style={{ paddingLeft: '1.5rem' }}>
          <li style={{ marginBottom: '0.5rem' }}>
            <strong>cmd/api/main.go:</strong> Acts purely as the application entry point. It contains no business logic.
          </li>
          <li style={{ marginBottom: '0.5rem' }}>
            <strong>internal/app/app.go & internal/di/container.go:</strong> Manages the Dependency Injection (DI). It instantiates stores, injects them into services, and wires up the HTTP routers.
          </li>
          <li style={{ marginBottom: '0.5rem' }}>
            <strong>internal/solutions/:</strong> Contains the core business logic domains (e.g., the <code>payment</code> and <code>palindrome</code> packages). They define their own interfaces for data storage to invert dependencies.
          </li>
          <li style={{ marginBottom: '0.5rem' }}>
            <strong>pkg/store/:</strong> A highly reusable, generic, thread-safe in-memory key-value store (<code>Store[K, V]</code>). It is unaware of the business logic and simply provides persistence.
          </li>
          <li style={{ marginBottom: '0.5rem' }}>
            <strong>data/solutions.json:</strong> Contains the dynamic UI metadata. By keeping the expected behaviors, input schemas, and descriptions in JSON, the Go binary doesn't need to be recompiled just to fix a typo in the UI.
          </li>
        </ul>
      </div>

      <div style={{ marginTop: '2rem' }}>
        <h3>Frontend Structure (React)</h3>
        <p>The frontend is a lightweight React application that dynamically generates its UI based on the backend data.</p>
        <ul style={{ paddingLeft: '1.5rem' }}>
          <li style={{ marginBottom: '0.5rem' }}>
            <strong>Dynamic Rendering:</strong> It fetches <code>/api/solutions</code> on mount. It doesn't hardcode "Payment" or "Palindrome" forms. Instead, it reads the <code>input_fields</code> array from the JSON and dynamically builds the form inputs.
          </li>
          <li style={{ marginBottom: '0.5rem' }}>
            <strong>Service Layer:</strong> <code>src/services/api.ts</code> abstracts all HTTP fetch logic, ensuring components like <code>SolutionCard</code> remain focused strictly on UI state.
          </li>
        </ul>
      </div>

      <div style={{ marginTop: '2rem' }}>
        <h3>Key Design Decisions</h3>
        
        <div style={{ marginBottom: '1rem', padding: '1rem', backgroundColor: '#f8f9fa', borderRadius: '4px', borderLeft: '4px solid #0066cc' }}>
          <h4 style={{ margin: '0 0 0.5rem 0' }}>1. Idempotency & Simulated Delay</h4>
          <p style={{ margin: 0 }}>
            The payment system simulates real-world processing. A new transaction is immediately recorded as <code>pending</code>. A goroutine waits for a configurable duration (<code>PAYMENT_PROCESSING_DELAY</code>) before updating the state to <code>processed</code>. If the client retries the same <code>transactionID</code> during this window, it safely returns the existing state without creating a duplicate.
          </p>
        </div>

        <div style={{ marginBottom: '1rem', padding: '1rem', backgroundColor: '#f8f9fa', borderRadius: '4px', borderLeft: '4px solid #0066cc' }}>
          <h4 style={{ margin: '0 0 0.5rem 0' }}>2. Dependency Inversion</h4>
          <p style={{ margin: 0 }}>
            The <code>payment.Service</code> does not import the generic <code>pkg/store</code> directly. Instead, it declares a local <code>Store interface</code>. This means the payment logic is entirely decoupled from the database implementation and can easily be unit tested with a mock store.
          </p>
        </div>
      </div>
    </div>
  );
}
