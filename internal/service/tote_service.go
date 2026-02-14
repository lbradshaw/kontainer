package service

import (
	"database/sql"
	"fmt"
	"time"

	"totetrax/internal/models"
)

type ToteService struct {
	db *sql.DB
}

func NewToteService(db *sql.DB) *ToteService {
	return &ToteService{db: db}
}

// GetAll retrieves all totes
func (s *ToteService) GetAll() ([]models.Tote, error) {
	query := `
		SELECT id, name, description, items, image_path, qr_code, created_at, updated_at
		FROM totes
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
			&t.ID, &t.Name, &t.Description, &t.Items, &t.ImagePath,
			&t.QRCode, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		totes = append(totes, t)
	}

	return totes, nil
}

// GetByID retrieves a tote by ID
func (s *ToteService) GetByID(id int) (*models.Tote, error) {
	query := `
		SELECT id, name, description, items, image_path, qr_code, created_at, updated_at
		FROM totes
		WHERE id = ?
	`

	var t models.Tote
	err := s.db.QueryRow(query, id).Scan(
		&t.ID, &t.Name, &t.Description, &t.Items, &t.ImagePath,
		&t.QRCode, &t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tote not found")
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// GetByQRCode retrieves a tote by QR code
func (s *ToteService) GetByQRCode(qrCode string) (*models.Tote, error) {
	query := `
		SELECT id, name, description, items, image_path, qr_code, created_at, updated_at
		FROM totes
		WHERE qr_code = ?
	`

	var t models.Tote
	err := s.db.QueryRow(query, qrCode).Scan(
		&t.ID, &t.Name, &t.Description, &t.Items, &t.ImagePath,
		&t.QRCode, &t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tote not found")
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Create creates a new tote
func (s *ToteService) Create(req models.ToteCreateRequest) (*models.Tote, error) {
	// Generate QR code (TOTE-XXXXX format)
	var maxID int
	err := s.db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM totes").Scan(&maxID)
	if err != nil {
		return nil, err
	}
	qrCode := fmt.Sprintf("TOTE-%05d", maxID+1)

	query := `
		INSERT INTO totes (name, description, items, image_path, qr_code, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := s.db.Exec(query, req.Name, req.Description, req.Items, req.ImagePath, qrCode, now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetByID(int(id))
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

	return s.GetByID(id)
}

// Delete deletes a tote by ID
func (s *ToteService) Delete(id int) error {
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
