package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

// DB wraps the sql.DB connection
type DB struct {
	conn *sql.DB
}

// New creates a new database connection
func New(dbPath string) (*DB, error) {
	log.Infof("Opening database at %s", dbPath)

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable foreign keys
	if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	log.Info("Database connection established")

	return &DB{conn: conn}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// Conn returns the underlying sql.DB connection
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// Migrate runs database migrations
func (db *DB) Migrate() error {
	log.Info("Running database migrations...")

	migrations := []string{
		createStudentsTable,
		createCoursesTable,
		createAssignmentsTable,
		createGradesTable,
		createEventsTable,
		createMessagesTable,
		createConfigTable,
		createSyncStatusTable,
		createIndexes,
	}

	for i, migration := range migrations {
		log.Debugf("Running migration %d/%d", i+1, len(migrations))
		if _, err := db.conn.Exec(migration); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	log.Info("Database migrations completed successfully")
	return nil
}

const createStudentsTable = `
CREATE TABLE IF NOT EXISTS students (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    grade_level INTEGER,
    schoology_user_id TEXT,
    powerschool_student_id TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createCoursesTable = `
CREATE TABLE IF NOT EXISTS courses (
    id TEXT PRIMARY KEY,
    student_id TEXT NOT NULL,
    name TEXT NOT NULL,
    teacher TEXT,
    period TEXT,
    source TEXT NOT NULL,
    external_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
);
`

const createAssignmentsTable = `
CREATE TABLE IF NOT EXISTS assignments (
    id TEXT PRIMARY KEY,
    student_id TEXT NOT NULL,
    course_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    due_date TIMESTAMP,
    status TEXT,
    source TEXT NOT NULL,
    external_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);
`

const createGradesTable = `
CREATE TABLE IF NOT EXISTS grades (
    id TEXT PRIMARY KEY,
    student_id TEXT NOT NULL,
    course_id TEXT NOT NULL,
    assignment_id TEXT,
    score REAL,
    max_score REAL,
    percentage REAL,
    letter_grade TEXT,
    grading_period TEXT,
    source TEXT NOT NULL,
    external_id TEXT,
    recorded_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    FOREIGN KEY (assignment_id) REFERENCES assignments(id) ON DELETE CASCADE
);
`

const createEventsTable = `
CREATE TABLE IF NOT EXISTS events (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    event_date DATE NOT NULL,
    event_type TEXT,
    source TEXT NOT NULL,
    external_id TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createMessagesTable = `
CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    student_id TEXT NOT NULL,
    from_name TEXT NOT NULL,
    from_email TEXT,
    subject TEXT,
    body TEXT,
    received_at TIMESTAMP,
    read BOOLEAN DEFAULT FALSE,
    source TEXT NOT NULL,
    external_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
);
`

const createConfigTable = `
CREATE TABLE IF NOT EXISTS config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    encrypted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createSyncStatusTable = `
CREATE TABLE IF NOT EXISTS sync_status (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    student_id TEXT NOT NULL,
    source TEXT NOT NULL,
    data_type TEXT NOT NULL,
    last_sync TIMESTAMP,
    status TEXT,
    error_message TEXT,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    UNIQUE(student_id, source, data_type)
);
`

const createIndexes = `
CREATE INDEX IF NOT EXISTS idx_assignments_student_due ON assignments(student_id, due_date);
CREATE INDEX IF NOT EXISTS idx_assignments_status ON assignments(status);
CREATE INDEX IF NOT EXISTS idx_grades_student_course ON grades(student_id, course_id);
CREATE INDEX IF NOT EXISTS idx_events_date ON events(event_date);
CREATE INDEX IF NOT EXISTS idx_messages_student_unread ON messages(student_id, read);
CREATE INDEX IF NOT EXISTS idx_sync_status_lookup ON sync_status(student_id, source, data_type);
`
