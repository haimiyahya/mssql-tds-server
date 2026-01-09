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
	// Normalize SQL
	sql = strings.TrimSpace(sql)
	sql = strings.ToUpper(sql)

	// Check if it starts with CREATE PROCEDURE or CREATE PROC
	if !strings.HasPrefix(sql, "CREATE PROCEDURE") && !strings.HasPrefix(sql, "CREATE PROC") {
		return nil, fmt.Errorf("not a CREATE PROCEDURE statement")
	}

	// Extract procedure name and body using regex
	// Pattern: CREATE PROC[EDURE] proc_name (@params) AS body
	re := regexp.MustCompile(`CREATE\s+PROC(?:EDURE)?\s+(\w+)\s*\((.*)\)\s+AS\s+(.*)$`)
	matches := re.FindStringSubmatch(sql)

	if len(matches) != 4 {
		return nil, fmt.Errorf("invalid CREATE PROCEDURE syntax: %s", sql)
	}

	procName := matches[1]
	paramsStr := strings.TrimSpace(matches[2])
	bodyStr := strings.TrimSpace(matches[3])

	// Parse parameters
	params, err := parseParameters(paramsStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	return &Procedure{
		Name:       procName,
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
