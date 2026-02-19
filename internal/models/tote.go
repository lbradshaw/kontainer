package models

import "time"

// ToteImage represents an image attached to a tote
type ToteImage struct {
	ID           int       `json:"id"`
	ToteID       int       `json:"tote_id"`
	ImageData    string    `json:"image_data"`      // Base64 encoded image data
	ImageType    string    `json:"image_type"`      // MIME type (image/jpeg, image/png, etc.)
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
}

// Tote represents a storage container/box/tote
type Tote struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Items       string      `json:"items"`        // Text list of items (newline separated)
	Location    string      `json:"location"`     // Physical location (e.g., "Garage", "Basement")
	ImagePath   string      `json:"image_path"`   // Legacy: Path to first uploaded image (for backward compatibility)
	Images      []ToteImage `json:"images"`       // Array of images
	QRCode      string      `json:"qr_code"`      // QR code identifier (TOTE-XXXXX)
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// ToteCreateRequest represents the data needed to create a new tote
type ToteCreateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Items       string   `json:"items"`
	Location    string   `json:"location"`
	ImagePath   string   `json:"image_path"`    // Legacy single image (base64)
	ImagePaths  []string `json:"image_paths"`   // Multiple images (base64)
	ImageTypes  []string `json:"image_types"`   // MIME types for images
}

// ToteUpdateRequest represents the data that can be updated
type ToteUpdateRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Items       *string  `json:"items,omitempty"`
	Location    *string  `json:"location,omitempty"`
	ImagePath   *string  `json:"image_path,omitempty"`
	ImagePaths  []string `json:"image_paths,omitempty"`   // New images to add (base64)
	ImageTypes  []string `json:"image_types,omitempty"`   // MIME types for new images
}
