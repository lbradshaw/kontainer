package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrateDatabase(t *testing.T) {
	// Create temporary test directory
	tempDir := t.TempDir()
	
	// Create a test database file
	oldPath := filepath.Join(tempDir, "old.db")
	newPath := filepath.Join(tempDir, "new", "database.db")
	
	// Write test data to old database
	testData := []byte("test database content")
	if err := os.WriteFile(oldPath, testData, 0644); err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	
	// Create service
	service := NewSettingsService()
	
	// Test migration
	if err := service.MigrateDatabase(oldPath, newPath); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}
	
	// Verify new file exists
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Fatal("New database file was not created")
	}
	
	// Verify content matches
	newData, err := os.ReadFile(newPath)
	if err != nil {
		t.Fatalf("Failed to read new database: %v", err)
	}
	
	if string(newData) != string(testData) {
		t.Errorf("Database content mismatch. Expected %s, got %s", testData, newData)
	}
	
	// Verify old file still exists (copy, not move)
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		t.Error("Old database file should still exist after migration")
	}
}

func TestMigrateDatabaseSamePath(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	// Create test database
	if err := os.WriteFile(dbPath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	
	service := NewSettingsService()
	
	// Migrating to same path should be a no-op
	if err := service.MigrateDatabase(dbPath, dbPath); err != nil {
		t.Fatalf("Migration to same path should not error: %v", err)
	}
}

func TestMigrateDatabaseNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	oldPath := filepath.Join(tempDir, "nonexistent.db")
	newPath := filepath.Join(tempDir, "new.db")
	
	service := NewSettingsService()
	
	// Migrating non-existent database should not error (nothing to migrate)
	if err := service.MigrateDatabase(oldPath, newPath); err != nil {
		t.Fatalf("Migration of non-existent database should not error: %v", err)
	}
}

func TestMigrateDatabaseAlreadyExists(t *testing.T) {
	tempDir := t.TempDir()
	oldPath := filepath.Join(tempDir, "old.db")
	newPath := filepath.Join(tempDir, "new.db")
	
	// Create both files
	if err := os.WriteFile(oldPath, []byte("old"), 0644); err != nil {
		t.Fatalf("Failed to create old database: %v", err)
	}
	if err := os.WriteFile(newPath, []byte("new"), 0644); err != nil {
		t.Fatalf("Failed to create new database: %v", err)
	}
	
	service := NewSettingsService()
	
	// Should error because new file already exists
	err := service.MigrateDatabase(oldPath, newPath)
	if err == nil {
		t.Fatal("Expected error when destination file exists, got nil")
	}
}
