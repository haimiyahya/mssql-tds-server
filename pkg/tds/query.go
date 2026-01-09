package tds

import (
	"fmt"
	"strings"

	"github.com/factory/mssql-tds-server/pkg/sqlexecutor"
)

// QueryProcessor handles SQL query processing
type QueryProcessor struct {
	executor *sqlexecutor.Executor
}

// NewQueryProcessor creates a new query processor
func NewQueryProcessor() *QueryProcessor {
	return &QueryProcessor{}
}

// SetExecutor sets the SQL executor for the query processor
func (qp *QueryProcessor) SetExecutor(executor *sqlexecutor.Executor) {
	qp.executor = executor
}

// ProcessQuery processes a SQL query and returns results
// This is the main entry point for SQL query execution
func (qp *QueryProcessor) ProcessQuery(query string) ([][]string, error) {
	return qp.ExecuteSQLBatch(query)
}

// ExecuteSQLBatch executes a SQL batch command
func (qp *QueryProcessor) ExecuteSQLBatch(batch string) ([][]string, error) {
	batch = strings.TrimSpace(batch)
	if batch == "" {
		return nil, fmt.Errorf("empty query")
	}

	// Check if executor is set
	if qp.executor == nil {
		return nil, fmt.Errorf("SQL executor not initialized")
	}

	// Execute the query using the SQL executor
	result, err := qp.executor.Execute(batch)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// If it's a query (SELECT), return the rows
	if result.IsQuery {
		return convertResultRows(result), nil
	}

	// If it's a non-query (INSERT, UPDATE, DELETE, DDL), return a message row
	return [][]string{{result.Message}}, nil
}

// convertResultRows converts ExecuteResult rows to [][]string format
func convertResultRows(result *sqlexecutor.ExecuteResult) [][]string {
	if result == nil {
		return [][]string{}
	}

	rows := make([][]string, 0, len(result.Rows)+1)

	// Add column headers as the first row
	if len(result.Columns) > 0 {
		rows = append(rows, result.Columns)
	}

	// Add data rows
	for _, row := range result.Rows {
		stringRow := make([]string, len(row))
		for i, value := range row {
			stringRow[i] = sqlexecutor.ConvertValueToString(value)
		}
		rows = append(rows, stringRow)
	}

	// If no rows, return empty result set
	if len(rows) == 0 {
		rows = append(rows, []string{"No rows found"})
	}

	return rows
}
