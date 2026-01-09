package auth

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	// SQL Server authentication types
	AuthTypeSQLServer = "SQL_SERVER" // SQL Server Authentication
	AuthTypeWindows   = "WINDOWS"    // Windows Authentication
	AuthTypeMixed     = "MIXED"      // Mixed Authentication
)

// AuthManager handles user authentication and login management
type AuthManager struct {
	masterDB *sql.DB // Master database connection
}

// NewAuthManager creates a new authentication manager
func NewAuthManager(masterDB *sql.DB) (*AuthManager, error) {
	// Ensure syslogins table exists
	err := createSysLoginsTable(masterDB)
	if err != nil {
		return nil, fmt.Errorf("failed to create syslogins table: %w", err)
	}

	return &AuthManager{
		masterDB: masterDB,
	}, nil
}

// createSysLoginsTable creates the syslogins table
func createSysLoginsTable(db *sql.DB) error {
	tableSQL := `
	CREATE TABLE IF NOT EXISTS syslogins (
		sid INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		type TEXT DEFAULT 'SQL_SERVER',
		default_database_name TEXT DEFAULT 'master',
		default_language TEXT DEFAULT 'english',
		created_date DATETIME DEFAULT CURRENT_TIMESTAMP,
		modified_date DATETIME DEFAULT CURRENT_TIMESTAMP,
		is_disabled BOOLEAN DEFAULT 0,
		is_locked BOOLEAN DEFAULT 0,
		login_count INTEGER DEFAULT 0,
		last_login_date DATETIME,
		description TEXT
	)`

	_, err := db.Exec(tableSQL)
	if err != nil {
		return fmt.Errorf("failed to create syslogins table: %w", err)
	}

	return nil
}

// Login represents a SQL Server login
type Login struct {
	SID                  int
	Name                  string
	PasswordHash          string
	Type                  string
	DefaultDatabaseName   string
	DefaultLanguage      string
	CreatedDate          time.Time
	ModifiedDate         time.Time
	IsDisabled           bool
	IsLocked             bool
	LoginCount           int
	LastLoginDate        time.Time
	Description           string
}

// CreateLogin creates a new SQL Server login
func (am *AuthManager) CreateLogin(name, password, loginType string) (*Login, error) {
	// Validate inputs
	if name == "" {
		return nil, fmt.Errorf("login name cannot be empty")
	}

	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	if loginType == "" {
		loginType = AuthTypeSQLServer
	}

	// Hash password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert into syslogins table
	query := `
		INSERT INTO syslogins (name, password_hash, type, default_database_name, default_language)
		VALUES (?, ?, ?, 'master', 'english')
	`

	result, err := am.masterDB.Exec(query, name, hash, loginType)
	if err != nil {
		return nil, fmt.Errorf("failed to create login: %w", err)
	}

	// Get login ID
	id, _ := result.LastInsertId()

	// Return created login
	login := &Login{
		SID:                int(id),
		Name:               name,
		PasswordHash:       string(hash),
		Type:               loginType,
		DefaultDatabaseName: "master",
		DefaultLanguage:    "english",
		CreatedDate:       time.Now(),
		ModifiedDate:      time.Now(),
		IsDisabled:        false,
		IsLocked:          false,
		LoginCount:        0,
	}

	return login, nil
}

