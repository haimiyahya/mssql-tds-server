package procedure

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/factory/mssql-tds-server/pkg/sqlite"
)

// Executor handles stored procedure execution
type Executor struct {
	db       *sql.DB
	storage  *Storage
}

// NewExecutor creates a new procedure executor
func NewExecutor(db *sqlite.Database, storage *Storage) (*Executor, error) {
	return &Executor{
		db:      db.GetDB(),
		storage: storage,
	}, nil
}

// Execute executes a stored procedure with given parameters
func (e *Executor) Execute(name string, paramValues map[string]interface{}) ([][]string, error) {
	// Retrieve procedure
	proc, err := e.storage.Get(name)
	if err != nil {
		return nil, err
	}

	// Validate parameters
	if err := e.validateParameters(proc, paramValues); err != nil {
		return nil, fmt.Errorf("parameter validation failed: %w", err)
	}

	// Replace parameters in SQL
	sql, err := e.replaceParameters(proc.Body, paramValues)
	if err != nil {
		return nil, fmt.Errorf("parameter replacement failed: %w", err)
	}

	// Execute SQL
	rows, err := e.db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("failed to execute procedure: %w", err)
	}
	defer rows.Close()

	// Read results
	results, err := e.readResults(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to read results: %w", err)
	}

	return results, nil
}

// validateParameters validates that all required parameters are provided
func (e *Executor) validateParameters(proc *Procedure, paramValues map[string]interface{}) error {
	// Check if all required parameters are provided
	for _, param := range proc.Parameters {
		if !param.HasDefault {
			if _, exists := paramValues["@"+param.Name]; !exists {
				return fmt.Errorf("missing required parameter: @%s", param.Name)
			}
		}
	}

	// Check that no extra parameters are provided
	for paramName := range paramValues {
		found := false
		for _, param := range proc.Parameters {
			if "@"+param.Name == paramName {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("unexpected parameter: %s", paramName)
		}
	}

	return nil
}

// replaceParameters replaces parameter placeholders with actual values
func (e *Executor) replaceParameters(sql string, paramValues map[string]interface{}) (string, error) {
	// Replace parameters in format: @param
	result := sql

	// Build regex pattern to find all @param references
	re := regexp.MustCompile(`@\w+`)

	// Replace each parameter
	result = re.ReplaceAllStringFunc(result, func(match string) string {
		// Check if this is a parameter we have a value for
		if value, exists := paramValues[match]; exists {
			return e.formatValue(value)
		}
		// If no value provided, keep as is (might be a default or error will be caught later)
		return match
	})

	return result, nil
}

// formatValue formats a value for SQL
func (e *Executor) formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		// Escape single quotes and wrap in quotes
		escaped := strings.ReplaceAll(v, "'", "''")
		return fmt.Sprintf("'%s'", escaped)
	case int, int32, int64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	case nil:
		return "NULL"
	default:
		// Default to string format
		escaped := strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''")
		return fmt.Sprintf("'%s'", escaped)
	}
}

// readResults reads all rows from a query result
func (e *Executor) readResults(rows *sql.Rows) ([][]string, error) {
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results [][]string

	// Add column names as first row
	results = append(results, columns)

	// Read all rows
	for rows.Next() {
		// Create slice of interfaces for row values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan row
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert values to strings
		row := make([]string, len(columns))
		for i, val := range values {
			if val == nil {
				row[i] = "NULL"
			} else {
				row[i] = fmt.Sprintf("%v", val)
			}
		}

		results = append(results, row)
	}

	// Check for errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return results, nil
}
