export default function ArchitecturePage() {
  return (
    <div style={{ lineHeight: '1.6', color: '#333' }}>
      <h2 style={{ borderBottom: '1px solid #eee', paddingBottom: '0.5rem' }}>Project Architecture & Design</h2>
      
      <p style={{ fontSize: '1.1rem' }}>
        Welcome to my interactive technical assessment platform. I built this project to demonstrate my full-stack software engineering capabilities to you. It specifically focuses on clean architecture, dependency injection, and decoupled systems to show how I write scalable, production-ready code.
      </p>

      <div style={{ marginTop: '2rem' }}>
        <h3>Backend Structure (Go)</h3>
        <p>The Go backend follows a domain-driven, clean architecture approach. I designed it to eliminate global state and hardcoded dependencies.</p>
        <ul style={{ paddingLeft: '1.5rem' }}>
          <li style={{ marginBottom: '0.5rem' }}>
            <strong>cmd/api/main.go:</strong> Acts purely as the application entry point. It contains no business logic.
          </li>
          <li style={{ marginBottom: '0.5rem' }}>
            <strong>internal/app/app.go & internal/di/container.go:</strong> Manages the Dependency Injection (DI). It instantiates stores, injects them into services, and wires up the HTTP routers.
          </li>
          <li style={{ marginBottom: '0.5rem' }}>
            <strong>internal/solutions/:</strong> Contains the core business logic domains (e.g., <code>payment</code>, <code>workerpool</code>, <code>sql1</code>). They define their own interfaces for data storage to invert dependencies.
          </li>
          <li style={{ marginBottom: '0.5rem' }}>
            <strong>In-Memory SQLite (modernc.org/sqlite):</strong> For database questions, we use a CGO-free pure Go SQLite driver. This allows us to spin up isolated, in-memory databases (<code>:memory:</code>) on the fly without requiring developers to run separate Postgres/MySQL containers, while still executing standard SQL.
          </li>
          <li style={{ marginBottom: '0.5rem' }}>
            <strong>data/solutions.json:</strong> Contains the dynamic UI metadata. By keeping the expected behaviors, input schemas, and descriptions in JSON, the Go binary doesn't need to be recompiled to adjust the test definitions.
          </li>
        </ul>
      </div>

      <div style={{ marginTop: '2rem' }}>
        <h3>Frontend Structure (React)</h3>
        <p>The frontend is a React application utilizing Vite. I designed it to dynamically generate its UI based on the backend data, but intelligently swap to custom interactive components for complex questions so you can test my solutions directly in the browser.</p>
        <ul style={{ paddingLeft: '1.5rem' }}>
          <li style={{ marginBottom: '0.5rem' }}>
            <strong>Dynamic Rendering:</strong> It fetches <code>/api/solutions</code> on mount and loops through the questions.
          </li>
          <li style={{ marginBottom: '0.5rem' }}>
            <strong>Polymorphic Components:</strong> While standard questions use the generic <code>SolutionCard</code> form builder, questions requiring deep interaction (like Code Reviews or SQL execution) automatically render highly customized components (<code>CodeReviewCard</code>, <code>SqlOptimizationCard</code>).
          </li>
        </ul>
      </div>

      <div style={{ marginTop: '3rem' }}>
        <h2 style={{ borderBottom: '1px solid #eee', paddingBottom: '0.5rem' }}>Solution Approaches</h2>
        <p>Below is a breakdown of the architectural approach I took to solve each question in the assessment.</p>

        {/* Q1 */}
        <div style={{ marginBottom: '1.5rem', padding: '1.5rem', backgroundColor: '#f8f9fa', borderRadius: '8px', borderLeft: '4px solid #0066cc' }}>
          <h4 style={{ margin: '0 0 0.5rem 0', fontSize: '1.2rem' }}>Q1: Idempotent Payment System</h4>
          <p style={{ margin: '0 0 1rem 0' }}><strong>Approach:</strong> In-Memory Key-Value Store with Concurrency Control</p>
          <ul style={{ margin: 0, paddingLeft: '1.2rem' }}>
            <li>I created a highly reusable, generic, thread-safe <code>Store[K, V]</code> package using <code>sync.RWMutex</code>.</li>
            <li><strong>Dependency Inversion:</strong> I designed the payment service to depend on a local <code>Store interface</code> rather than importing the database implementation directly.</li>
            <li>I simulated async processing by spawning a goroutine that waits before transitioning the state from <code>pending</code> to <code>processed</code>. If the exact same <code>transactionID</code> is submitted during or after this window, it returns the existing state safely.</li>
          </ul>
        </div>

        {/* Q2 */}
        <div style={{ marginBottom: '1.5rem', padding: '1.5rem', backgroundColor: '#f8f9fa', borderRadius: '8px', borderLeft: '4px solid #0066cc' }}>
          <h4 style={{ margin: '0 0 0.5rem 0', fontSize: '1.2rem' }}>Q2: Worker Pool Task Processor</h4>
          <p style={{ margin: '0 0 1rem 0' }}><strong>Approach:</strong> Channel-based Concurrency</p>
          <ul style={{ margin: 0, paddingLeft: '1.2rem' }}>
            <li>I utilized Go channels (<code>tasks</code> and <code>results</code>) to safely distribute work across exactly 5 worker goroutines.</li>
            <li>Instead of waiting sequentially, I fire all tasks into a buffered channel. The workers pull from the channel as fast as they can process.</li>
            <li>I collect the results asynchronously, and then sort them back into the original ID order before returning the JSON payload, ensuring correctness alongside concurrency speed.</li>
          </ul>
        </div>

        {/* Q3 & Q6 */}
        <div style={{ marginBottom: '1.5rem', padding: '1.5rem', backgroundColor: '#f8f9fa', borderRadius: '8px', borderLeft: '4px solid #0066cc' }}>
          <h4 style={{ margin: '0 0 0.5rem 0', fontSize: '1.2rem' }}>Q3 & Q6: Code Review (HTTP Handlers)</h4>
          <p style={{ margin: '0 0 1rem 0' }}><strong>Approach:</strong> Interactive Simulation Environment</p>
          <ul style={{ margin: 0, paddingLeft: '1.2rem' }}>
            <li>Rather than just providing text answers, I built actual running implementations of the "Bad" code and my "Good" refactored code.</li>
            <li>I created a <code>SimulatorService</code> that intentionally hammers the Bad endpoints using <code>sync.WaitGroup</code> to artificially induce and catch the race conditions caused by the original global variables.</li>
            <li>I built a custom React component (<code>CodeReviewCard</code>) with tabs allowing you to manually trigger XSS payloads against the bad endpoint, or view my written architectural analysis outlining the memory exhaustion and concurrency vulnerabilities.</li>
          </ul>
        </div>

        {/* Q4 */}
        <div style={{ marginBottom: '1.5rem', padding: '1.5rem', backgroundColor: '#f8f9fa', borderRadius: '8px', borderLeft: '4px solid #0066cc' }}>
          <h4 style={{ margin: '0 0 0.5rem 0', fontSize: '1.2rem' }}>Q4: SQL LEFT JOIN</h4>
          <p style={{ margin: '0 0 1rem 0' }}><strong>Approach:</strong> Dynamic In-Memory SQL Playground</p>
          <ul style={{ margin: 0, paddingLeft: '1.2rem' }}>
            <li>I used a CGO-free SQLite driver to bootstrap an ephemeral in-memory database on application startup.</li>
            <li>I seed the database programmatically with the exact <code>users</code> and <code>orders</code> table structures required for the assessment.</li>
            <li>I built a custom query editor UI that allows you to execute raw SQL against the in-memory database and view the parsed rows in a dynamic HTML table to prove my query works.</li>
          </ul>
        </div>

        {/* Q5 */}
        <div style={{ marginBottom: '1.5rem', padding: '1.5rem', backgroundColor: '#f8f9fa', borderRadius: '8px', borderLeft: '4px solid #0066cc' }}>
          <h4 style={{ margin: '0 0 0.5rem 0', fontSize: '1.2rem' }}>Q5: SQL Optimization / Indexing</h4>
          <p style={{ margin: '0 0 1rem 0' }}><strong>Approach:</strong> Explain Query Plan Visualizer</p>
          <ul style={{ margin: 0, paddingLeft: '1.2rem' }}>
            <li>I seeded a <code>transactions</code> table with 10,000 randomized rows to simulate a large dataset.</li>
            <li>I exposed specialized API actions to run <code>EXPLAIN QUERY PLAN</code> strings.</li>
            <li>I built a 3-step interactive UI to demonstrate my understanding of database indexing: (1) Run the raw query, (2) See the Database Engine fall back to a Full Table Scan (SCAN TABLE), and (3) Execute my <code>CREATE INDEX</code> DDL statement and watch the execution plan dynamically upgrade to an Index Search (SEARCH TABLE USING INDEX).</li>
          </ul>
        </div>

      </div>
    </div>
  );
}
