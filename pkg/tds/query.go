package tds

import (
	"fmt"
	"strings"
)

// QueryProcessor handles simple query processing
type QueryProcessor struct{}

// NewQueryProcessor creates a new query processor
func NewQueryProcessor() *QueryProcessor {
	return &QueryProcessor{}
}

// ProcessQuery processes a SQL query and returns results
// For Phase 2, this implements a simple echo functionality that returns
// the query in uppercase as a result set
func (qp *QueryProcessor) ProcessQuery(query string) ([][]string, error) {
	// For Phase 2, implement simple echo functionality
	// Convert the query to uppercase and return as a result set
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("empty query")
	}

	// Echo functionality: return the query in uppercase
	result := [][]string{
		{strings.ToUpper(query)},
	}

	return result, nil
}

// ExecuteSQLBatch executes a SQL batch command
func (qp *QueryProcessor) ExecuteSQLBatch(batch string) ([][]string, error) {
	return qp.ProcessQuery(batch)
}
