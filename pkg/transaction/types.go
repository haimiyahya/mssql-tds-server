package transaction

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
)

// TransactionType represents type of transaction statement
type TransactionType int

const (
	TransactionUnknown TransactionType = iota
	TransactionBegin
	TransactionCommit
	TransactionRollback
)

// Context represents a transaction context
type Context struct {
	mu           sync.RWMutex
	transactions []*sql.Tx
	nestedLevel  int
}

// NewContext creates a new transaction context
func NewContext() *Context {
	return &Context{
		transactions: []*sql.Tx{},
		nestedLevel:  0,
	}
}

// Begin begins a new transaction
func (c *Context) Begin(db *sql.DB) (*sql.Tx, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	c.transactions = append(c.transactions, tx)
	c.nestedLevel++

	return tx, nil
}

// Commit commits the current transaction
func (c *Context) Commit() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.nestedLevel == 0 {
		return fmt.Errorf("no transaction to commit")
	}

	// Get current transaction
	tx := c.transactions[len(c.transactions)-1]

	err := tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Remove from context
	c.transactions = c.transactions[:len(c.transactions)-1]
	c.nestedLevel--

	return nil
}

// Rollback rolls back the current transaction
func (c *Context) Rollback() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.nestedLevel == 0 {
		return fmt.Errorf("no transaction to rollback")
	}

	// Get current transaction
	tx := c.transactions[len(c.transactions)-1]

	err := tx.Rollback()
	if err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	// Remove from context
	c.transactions = c.transactions[:len(c.transactions)-1]
	c.nestedLevel--

	return nil
}

// RollbackAll rolls back all transactions
func (c *Context) RollbackAll() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var lastErr error

	// Rollback all transactions in reverse order
	for i := len(c.transactions) - 1; i >= 0; i-- {
		tx := c.transactions[i]
		err := tx.Rollback()
		if err != nil {
			lastErr = fmt.Errorf("failed to rollback transaction: %w", err)
		}
	}

	// Clear context
	c.transactions = []*sql.Tx{}
	c.nestedLevel = 0

	return lastErr
}

// IsActive checks if a transaction is active
func (c *Context) IsActive() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.nestedLevel > 0
}

// GetLevel returns the current nesting level
func (c *Context) GetLevel() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.nestedLevel
}

// GetCurrentTx returns the current transaction
func (c *Context) GetCurrentTx() *sql.Tx {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.nestedLevel == 0 {
		return nil
	}

	return c.transactions[len(c.transactions)-1]
}

// ParseStatement parses a transaction statement
func ParseStatement(sql string) TransactionType {
	sql = strings.TrimSpace(sql)
	sqlUpper := strings.ToUpper(sql)

	// Check for BEGIN TRANSACTION variations
	if strings.HasPrefix(sqlUpper, "BEGIN TRAN") ||
		strings.HasPrefix(sqlUpper, "BEGIN TRANSACTION") ||
		strings.HasPrefix(sqlUpper, "START TRANSACTION") {
		return TransactionBegin
	}

	// Check for COMMIT
	if strings.HasPrefix(sqlUpper, "COMMIT") ||
		strings.HasPrefix(sqlUpper, "COMMIT TRAN") ||
		strings.HasPrefix(sqlUpper, "COMMIT TRANSACTION") {
		return TransactionCommit
	}

	// Check for ROLLBACK
	if strings.HasPrefix(sqlUpper, "ROLLBACK") ||
		strings.HasPrefix(sqlUpper, "ROLLBACK TRAN") ||
		strings.HasPrefix(sqlUpper, "ROLLBACK TRANSACTION") {
		return TransactionRollback
	}

	return TransactionUnknown
}

// IsTransactionStatement checks if a statement is a transaction statement
func IsTransactionStatement(sql string) bool {
	return ParseStatement(sql) != TransactionUnknown
}

// DetectTransactionUsage detects if SQL contains transaction statements
func DetectTransactionUsage(sql string) bool {
	sqlUpper := strings.ToUpper(sql)
	return strings.Contains(sqlUpper, "BEGIN TRAN") ||
		strings.Contains(sqlUpper, "BEGIN TRANSACTION") ||
		strings.Contains(sqlUpper, "START TRANSACTION") ||
		strings.Contains(sqlUpper, "COMMIT") ||
		strings.Contains(sqlUpper, "ROLLBACK")
}
