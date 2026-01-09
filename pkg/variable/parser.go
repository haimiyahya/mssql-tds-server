package variable

import (
	"fmt"
	"regexp"
	"strings"
)

// StatementType represents the type of SQL statement
type StatementType int

const (
	StatementUnknown StatementType = iota
	StatementDeclare
	StatementSet
	StatementSelectAssignment
	StatementQuery
)

// ParseStatement determines the type of statement
func ParseStatement(sql string) StatementType {
	sql = strings.TrimSpace(sql)
	sqlUpper := strings.ToUpper(sql)

	// Check for DECLARE
	if strings.HasPrefix(sqlUpper, "DECLARE") {
		return StatementDeclare
	}

	// Check for SET
	if strings.HasPrefix(sqlUpper, "SET ") {
		return StatementSet
	}

	// Check for SELECT with assignment
	// Format: SELECT @var = expression
	re := regexp.MustCompile(`SELECT\s+@\w+\s*=\s+`)
	if re.MatchString(sqlUpper) {
		return StatementSelectAssignment
	}

	// Default to query
	return StatementQuery
}

// ParseProcedureBody parses a procedure body and extracts individual statements
func ParseProcedureBody(body string) ([]string, error) {
	// Split by semicolon, but handle string literals and nested structures
	statements := []string{}
	current := ""
	inQuotes := false
	inParentheses := 0

	for _, ch := range body {
		switch ch {
		case '\'':
			inQuotes = !inQuotes
			current += string(ch)
		case '(':
			if !inQuotes {
				inParentheses++
			}
			current += string(ch)
		case ')':
			if !inQuotes {
				inParentheses--
			}
			current += string(ch)
		case ';':
			if !inQuotes && inParentheses == 0 {
				// End of statement
				stmt := strings.TrimSpace(current)
				if stmt != "" {
					statements = append(statements, stmt)
				}
				current = ""
			} else {
				current += string(ch)
			}
		default:
			current += string(ch)
		}
	}

	// Add last statement if exists
	if strings.TrimSpace(current) != "" {
		statements = append(statements, strings.TrimSpace(current))
	}

	return statements, nil
}

// FindVariableReferences finds all variable references in a SQL statement
func FindVariableReferences(sql string) []string {
	// Find all @variable references
	re := regexp.MustCompile(`@\w+`)
	matches := re.FindAllString(sql, -1)

	// Deduplicate
	seen := make(map[string]bool)
	variables := []string{}

	for _, match := range matches {
		if !seen[match] {
			seen[match] = true
			variables = append(variables, match)
		}
	}

	return variables
}

// ReplaceVariables replaces variable references with their values in SQL
func ReplaceVariables(sql string, context *Context) (string, error) {
	// Find all variable references
	references := FindVariableReferences(sql)

	result := sql
	for _, ref := range references {
		// Get variable from context
		variable, exists := context.Get(ref)
		if !exists {
			return "", fmt.Errorf("variable '%s' not declared", ref)
		}

		// Skip if variable is NULL (we handle this in SQL generation)
		if variable.IsNull {
			return "", fmt.Errorf("variable '%s' is NULL", ref)
		}

		// Format value based on type
		formattedValue := FormatValue(variable.Value, variable.Type)

		// Replace in SQL (case-sensitive replacement to avoid partial matches)
		result = strings.ReplaceAll(result, ref, formattedValue)
	}

	return result, nil
}
