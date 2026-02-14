package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"totetrax/internal/models"
	"totetrax/internal/service"
)

type Handler struct {
	toteService     *service.ToteService
	settingsService *service.SettingsService
}

func NewHandler(toteService *service.ToteService, settingsService *service.SettingsService) *Handler {
	return &Handler{
		toteService:     toteService,
		settingsService: settingsService,
	}
}

// IndexHandler serves the main web interface
func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	html := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>ToteTrax - Storage Container Inventory</title>
	<link rel="stylesheet" href="/static/css/style.css">
	<script>
		(function() {
			try {
				const settings = JSON.parse(localStorage.getItem('totetrax_settings') || '{}');
				if (settings.theme === 'dark') {
					document.documentElement.classList.add('dark-mode');
				}
			} catch(e) {}
		})();
	</script>
</head>
<body>
	<header>
		<div class="container">
			<div class="header-content">
				<div class="header-title">
					<h1>📦 ToteTrax</h1>
				</div>
				<div class="header-actions">
					<button class="btn btn-secondary" onclick="window.location.href='/scan'">
						📱 Scan QR
					</button>
					<button class="btn btn-primary" onclick="window.location.href='/add'">
						➕ Add Tote
					</button>
				</div>
			</div>
		</div>
	</header>

	<main class="container">
		<div class="stats-grid">
			<div class="stat-card">
				<div class="stat-header">
					<span class="stat-label">Total Totes</span>
					<span class="stat-icon">📦</span>
				</div>
				<div class="stat-value" id="total-totes">0</div>
				<div class="stat-description">Storage containers tracked</div>
			</div>
		</div>

		<div class="search-section">
			<input type="text" id="search" class="search-input" placeholder="Search totes by name or items...">
		</div>

		<div id="totes-grid" class="totes-grid">
			<div class="loading">Loading totes...</div>
		</div>

		<div id="empty-state" class="empty-state" style="display: none;">
			<div class="empty-icon">📦</div>
			<h2>No Totes Yet</h2>
			<p>Start organizing by adding your first storage tote</p>
			<button class="btn btn-primary" onclick="window.location.href='/add'">➕ Add First Tote</button>
		</div>
	</main>

	<script src="/static/js/app.js"></script>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// AddToteHandler serves the add tote form
func (h *Handler) AddToteHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Add Tote - ToteTrax</title>
	<link rel="stylesheet" href="/static/css/style.css">
	<script>
		(function() {
			try {
				const settings = JSON.parse(localStorage.getItem('totetrax_settings') || '{}');
				if (settings.theme === 'dark') {
					document.documentElement.classList.add('dark-mode');
				}
			} catch(e) {}
		})();
	</script>
</head>
<body>
	<header>
		<div class="container">
			<div class="header-content">
				<div class="header-title">
					<h1>📦 ToteTrax</h1>
				</div>
				<div class="header-actions">
					<button class="btn btn-secondary" onclick="window.location.href='/'">
						← Back to Home
					</button>
				</div>
			</div>
		</div>
	</header>

	<main class="container">
		<div class="form-container">
			<h2>Add New Tote</h2>
			<form id="tote-form">
				<div class="form-group">
					<label for="name">Tote Name *</label>
					<input type="text" id="name" name="name" required placeholder="e.g., Kitchen Supplies, Holiday Decorations">
				</div>

				<div class="form-group">
					<label for="description">Description</label>
					<input type="text" id="description" name="description" placeholder="Brief description of contents">
				</div>

				<div class="form-group">
					<label for="items">Items List</label>
					<textarea id="items" name="items" rows="6" placeholder="Enter items (one per line)&#10;Example:&#10;- 4x Dish towels&#10;- 2x Pot holders&#10;- 1x Apron"></textarea>
				</div>

				<div class="form-group">
					<label for="image">Tote Image</label>
					<input type="file" id="image" name="image" accept="image/*">
					<div id="image-preview" style="margin-top: 10px; display: none;">
						<img id="preview-img" style="max-width: 300px; max-height: 300px; border: 1px solid #ddd; border-radius: 4px;">
					</div>
				</div>

				<div class="form-actions">
					<button type="button" class="btn btn-secondary" onclick="window.location.href='/'">Cancel</button>
					<button type="submit" class="btn btn-primary">Save Tote</button>
				</div>
			</form>
		</div>
	</main>

	<script src="/static/js/form.js"></script>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// EditToteHandler serves the edit tote form
