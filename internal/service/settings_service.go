package service

import (
	"encoding/json"
	"fmt"
	"os"

	"totetrax/internal/models"
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
