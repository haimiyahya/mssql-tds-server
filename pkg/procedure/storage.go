package procedure

import (
	"database/sql"
	"fmt"

	"github.com/factory/mssql-tds-server/pkg/sqlite"
)

// Storage handles stored procedure storage in SQLite
type Storage struct {
	db *sql.DB
}

// NewStorage creates a new procedure storage
func NewStorage(db *sqlite.Database) (*Storage, error) {
	return &Storage{
		db: db.GetDB(),
	}, nil
}

// Create stores a new procedure in the database
func (s *Storage) Create(proc *Procedure) error {
	paramsJSON, err := ParametersToJSON(proc.Parameters)
	if err != nil {
		return fmt.Errorf("failed to serialize parameters: %w", err)
	}

	query := `
	INSERT INTO procedures (name, body, parameters)
	VALUES (?, ?, ?)
	`

	result, err := s.db.Exec(query, proc.Name, proc.Body, paramsJSON)
	if err != nil {
		return fmt.Errorf("failed to create procedure: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	proc.ID = int(id)
	return nil
}

// Get retrieves a procedure by name
func (s *Storage) Get(name string) (*Procedure, error) {
	query := `
	SELECT id, name, body, parameters, created_at
	FROM procedures
	WHERE name = ?
	`

	var paramsJSON string
	proc := &Procedure{}

	err := s.db.QueryRow(query, name).Scan(
		&proc.ID,
		&proc.Name,
		&proc.Body,
		&paramsJSON,
		&proc.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("procedure '%s' not found", name)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve procedure: %w", err)
	}

	// Parse parameters
	proc.Parameters, err = ParametersFromJSON(paramsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	return proc, nil
}

// List returns all procedures
func (s *Storage) List() ([]*Procedure, error) {
	query := `
	SELECT id, name, body, parameters, created_at
	FROM procedures
	ORDER BY name
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list procedures: %w", err)
	}
	defer rows.Close()

	var procedures []*Procedure

	for rows.Next() {
		var paramsJSON string
		proc := &Procedure{}

		err := rows.Scan(
			&proc.ID,
			&proc.Name,
			&proc.Body,
			&paramsJSON,
			&proc.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan procedure: %w", err)
		}

		// Parse parameters
		proc.Parameters, err = ParametersFromJSON(paramsJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}

		procedures = append(procedures, proc)
	}

	return procedures, nil
}

// Drop removes a procedure by name
func (s *Storage) Drop(name string) error {
	query := `DELETE FROM procedures WHERE name = ?`

	result, err := s.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to drop procedure: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("procedure '%s' not found", name)
	}

	return nil
}

// Exists checks if a procedure exists
func (s *Storage) Exists(name string) (bool, error) {
	query := `SELECT COUNT(*) FROM procedures WHERE name = ?`

	var count int
	err := s.db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check procedure existence: %w", err)
	}

	return count > 0, nil
}
