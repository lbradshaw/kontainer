package models

// Settings represents application settings
type Settings struct {
	Port int `json:"port"`
}

// DefaultSettings returns the default application settings
func DefaultSettings() *Settings {
	return &Settings{
		Port: 3818,
	}
}