func (h *Handler) EditToteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing tote ID", http.StatusBadRequest)
		return
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Edit Tote - ToteTrax</title>
	<link rel="stylesheet" href="/static/css/style.css">
	<script>
		(function() {
			try {
				const settings = JSON.parse(localStorage.getItem('totetrax_settings') || '{}');
				if (settings.theme === 'dark') {
					document.documentElement.classList.add('dark-mode');
				}
			} catch(e) {}
		})();
	</script>
</head>
<body>
	<header>
		<div class="container">
			<div class="header-content">
				<div class="header-title">
					<h1>📦 ToteTrax</h1>
				</div>
				<div class="header-actions">
					<button class="btn btn-secondary" onclick="window.location.href='/tote/%s'">
						← Back to Details
					</button>
				</div>
			</div>
		</div>
	</header>

	<main class="container">
		<div class="form-container">
			<h2>Edit Tote</h2>
			<form id="tote-form">
				<input type="hidden" id="tote-id" value="%s">
				
				<div class="form-group">
					<label for="name">Tote Name *</label>
					<input type="text" id="name" name="name" required>
				</div>

				<div class="form-group">
					<label for="description">Description</label>
					<input type="text" id="description" name="description">
				</div>

				<div class="form-group">
					<label for="items">Items List</label>
					<textarea id="items" name="items" rows="6"></textarea>
				</div>

				<div class="form-group">
					<label for="image">Change Tote Image</label>
					<input type="file" id="image" name="image" accept="image/*">
					<div id="current-image" style="margin-top: 10px;"></div>
					<div id="image-preview" style="margin-top: 10px; display: none;">
						<img id="preview-img" style="max-width: 300px; max-height: 300px; border: 1px solid #ddd; border-radius: 4px;">
					</div>
				</div>

				<div class="form-actions">
					<button type="button" class="btn btn-secondary" onclick="window.location.href='/tote/%s'">Cancel</button>
					<button type="submit" class="btn btn-primary">Update Tote</button>
				</div>
			</form>
		</div>
	</main>

	<script src="/static/js/form.js"></script>
</body>
</html>
`, idStr, idStr, idStr)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// ToteDetailHandler serves the tote detail page
func (h *Handler) ToteDetailHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, "Invalid tote ID", http.StatusBadRequest)
		return
	}
	id := parts[1]

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Tote Details - ToteTrax</title>
	<link rel="stylesheet" href="/static/css/style.css">
	<script src="/static/js/qrcode.min.js"></script>
	<script>
		(function() {
			try {
				const settings = JSON.parse(localStorage.getItem('totetrax_settings') || '{}');
				if (settings.theme === 'dark') {
					document.documentElement.classList.add('dark-mode');
				}
			} catch(e) {}
		})();
	</script>
</head>
<body>
	<header>
		<div class="container">
			<div class="header-content">
				<div class="header-title">
					<h1>📦 ToteTrax</h1>
				</div>
				<div class="header-actions">
					<button class="btn btn-secondary" onclick="window.location.href='/'">
						← Back to Home
					</button>
					<button class="btn btn-secondary" onclick="window.location.href='/edit?id=%s'">
						✏️ Edit
					</button>
					<button class="btn btn-danger" onclick="deleteTote()">
						🗑️ Delete
					</button>
				</div>
			</div>
		</div>
	</header>

	<main class="container">
		<div id="tote-detail" class="detail-container">
			<div class="loading">Loading tote details...</div>
		</div>
	</main>

	<script src="/static/js/detail.js"></script>
	<script>
		const TOTE_ID = %s;
	</script>
</body>
</html>
`, id, id)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// ScanHandler serves the QR code scanning page
func (h *Handler) ScanHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Scan QR Code - ToteTrax</title>
	<link rel="stylesheet" href="/static/css/style.css">
	<script src="/static/js/html5-qrcode.min.js"></script>
	<script>
		(function() {
			try {
				const settings = JSON.parse(localStorage.getItem('totetrax_settings') || '{}');
				if (settings.theme === 'dark') {
					document.documentElement.classList.add('dark-mode');
				}
			} catch(e) {}
		})();
	</script>
