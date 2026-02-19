package models

// Settings represents application settings
type Settings struct {
	Port       int    `json:"port"`
	Theme      string `json:"theme"`       // "light" or "dark"
	DatabasePath string `json:"database_path"` // Path to database file
}

// DefaultSettings returns the default application settings
func DefaultSettings() *Settings {
	return &Settings{
		Port:       3818,
		Theme:      "dark",
		DatabasePath: "kontainer.db",
	}
}