// AuthenticateLogin authenticates a user with username and password
func (am *AuthManager) AuthenticateLogin(name, password string) (*Login, error) {
	// Get login by name
	login, err := am.GetLoginByName(name)
	if err != nil {
		return nil, fmt.Errorf("login not found: %w", err)
	}

	// Check if login is disabled
	if login.IsDisabled {
		return nil, fmt.Errorf("login is disabled")
	}

	// Check if login is locked
	if login.IsLocked {
		return nil, fmt.Errorf("login is locked")
	}

	// Verify password hash
	err = bcrypt.CompareHashAndPassword([]byte(login.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Update login statistics
	err = am.updateLoginStats(login.SID)
	if err != nil {
		// Log error but don't fail authentication
		fmt.Printf("Warning: failed to update login stats: %v\n", err)
	}

	// Return authenticated login
	return login, nil
}

// GetLoginByName retrieves a login by username
func (am *AuthManager) GetLoginByName(name string) (*Login, error) {
	query := `
		SELECT sid, name, password_hash, type, default_database_name, default_language,
		       created_date, modified_date, is_disabled, is_locked,
		       login_count, last_login_date, description
		FROM syslogins
		WHERE name = ?
	`

	login := &Login{}
	var createdDate, modifiedDate, lastLoginDate sql.NullString

	err := am.masterDB.QueryRow(query, name).Scan(
		&login.SID,
		&login.Name,
		&login.PasswordHash,
		&login.Type,
		&login.DefaultDatabaseName,
		&login.DefaultLanguage,
		&createdDate,
		&modifiedDate,
		&login.IsDisabled,
		&login.IsLocked,
		&login.LoginCount,
		&lastLoginDate,
		&login.Description,
	)

	if err != nil {
		return nil, fmt.Errorf("login not found: %w", err)
	}

	// Parse dates
	if createdDate.Valid {
		login.CreatedDate, _ = time.Parse(time.RFC3339, createdDate.String)
	}

	if modifiedDate.Valid {
		login.ModifiedDate, _ = time.Parse(time.RFC3339, modifiedDate.String)
	}

	if lastLoginDate.Valid {
		login.LastLoginDate, _ = time.Parse(time.RFC3339, lastLoginDate.String)
	}

	return login, nil
}

// GetLoginBySID retrieves a login by SID
func (am *AuthManager) GetLoginBySID(sid int) (*Login, error) {
	query := `
		SELECT sid, name, password_hash, type, default_database_name, default_language,
		       created_date, modified_date, is_disabled, is_locked,
		       login_count, last_login_date, description
		FROM syslogins
		WHERE sid = ?
	`

	login := &Login{}
	var createdDate, modifiedDate, lastLoginDate sql.NullString

	err := am.masterDB.QueryRow(query, sid).Scan(
		&login.SID,
		&login.Name,
		&login.PasswordHash,
		&login.Type,
		&login.DefaultDatabaseName,
		&login.DefaultLanguage,
		&createdDate,
		&modifiedDate,
		&login.IsDisabled,
		&login.IsLocked,
		&login.LoginCount,
		&lastLoginDate,
		&login.Description,
	)

	if err != nil {
		return nil, fmt.Errorf("login not found: %w", err)
	}

	// Parse dates
	if createdDate.Valid {
		login.CreatedDate, _ = time.Parse(time.RFC3339, createdDate.String)
	}

	if modifiedDate.Valid {
		if parsed, err := time.Parse(time.RFC3339, modifiedDate.String); err == nil {
			login.ModifiedDate = parsed
		}
	}

	if lastLoginDate.Valid {
		if parsed, err := time.Parse(time.RFC3339, lastLoginDate.String); err == nil {
			login.LastLoginDate = parsed
		}
	}

	return login, nil
}

// ListLogins lists all logins
func (am *AuthManager) ListLogins() ([]*Login, error) {
	query := `
		SELECT sid, name, password_hash, type, default_database_name, default_language,
		       created_date, modified_date, is_disabled, is_locked,
		       login_count, last_login_date, description
		FROM syslogins
		ORDER BY name ASC
	`

	rows, err := am.masterDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list logins: %w", err)
	}
	defer rows.Close()

	logins := []*Login{}

	for rows.Next() {
		login := &Login{}
		var createdDate, modifiedDate, lastLoginDate sql.NullString

		err := rows.Scan(
			&login.SID,
			&login.Name,
			&login.PasswordHash,
			&login.Type,
			&login.DefaultDatabaseName,
			&login.DefaultLanguage,
			&createdDate,
			&modifiedDate,
			&login.IsDisabled,
			&login.IsLocked,
			&login.LoginCount,
			&lastLoginDate,
			&login.Description,
		)

		if err != nil {
			continue
		}

		// Parse dates
		if createdDate.Valid {
			login.CreatedDate, _ = time.Parse(time.RFC3339, createdDate.String)
		}

		if modifiedDate.Valid {
			if parsed, err := time.Parse(time.RFC3339, modifiedDate.String); err == nil {
				login.ModifiedDate = parsed
			}
		}

		if lastLoginDate.Valid {
			if parsed, err := time.Parse(time.RFC3339, lastLoginDate.String); err == nil {
				login.LastLoginDate = parsed
			}
		}

		logins = append(logins, login)
	}

	return logins, nil
}

// DropLogin removes a login
func (am *AuthManager) DropLogin(name string) error {
	// Prevent dropping sa login
	if name == "sa" {
		return fmt.Errorf("cannot drop 'sa' login")
	}

	query := `DELETE FROM syslogins WHERE name = ?`
	_, err := am.masterDB.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to drop login: %w", err)
	}

	return nil
}

// ChangePassword changes a login's password
func (am *AuthManager) ChangePassword(name, newPassword string) error {
	// Validate password
	if newPassword == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password hash
	query := `UPDATE syslogins SET password_hash = ?, modified_date = CURRENT_TIMESTAMP WHERE name = ?`
	_, err = am.masterDB.Exec(query, hash, name)
	if err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}

	return nil
}

// SetDefaultDatabase sets the default database for a login
func (am *AuthManager) SetDefaultDatabase(loginName, databaseName string) error {
	query := `UPDATE syslogins SET default_database_name = ? WHERE name = ?`
	_, err := am.masterDB.Exec(query, databaseName, loginName)
	if err != nil {
		return fmt.Errorf("failed to set default database: %w", err)
	}

	return nil
}

// DisableLogin disables a login
func (am *AuthManager) DisableLogin(name string) error {
	query := `UPDATE syslogins SET is_disabled = 1 WHERE name = ?`
	_, err := am.masterDB.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to disable login: %w", err)
	}

	return nil
}

// EnableLogin enables a login
func (am *AuthManager) EnableLogin(name string) error {
	query := `UPDATE syslogins SET is_disabled = 0 WHERE name = ?`
	_, err := am.masterDB.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to enable login: %w", err)
	}

	return nil
}

// updateLoginStats updates login statistics (login count and last login date)
func (am *AuthManager) updateLoginStats(sid int) error {
	query := `
		UPDATE syslogins
		SET login_count = login_count + 1,
		    last_login_date = CURRENT_TIMESTAMP
		WHERE sid = ?
	`

	_, err := am.masterDB.Exec(query, sid)
	return err
}

// InitializeDefaultLogins creates default logins (sa)
func (am *AuthManager) InitializeDefaultLogins() error {
	// Check if sa login exists
	_, err := am.GetLoginByName("sa")
	if err == nil {
		// sa login already exists
		return nil
	}

	// Create sa login with default password (empty password for now)
	// In production, require setting password
	_, err = am.CreateLogin("sa", "", AuthTypeSQLServer)
	if err != nil {
		return fmt.Errorf("failed to create sa login: %w", err)
	}

	return nil
}