</head>
<body>
	<header>
		<div class="container">
			<div class="header-content">
				<div class="header-title">
					<h1>📦 ToteTrax</h1>
				</div>
				<div class="header-actions">
					<button class="btn btn-secondary" onclick="window.location.href='/'">
						← Back to Home
					</button>
				</div>
			</div>
		</div>
	</header>

	<main class="container">
		<div class="scan-container">
			<h2>Scan Tote QR Code</h2>
			
			<div class="scan-methods">
				<div class="scan-method">
					<h3>📷 Method 1: Camera Scan</h3>
					<div id="qr-reader" style="width: 100%; max-width: 500px; margin: 20px auto;"></div>
					<div id="qr-reader-results"></div>
				</div>

				<div class="scan-method">
					<h3>🖼️ Method 2: Upload Image</h3>
					<input type="file" id="qr-file" accept="image/*" style="margin: 10px 0;">
					<div id="file-reader-results"></div>
				</div>

				<div class="scan-method">
					<h3>⌨️ Method 3: Manual Entry</h3>
					<input type="text" id="manual-code" placeholder="Enter QR code (e.g., TOTE-00001)" style="width: 100%; max-width: 300px; padding: 10px;">
					<button class="btn btn-primary" onclick="manualLookup()" style="margin-top: 10px;">Look Up</button>
				</div>
			</div>
		</div>
	</main>

	<script>
		let html5QrCode;
		
		function onScanSuccess(decodedText) {
			if (decodedText.startsWith('TOTE-')) {
				fetch('/api/tote/qr/' + decodedText)
					.then(response => response.json())
					.then(data => {
						if (data.id) {
							window.location.href = '/tote/' + data.id;
						}
					})
					.catch(error => {
						alert('Tote not found: ' + decodedText);
					});
			}
		}

		// Start camera scanning
		html5QrCode = new Html5Qrcode("qr-reader");
		html5QrCode.start(
			{ facingMode: "environment" },
			{ fps: 10, qrbox: 250 },
			onScanSuccess
		).catch(err => {
			document.getElementById('qr-reader-results').innerHTML = 
				'<p style="color: orange;">Camera not available. Use image upload or manual entry.</p>';
		});

		// File upload scanning
		document.getElementById('qr-file').addEventListener('change', function(e) {
			const file = e.target.files[0];
			if (file) {
				const tempReader = new Html5Qrcode("temp-reader");
				tempReader.scanFile(file, true)
					.then(decodedText => {
						onScanSuccess(decodedText);
					})
					.catch(err => {
						document.getElementById('file-reader-results').innerHTML = 
							'<p style="color: red;">Could not read QR code from image.</p>';
					});
			}
		});

		// Manual entry
		function manualLookup() {
			const code = document.getElementById('manual-code').value.trim();
			if (code) {
				onScanSuccess(code);
			}
		}
	</script>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// PrintLabelHandler serves the printable label page
func (h *Handler) PrintLabelHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, "Invalid tote ID", http.StatusBadRequest)
		return
	}
	id := parts[1]

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Print Label - ToteTrax</title>
	<script src="/static/js/qrcode.min.js"></script>
	<style>
		@page {
			size: 4in 6in;
			margin: 0;
		}
		* {
			margin: 0;
			padding: 0;
			box-sizing: border-box;
		}
		body {
			font-family: Arial, sans-serif;
			background: white;
			color: black;
		}
		.print-header {
			background: #333;
			color: white;
			padding: 15px;
			text-align: center;
		}
		.print-header button {
			background: white;
			color: #333;
			border: none;
			padding: 10px 20px;
			margin: 5px;
			cursor: pointer;
			border-radius: 4px;
		}
		.label {
			width: 4in;
			padding: 0.25in;
			background: white;
			color: black;
		}
		.label-header {
			text-align: center;
			margin-bottom: 10px;
			border-bottom: 2px solid #333;
			padding-bottom: 10px;
		}
		.label-qr {
			text-align: center;
			margin: 15px 0;
		}
		.label-info {
			font-size: 11pt;
		}
		.label-info h3 {
			font-size: 13pt;
			margin-bottom: 5px;
		}
		.label-items {
			margin-top: 10px;
			font-size: 10pt;
			white-space: pre-wrap;
		}
		@media print {
			.print-header {
				display: none;
			}
			body {
				print-color-adjust: exact;
				-webkit-print-color-adjust: exact;
			}
		}
	</style>
