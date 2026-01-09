package sqlexecutor

import (
	"database/sql"
	"fmt"

	"github.com/factory/mssql-tds-server/pkg/sqlparser"
)

// Executor handles SQL statement execution
type Executor struct {
	db    *sql.DB
	views map[string]string // Store view name -> SELECT query mapping
}

// NewExecutor creates a new SQL executor
func NewExecutor(db *sql.DB) *Executor {
	return &Executor{
		db:    db,
		views: make(map[string]string),
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

	case sqlparser.StatementTypeCreateView:
		return e.executeCreateView(query)

	case sqlparser.StatementTypeDropView:
		return e.executeDropView(query)

	case sqlparser.StatementTypeCreateIndex:
		return e.executeCreateIndex(query)

	case sqlparser.StatementTypeDropIndex:
		return e.executeDropIndex(query)

	case sqlparser.StatementTypeBeginTransaction:
		return e.executeBeginTransaction(query)

	case sqlparser.StatementTypeCommit:
		return e.executeCommit(query)

	case sqlparser.StatementTypeRollback:
		return e.executeRollback(query)

	default:
		// Try to execute as raw SQL (for unsupported statements)
		return e.executeRaw(query)
	}
}

// executeSelect executes a SELECT query
func (e *Executor) executeSelect(query string) (*ExecuteResult, error) {
	// Parse the query to get ORDER BY and DISTINCT information
	stmt, err := sqlparser.NewParser().Parse(query)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SELECT query: %w", err)
	}

	// If not a SELECT statement, execute as raw SQL
	if stmt.Type != sqlparser.StatementTypeSelect || stmt.Select == nil {
		return e.executeRaw(query)
	}

	// Remove ORDER BY and DISTINCT from query if present
	// SQLite supports them natively, but we want to apply them ourselves
	// For now, let SQLite handle them (simpler approach)
	// In production, we would implement custom ORDER BY and DISTINCT logic

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

	// Apply ORDER BY (if parsed and we want to handle it manually)
	// For Phase 11 Iteration 1, we'll let SQLite handle ORDER BY
	// In a more advanced implementation, we would:
	// resultRows = e.sortRows(resultRows, stmt.Select.OrderBy, columns)

	// Apply DISTINCT (if parsed and we want to handle it manually)
	// For Phase 11 Iteration 1, we'll let SQLite handle DISTINCT
	// In a more advanced implementation, we would:
	// if stmt.Select.Distinct {
	// 	resultRows = e.removeDuplicates(resultRows)
	// }

	// Handle aggregate functions
	// For Phase 11 Iteration 2, we'll let SQLite handle aggregate functions
	// SQLite supports COUNT, SUM, AVG, MIN, MAX natively
	// In a more advanced implementation, we would:
	// if stmt.Select.IsAggregateQuery {
	// 	resultRows = e.calculateAggregates(resultRows, stmt.Select.Aggregates)
	// }
	// For now, SQLite handles all aggregate calculations

	// Handle GROUP BY and HAVING
	// For Phase 11 Iteration 3, we'll let SQLite handle GROUP BY and HAVING
	// SQLite supports GROUP BY and HAVING natively
	// In a more advanced implementation, we would:
	// if len(stmt.Select.GroupBy) > 0 {
	// 	resultRows = e.groupBy(resultRows, stmt.Select.GroupBy, stmt.Select.Aggregates)
	// 	if stmt.Select.HavingClause != "" {
	// 		resultRows = e.filterHaving(resultRows, stmt.Select.HavingClause)
	// 	}
	// }
	// For now, SQLite handles all GROUP BY and HAVING calculations

	// Handle JOINs
	// For Phase 11 Iteration 4, we'll let SQLite handle JOINs
	// SQLite supports INNER, LEFT, and CROSS JOINs natively
	// RIGHT JOIN and FULL JOIN are not supported by SQLite directly
	// Workarounds would be needed for RIGHT and FULL JOINs:
	// - RIGHT JOIN: Can be emulated by swapping table order and using LEFT JOIN
	// - FULL JOIN: Can be emulated by combining LEFT JOIN and RIGHT JOIN with UNION
	// In a more advanced implementation, we would:
	// if len(stmt.Select.Joins) > 0 {
	// 	for _, join := range stmt.Select.Joins {
	// 		if join.Type == "RIGHT" || join.Type == "FULL" {
	// 			// Implement JOIN workaround
	// 		}
	// 	}
	// }
	// For now, let SQLite handle all JOINs (will fail on RIGHT and FULL with error)

	// Handle Subqueries
	// For Phase 11 Iteration 5, we'll let SQLite handle subqueries
	// SQLite supports subqueries natively in WHERE clause (IN, EXISTS, =, !=, >, <, >=, <=)
	// SQLite supports subqueries in SELECT list
	// SQLite supports subqueries in FROM clause (derived tables)
	// In a more advanced implementation, we would:
	// if stmt.Select.HasSubqueries {
	// 	// Parse and execute subqueries separately
	// 	// Substitute subquery results into main query
	// }
	// For now, let SQLite handle all subqueries

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

