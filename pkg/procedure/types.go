package procedure

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Parameter represents a stored procedure parameter
type Parameter struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Length     int    `json:"length,omitempty"`
	HasDefault bool   `json:"has_default"`
	Default    string `json:"default,omitempty"`
}

// Procedure represents a stored procedure definition
type Procedure struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Body       string     `json:"body"`
	Parameters []Parameter `json:"parameters"`
	CreatedAt  string     `json:"created_at"`
}

// ParseCreateProcedure parses a CREATE PROCEDURE statement
func ParseCreateProcedure(sql string) (*Procedure, error) {
	// Store original SQL string for body extraction
	originalSQL := strings.TrimSpace(sql)
	
	// Normalize SQL for parsing (convert to uppercase for keyword matching)
	normalizedSQL := strings.ToUpper(originalSQL)

	// Check if it starts with CREATE PROCEDURE or CREATE PROC
	if !strings.HasPrefix(normalizedSQL, "CREATE PROCEDURE") && !strings.HasPrefix(normalizedSQL, "CREATE PROC") {
		return nil, fmt.Errorf("not a CREATE PROCEDURE statement")
	}

	// Find AS keyword to separate name+params from body
	asPos := strings.Index(normalizedSQL, " AS ")
	if asPos < 0 {
		return nil, fmt.Errorf("invalid CREATE PROCEDURE syntax: missing AS keyword")
	}

	beforeAS := strings.TrimSpace(normalizedSQL[:asPos])
	
	// Extract body from original SQL string to preserve case
	bodyStr := strings.TrimSpace(originalSQL[asPos+4:])

	// Parse name and parameters from beforeAS
	// Format 1: CREATE PROCEDURE name (@params) AS body
	// Format 2: CREATE PROCEDURE name @param1 type, @param2 type AS body
	// Format 3: CREATE PROC name @param1 type AS body
	name := ""
	paramsStr := ""

	// Remove "CREATE PROCEDURE" or "CREATE PROC" prefix
	withoutPrefix := beforeAS
	if strings.HasPrefix(beforeAS, "CREATE PROCEDURE ") {
		withoutPrefix = strings.TrimPrefix(beforeAS, "CREATE PROCEDURE ")
	} else if strings.HasPrefix(beforeAS, "CREATE PROC ") {
		withoutPrefix = strings.TrimPrefix(beforeAS, "CREATE PROC ")
	} else {
		// Fallback - remove any whitespace after keyword
		withoutPrefix = strings.TrimPrefix(beforeAS, "CREATE PROCEDURE")
		withoutPrefix = strings.TrimPrefix(withoutPrefix, "CREATE PROC")
		withoutPrefix = strings.TrimSpace(withoutPrefix)
	}
	
	// Check if format has parentheses (for parameter list): name (...)
	// The name should be followed immediately by space and (
	openParenAfterNamePos := strings.Index(withoutPrefix, " (")
	if openParenAfterNamePos > 0 {
		// Format with parentheses: name (@params)
		name = strings.TrimSpace(withoutPrefix[:openParenAfterNamePos])
		// Extract params between ( and )
		openParen := strings.Index(withoutPrefix[openParenAfterNamePos:], "(")
		closeParen := strings.LastIndex(withoutPrefix, ")")
		if openParen >= 0 && closeParen >= 0 {
			paramsStr = strings.TrimSpace(withoutPrefix[openParenAfterNamePos+openParen+1 : closeParen])
		}
	} else {
		// Format without parentheses: name @param1 type, @param2 type
		// Find first @ to get name and params
		atPos := strings.Index(withoutPrefix, "@")
		if atPos >= 0 {
			name = strings.TrimSpace(withoutPrefix[:atPos])
			paramsStr = strings.TrimSpace(withoutPrefix[atPos:])
		} else {
			name = strings.TrimSpace(withoutPrefix)
		}
	}

	// Parse parameters
	params, err := parseParameters(paramsStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	return &Procedure{
		Name:       name,
		Body:       bodyStr,
		Parameters: params,
	}, nil
}

// parseParameters parses the parameter list from CREATE PROCEDURE
func parseParameters(paramsStr string) ([]Parameter, error) {
	params := []Parameter{}

	if paramsStr == "" {
		return params, nil
	}

	// Split by comma (but not inside parentheses)
	paramParts := splitByCommaOutsideParentheses(paramsStr)

	for _, paramPart := range paramParts {
		paramPart = strings.TrimSpace(paramPart)
		if paramPart == "" {
			continue
		}

		// Parse individual parameter
		// Format: @param TYPE [DEFAULT value]
		paramRegex := regexp.MustCompile(`@(\w+)\s+(\w+)(?:\((\d+)\))?\s*(?:DEFAULT\s+(.+))?`)
		paramMatches := paramRegex.FindStringSubmatch(paramPart)

		if len(paramMatches) < 3 {
			return nil, fmt.Errorf("invalid parameter syntax: %s", paramPart)
		}

		param := Parameter{
			Name: paramMatches[1],
			Type: paramMatches[2],
		}

		if len(paramMatches) > 3 && paramMatches[3] != "" {
			// Parse length
			fmt.Sscanf(paramMatches[3], "%d", &param.Length)
		}

		if len(paramMatches) > 4 && paramMatches[4] != "" {
			param.HasDefault = true
			param.Default = paramMatches[4]
		}

		params = append(params, param)
	}

	return params, nil
}

// splitByCommaOutsideParentheses splits a string by commas, but not inside parentheses
func splitByCommaOutsideParentheses(s string) []string {
	var result []string
	var current strings.Builder
	depth := 0

	for _, ch := range s {
		switch ch {
		case '(':
			depth++
			current.WriteRune(ch)
		case ')':
			depth--
			current.WriteRune(ch)
		case ',':
			if depth == 0 {
				result = append(result, current.String())
				current.Reset()
			} else {
				current.WriteRune(ch)
			}
		default:
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

// ParametersToJSON converts parameters to JSON for storage
func ParametersToJSON(params []Parameter) (string, error) {
	if params == nil {
		return "[]", nil
	}

	data, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("failed to marshal parameters: %w", err)
	}
	return string(data), nil
}

// ParametersFromJSON converts JSON to parameters
func ParametersFromJSON(jsonStr string) ([]Parameter, error) {
	if jsonStr == "" || jsonStr == "[]" {
		return []Parameter{}, nil
	}

	var params []Parameter
	err := json.Unmarshal([]byte(jsonStr), &params)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal parameters: %w", err)
	}
	return params, nil
}
