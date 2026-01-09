package controlflow

import (
	"fmt"
	"regexp"
	"strings"
)

// StatementType represents type of statement (extends variable.StatementType)
type StatementType int

const (
	StatementUnknown StatementType = iota
	StatementQuery
	StatementDeclare
	StatementSet
	StatementSelectAssignment
	StatementIF
	StatementELSE
	StatementEND
)

// ParseStatement determines the type of SQL statement
func ParseStatement(sql string) StatementType {
	sql = strings.TrimSpace(sql)
	sqlUpper := strings.ToUpper(sql)

	// Check for IF
	if strings.HasPrefix(sqlUpper, "IF") {
		return StatementIF
	}

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

// ParseIFBlock parses an IF statement block including nested statements
// Format: IF condition THEN statements [ELSE statements] END
func ParseIFBlock(sql string) (*Block, []string, error) {
	sql = strings.TrimSpace(sql)
	sqlUpper := strings.ToUpper(sql)

	// Must start with IF
	if !strings.HasPrefix(sqlUpper, "IF") {
		return nil, nil, fmt.Errorf("not an IF block")
	}

	// Find THEN keyword
	thenPos := strings.Index(sqlUpper, "THEN")
	if thenPos < 0 {
		return nil, nil, fmt.Errorf("IF block missing THEN keyword")
	}

	// Extract condition
	conditionStr := strings.TrimSpace(sql[2:thenPos])

	// Parse condition
	condition, err := ParseCondition(conditionStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse IF condition: %w", err)
	}

	// Find END keyword (outermost)
	endPos := findOutermostEND(sqlUpper)
	if endPos < 0 {
		return nil, nil, fmt.Errorf("IF block missing END keyword")
	}

	// Extract body (between THEN and END)
	bodyStr := strings.TrimSpace(sql[thenPos+4:endPos])

	// Check for ELSE
	elsePos := findELSEInBody(bodyStr)

	var ifBlock *Block
	var elseBlock *Block

	if elsePos >= 0 {
		// Split into IF and ELSE blocks
		ifBody := strings.TrimSpace(bodyStr[:elsePos])
		elseBody := strings.TrimSpace(bodyStr[elsePos+4:])

		ifBlock = &Block{
			Type:     BlockIf,
			Condition: condition,
			Body:     []string{ifBody},
		}

		elseBlock = &Block{
			Type:     BlockElse,
			Condition: nil,
			Body:     []string{elseBody},
		}

		// Return both blocks
		return ifBlock, []string{"ELSE", elseBody}, nil
	}

	// Only IF block
	ifBlock = &Block{
		Type:     BlockIf,
		Condition: condition,
		Body:     []string{bodyStr},
	}

	return ifBlock, nil, nil
}

// findOutermostEND finds the outermost END keyword (not inside IF blocks)
func findOutermostEND(sql string) int {
	depth := 0
	inQuotes := false

	for i := 0; i < len(sql); i++ {
		ch := sql[i]
		switch ch {
		case '\'':
			inQuotes = !inQuotes
		case '(':
			if !inQuotes {
				depth++
			}
		case ')':
			if !inQuotes {
				depth--
			}
		}

		if !inQuotes && depth == 0 {
			// Check for END keyword
			if strings.ToUpper(sql[i:]) == "END" || (i+3 < len(sql) && strings.ToUpper(sql[i:i+3]) == "END ") {
				return i
			}
		}
	}

	return -1
}

// findELSEInBody finds ELSE keyword in body (not inside nested IF blocks)
func findELSEInBody(body string) int {
	ifDepth := 0
	inQuotes := false

	for i := 0; i < len(body); i++ {
		ch := body[i]
		switch ch {
		case '\'':
			inQuotes = !inQuotes
		case '(':
			if !inQuotes {
				ifDepth++
			}
		case ')':
			if !inQuotes {
				ifDepth--
			}
		}

		if !inQuotes && ifDepth == 0 {
			// Check for ELSE keyword
			if strings.ToUpper(body[i:]) == "ELSE" || (i+4 < len(body) && strings.ToUpper(body[i:i+4]) == "ELSE ") {
				return i
			}
		}
	}

	return -1
}

// SplitStatements splits procedure body into individual statements
// Handles semicolon, IF blocks, etc.
func SplitStatements(body string) ([]string, error) {
	statements := []string{}
	current := ""
	inQuotes := false
	inParentheses := 0
	inIFBlock := 0

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
			if !inQuotes && inParentheses == 0 && inIFBlock == 0 {
				// End of statement
				stmt := strings.TrimSpace(current)
				if stmt != "" {
					statements = append(statements, stmt)
				}
				current = ""
			} else {
				current += string(ch)
			}
		case ' ', '\t', '\n', '\r':
			// Track IF/END keywords
			currentUpper := strings.ToUpper(current)
			if strings.HasSuffix(currentUpper, "IF") {
				// Check if it's IF (not part of another word)
				if len(current) == 2 || (len(current) > 2 && (current[len(current)-3] == ' ' || current[len(current)-3] == '(')) {
					inIFBlock++
				}
			}
			if strings.HasSuffix(currentUpper, "END") {
				// Check if it's END (not part of another word)
				if len(current) == 3 || (len(current) > 3 && (current[len(current)-4] == ' ' || current[len(current)-4] == ';')) {
					if inIFBlock > 0 {
						inIFBlock--
					}
				}
			}
			current += string(ch)
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
