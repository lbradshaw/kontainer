package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// InitDB initializes the SQLite database and creates tables
func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Enable foreign key constraints (required for CASCADE DELETE)
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("error enabling foreign keys: %w", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("error creating tables: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("error running migrations: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS totes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		items TEXT,
		image_path TEXT,
		qr_code TEXT UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tote_images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tote_id INTEGER NOT NULL,
		image_data BLOB NOT NULL,
		image_type TEXT NOT NULL,
		display_order INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (tote_id) REFERENCES totes(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_name ON totes(name);
	CREATE INDEX IF NOT EXISTS idx_qr_code ON totes(qr_code);
	CREATE INDEX IF NOT EXISTS idx_tote_images_tote_id ON tote_images(tote_id);
	`

	if _, err := db.Exec(schema); err != nil {
		return err
	}

	return nil
}

func runMigrations(db *sql.DB) error {
	// Add location column if it doesn't exist
	var columnExists bool
	err := db.QueryRow(`
		SELECT COUNT(*) > 0 
		FROM pragma_table_info('totes') 
		WHERE name = 'location'
	`).Scan(&columnExists)
	
	if err != nil {
		return fmt.Errorf("error checking for location column: %w", err)
	}

	if !columnExists {
		_, err = db.Exec("ALTER TABLE totes ADD COLUMN location TEXT")
		if err != nil {
			return fmt.Errorf("error adding location column: %w", err)
		}
	}

	// Add parent_id column if it doesn't exist
	err = db.QueryRow(`
		SELECT COUNT(*) > 0 
		FROM pragma_table_info('totes') 
		WHERE name = 'parent_id'
	`).Scan(&columnExists)
	
	if err != nil {
		return fmt.Errorf("error checking for parent_id column: %w", err)
	}

	if !columnExists {
		_, err = db.Exec("ALTER TABLE totes ADD COLUMN parent_id INTEGER REFERENCES totes(id) ON DELETE CASCADE")
		if err != nil {
			return fmt.Errorf("error adding parent_id column: %w", err)
		}
		_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_parent_id ON totes(parent_id)")
		if err != nil {
			return fmt.Errorf("error creating parent_id index: %w", err)
		}
	}

	// Add depth column if it doesn't exist
	err = db.QueryRow(`
		SELECT COUNT(*) > 0 
		FROM pragma_table_info('totes') 
		WHERE name = 'depth'
	`).Scan(&columnExists)
	
	if err != nil {
		return fmt.Errorf("error checking for depth column: %w", err)
	}

	if !columnExists {
		_, err = db.Exec("ALTER TABLE totes ADD COLUMN depth INTEGER NOT NULL DEFAULT 0 CHECK (depth IN (0, 1))")
		if err != nil {
			return fmt.Errorf("error adding depth column: %w", err)
		}
	}

	return nil
}
