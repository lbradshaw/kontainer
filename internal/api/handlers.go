package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"kontainer/internal/models"
	"kontainer/internal/service"
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

// convertImagePathsToURLs is no longer needed since images are base64 encoded in database
// but keeping for backward compatibility with legacy image_path field
func convertImagePathsToURLs(tote *models.Tote) {
	// Images are now base64 data URIs, no conversion needed
	// Legacy image_path can remain as-is for backward compatibility
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
	<title>Kontainer - Storage Container Inventory</title>
	<link rel="stylesheet" href="/static/css/style.css">
	<script>
		(function() {
			try {
				const settings = JSON.parse(localStorage.getItem('kontainer_settings') || '{}');
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
					<h1>📦 Kontainer</h1>
				</div>
				<div class="header-actions">
					<button class="btn btn-secondary" onclick="window.location.href='/scan'">
						🔍 Look up
					</button>
					<button class="btn btn-secondary" onclick="exportData()">
						📥 Export
					</button>
					<button class="btn btn-secondary" onclick="document.getElementById('import-file').click()">
						📤 Import
					</button>
					<button class="btn btn-primary" onclick="window.location.href='/add'">
						➕ Add Kontainer
					</button>
					<button class="btn btn-secondary" onclick="window.location.href='/settings'">
						⚙️ Settings
					</button>
				</div>
				<input type="file" id="import-file" accept=".json" style="display: none;" onchange="importData(event)">
			</div>
		</div>
	</header>

	<main class="container">
		<div class="stats-grid">
			<div class="stat-card">
				<div class="stat-header">
					<span class="stat-label">Total Containers</span>
					<span class="stat-icon">📦</span>
				</div>
				<div class="stat-value" id="total-totes">0</div>
				<div class="stat-description">Storage containers tracked</div>
			</div>
		</div>

		<div class="search-section">
			<input type="text" id="search" class="search-input" placeholder="Search kontainers by name or items...">
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
	<title>Add Kontainer</title>
	<link rel="stylesheet" href="/static/css/style.css">
	<script>
		(function() {
			try {
				const settings = JSON.parse(localStorage.getItem('kontainer_settings') || '{}');
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
					<h1>📦 Kontainer</h1>
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
			<h2>Add New Kontainer</h2>
			<form id="tote-form">
				<div class="form-group">
					<label for="name">Name (required)</label>
					<input type="text" id="name" name="name" required placeholder="e.g., Kitchen Supplies, Holiday Decorations">
				</div>

				<div class="form-group">
					<label for="description">Description</label>
					<input type="text" id="description" name="description" placeholder="Brief description of contents">
				</div>

				<div class="form-group">
					<label for="location">Location</label>
					<input type="text" id="location" name="location" placeholder="e.g., Garage, Basement, Storage Unit A">
				</div>

				<div class="form-group">
					<label for="items">Items List</label>
					<textarea id="items" name="items" rows="6" placeholder="Enter items (one per line)&#10;Example:&#10;- 4x Dish towels&#10;- 2x Pot holders&#10;- 1x Apron"></textarea>
				</div>

				<div class="form-group">
					<label for="image">Images</label>
					<input type="file" id="image" name="image" accept="image/*" multiple>
					<p style="font-size: 0.85rem; opacity: 0.7; margin-top: 0.3rem;">Select one or more images (Ctrl+Click or Shift+Click)</p>
					<div id="image-preview" style="margin-top: 10px; display: none;"></div>
				</div>

				<div class="form-actions">
					<button type="button" class="btn btn-secondary" onclick="window.location.href='/'">Cancel</button>
					<button type="submit" class="btn btn-primary">Save</button>
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
	<title>Edit Kontainer</title>
	<link rel="stylesheet" href="/static/css/style.css">
	<script>
		(function() {
			try {
				const settings = JSON.parse(localStorage.getItem('kontainer_settings') || '{}');
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
					<h1>📦 Kontainer</h1>
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
			<h2>Edit Kontainer</h2>
			<form id="tote-form">
				<input type="hidden" id="tote-id" value="%s">
				
				<div class="form-group">
					<label for="name">Name (required)</label>
					<input type="text" id="name" name="name" required>
				</div>

				<div class="form-group">
					<label for="description">Description</label>
					<input type="text" id="description" name="description">
				</div>

				<div class="form-group">
					<label for="location">Location</label>
					<input type="text" id="location" name="location" placeholder="e.g., Garage, Basement, Storage Unit A">
				</div>

				<div class="form-group">
					<label for="items">Items List</label>
					<textarea id="items" name="items" rows="6"></textarea>
				</div>

				<div class="form-group">
					<label for="image">Add More Images</label>
					<input type="file" id="image" name="image" accept="image/*" multiple>
					<p style="font-size: 0.85rem; color: #28a745; margin-top: 0.3rem;">
						✓ Select additional images to add (Ctrl+Click or Shift+Click)<br>
						✓ Existing images will be kept - new images will be added to the gallery
					</p>
					<div id="current-image" style="margin-top: 10px;"></div>
					<div id="image-preview" style="margin-top: 10px; display: none;"></div>
				</div>

				<div class="form-actions">
					<button type="button" class="btn btn-secondary" onclick="window.location.href='/tote/%s'">Cancel</button>
					<button type="submit" class="btn btn-primary">Update</button>
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
	<title>Kontainer Details</title>
	<link rel="stylesheet" href="/static/css/style.css">
	<script src="/static/js/qrcode.min.js"></script>
	<script>
		(function() {
			try {
				const settings = JSON.parse(localStorage.getItem('kontainer_settings') || '{}');
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
					<h1>📦 Kontainer</h1>
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
	<title>Scan QR Code - Kontainer</title>
	<link rel="stylesheet" href="/static/css/style.css">
	<script src="/static/js/html5-qrcode.min.js"></script>
	<script>
		(function() {
			try {
				const settings = JSON.parse(localStorage.getItem('kontainer_settings') || '{}');
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
					<h1>📦 Kontainer</h1>
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
			<h2>Look up Tote</h2>
			
			<div class="scan-methods">
				<div class="scan-method">
					<h3>🖼️ Method 1: Upload Image</h3>
					<input type="file" id="qr-file" accept="image/*" style="margin: 10px 0;">
					<div id="file-reader-results"></div>
				</div>

				<div class="scan-method">
					<h3>⌨️ Method 2: Manual Entry</h3>
					<input type="text" id="manual-code" placeholder="Enter QR code (e.g., TOTE-00001)" style="width: 100%; max-width: 300px; padding: 10px;">
					<button class="btn btn-primary" onclick="manualLookup()" style="margin-top: 10px;">Look Up</button>
				</div>
			</div>
		</div>
	</main>

	<script>
		document.addEventListener('DOMContentLoaded', function() {
			console.log('Page loaded, initializing QR scanner');
			
			function onScanSuccess(decodedText) {
				console.log('Decoded QR code:', decodedText);
				if (decodedText.startsWith('TOTE-')) {
					fetch('/api/tote/qr/' + decodedText)
						.then(response => {
							console.log('API response status:', response.status);
							return response.json();
						})
						.then(data => {
							console.log('API data:', data);
							if (data.id) {
								window.location.href = '/tote/' + data.id;
							} else {
								alert('Tote not found: ' + decodedText);
							}
						})
						.catch(error => {
							console.error('API error:', error);
							alert('Tote not found: ' + decodedText);
						});
				} else {
					alert('Invalid QR code format. Expected TOTE-XXXXX, got: ' + decodedText);
				}
			}

			// File upload scanning - using html5-qrcode library
			const fileInput = document.getElementById('qr-file');
			console.log('File input element:', fileInput);
			
			fileInput.addEventListener('change', function(e) {
				console.log('File input changed');
				const file = e.target.files[0];
				const resultsDiv = document.getElementById('file-reader-results');
				
				if (file) {
					console.log('File selected:', file.name, file.type, file.size);
					resultsDiv.innerHTML = '<p style="color: blue;">Processing image...</p>';
					
					// Create a temporary scanner instance
					const html5QrCode = new Html5Qrcode("file-reader-results");
					
					html5QrCode.scanFile(file, true)
						.then(decodedText => {
							console.log('Scan success:', decodedText);
							resultsDiv.innerHTML = '<p style="color: green;">QR Code found: ' + decodedText + '</p>';
							// Clear and redirect
							try {
								html5QrCode.clear();
							} catch(e) {
								console.log('Clear error (ignorable):', e);
							}
							setTimeout(() => onScanSuccess(decodedText), 500);
						})
						.catch(err => {
							console.error('Scan error:', err);
							resultsDiv.innerHTML = '<p style="color: red;">Could not read QR code from this image. Please try a clearer image or use Manual Entry below.</p>';
							try {
								html5QrCode.clear();
							} catch(e) {
								console.log('Clear error (ignorable):', e);
							}
						});
				} else {
					console.log('No file selected');
				}
			});

			// Manual entry
			window.manualLookup = function() {
				const code = document.getElementById('manual-code').value.trim();
				if (code) {
					onScanSuccess(code);
				} else {
					alert('Please enter a QR code');
				}
			};

			// Allow Enter key on manual input
			document.getElementById('manual-code').addEventListener('keypress', function(e) {
				if (e.key === 'Enter') {
					manualLookup();
				}
			});
		});
	</script>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// SettingsPageHandler serves the settings configuration page
func (h *Handler) SettingsPageHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Settings - Kontainer</title>
	<link rel="stylesheet" href="/static/css/style.css">
	<script>
		(function() {
			try {
				const settings = JSON.parse(localStorage.getItem('kontainer_settings') || '{}');
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
					<h1>📦 Kontainer</h1>
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
		<div class="settings-container" style="max-width: 800px; margin: 0 auto;">
			<h2>⚙️ Settings</h2>

			<div class="settings-section">
				<h3>Application Settings</h3>
				<p style="color: #f39c12; margin-bottom: 1rem;">⚠️ Changing these settings requires restarting the application</p>

				<div class="form-group">
					<label for="port">Server Port</label>
					<input type="number" id="port" name="port" min="1024" max="65535" style="max-width: 200px;">
					<p style="font-size: 0.85rem; opacity: 0.7; margin-top: 0.3rem;">
						Current port. Change and restart to apply. (Default: 3818)
					</p>
				</div>

				<div class="form-group">
					<label for="database_path">Database File Path</label>
					<input type="text" id="database_path" name="database_path" placeholder="kontainer.db">
					<p style="font-size: 0.85rem; opacity: 0.7; margin-top: 0.3rem;">
						Local path to SQLite database file. For NAS storage, run Kontainer on the NAS itself via Docker.
					</p>
				</div>
			</div>

			<div class="settings-section" style="margin-top: 2rem;">
				<h3>Appearance</h3>

				<div class="form-group">
					<label for="theme">Theme</label>
					<select id="theme" name="theme" style="max-width: 200px;">
						<option value="light">Light</option>
						<option value="dark">Dark</option>
					</select>
					<p style="font-size: 0.85rem; opacity: 0.7; margin-top: 0.3rem;">
						Choose your preferred color theme
					</p>
				</div>
			</div>

			<div class="settings-section" style="margin-top: 2rem;">
				<h3>Data Management</h3>

				<div class="form-group">
					<label>Backup & Restore</label>
					<div style="display: flex; gap: 1rem; flex-wrap: wrap; margin-top: 0.5rem;">
						<button class="btn btn-primary" onclick="window.location.href='/api/export'">
							📥 Export All Data
						</button>
						<button class="btn btn-secondary" onclick="document.getElementById('settings-import-file').click()">
							📤 Import Data
						</button>
						<button class="btn btn-danger" onclick="deleteAllData()">
							🗑️ Delete All Totes
						</button>
					</div>
					<input type="file" id="settings-import-file" accept=".json" style="display: none;" onchange="importDataFromSettings(event)">
					<p style="font-size: 0.85rem; opacity: 0.7; margin-top: 0.5rem;">
						Export creates a JSON backup. Import adds totes from backup file.
					</p>
				</div>
			</div>

			<div class="form-actions" style="margin-top: 2rem; display: flex; gap: 1rem;">
				<button class="btn btn-primary" onclick="saveSettings()">💾 Save Settings</button>
				<button class="btn btn-secondary" onclick="resetSettings()">🔄 Reset to Defaults</button>
			</div>

			<div class="settings-section" style="margin-top: 3rem; padding: 1.5rem; background: var(--card-bg); border-radius: 8px; border-left: 4px solid #3498db;">
				<h3>ℹ️ Current Configuration</h3>
				<div id="current-config" style="font-family: monospace; font-size: 0.9rem; margin-top: 1rem;">
					Loading...
				</div>
			</div>
		</div>
	</main>

	<script src="/static/js/settings-page.js"></script>
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
	<title>Print Label - Kontainer</title>
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
		<h2>QR Label</h2>
		<button onclick="window.print()">🖨️ Print</button>
		<button onclick="history.back()">← Back</button>
	</div>

	<div class="label" id="label">
		<div class="label-header">
			<h2>📦 Kontainer</h2>
		</div>
		<div id="tote-name" style="font-size: 14pt; font-weight: bold; text-align: center; margin: 15px 0 10px 0;">Loading...</div>
		<div class="label-qr">
			<div id="qrcode" style="display: inline-block;"></div>
			<div id="qr-code-text" style="font-size: 11pt; font-weight: bold; margin-top: 8px; text-align: center;"></div>
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

				// Display QR code text below QR code
				document.getElementById('qr-code-text').textContent = tote.qr_code;

				// Display the name above QR code
				document.getElementById('tote-name').textContent = tote.name;
			})
			.catch(error => {
				document.getElementById('tote-name').textContent = 'Error loading kontainer';
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

	// Convert file paths to web URLs
	for i := range totes {
		convertImagePathsToURLs(&totes[i])
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totes)
}

// AllTotesHandler handles GET /api/totes/all - returns ALL totes including sub-containers
func (h *Handler) AllTotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	totes, err := h.toteService.GetAllIncludingChildren()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert file paths to web URLs
	for i := range totes {
		convertImagePathsToURLs(&totes[i])
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

	convertImagePathsToURLs(tote)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tote)
}

// ToteHandler handles GET/PUT/DELETE /api/tote/{id} and POST /api/tote/{id}/add-image
func (h *Handler) ToteHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid tote ID", http.StatusBadRequest)
		return
	}

	// Check if it's an add-image request
	if len(parts) == 4 && parts[3] == "add-image" && r.Method == http.MethodPost {
		h.AddImageToToteHandler(w, r)
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
		convertImagePathsToURLs(tote)
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

		convertImagePathsToURLs(tote)
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

	convertImagePathsToURLs(tote)
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

	filename := fmt.Sprintf("kontainer-export-%s.json", time.Now().Format("2006-01-02"))
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
		// Prepare image data arrays from the imported tote
		var imagePaths []string
		var imageTypes []string
		
		if tote.Images != nil && len(tote.Images) > 0 {
			// Use images array if available (base64 data URIs)
			for _, img := range tote.Images {
				imagePaths = append(imagePaths, img.ImageData)
				imageTypes = append(imageTypes, img.ImageType)
			}
		} else if tote.ImagePath != "" {
			// Fallback to single image_path for backward compatibility
			imagePaths = []string{tote.ImagePath}
		}

		req := models.ToteCreateRequest{
			Name:        tote.Name,
			Description: tote.Description,
			Items:       tote.Items,
			Location:    tote.Location,
			ImagePath:   tote.ImagePath,
			ImagePaths:  imagePaths,
			ImageTypes:  imageTypes,
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
		var newSettings models.Settings
		if err := json.NewDecoder(r.Body).Decode(&newSettings); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Load current settings to check if database path changed
		currentSettings, err := h.settingsService.LoadSettings()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if database path changed
		if currentSettings.DatabasePath != newSettings.DatabasePath && newSettings.DatabasePath != "" {
			// Migrate database to new location
			oldPath := currentSettings.DatabasePath
			if oldPath == "" {
				oldPath = "kontainer.db" // Default path
			}
			
			if err := h.settingsService.MigrateDatabase(oldPath, newSettings.DatabasePath); err != nil {
				http.Error(w, fmt.Sprintf("Failed to migrate database: %v", err), http.StatusInternalServerError)
				return
			}
		}

		// Save new settings
		if err := h.settingsService.SaveSettings(&newSettings); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newSettings)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// UploadImageHandler is now deprecated - images are sent as base64 in request body
// Keeping for backward compatibility
func (h *Handler) UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Image upload deprecated - send base64 data in request body", http.StatusGone)
}

// AddImageToToteHandler handles POST /api/tote/{id}/add-image
func (h *Handler) AddImageToToteHandler(w http.ResponseWriter, r *http.Request) {
	// Extract tote ID from path
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid tote ID", http.StatusBadRequest)
		return
	}

	toteID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid tote ID", http.StatusBadRequest)
		return
	}

	// Get image data from request body (base64 data URI)
	var req struct {
		ImageData string `json:"image_data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add image to tote
	image, err := h.toteService.AddImage(toteID, req.ImageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(image)
}

// ToteImageHandler handles DELETE /api/tote-image/{id}
func (h *Handler) ToteImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	imageID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	if err := h.toteService.DeleteImage(imageID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ServeImageHandler is now deprecated - images are served as base64 data URIs
// Keeping for backward compatibility
func (h *Handler) ServeImageHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Image serving deprecated - images are base64 encoded", http.StatusGone)
}
