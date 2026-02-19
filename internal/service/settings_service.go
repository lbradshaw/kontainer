package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"kontainer/internal/models"
)

type SettingsService struct {
	filepath string
}

func NewSettingsService() *SettingsService {
	return &SettingsService{
		filepath: "settings.json",
	}
}

// LoadSettings loads settings from file or returns defaults
func (s *SettingsService) LoadSettings() (*models.Settings, error) {
	// Check if file exists
	if _, err := os.Stat(s.filepath); os.IsNotExist(err) {
		// Create default settings file
		settings := models.DefaultSettings()
		if err := s.SaveSettings(settings); err != nil {
			return nil, fmt.Errorf("failed to create default settings: %w", err)
		}
		return settings, nil
	}

	// Read existing file
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings file: %w", err)
	}

	var settings models.Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse settings: %w", err)
	}

	return &settings, nil
}

// SaveSettings saves settings to file
func (s *SettingsService) SaveSettings(settings *models.Settings) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(s.filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	return nil
}

// MigrateDatabase moves the database file from oldPath to newPath
func (s *SettingsService) MigrateDatabase(oldPath, newPath string) error {
	// Validate paths
	if oldPath == "" || newPath == "" {
		return fmt.Errorf("database paths cannot be empty")
	}

	// If paths are the same, no migration needed
	if oldPath == newPath {
		return nil
	}

	// Check if old database exists
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		// Old database doesn't exist, nothing to migrate
		return nil
	}

	// Create directory for new database if needed
	newDir := filepath.Dir(newPath)
	if newDir != "" && newDir != "." {
		if err := os.MkdirAll(newDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for new database: %w", err)
		}
	}

	// Check if file already exists at new location
	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("database file already exists at new location: %s", newPath)
	}

	// Copy database file to new location
	if err := copyFile(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to copy database: %w", err)
	}

	// Copy associated files (journal, wal, shm if they exist)
	for _, ext := range []string{"-journal", "-wal", "-shm"} {
		oldFile := oldPath + ext
		newFile := newPath + ext
		if _, err := os.Stat(oldFile); err == nil {
			if err := copyFile(oldFile, newFile); err != nil {
				// Non-critical, just log
				fmt.Printf("Warning: failed to copy %s: %v\n", oldFile, err)
			}
		}
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Sync to ensure data is written to disk
	return destFile.Sync()
}
