package sql1

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Service struct {
	db *sql.DB
}

func NewService() (*Service, error) {
	// Use an in-memory database that is shared within this connection
	// Using a unique URI ensures the memory DB persists across queries for this db object
	db, err := sql.Open("sqlite", "file:sql1?mode=memory&cache=shared")
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}

	// Ensure connection is established
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite db: %w", err)
	}

	svc := &Service{db: db}
	if err := svc.seedDatabase(); err != nil {
		return nil, fmt.Errorf("failed to seed database: %w", err)
	}

	return svc, nil
}

func (s *Service) seedDatabase() error {
	queries := []string{
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL
		)`,
		`CREATE TABLE orders (
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			amount DECIMAL(10, 2),
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`,
		`INSERT INTO users (id, name) VALUES 
			(101, 'Alice'),
			(102, 'Bob'),
			(103, 'Charlie')`, // Charlie added to explicitly test LEFT JOIN
		`INSERT INTO orders (id, user_id, amount) VALUES 
			(1, 101, 50.00),
			(2, 101, 75.00),
			(3, 102, 30.00)`,
	}

	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return fmt.Errorf("error executing seed query '%s': %w", q, err)
		}
	}
	return nil
}

// ExecuteQuery runs a raw SQL query and returns the column names and rows.
func (s *Service) ExecuteQuery(query string) ([]string, [][]interface{}, error) {
	// Execute the query
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	// Get column names
	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var result [][]interface{}

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Extract the actual values
		rowData := make([]interface{}, len(cols))
		for i, col := range columns {
			val := col
			// SQLite driver sometimes returns []byte for strings/numbers
			if b, ok := val.([]byte); ok {
				rowData[i] = string(b)
			} else {
				rowData[i] = val
			}
		}
		result = append(result, rowData)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("row iteration error: %w", err)
	}

	return cols, result, nil
}
