package models

import "time"

// Tote represents a storage container/box/tote
type Tote struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Items       string    `json:"items"`        // Text list of items (newline separated)
	ImagePath   string    `json:"image_path"`   // Path to uploaded image
	QRCode      string    `json:"qr_code"`      // QR code identifier (TOTE-XXXXX)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToteCreateRequest represents the data needed to create a new tote
type ToteCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Items       string `json:"items"`
	ImagePath   string `json:"image_path"`
}

// ToteUpdateRequest represents the data that can be updated
type ToteUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Items       *string `json:"items,omitempty"`
	ImagePath   *string `json:"image_path,omitempty"`
}