// executeCreateView executes a CREATE VIEW statement
func (e *Executor) executeCreateView(query string) (*ExecuteResult, error) {
	// Parse query to get view information
	stmt, err := sqlparser.NewParser().Parse(query)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CREATE VIEW: %w", err)
	}

	if stmt.Type != sqlparser.StatementTypeCreateView || stmt.CreateView == nil {
		return nil, fmt.Errorf("invalid CREATE VIEW statement")
	}

	// Store view definition
	e.views[stmt.CreateView.ViewName] = stmt.CreateView.SelectQuery

	// Execute CREATE VIEW on SQLite (SQLite supports CREATE VIEW natively)
	_, err = e.db.Exec(query)
	if err != nil {
		// If SQLite fails, we still have the view definition stored
		// This allows us to handle queries against the view
		return nil, fmt.Errorf("failed to create view in SQLite: %w (view definition stored)", err)
	}

	return &ExecuteResult{
		RowCount: 0,
		IsQuery:  false,
		Message:  fmt.Sprintf("View '%s' created successfully", stmt.CreateView.ViewName),
	}, nil
}

// executeDropView executes a DROP VIEW statement
func (e *Executor) executeDropView(query string) (*ExecuteResult, error) {
	// Parse query to get view name
	stmt, err := sqlparser.NewParser().Parse(query)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DROP VIEW: %w", err)
	}

	if stmt.Type != sqlparser.StatementTypeDropView || stmt.DropView == nil {
		return nil, fmt.Errorf("invalid DROP VIEW statement")
	}

	// Remove view definition
	delete(e.views, stmt.DropView.ViewName)

	// Execute DROP VIEW on SQLite (SQLite supports DROP VIEW natively)
	_, err = e.db.Exec(query)
	if err != nil {
		// If SQLite fails, we still removed the view definition
		return nil, fmt.Errorf("failed to drop view in SQLite: %w (view definition removed)", err)
	}

	return &ExecuteResult{
		RowCount: 0,
		IsQuery:  false,
		Message:  fmt.Sprintf("View '%s' dropped successfully", stmt.DropView.ViewName),
	}, nil
}

// executeCreateIndex executes a CREATE INDEX statement
func (e *Executor) executeCreateIndex(query string) (*ExecuteResult, error) {
	// Parse query to get index information
	stmt, err := sqlparser.NewParser().Parse(query)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CREATE INDEX: %w", err)
	}

	if stmt.Type != sqlparser.StatementTypeCreateIndex || stmt.CreateIndex == nil {
		return nil, fmt.Errorf("invalid CREATE INDEX statement")
	}

	// Execute CREATE INDEX on SQLite (SQLite supports CREATE INDEX natively)
	_, err = e.db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to create index in SQLite: %w", err)
	}

	uniqueFlag := ""
	if stmt.CreateIndex.Unique {
		uniqueFlag = "UNIQUE "
	}

	return &ExecuteResult{
		RowCount: 0,
		IsQuery:  false,
		Message:  fmt.Sprintf("%sIndex '%s' created successfully on table '%s'", uniqueFlag, stmt.CreateIndex.IndexName, stmt.CreateIndex.TableName),
	}, nil
}

// executeDropIndex executes a DROP INDEX statement
func (e *Executor) executeDropIndex(query string) (*ExecuteResult, error) {
	// Parse query to get index name
	stmt, err := sqlparser.NewParser().Parse(query)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DROP INDEX: %w", err)
	}

	if stmt.Type != sqlparser.StatementTypeDropIndex || stmt.DropIndex == nil {
		return nil, fmt.Errorf("invalid DROP INDEX statement")
	}

	// Execute DROP INDEX on SQLite (SQLite supports DROP INDEX natively)
	_, err = e.db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to drop index in SQLite: %w", err)
	}

	return &ExecuteResult{
		RowCount: 0,
		IsQuery:  false,
		Message:  fmt.Sprintf("Index '%s' dropped successfully", stmt.DropIndex.IndexName),
	}, nil
}

// executeBeginTransaction executes a BEGIN TRANSACTION statement
func (e *Executor) executeBeginTransaction(query string) (*ExecuteResult, error) {
	// Begin a transaction
	_, err := e.db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &ExecuteResult{
		RowCount: 0,
		IsQuery:  false,
		Message:  "Transaction started successfully",
	}, nil
}

// executeCommit executes a COMMIT statement
func (e *Executor) executeCommit(query string) (*ExecuteResult, error) {
	// Commit the transaction
	_, err := e.db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &ExecuteResult{
		RowCount: 0,
		IsQuery:  false,
		Message:  "Transaction committed successfully",
	}, nil
}

// executeRollback executes a ROLLBACK statement
func (e *Executor) executeRollback(query string) (*ExecuteResult, error) {
	// Rollback the transaction
	_, err := e.db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to rollback transaction: %w", err)
	}

	return &ExecuteResult{
		RowCount: 0,
		IsQuery:  false,
		Message:  "Transaction rolled back successfully",
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
