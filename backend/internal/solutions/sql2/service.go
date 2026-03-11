package sql2

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "modernc.org/sqlite"
)

type Service struct {
	db *sql.DB
}

func NewService() (*Service, error) {
	// Unique shared memory DB for sql2
	db, err := sql.Open("sqlite", "file:sql2?mode=memory&cache=shared")
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite db: %w", err)
	}

	svc := &Service{db: db}
	if err := svc.ResetDatabase(); err != nil {
		return nil, fmt.Errorf("failed to seed database: %w", err)
	}

	return svc, nil
}

func (s *Service) ResetDatabase() error {
	// Drop existing table if any
	_, _ = s.db.Exec(`DROP TABLE IF EXISTS transactions;`)

	createTable := `
	CREATE TABLE transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		amount DECIMAL(10, 2),
		created_at TIMESTAMP
	);`

	if _, err := s.db.Exec(createTable); err != nil {
		return err
	}

	// Seed with a decent amount of data (e.g., 10,000 rows) to simulate a "millions of rows" table
	// We ensure user_id = 123 has multiple transactions spread out.
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO transactions (user_id, amount, created_at) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	r := rand.New(rand.NewSource(42))
	baseTime := time.Now().Add(-30 * 24 * time.Hour) // 30 days ago

	for i := 0; i < 10000; i++ {
		userID := r.Intn(1000) + 1 // random user 1-1000

		// Force some rows for user 123
		if i%50 == 0 {
			userID = 123
		}

		amount := float64(r.Intn(50000)) / 100.0
		createdAt := baseTime.Add(time.Duration(i) * time.Minute)

		if _, err := stmt.Exec(userID, amount, createdAt); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// ExecuteQuery runs a raw SQL query and returns the columns and rows.
func (s *Service) ExecuteQuery(query string) ([]string, [][]interface{}, error) {
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var result [][]interface{}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}

		rowData := make([]interface{}, len(cols))
		for i, col := range columns {
			val := col
			if b, ok := val.([]byte); ok {
				rowData[i] = string(b)
			} else {
				rowData[i] = val
			}
		}
		result = append(result, rowData)
	}

	return cols, result, nil
}

// ExecQuery runs a raw SQL query (like CREATE INDEX) without returning rows.
func (s *Service) ExecQuery(query string) error {
	_, err := s.db.Exec(query)
	return err
}

// ExplainQuery runs EXPLAIN QUERY PLAN and returns the plan description string.
func (s *Service) ExplainQuery(query string) (string, error) {
	explainQuery := fmt.Sprintf("EXPLAIN QUERY PLAN %s", query)

	rows, err := s.db.Query(explainQuery)
	if err != nil {
		return "", fmt.Errorf("explain failed: %w", err)
	}
	defer rows.Close()

	var fullPlan string
	for rows.Next() {
		var id, parent, notused int
		var detail string
		// SQLite EXPLAIN QUERY PLAN returns: id, parent, notused, detail
		if err := rows.Scan(&id, &parent, &notused, &detail); err != nil {
			return "", err
		}
		fullPlan += detail + "\n"
	}

	return fullPlan, nil
}
