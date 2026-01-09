package controlflow

import (
	"fmt"
	"strings"
)

// Operator represents comparison operators
type Operator string

const (
	OpEqual      Operator = "="
	OpNotEqual   Operator = "<>"
	OpLess       Operator = "<"
	OpGreater    Operator = ">"
	OpLessEqual  Operator = "<="
	OpGreaterEqual Operator = ">="
)

// LogicalOperator represents logical operators
type LogicalOperator string

const (
	LogOpAnd LogicalOperator = "AND"
	LogOpOr  LogicalOperator = "OR"
)

// Condition represents a single condition
type Condition struct {
	Left     interface{}
	Operator  Operator
	Right    interface{}
}

// LogicalCondition represents conditions with AND/OR
type LogicalCondition struct {
	Left      interface{}
	Operator  LogicalOperator
	Right     interface{}
}

// BlockType represents type of control block
type BlockType int

const (
	BlockIf BlockType = iota
	BlockElse
	BlockEnd
)

// Block represents a control flow block
type Block struct {
	Type     BlockType
	Condition *LogicalCondition
	Body     []string
}

// ParseCondition parses a condition expression
// Format: expr operator expr [LOGICAL_OPERATOR expr operator expr ...]
func ParseCondition(expr string) (*LogicalCondition, error) {
	expr = strings.TrimSpace(expr)

	// Try to find AND/OR operators (outermost only)
	andOrPos := findOutermostLogicalOp(expr, "AND", "OR")
	if andOrPos >= 0 {
		// Get operator
		parts := strings.Fields(expr[andOrPos:andOrPos+5])
		if len(parts) == 0 {
			return nil, fmt.Errorf("invalid logical operator at position %d", andOrPos)
		}
		logOp := LogicalOperator(strings.ToUpper(parts[0]))

		// Recursively parse left and right
		leftExpr := strings.TrimSpace(expr[:andOrPos])
		rightExpr := strings.TrimSpace(expr[andOrPos+len(logOp)+1:])

		left, err := ParseCondition(leftExpr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse left condition: %w", err)
		}

		right, err := ParseCondition(rightExpr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse right condition: %w", err)
		}

		return &LogicalCondition{
			Left:     left,
			Operator:  logOp,
			Right:    right,
		}, nil
	}

	// Simple condition without AND/OR
	return parseSimpleCondition(expr)
}

// parseSimpleCondition parses a simple condition (no AND/OR)
func parseSimpleCondition(expr string) (*LogicalCondition, error) {
	expr = strings.TrimSpace(expr)

	// Try each operator
	operators := []Operator{
		OpGreaterEqual,
		OpLessEqual,
		OpNotEqual,
		OpEqual,
		OpLess,
		OpGreater,
	}

	for _, op := range operators {
		pos := findOperatorPos(expr, string(op))
		if pos >= 0 {
			left := strings.TrimSpace(expr[:pos])
			right := strings.TrimSpace(expr[pos+len(op):])

			return &LogicalCondition{
				Left: &Condition{
					Left:    left,
					Operator: op,
					Right:   right,
				},
				Operator: "",
				Right:    nil,
			}, nil
		}
	}

	return nil, fmt.Errorf("no comparison operator found in: %s", expr)
}

// findOperatorPos finds the position of an operator, ignoring quotes
func findOperatorPos(expr, op string) int {
	inQuotes := false
	i := 0

	for i < len(expr) {
		ch := expr[i]
		switch ch {
		case '\'':
			inQuotes = !inQuotes
		case ' ', '\t':
			i++
			continue
		}

		if !inQuotes && i+len(op) <= len(expr) {
			if strings.ToUpper(expr[i:i+len(op)]) == strings.ToUpper(op) {
				return i
			}
		}

		i++
	}

	return -1
}

// findOutermostLogicalOp finds the outermost AND/OR operator (not inside parentheses)
func findOutermostLogicalOp(expr, op1, op2 string) int {
	expr = strings.TrimSpace(expr)
	depth := 0
	inQuotes := false

	for i := 0; i < len(expr); i++ {
		ch := expr[i]
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
			// Check for OR (should be checked before AND)
			if strings.ToUpper(expr[i:]) == "OR" || (i+2 < len(expr) && strings.ToUpper(expr[i:i+3]) == "OR ") {
				return i
			}
			// Check for AND
			if strings.ToUpper(expr[i:]) == "AND" || (i+3 < len(expr) && strings.ToUpper(expr[i:i+4]) == "AND ") {
				return i
			}
		}
	}

	return -1
}

// Evaluate evaluates a condition with given variable values
func Evaluate(condition *LogicalCondition, variables map[string]interface{}) (bool, error) {
	if condition == nil {
		return true, nil
	}

	// Logical condition (AND/OR)
	if condition.Operator != "" {
		left, err := Evaluate(getLogicalOperand(condition.Left), variables)
		if err != nil {
			return false, err
		}

		right, err := Evaluate(getLogicalOperand(condition.Right), variables)
		if err != nil {
			return false, err
		}

		switch condition.Operator {
		case LogOpAnd:
			return left && right, nil
		case LogOpOr:
			return left || right, nil
		default:
			return false, fmt.Errorf("unknown logical operator: %s", condition.Operator)
		}
	}

	// Simple condition
	cond, ok := condition.Left.(*Condition)
	if !ok {
		return false, fmt.Errorf("invalid condition type")
	}

	return evaluateSimpleCondition(cond, variables)
}

