package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Database wraps the SQLite connection and provides methods for data operations
type Database struct {
	conn *sql.DB
}

// Profile represents a LinkedIn profile scraped from search
type Profile struct {
	ID         string
	Name       string
	Title      string
	Company    string
	Location   string
	ProfileURL string
	VisitedAt  time.Time
	CreatedAt  time.Time
}

// ConnectionRequest tracks sent connection requests
type ConnectionRequest struct {
	ID        int
	ProfileID string
	SentAt    time.Time
	NoteUsed  string
	Status    string // 'pending', 'accepted', 'rejected', 'withdrawn'
	CreatedAt time.Time
}

// Message tracks sent messages to connections
type Message struct {
	ID             int
	ConnectionID   string
	TemplateName   string
	MessageContent string
	SentAt         time.Time
	CreatedAt      time.Time
}

// RateLimit tracks daily action limits
type RateLimit struct {
	Date            string // YYYY-MM-DD format
	ConnectionCount int
	MessageCount    int
	SearchCount     int
	LastUpdated     time.Time
}

// InitDB creates a new database connection and initializes tables
func InitDB(dbPath string) (*Database, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &Database{conn: conn}

	// Create tables
	if err := db.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// createTables initializes all required database tables
func (db *Database) createTables() error {
	schema := `
	-- Profiles table: stores scraped LinkedIn profiles
	CREATE TABLE IF NOT EXISTS profiles (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		title TEXT,
		company TEXT,
		location TEXT,
		profile_url TEXT NOT NULL UNIQUE,
		visited_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Connection requests table: tracks all sent connection requests
	CREATE TABLE IF NOT EXISTS connection_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_id TEXT NOT NULL,
		sent_at DATETIME NOT NULL,
		note_used TEXT,
		status TEXT DEFAULT 'pending',
		has_replied BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (profile_id) REFERENCES profiles(id)
	);

	-- Messages table: tracks all sent messages
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		connection_id TEXT NOT NULL,
		template_name TEXT,
		message_content TEXT NOT NULL,
		sent_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Rate limits table: tracks daily action quotas
	CREATE TABLE IF NOT EXISTS rate_limits (
		date TEXT PRIMARY KEY,
		connection_count INTEGER DEFAULT 0,
		message_count INTEGER DEFAULT 0,
		search_count INTEGER DEFAULT 0,
		last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_profiles_visited ON profiles(visited_at);
	CREATE INDEX IF NOT EXISTS idx_connection_requests_profile ON connection_requests(profile_id);
	CREATE INDEX IF NOT EXISTS idx_connection_requests_sent ON connection_requests(sent_at);
	CREATE INDEX IF NOT EXISTS idx_messages_connection ON messages(connection_id);
	CREATE INDEX IF NOT EXISTS idx_messages_sent ON messages(sent_at);
	`

	_, err := db.conn.Exec(schema)
	return err
}

// Close closes the database connection
func (db *Database) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// --- Profile Operations ---

// SaveProfile saves a profile to the database
func (db *Database) SaveProfile(profile Profile) error {
	query := `
		INSERT INTO profiles (id, name, title, company, location, profile_url, visited_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			title = excluded.title,
			company = excluded.company,
			location = excluded.location,
			visited_at = excluded.visited_at
	`

	_, err := db.conn.Exec(query,
		profile.ID,
		profile.Name,
		profile.Title,
		profile.Company,
		profile.Location,
		profile.ProfileURL,
		profile.VisitedAt,
		profile.CreatedAt,
	)

	return err
}

// IsDuplicateProfile checks if a profile was visited recently (within 30 days)
func (db *Database) IsDuplicateProfile(profileID string, daysSince int) (bool, error) {
	query := `
		SELECT COUNT(*) FROM profiles
		WHERE id = ? AND datetime(visited_at, 'utc') > datetime('now', '-' || ? || ' days')
	`

	var count int
	err := db.conn.QueryRow(query, profileID, daysSince).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetProfile retrieves a profile by ID
func (db *Database) GetProfile(profileID string) (*Profile, error) {
	query := `
		SELECT id, name, title, company, location, profile_url, visited_at, created_at
		FROM profiles WHERE id = ?
	`

	var profile Profile
	err := db.conn.QueryRow(query, profileID).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Title,
		&profile.Company,
		&profile.Location,
		&profile.ProfileURL,
		&profile.VisitedAt,
		&profile.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &profile, nil
}

// --- Connection Request Operations ---

// SaveConnectionRequest records a sent connection request
func (db *Database) SaveConnectionRequest(req ConnectionRequest) error {
	query := `
		INSERT INTO connection_requests (profile_id, sent_at, note_used, status, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := db.conn.Exec(query,
		req.ProfileID,
		req.SentAt,
		req.NoteUsed,
		req.Status,
		req.CreatedAt,
	)

	return err
}

// UpdateConnectionStatus updates the status of a connection request
func (db *Database) UpdateConnectionStatus(profileID, status string) error {
	query := `
		UPDATE connection_requests
		SET status = ?
		WHERE profile_id = ? AND status = 'pending'
	`

	_, err := db.conn.Exec(query, status, profileID)
	return err
}

// GetPendingConnections retrieves all pending connection requests
func (db *Database) GetPendingConnections() ([]ConnectionRequest, error) {
	query := `
		SELECT id, profile_id, sent_at, note_used, status, created_at
		FROM connection_requests
		WHERE status = 'pending'
		ORDER BY sent_at DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []ConnectionRequest
	for rows.Next() {
		var req ConnectionRequest
		err := rows.Scan(
			&req.ID,
			&req.ProfileID,
			&req.SentAt,
			&req.NoteUsed,
			&req.Status,
			&req.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	return requests, nil
}

// HasSentConnectionRequest checks if a connection request was already sent to a profile
func (db *Database) HasSentConnectionRequest(profileID string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM connection_requests
		WHERE profile_id = ?
	`

	var count int
	err := db.conn.QueryRow(query, profileID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// --- Message Operations ---

// SaveMessage records a sent message
func (db *Database) SaveMessage(msg Message) error {
	query := `
		INSERT INTO messages (connection_id, template_name, message_content, sent_at, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := db.conn.Exec(query,
		msg.ConnectionID,
		msg.TemplateName,
		msg.MessageContent,
		msg.SentAt,
		msg.CreatedAt,
	)

	return err
}

// HasSentMessage checks if a message was already sent to a connection
func (db *Database) HasSentMessage(connectionID, templateName string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM messages
		WHERE connection_id = ? AND template_name = ?
	`

	var count int
	err := db.conn.QueryRow(query, connectionID, templateName).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetMessageHistory retrieves all messages sent to a connection
func (db *Database) GetMessageHistory(connectionID string) ([]Message, error) {
	query := `
		SELECT id, connection_id, template_name, message_content, sent_at, created_at
		FROM messages
		WHERE connection_id = ?
		ORDER BY sent_at ASC
	`

	rows, err := db.conn.Query(query, connectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(
			&msg.ID,
			&msg.ConnectionID,
			&msg.TemplateName,
			&msg.MessageContent,
			&msg.SentAt,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// --- Rate Limit Operations ---

// GetTodayRateLimit retrieves or creates today's rate limit record
func (db *Database) GetTodayRateLimit() (*RateLimit, error) {
	today := time.Now().Format("2006-01-02")

	query := `
		SELECT date, connection_count, message_count, search_count, last_updated
		FROM rate_limits WHERE date = ?
	`

	var limit RateLimit
	err := db.conn.QueryRow(query, today).Scan(
		&limit.Date,
		&limit.ConnectionCount,
		&limit.MessageCount,
		&limit.SearchCount,
		&limit.LastUpdated,
	)

	if err == sql.ErrNoRows {
		// Create new record for today
		insertQuery := `
			INSERT INTO rate_limits (date, connection_count, message_count, search_count, last_updated)
			VALUES (?, 0, 0, 0, ?)
		`
		_, err := db.conn.Exec(insertQuery, today, time.Now())
		if err != nil {
			return nil, err
		}

		// Return fresh limit
		return &RateLimit{
			Date:            today,
			ConnectionCount: 0,
			MessageCount:    0,
			SearchCount:     0,
			LastUpdated:     time.Now(),
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &limit, nil
}

// IncrementConnectionCount increments today's connection request count
func (db *Database) IncrementConnectionCount() error {
	today := time.Now().Format("2006-01-02")

	query := `
		INSERT INTO rate_limits (date, connection_count, message_count, search_count, last_updated)
		VALUES (?, 1, 0, 0, ?)
		ON CONFLICT(date) DO UPDATE SET
			connection_count = connection_count + 1,
			last_updated = ?
	`

	now := time.Now()
	_, err := db.conn.Exec(query, today, now, now)
	return err
}

// IncrementMessageCount increments today's message count
func (db *Database) IncrementMessageCount() error {
	today := time.Now().Format("2006-01-02")

	query := `
		INSERT INTO rate_limits (date, connection_count, message_count, search_count, last_updated)
		VALUES (?, 0, 1, 0, ?)
		ON CONFLICT(date) DO UPDATE SET
			message_count = message_count + 1,
			last_updated = ?
	`

	now := time.Now()
	_, err := db.conn.Exec(query, today, now, now)
	return err
}

// IncrementSearchCount increments today's search count
func (db *Database) IncrementSearchCount() error {
	today := time.Now().Format("2006-01-02")

	query := `
		INSERT INTO rate_limits (date, connection_count, message_count, search_count, last_updated)
		VALUES (?, 0, 0, 1, ?)
		ON CONFLICT(date) DO UPDATE SET
			search_count = search_count + 1,
			last_updated = ?
	`

	now := time.Now()
	_, err := db.conn.Exec(query, today, now, now)
	return err
}

// GetRecentProfiles retrieves recent profiles that haven't been contacted
func (db *Database) GetRecentProfiles(limit int, daysBack int) ([]Profile, error) {
	query := `
		SELECT DISTINCT p.id, p.name, p.title, p.company, p.location, p.profile_url, p.visited_at, p.created_at
		FROM profiles p
		WHERE datetime(p.visited_at, 'utc') >= datetime('now', '-' || ? || ' days')
		AND p.id NOT IN (
			SELECT profile_id FROM connection_requests
			WHERE datetime(sent_at, 'utc') >= datetime('now', '-' || ? || ' days')
		)
		ORDER BY p.visited_at DESC
		LIMIT ?
	`

	rows, err := db.conn.Query(query, daysBack, daysBack, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var profile Profile
		err := rows.Scan(
			&profile.ID,
			&profile.Name,
			&profile.Title,
			&profile.Company,
			&profile.Location,
			&profile.ProfileURL,
			&profile.VisitedAt,
			&profile.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, rows.Err()
}

// GetDailyStats retrieves statistics for a specific date
func (db *Database) GetDailyStats(date string) (*RateLimit, error) {
	query := `
		SELECT date, connection_count, message_count, search_count, last_updated
		FROM rate_limits WHERE date = ?
	`

	var limit RateLimit
	err := db.conn.QueryRow(query, date).Scan(
		&limit.Date,
		&limit.ConnectionCount,
		&limit.MessageCount,
		&limit.SearchCount,
		&limit.LastUpdated,
	)

	if err == sql.ErrNoRows {
		return &RateLimit{
			Date:            date,
			ConnectionCount: 0,
			MessageCount:    0,
			SearchCount:     0,
			LastUpdated:     time.Now(),
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &limit, nil
}

// GetAcceptedConnectionProfiles retrieves profiles where connection was accepted and haven't been messaged yet
// This is used for messaging automation to only message actual connections
func (db *Database) GetAcceptedConnectionProfiles(limit int, daysBack int) ([]Profile, error) {
	query := `
		SELECT DISTINCT p.id, p.name, p.title, p.company, p.location, p.profile_url, p.visited_at, p.created_at
		FROM profiles p
		INNER JOIN connection_requests cr ON p.id = cr.profile_id
		WHERE cr.status = 'accepted'
		AND (cr.has_replied IS NULL OR cr.has_replied = 0)
		AND datetime(cr.sent_at, 'utc') >= datetime('now', '-' || ? || ' days')
		AND p.id NOT IN (
			SELECT connection_id FROM messages
			WHERE datetime(sent_at, 'utc') >= datetime('now', '-' || ? || ' days')
		)
		ORDER BY cr.sent_at DESC
		LIMIT ?
	`

	rows, err := db.conn.Query(query, daysBack, daysBack, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var profile Profile
		err := rows.Scan(
			&profile.ID,
			&profile.Name,
			&profile.Title,
			&profile.Company,
			&profile.Location,
			&profile.ProfileURL,
			&profile.VisitedAt,
			&profile.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	return profiles, rows.Err()
}

// UpdateConnectionReplyStatus updates the has_replied status for a connection
func (db *Database) UpdateConnectionReplyStatus(profileID string, hasReplied bool) error {
	query := `
		UPDATE connection_requests
		SET has_replied = ?
		WHERE profile_id = ?
	`
	_, err := db.conn.Exec(query, hasReplied, profileID)
	return err
}
