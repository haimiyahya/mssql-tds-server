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
	StatementWHILE
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

	// Check for WHILE
	if strings.HasPrefix(sqlUpper, "WHILE") {
		return StatementWHILE
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

	if elsePos >= 0 {
		// Split into IF and ELSE blocks
		ifBody := strings.TrimSpace(bodyStr[:elsePos])
		elseBody := strings.TrimSpace(bodyStr[elsePos+4:])

		ifBlock = &Block{
			Type:     BlockIf,
			Condition: condition,
			Body:     []string{ifBody},
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
				// Return position of separator before END (space or beginning)
				if i > 0 && (sql[i-1] == ' ' || sql[i-1] == '\t') {
					return i - 1
				}
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
			if strings.ToUpper(body[i:]) == "ELSE" || (i+4 < len(body) && strings.ToUpper(body[i:i+4]) == "ELSE") {
				// Return position of separator before ELSE (space or beginning)
				if i > 0 && (body[i-1] == ' ' || body[i-1] == '\t') {
					return i - 1
				}
				return i
			}
		}
	}

	return -1
}

// ParseWHILEBlock parses a WHILE statement block
// Format: WHILE condition statements END
func ParseWHILEBlock(sql string) (*Block, error) {
	sql = strings.TrimSpace(sql)
	sqlUpper := strings.ToUpper(sql)

	// Must start with WHILE
	if !strings.HasPrefix(sqlUpper, "WHILE") {
		return nil, fmt.Errorf("not a WHILE block")
	}

	// Find statements after WHILE keyword
	afterWHILE := strings.TrimSpace(sql[5:])

	// Find END keyword (outermost)
	endPos := findOutermostEND(afterWHILE)
	if endPos < 0 {
		return nil, fmt.Errorf("WHILE block missing END keyword")
	}

	bodyStr := strings.TrimSpace(afterWHILE[:endPos])

	// Parse condition from beginning of body
	// Format: condition statements...
	// We need to split condition from statements
	// For simplicity, assume first part is condition until we see a statement keyword
	conditionStr := extractConditionFromBody(bodyStr)
	body := bodyStr[len(conditionStr):]

	// Parse condition
	condition, err := ParseCondition(conditionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse WHILE condition: %w", err)
	}

	return &Block{
		Type:     BlockWhile,
		Condition: condition,
		Body:     []string{body},
	}, nil
}

// extractConditionFromBody extracts condition from WHILE body
// Format: condition statements...
func extractConditionFromBody(body string) string {
	// Simple approach: find first keyword that indicates start of statements
	// Keywords: SELECT, INSERT, UPDATE, DELETE, IF, WHILE, DECLARE, SET, PRINT
	body = strings.TrimSpace(body)
	bodyUpper := strings.ToUpper(body)

	// Check if body starts with a keyword
	keywords := []string{"SELECT ", "INSERT ", "UPDATE ", "DELETE ", "IF ", "WHILE ", "DECLARE ", "SET ", "PRINT ", "BREAK", "CONTINUE"}

	for _, keyword := range keywords {
		if strings.HasPrefix(bodyUpper, keyword) {
			// Keyword found at start, condition is empty
			return ""
		}
	}

	// For simple WHILE loops, the condition is usually a single expression
	// We can parse it by finding the first statement keyword
	// But this is complex without a full SQL parser

	// Simple heuristic: split by first statement keyword
	for i := 0; i < len(body); i++ {
		checkUpper := strings.ToUpper(body[i:])
		for _, keyword := range keywords {
			if strings.HasPrefix(checkUpper, keyword) {
				return strings.TrimSpace(body[:i])
			}
		}
	}

	// No keyword found, entire body is condition
	return ""
}

// BlockWhile represents a WHILE block type
const BlockWhile BlockType = 5

// SplitStatements splits procedure body into individual statements
// Handles semicolon, IF blocks, WHILE blocks, etc.
func SplitStatements(body string) ([]string, error) {
	statements := []string{}
	current := ""
	inQuotes := false
	inParentheses := 0
	inIFBlock := 0
	inWHILEBlock := 0

	for i, ch := range body {
		// Handle escaped quotes
		if ch == '\'' && i > 0 && body[i-1] == '\\' {
			// It's an escaped quote, don't toggle inQuotes
			current += string(ch)
			continue
		}

		// Check for END followed by semicolon before processing whitespace
		if strings.HasSuffix(strings.ToUpper(current), "END") && ch == ';' && inQuotes == false && inParentheses == 0 {
			// END block ends with semicolon - split the END statement, then add semicolon to next statement
			stmt := strings.TrimSpace(strings.TrimSuffix(current, "END;"))
			stmt += "END;"
			if stmt != "" {
				statements = append(statements, stmt)
			}
			current = ""
			continue
		}

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
			if !inQuotes && inParentheses == 0 && inIFBlock == 0 && inWHILEBlock == 0 {
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
			// Track IF/WHILE/END keywords
			currentUpper := strings.ToUpper(current)
			if strings.HasSuffix(currentUpper, "IF") {
				// Check if it's IF (not part of another word)
				if len(current) == 2 || (len(current) > 2 && (current[len(current)-3] == ' ' || current[len(current)-3] == '(')) {
					inIFBlock++
				}
			}
			if strings.HasSuffix(currentUpper, "WHILE") {
				// Check if it's WHILE (not part of another word)
				if len(current) == 5 || (len(current) > 5 && (current[len(current)-6] == ' ' || current[len(current)-6] == '(')) {
					inWHILEBlock++
				}
			}
			if strings.HasSuffix(currentUpper, "END") {
				// Check if it's END (not part of another word)
				if len(current) == 3 || (len(current) > 3 && (current[len(current)-4] == ' ' || current[len(current)-4] == ';')) {
					if inIFBlock > 0 {
						inIFBlock--
					}
					if inWHILEBlock > 0 {
						inWHILEBlock--
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
