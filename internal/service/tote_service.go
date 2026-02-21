package service

import (
	"database/sql"
	"fmt"
	"time"

	"kontainer/internal/models"
)

type ToteService struct {
	db *sql.DB
}

func NewToteService(db *sql.DB) *ToteService {
	return &ToteService{db: db}
}

// loadImagesForTote loads all images for a specific tote
func (s *ToteService) loadImagesForTote(toteID int) ([]models.ToteImage, error) {
	query := `
		SELECT id, tote_id, image_data, image_type, display_order, created_at
		FROM tote_images
		WHERE tote_id = ?
		ORDER BY display_order ASC, created_at ASC
	`

	rows, err := s.db.Query(query, toteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.ToteImage
	for rows.Next() {
		var img models.ToteImage
		var imageData []byte
		err := rows.Scan(&img.ID, &img.ToteID, &imageData, &img.ImageType, &img.DisplayOrder, &img.CreatedAt)
		if err != nil {
			return nil, err
		}
		// Convert blob to base64 for JSON
		img.ImageData = "data:" + img.ImageType + ";base64," + base64Encode(imageData)
		images = append(images, img)
	}

	return images, nil
}

func base64Encode(data []byte) string {
	// Simple base64 encoding
	const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	
	result := make([]byte, ((len(data)+2)/3)*4)
	j := 0
	
	for i := 0; i < len(data); i += 3 {
		b := (uint32(data[i]) << 16)
		if i+1 < len(data) {
			b |= uint32(data[i+1]) << 8
		}
		if i+2 < len(data) {
			b |= uint32(data[i+2])
		}
		
		result[j] = base64Table[(b>>18)&0x3F]
		result[j+1] = base64Table[(b>>12)&0x3F]
		if i+1 < len(data) {
			result[j+2] = base64Table[(b>>6)&0x3F]
		} else {
			result[j+2] = '='
		}
		if i+2 < len(data) {
			result[j+3] = base64Table[b&0x3F]
		} else {
			result[j+3] = '='
		}
		j += 4
	}
	
	return string(result)
}

// loadImagesForTotes loads images for multiple totes
func (s *ToteService) loadImagesForTotes(totes []models.Tote) ([]models.Tote, error) {
	for i := range totes {
		images, err := s.loadImagesForTote(totes[i].ID)
		if err != nil {
			return nil, err
		}
		totes[i].Images = images
	}
	return totes, nil
}

// GetAll retrieves all top-level totes (parent_id IS NULL)
func (s *ToteService) GetAll() ([]models.Tote, error) {
	query := `
		SELECT id, name, description, items, location, image_path, qr_code, parent_id, depth, created_at, updated_at
		FROM totes
		WHERE parent_id IS NULL
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totes []models.Tote
	for rows.Next() {
		var t models.Tote
		err := rows.Scan(
			&t.ID, &t.Name, &t.Description, &t.Items, &t.Location, &t.ImagePath,
			&t.QRCode, &t.ParentID, &t.Depth, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		totes = append(totes, t)
	}

	// Load images for all totes
	return s.loadImagesForTotes(totes)
}

// GetAllIncludingChildren retrieves all totes (including sub-containers)
func (s *ToteService) GetAllIncludingChildren() ([]models.Tote, error) {
	query := `
		SELECT id, name, description, items, location, image_path, qr_code, parent_id, depth, created_at, updated_at
		FROM totes
		ORDER BY parent_id NULLS FIRST, created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totes []models.Tote
	for rows.Next() {
		var t models.Tote
		err := rows.Scan(
			&t.ID, &t.Name, &t.Description, &t.Items, &t.Location, &t.ImagePath,
			&t.QRCode, &t.ParentID, &t.Depth, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		totes = append(totes, t)
	}

	// Load images for all totes
	return s.loadImagesForTotes(totes)
}

// GetChildren retrieves all child totes for a parent tote
func (s *ToteService) GetChildren(parentID int) ([]models.Tote, error) {
	query := `
		SELECT id, name, description, items, location, image_path, qr_code, parent_id, depth, created_at, updated_at
		FROM totes
		WHERE parent_id = ?
		ORDER BY created_at ASC
	`

	rows, err := s.db.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []models.Tote
	for rows.Next() {
		var t models.Tote
		err := rows.Scan(
			&t.ID, &t.Name, &t.Description, &t.Items, &t.Location, &t.ImagePath,
			&t.QRCode, &t.ParentID, &t.Depth, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		children = append(children, t)
	}

	// Load images for all children
	return s.loadImagesForTotes(children)
}

// GetByID retrieves a tote by ID
func (s *ToteService) GetByID(id int) (*models.Tote, error) {
	query := `
		SELECT id, name, description, items, location, image_path, qr_code, parent_id, depth, created_at, updated_at
		FROM totes
		WHERE id = ?
	`

	var t models.Tote
	err := s.db.QueryRow(query, id).Scan(
		&t.ID, &t.Name, &t.Description, &t.Items, &t.Location, &t.ImagePath,
		&t.QRCode, &t.ParentID, &t.Depth, &t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tote not found")
	}
	if err != nil {
		return nil, err
	}

	// Load images for this tote
	images, err := s.loadImagesForTote(t.ID)
	if err != nil {
		return nil, err
	}
	t.Images = images

	// Load children for this tote
	children, err := s.GetChildren(t.ID)
	if err != nil {
		return nil, err
	}
	t.Children = children

	return &t, nil
}

// GetByQRCode retrieves a tote by QR code
func (s *ToteService) GetByQRCode(qrCode string) (*models.Tote, error) {
	query := `
		SELECT id, name, description, items, location, image_path, qr_code, parent_id, depth, created_at, updated_at
		FROM totes
		WHERE qr_code = ?
	`

	var t models.Tote
	err := s.db.QueryRow(query, qrCode).Scan(
		&t.ID, &t.Name, &t.Description, &t.Items, &t.Location, &t.ImagePath,
		&t.QRCode, &t.ParentID, &t.Depth, &t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tote not found")
	}
	if err != nil {
		return nil, err
	}

	// Load images for this tote
	images, err := s.loadImagesForTote(t.ID)
	if err != nil {
		return nil, err
	}
	t.Images = images

	// Load children for this tote
	children, err := s.GetChildren(t.ID)
	if err != nil {
		return nil, err
	}
	t.Children = children

	return &t, nil
}

// Create creates a new tote
func (s *ToteService) Create(req models.ToteCreateRequest) (*models.Tote, error) {
	// Validate parent_id and calculate depth
	depth := 0
	if req.ParentID != nil {
		// Check if parent exists
		var parentDepth int
		err := s.db.QueryRow("SELECT depth FROM totes WHERE id = ?", *req.ParentID).Scan(&parentDepth)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("parent tote not found")
		}
		if err != nil {
			return nil, err
		}
		
		// Enforce depth limit: sub-containers can only be created under top-level containers
		if parentDepth >= 1 {
			return nil, fmt.Errorf("cannot create sub-container: maximum nesting depth (2 levels) exceeded")
		}
		
		depth = parentDepth + 1
	}

	// Generate QR code (TOTE-XXXXX format)
	var maxID int
	err := s.db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM totes").Scan(&maxID)
	if err != nil {
		return nil, err
	}
	qrCode := fmt.Sprintf("TOTE-%05d", maxID+1)

	query := `
		INSERT INTO totes (name, description, items, location, image_path, qr_code, parent_id, depth, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	// Use first image from ImagePaths array as legacy image_path if available
	legacyImagePath := req.ImagePath
	if len(req.ImagePaths) > 0 && legacyImagePath == "" {
		legacyImagePath = req.ImagePaths[0]
	}

	result, err := s.db.Exec(query, req.Name, req.Description, req.Items, req.Location, legacyImagePath, qrCode, req.ParentID, depth, now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	toteID := int(id)

	// Insert images into tote_images table
	if len(req.ImagePaths) > 0 {
		for i, imageData := range req.ImagePaths {
			// Decode base64 image data
			imageBlob, imageType, err := decodeBase64Image(imageData)
			if err != nil {
				return nil, fmt.Errorf("failed to decode image %d: %w", i, err)
			}
			
			// Use provided image type or detected type
			mimeType := imageType
			if i < len(req.ImageTypes) && req.ImageTypes[i] != "" {
				mimeType = req.ImageTypes[i]
			}
			
			_, err = s.db.Exec(`
				INSERT INTO tote_images (tote_id, image_data, image_type, display_order, created_at)
				VALUES (?, ?, ?, ?, ?)
			`, toteID, imageBlob, mimeType, i, now)
			if err != nil {
				return nil, err
			}
		}
	}

	return s.GetByID(toteID)
}

// decodeBase64Image decodes a base64 data URI and returns the binary data and MIME type
func decodeBase64Image(dataURI string) ([]byte, string, error) {
	// Remove data URI prefix if present (e.g., "data:image/png;base64,")
	var base64Data string
	var mimeType string = "image/png" // default
	
	if len(dataURI) > 5 && dataURI[:5] == "data:" {
		// Extract MIME type and base64 data
		parts := splitString(dataURI, ",")
		if len(parts) != 2 {
			return nil, "", fmt.Errorf("invalid data URI format")
		}
		
		// Parse MIME type from "data:image/png;base64"
		header := parts[0]
		if len(header) > 5 {
			typeParts := splitString(header[5:], ";")
			if len(typeParts) > 0 {
				mimeType = typeParts[0]
			}
		}
		base64Data = parts[1]
	} else {
		base64Data = dataURI
	}
	
	// Decode base64
	decoded, err := base64Decode(base64Data)
	if err != nil {
		return nil, "", err
	}
	
	return decoded, mimeType, nil
}

func base64Decode(encoded string) ([]byte, error) {
	const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	
	// Create reverse lookup table
	lookup := make(map[byte]int)
	for i := 0; i < len(base64Table); i++ {
		lookup[base64Table[i]] = i
	}
	
	// Remove padding
	encoded = trimRight(encoded, "=")
	
	result := make([]byte, (len(encoded)*3)/4)
	j := 0
	
	for i := 0; i < len(encoded); i += 4 {
		b0, ok := lookup[encoded[i]]
		if !ok {
			return nil, fmt.Errorf("invalid base64 character")
		}
		b1, ok := lookup[encoded[i+1]]
		if !ok {
			return nil, fmt.Errorf("invalid base64 character")
		}
		
		result[j] = byte(b0<<2 | b1>>4)
		j++
		
		if i+2 < len(encoded) {
			b2, ok := lookup[encoded[i+2]]
			if !ok {
				return nil, fmt.Errorf("invalid base64 character")
			}
			result[j] = byte(b1<<4 | b2>>2)
			j++
			
			if i+3 < len(encoded) {
				b3, ok := lookup[encoded[i+3]]
				if !ok {
					return nil, fmt.Errorf("invalid base64 character")
				}
				result[j] = byte(b2<<6 | b3)
				j++
			}
		}
	}
	
	return result[:j], nil
}

func splitString(s, sep string) []string {
	if sep == "" {
		return []string{s}
	}
	
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimRight(s, cutset string) string {
	for len(s) > 0 {
		found := false
		for i := 0; i < len(cutset); i++ {
			if s[len(s)-1] == cutset[i] {
				s = s[:len(s)-1]
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	return s
}

// Update updates an existing tote
func (s *ToteService) Update(id int, req models.ToteUpdateRequest) (*models.Tote, error) {
	// Check if tote exists
	existing, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}

	if req.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *req.Description)
	}
	if req.Items != nil {
		updates = append(updates, "items = ?")
		args = append(args, *req.Items)
	}
	if req.Location != nil {
		updates = append(updates, "location = ?")
		args = append(args, *req.Location)
	}
	if req.ImagePath != nil {
		updates = append(updates, "image_path = ?")
		args = append(args, *req.ImagePath)
	}

	if len(updates) == 0 {
		return existing, nil
	}

	updates = append(updates, "updated_at = ?")
	args = append(args, time.Now())
	args = append(args, id)

	query := fmt.Sprintf("UPDATE totes SET %s WHERE id = ?", joinStrings(updates, ", "))
	_, err = s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// Add new images if provided
	if len(req.ImagePaths) > 0 {
		for _, imagePath := range req.ImagePaths {
			_, err = s.AddImage(id, imagePath)
			if err != nil {
				// Log error but don't fail the update
				fmt.Printf("Warning: failed to add image: %v\n", err)
			}
		}
	}

	return s.GetByID(id)
}

// Delete deletes a tote by ID
func (s *ToteService) Delete(id int) error {
	// First, delete all sub-containers (children) of this tote
	// This ensures cascading delete works even if foreign keys aren't enabled
	_, err := s.db.Exec("DELETE FROM totes WHERE parent_id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting sub-containers: %w", err)
	}

	// Now delete the tote itself
	result, err := s.db.Exec("DELETE FROM totes WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tote not found")
	}

	return nil
}

// DeleteAll deletes all totes
func (s *ToteService) DeleteAll() (int, error) {
	result, err := s.db.Exec("DELETE FROM totes")
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// AddImage adds a new image to a tote
func (s *ToteService) AddImage(toteID int, imageData string) (*models.ToteImage, error) {
	// Get current max display order
	var maxOrder int
	err := s.db.QueryRow("SELECT COALESCE(MAX(display_order), -1) FROM tote_images WHERE tote_id = ?", toteID).Scan(&maxOrder)
	if err != nil {
		return nil, err
	}

	// Decode base64 image
	imageBlob, imageType, err := decodeBase64Image(imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	query := `
		INSERT INTO tote_images (tote_id, image_data, image_type, display_order, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := s.db.Exec(query, toteID, imageBlob, imageType, maxOrder+1, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Return the created image
	var img models.ToteImage
	var imgData []byte
	err = s.db.QueryRow(`
		SELECT id, tote_id, image_data, image_type, display_order, created_at
		FROM tote_images
		WHERE id = ?
	`, id).Scan(&img.ID, &img.ToteID, &imgData, &img.ImageType, &img.DisplayOrder, &img.CreatedAt)

	if err != nil {
		return nil, err
	}

	img.ImageData = "data:" + img.ImageType + ";base64," + base64Encode(imgData)
	return &img, nil
}

// DeleteImage deletes an image from a tote
func (s *ToteService) DeleteImage(imageID int) error {
	result, err := s.db.Exec("DELETE FROM tote_images WHERE id = ?", imageID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("image not found")
	}

	return nil
}

// GetImage retrieves a specific image
func (s *ToteService) GetImage(imageID int) (*models.ToteImage, error) {
	var img models.ToteImage
	var imageData []byte
	err := s.db.QueryRow(`
		SELECT id, tote_id, image_data, image_type, display_order, created_at
		FROM tote_images
		WHERE id = ?
	`, imageID).Scan(&img.ID, &img.ToteID, &imageData, &img.ImageType, &img.DisplayOrder, &img.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("image not found")
	}
	if err != nil {
		return nil, err
	}

	img.ImageData = "data:" + img.ImageType + ";base64," + base64Encode(imageData)
	return &img, nil
}
