package api

import (
	"net/http"

	"totetrax/internal/service"
)

// NewRouter creates and configures the HTTP router
func NewRouter(toteService *service.ToteService, settingsService *service.SettingsService) http.Handler {
	mux := http.NewServeMux()

	// Create handler
	handler := NewHandler(toteService, settingsService)

	// API routes - Totes
	mux.HandleFunc("/api/totes", handler.TotesHandler)              // GET all totes
	mux.HandleFunc("/api/tote", handler.ToteCreateHandler)          // POST single tote
	mux.HandleFunc("/api/tote/", handler.ToteHandler)               // GET/PUT/DELETE by ID
	mux.HandleFunc("/api/tote/qr/", handler.ToteByQRCodeHandler)    // GET by QR code

	// API routes - Settings
	mux.HandleFunc("/api/settings", handler.APISettingsHandler)     // GET/PUT settings

	// API routes - Import/Export
	mux.HandleFunc("/api/export", handler.ExportHandler)            // GET export inventory
	mux.HandleFunc("/api/import", handler.ImportHandler)            // POST import inventory
	mux.HandleFunc("/api/totes/delete-all", handler.DeleteAllTotesHandler) // DELETE all totes

	// API routes - Image upload
	mux.HandleFunc("/api/upload-image", handler.UploadImageHandler) // POST upload image

	// Web UI routes
	mux.HandleFunc("/", handler.IndexHandler)
	mux.HandleFunc("/add", handler.AddToteHandler)
	mux.HandleFunc("/edit", handler.EditToteHandler)
	mux.HandleFunc("/tote/", handler.ToteDetailHandler)
	mux.HandleFunc("/scan", handler.ScanHandler)
	mux.HandleFunc("/print-label/", handler.PrintLabelHandler)

	// Static files
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	return mux
}
