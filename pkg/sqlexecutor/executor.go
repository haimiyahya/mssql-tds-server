package sqlexecutor

import (
	"database/sql"
	"fmt"

	"github.com/factory/mssql-tds-server/pkg/sqlparser"
)

// Executor handles SQL statement execution
type Executor struct {
	db *sql.DB
}

// NewExecutor creates a new SQL executor
func NewExecutor(db *sql.DB) *Executor {
	return &Executor{
		db: db,
	}
}

// ExecuteResult represents the result of SQL execution
type ExecuteResult struct {
	Columns    []string
	Rows       [][]interface{}
	RowCount   int64
	IsQuery    bool
	Message    string
}

// Execute executes a SQL query and returns results
func (e *Executor) Execute(query string) (*ExecuteResult, error) {
	// Strip comments from query
	query = sqlparser.StripComments(query)

	// Parse the query to determine statement type
	stmt, err := sqlparser.NewParser().Parse(query)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}

	switch stmt.Type {
	case sqlparser.StatementTypeSelect:
		return e.executeSelect(query)

	case sqlparser.StatementTypeInsert:
		return e.executeInsert(query)

	case sqlparser.StatementTypeUpdate:
		return e.executeUpdate(query)

	case sqlparser.StatementTypeDelete:
		return e.executeDelete(query)

	case sqlparser.StatementTypeCreateTable:
		return e.executeCreateTable(query)

	case sqlparser.StatementTypeDropTable:
		return e.executeDropTable(query)

	default:
		// Try to execute as raw SQL (for unsupported statements)
		return e.executeRaw(query)
	}
}

// executeSelect executes a SELECT query
func (e *Executor) executeSelect(query string) (*ExecuteResult, error) {
	rows, err := e.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SELECT: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	// Read all rows
	var resultRows [][]interface{}
	for rows.Next() {
		// Create slice for values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// Scan row
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		resultRows = append(resultRows, values)
	}

	// Check for errors after scanning
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning rows: %w", err)
	}

	return &ExecuteResult{
		Columns:  columns,
		Rows:     resultRows,
		RowCount: int64(len(resultRows)),
		IsQuery:  true,
	}, nil
}

// executeInsert executes an INSERT statement
func (e *Executor) executeInsert(query string) (*ExecuteResult, error) {
	result, err := e.db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute INSERT: %w", err)
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return &ExecuteResult{
		RowCount: rowCount,
		IsQuery:  false,
		Message:  fmt.Sprintf("%d row(s) inserted", rowCount),
	}, nil
}

// executeUpdate executes an UPDATE statement
func (e *Executor) executeUpdate(query string) (*ExecuteResult, error) {
	result, err := e.db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute UPDATE: %w", err)
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return &ExecuteResult{
		RowCount: rowCount,
		IsQuery:  false,
		Message:  fmt.Sprintf("%d row(s) updated", rowCount),
	}, nil
}

// executeDelete executes a DELETE statement
func (e *Executor) executeDelete(query string) (*ExecuteResult, error) {
	result, err := e.db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute DELETE: %w", err)
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return &ExecuteResult{
		RowCount: rowCount,
		IsQuery:  false,
		Message:  fmt.Sprintf("%d row(s) deleted", rowCount),
	}, nil
}

// executeCreateTable executes a CREATE TABLE statement
func (e *Executor) executeCreateTable(query string) (*ExecuteResult, error) {
	// Convert T-SQL CREATE TABLE to SQLite-compatible SQL
	sqliteQuery := e.convertCreateTable(query)

	_, err := e.db.Exec(sqliteQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute CREATE TABLE: %w", err)
	}

	return &ExecuteResult{
		RowCount: 0,
		IsQuery:  false,
		Message:  "Table created successfully",
	}, nil
}

// executeDropTable executes a DROP TABLE statement
func (e *Executor) executeDropTable(query string) (*ExecuteResult, error) {
	_, err := e.db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute DROP TABLE: %w", err)
	}

	return &ExecuteResult{
		RowCount: 0,
		IsQuery:  false,
		Message:  "Table dropped successfully",
	}, nil
}

// executeRaw executes raw SQL (for unsupported statement types)
func (e *Executor) executeRaw(query string) (*ExecuteResult, error) {
	// Try to execute as query first
	rows, err := e.db.Query(query)
	if err == nil {
		defer rows.Close()

		// Get column names
		columns, err := rows.Columns()
		if err != nil {
			return nil, fmt.Errorf("failed to get columns: %w", err)
		}

		// Read all rows
		var resultRows [][]interface{}
		for rows.Next() {
			// Create slice for values
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range columns {
				valuePtrs[i] = &values[i]
			}

			// Scan row
			if err := rows.Scan(valuePtrs...); err != nil {
				return nil, fmt.Errorf("failed to scan row: %w", err)
			}

			resultRows = append(resultRows, values)
		}

		return &ExecuteResult{
			Columns:  columns,
			Rows:     resultRows,
			RowCount: int64(len(resultRows)),
			IsQuery:  true,
		}, nil
	}

	// Try to execute as non-query
	result, err := e.db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute raw SQL: %w", err)
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return &ExecuteResult{
		RowCount: rowCount,
		IsQuery:  false,
		Message:  fmt.Sprintf("%d row(s) affected", rowCount),
	}, nil
}

// convertCreateTable converts T-SQL CREATE TABLE to SQLite-compatible SQL
func (e *Executor) convertCreateTable(query string) string {
	// Basic conversion - for production, use a proper SQL parser/converter
	// This handles common differences:

	// 1. Remove square brackets
	result := query
	result = replaceAll(result, "[", "")
	result = replaceAll(result, "]", "")

	// 2. Convert IDENTITY to AUTOINCREMENT
	// This is a basic replacement - may need more sophisticated handling
	// result = convertIdentity(result)

	// 3. Other conversions could be added here

	return result
}

// replaceAll replaces all occurrences of old with new in s
func replaceAll(s, old, new string) string {
	result := ""
	for {
		idx := findSubstring(s, old)
		if idx == -1 {
			break
		}
		result += s[:idx] + new
		s = s[idx+len(old):]
	}
	result += s
	return result
}

// findSubstring finds the first occurrence of substr in s
func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// ConvertValueToString converts a database value to string for TDS
func ConvertValueToString(value interface{}) string {
	if value == nil {
		return "NULL"
	}

	switch v := value.(type) {
	case []byte:
		return string(v)
	default:
		return fmt.Sprintf("%v", value)
	}
}
