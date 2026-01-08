package tds

import (
	"fmt"
	"strings"
)

// StoredProcedureHandler manages stored procedure execution
type StoredProcedureHandler struct{}

// NewStoredProcedureHandler creates a new stored procedure handler
func NewStoredProcedureHandler() *StoredProcedureHandler {
	return &StoredProcedureHandler{}
}

// Execute executes a stored procedure with given parameters and returns results
func (h *StoredProcedureHandler) Execute(procName string, params []*RPCParameter) ([][]string, error) {
	// Normalize procedure name
	procName = strings.ToUpper(strings.TrimSpace(procName))

	// Handle different stored procedures
	switch procName {
	case "SP_HELLO", "HELLO_WORLD", "TEST_PROC":
		return h.executeHelloWorld(params)
	case "SP_ECHO", "ECHO_PROC":
		return h.executeEcho(params)
	case "SP_GET_DATA", "GET_DATA":
		return h.executeGetData(params)
	default:
		return nil, fmt.Errorf("stored procedure '%s' not found", procName)
	}
}

// executeHelloWorld is a simple test stored procedure
func (h *StoredProcedureHandler) executeHelloWorld(params []*RPCParameter) ([][]string, error) {
	var greeting string
	if len(params) > 0 {
		if name, ok := params[0].Value.(string); ok {
			greeting = fmt.Sprintf("Hello, %s!", name)
		} else {
			greeting = "Hello, World!"
		}
	} else {
		greeting = "Hello, World!"
	}

	return [][]string{
		{greeting},
	}, nil
}

// executeEcho echoes back the input parameters
func (h *StoredProcedureHandler) executeEcho(params []*RPCParameter) ([][]string, error) {
	if len(params) == 0 {
		return [][]string{
			{"No parameters provided"},
		}, nil
	}

	results := make([][]string, 0)
	for i, param := range params {
		var value string
		switch v := param.Value.(type) {
		case string:
			value = fmt.Sprintf("Param %d: %s", i+1, strings.ToUpper(v))
		case int32:
			value = fmt.Sprintf("Param %d: %d (INT)", i+1, v)
		case int64:
			value = fmt.Sprintf("Param %d: %d (BIGINT)", i+1, v)
		case int16:
			value = fmt.Sprintf("Param %d: %d (SMALLINT)", i+1, v)
		default:
			value = fmt.Sprintf("Param %d: %v (UNKNOWN TYPE)", i+1, v)
		}
		results = append(results, []string{value})
	}

	return results, nil
}

// executeGetData returns sample data for testing
func (h *StoredProcedureHandler) executeGetData(params []*RPCParameter) ([][]string, error) {
	// Return sample data set
	results := [][]string{
		{"1", "John", "Doe", "Engineering"},
		{"2", "Jane", "Smith", "Marketing"},
		{"3", "Bob", "Johnson", "Sales"},
		{"4", "Alice", "Williams", "HR"},
	}

	// Filter results if parameters are provided
	if len(params) > 0 {
		deptFilter := ""
		if dept, ok := params[0].Value.(string); ok {
			deptFilter = strings.ToUpper(dept)
		}

		if deptFilter != "" {
			filtered := make([][]string, 0)
			for _, row := range results {
				if len(row) >= 4 && strings.Contains(strings.ToUpper(row[3]), deptFilter) {
					filtered = append(filtered, row)
				}
			}
			results = filtered
		}
	}

	return results, nil
}

// GetProcedureInfo returns information about a stored procedure
func (h *StoredProcedureHandler) GetProcedureInfo(procName string) (string, []string) {
	procName = strings.ToUpper(strings.TrimSpace(procName))

	switch procName {
	case "SP_HELLO", "HELLO_WORLD", "TEST_PROC":
		return "SP_HELLO", []string{"Greeting (VARCHAR)"}
	case "SP_ECHO", "ECHO_PROC":
		return "SP_ECHO", []string{"Echo result (VARCHAR)"}
	case "SP_GET_DATA", "GET_DATA":
		return "SP_GET_DATA", []string{"ID (INT)", "FirstName (VARCHAR)", "LastName (VARCHAR)", "Department (VARCHAR)"}
	default:
		return "", []string{}
	}
}