</head>
<body>
	<div class="print-header">
		<h2>Tote Label</h2>
		<button onclick="window.print()">🖨️ Print</button>
		<button onclick="history.back()">← Back</button>
	</div>

	<div class="label" id="label">
		<div class="label-header">
			<h2>📦 ToteTrax</h2>
			<div id="qr-code-text" style="font-size: 12pt; font-weight: bold; margin-top: 5px;"></div>
		</div>
		<div class="label-qr">
			<div id="qrcode" style="display: inline-block;"></div>
		</div>
		<div class="label-info" id="tote-info">
			<div style="text-align: center;">Loading...</div>
		</div>
	</div>

	<script>
		const TOTE_ID = %s;
		
		fetch('/api/tote/' + TOTE_ID)
			.then(response => response.json())
			.then(tote => {
				// Generate QR code
				new QRCode(document.getElementById('qrcode'), {
					text: tote.qr_code,
					width: 120,
					height: 120,
					colorDark: '#000000',
					colorLight: '#ffffff',
					correctLevel: QRCode.CorrectLevel.H
				});

				document.getElementById('qr-code-text').textContent = tote.qr_code;

				// Display tote info
				let itemsHtml = '';
				if (tote.items) {
					itemsHtml = '<div class="label-items"><strong>Items:</strong><br>' + 
						tote.items.split('\n').slice(0, 8).join('\n') + '</div>';
				}

				document.getElementById('tote-info').innerHTML = 
					'<h3>' + tote.name + '</h3>' +
					(tote.description ? '<p>' + tote.description + '</p>' : '') +
					itemsHtml;
			})
			.catch(error => {
				document.getElementById('tote-info').innerHTML = 
					'<p style="color: red;">Error loading tote</p>';
			});
	</script>
</body>
</html>
`, id)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// API Handlers

// TotesHandler handles GET /api/totes
func (h *Handler) TotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	totes, err := h.toteService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totes)
}

// ToteCreateHandler handles POST /api/tote
func (h *Handler) ToteCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.ToteCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tote, err := h.toteService.Create(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tote)
}

// ToteHandler handles GET/PUT/DELETE /api/tote/{id}
func (h *Handler) ToteHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid tote ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid tote ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		tote, err := h.toteService.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tote)

	case http.MethodPut:
		var req models.ToteUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tote, err := h.toteService.Update(id, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tote)

	case http.MethodDelete:
		if err := h.toteService.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ToteByQRCodeHandler handles GET /api/tote/qr/{qr_code}
func (h *Handler) ToteByQRCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid QR code", http.StatusBadRequest)
		return
	}

	qrCode := parts[3]
	tote, err := h.toteService.GetByQRCode(qrCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tote)
}

// ExportHandler handles GET /api/export
func (h *Handler) ExportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	totes, err := h.toteService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("totetrax-export-%s.json", time.Now().Format("2006-01-02"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	json.NewEncoder(w).Encode(totes)
}

// ImportHandler handles POST /api/import
func (h *Handler) ImportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var totes []models.Tote
	if err := json.NewDecoder(r.Body).Decode(&totes); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	imported := 0
	for _, tote := range totes {
		req := models.ToteCreateRequest{
			Name:        tote.Name,
			Description: tote.Description,
			Items:       tote.Items,
			ImagePath:   tote.ImagePath,
		}
		_, err := h.toteService.Create(req)
		if err == nil {
			imported++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"imported": imported})
}

// DeleteAllTotesHandler handles DELETE /api/totes/delete-all
func (h *Handler) DeleteAllTotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count, err := h.toteService.DeleteAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"deleted": count})
}

// APISettingsHandler handles GET/PUT /api/settings
func (h *Handler) APISettingsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		settings, err := h.settingsService.LoadSettings()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)

	case http.MethodPut:
		var settings models.Settings
		if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.settingsService.SaveSettings(&settings); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(settings)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// UploadImageHandler handles POST /api/upload-image
func (h *Handler) UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create uploads directory if it doesn't exist
	uploadDir := "web/static/images/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Return the file path
	relativePath := "/static/images/uploads/" + filename
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"path": relativePath})
}