// getLogicalOperand extracts the condition from a logical operand
func getLogicalOperand(operand interface{}) *LogicalCondition {
	switch v := operand.(type) {
	case *LogicalCondition:
		return v
	case *Condition:
		return &LogicalCondition{
			Left:     v,
			Operator:  "",
			Right:    nil,
		}
	default:
		return nil
	}
}

// evaluateSimpleCondition evaluates a simple condition
func evaluateSimpleCondition(cond *Condition, variables map[string]interface{}) (bool, error) {
	// Get left value
	leftValue, err := getExpressionValue(cond.Left, variables)
	if err != nil {
		return false, err
	}

	// Get right value
	rightValue, err := getExpressionValue(cond.Right, variables)
	if err != nil {
		return false, err
	}

	// Compare based on operator
	switch cond.Operator {
	case OpEqual:
		return compareValues(leftValue, rightValue) == 0, nil
	case OpNotEqual:
		return compareValues(leftValue, rightValue) != 0, nil
	case OpLess:
		return compareValues(leftValue, rightValue) < 0, nil
	case OpGreater:
		return compareValues(leftValue, rightValue) > 0, nil
	case OpLessEqual:
		return compareValues(leftValue, rightValue) <= 0, nil
	case OpGreaterEqual:
		return compareValues(leftValue, rightValue) >= 0, nil
	default:
		return false, fmt.Errorf("unknown comparison operator: %s", cond.Operator)
	}
}

// getExpressionValue gets the value of an expression (variable or literal)
func getExpressionValue(expr interface{}, variables map[string]interface{}) (interface{}, error) {
	switch v := expr.(type) {
	case string:
		// Check if it's a variable reference
		if strings.HasPrefix(v, "@") {
			value, exists := variables[v]
			if !exists {
				return nil, fmt.Errorf("variable '%s' not found", v)
			}
			return value, nil
		}
		// String literal (remove quotes)
		if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
			return strings.ReplaceAll(v[1:len(v)-1], "''", "'"), nil
		}
		return v, nil
	default:
		return v, nil
	}
}

// compareValues compares two values
// Returns -1 if a < b, 0 if a == b, 1 if a > b
func compareValues(a, b interface{}) int {
	// Handle nil values
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// Try numeric comparison
	switch va := a.(type) {
	case int, int32, int64, uint, uint32, uint64:
		switch vb := b.(type) {
		case int, int32, int64, uint, uint32, uint64:
			// Convert to int64 for comparison
			var ia, ib int64
			fmt.Sscanf(fmt.Sprintf("%v", va), "%d", &ia)
			fmt.Sscanf(fmt.Sprintf("%v", vb), "%d", &ib)
			if ia < ib {
				return -1
			}
			if ia > ib {
				return 1
			}
			return 0
		}
	case float32, float64:
		switch vb := b.(type) {
		case float32, float64:
			fa := fmt.Sprintf("%v", va)
			fb := fmt.Sprintf("%v", vb)
			var faVal, fbVal float64
			fmt.Sscanf(fa, "%f", &faVal)
			fmt.Sscanf(fb, "%f", &fbVal)
			if faVal < fbVal {
				return -1
			}
			if faVal > fbVal {
				return 1
			}
			return 0
		}
	}

	// String comparison
	sa := fmt.Sprintf("%v", a)
	sb := fmt.Sprintf("%v", b)
	if sa < sb {
		return -1
	} else if sa > sb {
		return 1
	}
	return 0
}

// ParseIFStatement parses an IF statement
// Format: IF condition THEN/END [ELSE END]
func ParseIFStatement(sql string) (*Block, []string, error) {
	sql = strings.TrimSpace(sql)
	sqlUpper := strings.ToUpper(sql)

	// Must start with IF
	if !strings.HasPrefix(sqlUpper, "IF") {
		return nil, nil, fmt.Errorf("not an IF statement")
	}

	// Remove IF keyword
	afterIF := strings.TrimSpace(sql[2:])

	// Find THEN keyword
	thenPos := strings.Index(strings.ToUpper(afterIF), "THEN")
	if thenPos < 0 {
		return nil, nil, fmt.Errorf("IF statement missing THEN keyword")
	}

	// Extract condition
	conditionStr := strings.TrimSpace(afterIF[:thenPos])

	// Parse condition
	condition, err := ParseCondition(conditionStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse IF condition: %w", err)
	}

	// Extract body (after THEN)
	bodyStr := strings.TrimSpace(afterIF[thenPos+4:])

	// Check for ELSE
	elsePos := findELSEKeyword(bodyStr)

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

// findELSEKeyword finds the ELSE keyword (not inside quotes or sub-statements)
func findELSEKeyword(sql string) int {
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
			// Check for ELSE keyword
			if strings.ToUpper(sql[i:]) == "ELSE" || (i+4 < len(sql) && strings.ToUpper(sql[i:i+4]) == "ELSE ") {
				return i
			}
		}
	}

	return -1
}
