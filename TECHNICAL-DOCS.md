# Kontainer - Complete Technical Documentation

**Project Name:** Kontainer  
**Version:** 1.7.0  
**Created:** 2026-02-14  
**Last Updated:** 2026-02-19  
**Technology:** Go 1.24.0 + SQLite + HTML/CSS/JavaScript  
**Purpose:** Storage container inventory management with database-embedded images and mobile companion app

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture](#architecture)
3. [Database Schema](#database-schema)
4. [Backend Implementation](#backend-implementation)
5. [API Endpoints](#api-endpoints)
6. [Frontend Implementation](#frontend-implementation)
7. [Image Handling](#image-handling)
8. [Settings & Configuration](#settings--configuration)
9. [Development Guide](#development-guide)
10. [Deployment](#deployment)

---

## Project Overview

### Purpose
Kontainer helps users organize and track items stored in boxes, totes, bins, and containers. It's designed for home organization, moving, storage, and inventory management.

### Key Features
- ✅ Track storage containers with names, descriptions, and item lists
- ✅ **Location tracking** (optional physical location field for each container)
- ✅ **Multiple images per container** (unlimited)
- ✅ **Images stored in database as BLOBs** (no orphaned files, fully portable)
- ✅ Auto-generated QR codes (TOTE-XXXXX format)
- ✅ QR code lookup via image upload or manual entry (web UI)
- ✅ Printable labels with QR codes
- ✅ Search across names, descriptions, and item lists
- ✅ **Import/Export functionality with full image support** (JSON with base64 images)
- ✅ **Dashboard quick access** (Export/Import buttons)
- ✅ **Settings page data management** (Export/Import/Delete All)
- ✅ Image gallery with individual delete capability
- ✅ **Additive image uploads** (edit mode adds to existing images)
- ✅ **Configurable settings page** (port, database path, theme)
- ✅ **Light/Dark theme toggle** (instant application)
- ✅ **Network database support** (NAS/network drive compatible)
- ✅ **Automatic orphan prevention** (images deleted with tote via CASCADE)
- ✅ **Image hover gallery** (hover over card image to preview all images)
- ✅ **Modal image viewer** (click on gallery thumbnails for full-size view)
- ✅ **Full card clickable** (entire tote card navigates to edit page)

### Technology Stack
- **Backend:** Go 1.24.0 (single binary, no CGO dependencies)
- **Database:** SQLite with pure Go driver (`modernc.org/sqlite`)
- **Frontend:** Vanilla HTML/CSS/JavaScript
- **QR Codes:** `qrcode.min.js` (generation), `html5-qrcode.min.js` (scanning)
- **Web Server:** Go standard library `net/http`
- **Default Port:** 3818
- **Docker:** Multi-stage build with Alpine Linux (CGO_ENABLED=0 for pure Go)

---

## Architecture

### Project Structure
```
D:\projects\kontainer\
├── cmd\
│   └── totetrax\
│       └── main.go                 # Application entry point
├── internal\
│   ├── api\
│   │   ├── handlers.go            # HTTP request handlers (600+ lines)
│   │   └── router.go              # Route definitions
│   ├── database\
│   │   └── database.go            # SQLite initialization & schema
│   ├── models\
│   │   ├── tote.go                # Tote and ToteImage models
│   │   └── settings.go            # Application settings model
│   └── service\
│       ├── tote_service.go        # Business logic for totes & images
│       └── settings_service.go    # Settings management
├── web\
│   └── static\
│       ├── css\
│       │   └── style.css          # Application styles (dark mode support)
│       ├── js\
│       │   ├── app.js             # Dashboard functionality
│       │   ├── form.js            # Add/Edit forms with multi-image upload
│       │   ├── detail.js          # Tote detail page with image gallery
│       │   ├── settings-page.js   # Settings page functionality
│       │   ├── qrcode.min.js      # QR code generation
│       │   └── html5-qrcode.min.js # QR code scanning
│       └── images\
│           └── uploads\           # User-uploaded images
├── go.mod                          # Go module definition
├── go.sum                          # Dependency checksums
├── kontainer.exe                   # Compiled Windows binary
├── kontainer.db                    # SQLite database (auto-created)
├── settings.json                   # Runtime settings (auto-created)
├── README.md                       # User documentation
├── QUICKSTART.md                   # Quick reference guide
├── totetrax-project-info.md        # Project context (legacy name)
├── totetrax-technical-docs.md      # This file (legacy name)
└── .gitignore                      # Git ignore rules
```

### Design Patterns
- **Repository Pattern:** Service layer abstracts database access
- **Handler Pattern:** API handlers separated from business logic
- **Model-View-Controller (MVC):** Models, handlers (controllers), HTML (views)
- **API-First Design:** REST API with web UI as client

---

## Database Schema

### Tables

#### `totes` Table
Primary table for storage containers.

```sql
CREATE TABLE IF NOT EXISTS totes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    items TEXT,                              -- Newline-separated item list
    location TEXT,                           -- Optional physical location (e.g., "Garage", "Basement")
    image_path TEXT,                         -- Legacy field (first image for backward compatibility)
    qr_code TEXT UNIQUE NOT NULL,            -- Format: TOTE-00001, TOTE-00002, etc.
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_name ON totes(name);
CREATE INDEX IF NOT EXISTS idx_qr_code ON totes(qr_code);
```

**Key Points:**
- `qr_code` is auto-generated sequentially (TOTE-XXXXX)
- `location` is optional and stores where the container is physically located
- `image_path` kept for backward compatibility (stores first image)
- `items` is a multi-line text field (one item per line)

#### `tote_images` Table
Stores multiple images per tote as embedded binary data.

```sql
CREATE TABLE IF NOT EXISTS tote_images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tote_id INTEGER NOT NULL,
    image_data BLOB NOT NULL,                -- Binary image data
    image_type TEXT NOT NULL,                -- MIME type: image/jpeg, image/png, etc.
    display_order INTEGER NOT NULL DEFAULT 0, -- Controls sort order in gallery
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tote_id) REFERENCES totes(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_tote_images_tote_id ON tote_images(tote_id);
```

**Key Points:**
- **CASCADE DELETE:** Deleting a tote removes all its images (no orphans!)
- `image_data` stores binary BLOB data directly in database
- `image_type` stores MIME type for proper rendering
- `display_order` allows manual reordering (currently auto-assigned)
- Multiple images per tote (no limit)
- **No separate file storage needed** - fully self-contained database

### Database Migrations
- Schema auto-creates on first run via `database.InitDB()`
- **Automatic migration system** adds new columns to existing databases
- `location` column automatically added if missing (v1.7.0+)
- Backward compatible: Adding `tote_images` table doesn't break existing data
- No migration scripts needed - SQLite handles `CREATE TABLE IF NOT EXISTS`

---

## Backend Implementation

### Entry Point: `cmd/kontainer/main.go`

```go
func main() {
    // 1. Load settings (creates default settings.json if missing)
    settingsService := service.NewSettingsService()
    settings, err := settingsService.LoadSettings()
    
    // 2. Initialize SQLite database
    db, err := database.InitDB("totetrax.db")
    defer db.Close()
    
    // 3. Create service layer
    toteService := service.NewToteService(db)
    
    // 4. Initialize router with handlers
    router := api.NewRouter(toteService, settingsService)
    
    // 5. Start HTTP server
    http.ListenAndServe(fmt.Sprintf(":%d", settings.Port), router)
}
```

### Service Layer: `internal/service/tote_service.go`

**Core Methods:**

```go
type ToteService struct {
    db *sql.DB
}

// CRUD Operations
func (s *ToteService) GetAll() ([]models.Tote, error)
func (s *ToteService) GetByID(id int) (*models.Tote, error)
func (s *ToteService) GetByQRCode(qrCode string) (*models.Tote, error)
func (s *ToteService) Create(req models.ToteCreateRequest) (*models.Tote, error)
func (s *ToteService) Update(id int, req models.ToteUpdateRequest) (*models.Tote, error)
func (s *ToteService) Delete(id int) error
func (s *ToteService) DeleteAll() (int, error)

// Image Management
func (s *ToteService) AddImage(toteID int, imagePath string) (*models.ToteImage, error)
func (s *ToteService) DeleteImage(imageID int) error
func (s *ToteService) GetImage(imageID int) (*models.ToteImage, error)

// Private helpers
func (s *ToteService) loadImagesForTote(toteID int) ([]models.ToteImage, error)
func (s *ToteService) loadImagesForTotes(totes []models.Tote) ([]models.Tote, error)
```

**Key Implementation Details:**

1. **QR Code Generation:**
   ```go
   // Get max ID and generate next QR code
   var maxID int
   db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM totes").Scan(&maxID)
   qrCode := fmt.Sprintf("TOTE-%05d", maxID+1) // TOTE-00001
   ```

2. **Image Loading:**
   - `GetByID()`, `GetByQRCode()`, `GetAll()` all load associated images
   - Uses `loadImagesForTote()` to fetch images for each tote
   - Images sorted by `display_order ASC, created_at ASC`

3. **Create with Multiple Images:**
   ```go
   // Create tote first
   result := db.Exec("INSERT INTO totes ...")
   toteID := result.LastInsertId()
   
   // Insert each image
   for i, imagePath := range req.ImagePaths {
       db.Exec("INSERT INTO tote_images (tote_id, image_path, display_order) VALUES (?, ?, ?)",
           toteID, imagePath, i)
   }
   ```

4. **AddImage (for edit mode):**
   ```go
   // Get current max display_order
   var maxOrder int
   db.QueryRow("SELECT COALESCE(MAX(display_order), -1) FROM tote_images WHERE tote_id = ?", toteID).Scan(&maxOrder)
   
   // Insert with next order
   db.Exec("INSERT INTO tote_images (tote_id, image_path, display_order) VALUES (?, ?, ?)",
       toteID, imagePath, maxOrder+1)
   ```

### Models: `internal/models/tote.go`

```go
// ToteImage represents a single image
type ToteImage struct {
    ID           int       `json:"id"`
    ToteID       int       `json:"tote_id"`
    ImagePath    string    `json:"image_path"`
    DisplayOrder int       `json:"display_order"`
    CreatedAt    time.Time `json:"created_at"`
}

// Tote represents a storage container
type Tote struct {
    ID          int         `json:"id"`
    Name        string      `json:"name"`
    Description string      `json:"description"`
    Items       string      `json:"items"`
    ImagePath   string      `json:"image_path"`   // Legacy: first image
    Images      []ToteImage `json:"images"`       // NEW: all images
    QRCode      string      `json:"qr_code"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
}

// ToteCreateRequest for POST /api/tote
type ToteCreateRequest struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Items       string   `json:"items"`
    ImagePath   string   `json:"image_path"`    // Legacy: single image
    ImagePaths  []string `json:"image_paths"`   // NEW: multiple images
}

// ToteUpdateRequest for PUT /api/tote/{id}
type ToteUpdateRequest struct {
    Name        *string `json:"name,omitempty"`
    Description *string `json:"description,omitempty"`
    Items       *string `json:"items,omitempty"`
    ImagePath   *string `json:"image_path,omitempty"`
}
```

**Design Notes:**
- `Images` array always populated on GET requests
- `ImagePath` maintained for backward compatibility
- Update request uses pointers for optional fields (partial updates)

---

## API Endpoints

### Tote Management

#### `GET /api/totes`
List all totes with images.

**Response:**
```json
[
  {
    "id": 1,
    "name": "Kitchen Supplies",
    "description": "Extra kitchen items",
    "items": "4x Dish towels\n2x Pot holders\n1x Apron",
    "image_path": "/static/images/uploads/12345.jpg",
    "images": [
      {
        "id": 1,
        "tote_id": 1,
        "image_path": "/static/images/uploads/12345.jpg",
        "display_order": 0,
        "created_at": "2026-02-14T05:00:00Z"
      },
      {
        "id": 2,
        "tote_id": 1,
        "image_path": "/static/images/uploads/12346.jpg",
        "display_order": 1,
        "created_at": "2026-02-14T05:01:00Z"
      }
    ],
    "qr_code": "TOTE-00001",
    "created_at": "2026-02-14T05:00:00Z",
    "updated_at": "2026-02-14T05:01:00Z"
  }
]
```

#### `POST /api/tote`
Create new tote with multiple images.

**Request:**
```json
{
  "name": "Garage Tools",
  "description": "Hand tools and hardware",
  "items": "2x Hammers\n1x Screwdriver set\n1x Drill bits",
  "image_paths": [
    "/static/images/uploads/12347.jpg",
    "/static/images/uploads/12348.jpg",
    "/static/images/uploads/12349.jpg"
  ]
}
```

**Response:** Same as GET (newly created tote)

#### `GET /api/tote/{id}`
Get single tote by ID.

**Response:** Single tote object (same structure as array item above)

#### `GET /api/tote/qr/{qr_code}`
Get tote by QR code.

**Example:** `GET /api/tote/qr/TOTE-00001`

**Response:** Single tote object

#### `PUT /api/tote/{id}`
Update tote details (does NOT affect images).

**Request:**
```json
{
  "name": "Updated Name",
  "description": "Updated description",
  "items": "New item list"
}
```

**Response:** Updated tote object

#### `DELETE /api/tote/{id}`
Delete tote and all associated images (CASCADE).

**Response:** 204 No Content

### Image Management

#### `POST /api/upload-image`
Upload image file (returns path for use in tote creation).

**Request:** `multipart/form-data` with `image` field

**Response:**
```json
{
  "path": "/static/images/uploads/1739497234567.jpg"
}
```

**Implementation:**
- Generates unique filename using timestamp: `time.Now().UnixNano() + extension`
- Saves to `web/static/images/uploads/`
- Max file size: 10MB
- Accepts: `image/*`

#### `POST /api/tote/{id}/add-image`
Add image to existing tote (edit mode).

**Request:**
```json
{
  "image_path": "/static/images/uploads/12350.jpg"
}
```

**Response:**
```json
{
  "id": 5,
  "tote_id": 1,
  "image_path": "/static/images/uploads/12350.jpg",
  "display_order": 3,
  "created_at": "2026-02-14T06:00:00Z"
}
```

**Key Behavior:** This is **additive** - adds to existing images, never replaces.

#### `DELETE /api/tote-image/{id}`
Delete specific image by image ID.

**Response:** 204 No Content

**Note:** Does NOT delete the physical file (intentional - files may be referenced elsewhere)

### Import/Export

#### `GET /api/export`
Export all totes as JSON with full image arrays.

**Response:** Array of tote objects (downloads as `kontainer-export-YYYY-MM-DD.json`)

**Export Structure:**
```json
[
  {
    "id": 1,
    "name": "Kitchen Supplies",
    "description": "Kitchen items",
    "items": "4x Dish towels\n2x Pot holders",
    "image_path": "",
    "images": [
      {
        "id": 3,
        "tote_id": 1,
        "image_path": "/static/images/uploads/1771058813249973700.jpg",
        "display_order": 1,
        "created_at": "2026-02-14T00:00:00Z"
      },
      {
        "id": 4,
        "tote_id": 1,
        "image_path": "/static/images/uploads/1771058813255059400.png",
        "display_order": 2,
        "created_at": "2026-02-14T00:00:01Z"
      }
    ],
    "qr_code": "TOTE-00001",
    "created_at": "2026-02-14T00:00:00Z",
    "updated_at": "2026-02-14T00:00:00Z"
  }
]
```

**Implementation:**
- Uses `GetAll()` which loads all images via `loadImagesForTotes()`
- Includes complete images array with all metadata
- Filename format: `kontainer-export-2026-02-14.json`
- Content-Type: `application/json`
- Content-Disposition: `attachment` (auto-download)

#### `POST /api/import`
Import totes from JSON array with full image support.

**Request:** Array of tote objects (same format as export)

**Response:**
```json
{
  "imported": 5
}
```

**Behavior:** 
- **Additive operation** - Creates new totes, does not update existing
- **Multiple images support** - Extracts from `images` array
- **Backward compatible** - Falls back to single `image_path` if `images` is null
- **Image preservation** - All images recreated in new totes via `ImagePaths`

**Image Handling Logic:**
```go
var imagePaths []string
if tote.Images != nil && len(tote.Images) > 0 {
    // Use images array (preferred)
    for _, img := range tote.Images {
        imagePaths = append(imagePaths, img.ImagePath)
    }
} else if tote.ImagePath != "" {
    // Fallback to single image_path
    imagePaths = []string{tote.ImagePath}
}
```

**Important Notes:**
- Image files must exist at the paths specified in the JSON
- Import does NOT copy/move image files - assumes they exist
- For cross-system imports, ensure image files are copied to `web/static/images/uploads/`
- Image paths are relative: `/static/images/uploads/filename.ext`

#### `DELETE /api/totes/delete-all`
Delete all totes (with cascade to images).

**Response:**
```json
{
  "deleted": 10
}
```

**Safety:** Requires UI confirmation (see Data Management section below)

### Settings

#### `GET /api/settings`
Get application settings.

**Response:**
```json
{
  "port": 3818
}
```

#### `PUT /api/settings`
Update settings (requires app restart).

**Request:**
```json
{
  "port": 3819
}
```

---

## Frontend Implementation

### Page Structure

All pages served as inline HTML from `handlers.go`:

1. **Dashboard** (`/`) - `IndexHandler()`
2. **Add Tote** (`/add`) - `AddToteHandler()`
3. **Edit Tote** (`/edit?id={id}`) - `EditToteHandler()`
4. **Tote Detail** (`/tote/{id}`) - `ToteDetailHandler()`
5. **QR Scanner** (`/scan`) - `ScanHandler()`
6. **Settings** (`/settings`) - `SettingsPageHandler()`
7. **Print Label** (`/print-label/{id}`) - `PrintLabelHandler()`

### UI/UX Features

#### Image Hover Gallery
**Feature:** Hover over tote card image to preview all images in a modal gallery.

**Implementation:**
```javascript
// On dashboard, hovering over card image triggers gallery
<div class="tote-image" onmouseenter="showImageGallery(event, tote)" onmouseleave="hideImageGallery()">
    <img src="${firstImage}" alt="${tote.name}">
</div>

function showImageGallery(event, tote) {
    // Creates modal with all tote images as thumbnails
    // Positioned near mouse cursor
    // Scaled to fit viewport
}

function hideImageGallery() {
    // Removes gallery when mouse leaves
}
```

**Gallery Features:**
- Displays all tote images as thumbnails in a grid
- Click thumbnail to view full-size image in expanded modal
- Close button (×) or click outside to dismiss
- Smooth fade-in/fade-out animations
- Responsive layout (max 3 columns)

#### Full Card Navigation
**Feature:** Entire tote card is clickable and navigates to edit page.

**Implementation:**
```javascript
<div class="tote-card" onclick="window.location.href='/edit?id=${tote.id}'" style="cursor: pointer;">
    // Card contents
</div>
```

**Behavior:**
- Card has `cursor: pointer` to indicate clickability
- Hover effect on entire card
- Image gallery opens on image hover (doesn't interfere with card click)
- Click anywhere on card navigates to edit page

### JavaScript Files

#### `app.js` - Dashboard Logic

```javascript
let allTotes = [];

function loadTotes() {
    fetch('/api/totes')
        .then(response => response.json())
        .then(totes => {
            allTotes = totes || [];
            updateStats();
            displayTotes(allTotes);
        });
}

function displayTotes(totes) {
    // Renders tote cards
    // Shows first image if available
    // Truncates item list to 3 lines
}

function setupSearch() {
    // Real-time search across name, description, items, qr_code
}

// Export all totes as JSON
function exportData() {
    window.location.href = '/api/export';
    // Triggers download of kontainer-export-YYYY-MM-DD.json
}

// Import totes from JSON file
async function importData(event) {
    const file = event.target.files[0];
    if (!file) return;

    // Confirmation dialog
    if (!confirm(`Import data from ${file.name}?\n\nThis will ADD the totes from the file to your existing inventory.`)) {
        event.target.value = '';
        return;
    }

    try {
        const text = await file.text();
        const response = await fetch('/api/import', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: text
        });

        const result = await response.json();
        alert(`Successfully imported ${result.imported} tote(s)!`);
        window.location.reload();  // Refresh to show new totes
    } catch (error) {
        console.error('Error importing data:', error);
        alert('Error importing data. Please check the file format.');
    } finally {
        event.target.value = '';  // Reset file input
    }
}
```

**Dashboard UI Elements:**
```html
<!-- Export button -->
<button class="btn btn-secondary" onclick="exportData()">
    📥 Export
</button>

<!-- Import button -->
<button class="btn btn-secondary" onclick="document.getElementById('import-file').click()">
    📤 Import
</button>

<!-- Hidden file input -->
<input type="file" id="import-file" accept=".json" style="display: none;" onchange="importData(event)">
```

#### `form.js` - Add/Edit Form Logic

**Critical Implementation Details:**

```javascript
let uploadedImagePaths = [];  // Stores paths of uploaded images

// Multi-file handling
async function setupImagePreview() {
    const imageInput = document.getElementById('image');
    imageInput.addEventListener('change', async function(e) {
        const files = e.target.files;
        uploadedImagePaths = [];
        
        // Upload all files immediately
        for (let i = 0; i < files.length; i++) {
            await uploadImage(files[i]);
        }
    });
}

async function uploadImage(file) {
    const formData = new FormData();
    formData.append('image', file);
    
    const response = await fetch('/api/upload-image', {
        method: 'POST',
        body: formData
    });
    
    const data = await response.json();
    uploadedImagePaths.push(data.path);
}

// CREATE MODE
async function handleSubmit(e) {
    if (!isEditMode) {
        const toteData = {
            name,
            description,
            items,
            image_paths: uploadedImagePaths  // All uploaded images
        };
        
        await fetch('/api/tote', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(toteData)
        });
    }
}

// EDIT MODE - ADDITIVE BEHAVIOR
async function handleSubmit(e) {
    if (isEditMode) {
        // Update tote details first
        await fetch(`/api/tote/${toteId}`, {
            method: 'PUT',
            body: JSON.stringify({ name, description, items })
        });
        
        // ADD new images (does not replace existing)
        for (const imagePath of uploadedImagePaths) {
            await fetch(`/api/tote/${toteId}/add-image`, {
                method: 'POST',
                body: JSON.stringify({ image_path: imagePath })
            });
        }
    }
}
```

**Key Points:**
- Images upload **immediately** on file selection
- Preview shown while uploading
- Edit mode shows existing images with "✓ Will be kept" indicator
- New images **added** via `/add-image` endpoint

#### `detail.js` - Tote Detail Page

```javascript
function displayToteDetail(tote) {
    // Image gallery with delete buttons
    let imagesHtml = '<div class="images-gallery">';
    tote.images.forEach(img => {
        imagesHtml += `
            <div class="image-item">
                <img src="${img.image_path}" class="detail-image">
                <button onclick="deleteImage(${img.id})" class="btn btn-danger">
                    🗑️ Delete
                </button>
            </div>
        `;
    });
    
    // QR code generation
    new QRCode(document.getElementById('qrcode'), {
        text: tote.qr_code,
        width: 150,
        height: 150,
        correctLevel: QRCode.CorrectLevel.H
    });
}

function deleteImage(imageId) {
    if (!confirm('Delete this image?')) return;
    
    fetch(`/api/tote-image/${imageId}`, { method: 'DELETE' })
        .then(() => window.location.reload());
}
```

#### `settings-page.js` - Settings Configuration

**Critical Implementation:**

```javascript
// Load current settings on page load
async function loadSettings() {
    const response = await fetch('/api/settings');
    const settings = await response.json();
    
    // Populate form
    document.getElementById('port').value = settings.port || 3818;
    document.getElementById('database_path').value = settings.database_path || 'totetrax.db';
    document.getElementById('theme').value = settings.theme || 'light';
    
    // Apply theme from server settings
    applyTheme(settings.theme || 'light');
}

// Save settings to server
async function saveSettings() {
    const settings = {
        port: parseInt(document.getElementById('port').value),
        database_path: document.getElementById('database_path').value.trim(),
        theme: document.getElementById('theme').value
    };
    
    // Validate
    if (settings.port < 1024 || settings.port > 65535) {
        alert('Port must be between 1024 and 65535');
        return;
    }
    
    if (!settings.database_path) {
        alert('Database path cannot be empty');
        return;
    }
    
    // Save to server
    await fetch('/api/settings', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(settings)
    });
    
    // Save theme to localStorage for instant UI updates
    const localSettings = {theme: settings.theme};
    localStorage.setItem('totetrax_settings', JSON.stringify(localSettings));
    
    // Apply theme immediately (no restart)
    applyTheme(settings.theme);
    
    alert('Settings saved!\n\n⚠️ Restart required for port/database changes.');
}

// Apply theme without page reload
function applyTheme(theme) {
    if (theme === 'dark') {
        document.documentElement.classList.add('dark-mode');
    } else {
        document.documentElement.classList.remove('dark-mode');
    }
}

// Reset to defaults
async function resetSettings() {
    if (!confirm('Reset to defaults?')) return;
    
    const defaults = {
        port: 3818,
        database_path: 'totetrax.db',
        theme: 'light'
    };
    
    await fetch('/api/settings', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(defaults)
    });
    
    localStorage.removeItem('totetrax_settings');
    window.location.reload();
}

// Import data from JSON file (settings page version)
async function importDataFromSettings(event) {
    const file = event.target.files[0];
    if (!file) return;

    if (!confirm(`Import data from ${file.name}?\n\nThis will ADD the totes from the file to your existing inventory.`)) {
        event.target.value = '';
        return;
    }

    try {
        const text = await file.text();
        const response = await fetch('/api/import', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: text
        });

        if (!response.ok) throw new Error('Import failed');

        const result = await response.json();
        alert(`Successfully imported ${result.imported} tote(s)!`);
        window.location.href = '/';  // Go to dashboard
    } catch (error) {
        console.error('Error importing data:', error);
        alert('Error importing data. Please check the file format.');
    } finally {
        event.target.value = '';
    }
}

// Delete all totes with triple confirmation
async function deleteAllData() {
    if (!confirm('⚠️ WARNING: Delete ALL totes?\n\nThis will permanently delete all totes and their images from the database.\n\nThis action CANNOT be undone!')) {
        return;
    }

    if (!confirm('Are you ABSOLUTELY sure?\n\nType YES in the next prompt to confirm.')) {
        return;
    }

    const confirmation = prompt('Type YES to delete all data:');
    if (confirmation !== 'YES') {
        alert('Deletion cancelled.');
        return;
    }

    try {
        const response = await fetch('/api/totes/delete-all', {
            method: 'DELETE'
        });

        if (!response.ok) throw new Error('Failed to delete data');

        const result = await response.json();
        alert(`Successfully deleted ${result.deleted} tote(s).`);
        window.location.href = '/';
    } catch (error) {
        console.error('Error deleting data:', error);
        alert('Error deleting data.');
    }
}
    });
    
    localStorage.removeItem('totetrax_settings');
    window.location.reload();
}
```

**Key Behaviors:**
- Theme changes apply immediately via `applyTheme()`
- Theme saved to both server (`settings.json`) and client (`localStorage`)
- Port and database path changes require restart (user notified)
- Form validation prevents invalid port numbers
- Reset button clears all settings to defaults


### CSS Highlights

**Dark Mode Support:**
```css
:root {
    --bg-color: #f5f5f5;
    --text-color: #333;
    --card-bg: #fff;
}

.dark-mode {
    --bg-color: #1a1a1a;
    --text-color: #e0e0e0;
    --card-bg: #2d2d2d;
}
```

**Image Gallery Layout:**
```css
.images-gallery {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 15px;
}

.detail-image {
    width: 100%;
    height: 200px;
    object-fit: cover;
}
```

---

## Image Handling

### NEW: Database-Embedded Images (v1.4.0+)

Images are now stored **directly in the SQLite database as BLOBs**, eliminating the need for separate file storage.

### Upload Flow

1. **User selects files** (multiple via Ctrl+Click)
2. **JavaScript converts to base64** using `FileReader.readAsDataURL()`
3. **Stores in memory** as data URI (e.g., `data:image/jpeg;base64,/9j/4AAQ...`)
4. **On form submit:**
   - Create mode: Sends `image_paths` array (base64 data URIs) + `image_types` array
   - Edit mode: Calls `POST /api/tote/{id}/add-image` with base64 data
5. **Server decodes base64** to binary BLOB
6. **Inserts directly into `tote_images` table**
7. **On retrieval:** Server converts BLOB back to base64 data URI
8. **Frontend displays** using `<img src="data:image/jpeg;base64,...">`

### Base64 Encoding/Decoding

**Go Backend Implementation:**

```go
// Store: Decode base64 data URI → binary BLOB
func decodeBase64Image(dataURI string) ([]byte, string, error) {
    // Extract MIME type and base64 data
    // Input: "data:image/jpeg;base64,/9j/4AAQ..."
    // Output: []byte{binary data}, "image/jpeg"
    
    parts := strings.Split(dataURI, ",")
    header := parts[0] // "data:image/jpeg;base64"
    base64Data := parts[1]
    
    // Custom base64 decoder (no external dependencies)
    decoded := base64Decode(base64Data)
    mimeType := extractMimeType(header)
    
    return decoded, mimeType, nil
}

// Retrieve: Binary BLOB → base64 data URI
func loadImagesForTote(toteID int) ([]models.ToteImage, error) {
    // Read BLOB from database
    rows.Scan(&img.ID, &img.ToteID, &imageData, &img.ImageType, ...)
    
    // Convert to data URI
    img.ImageData = "data:" + img.ImageType + ";base64," + base64Encode(imageData)
    
    return images, nil
}
```

**JavaScript Frontend:**

```javascript
// Upload: File → base64 data URI
const reader = new FileReader();
reader.onload = function(e) {
    const base64Data = e.target.result; // "data:image/jpeg;base64,..."
    uploadedImages.push({
        data: base64Data,
        type: file.type
    });
};
reader.readAsDataURL(file);

// Submit to API
fetch('/api/tote', {
    method: 'POST',
    body: JSON.stringify({
        name: 'Kitchen Supplies',
        image_paths: uploadedImages.map(img => img.data),
        image_types: uploadedImages.map(img => img.type)
    })
});

// Display: Data URI → img src
tote.images.forEach(img => {
    html += `<img src="${img.image_data}">`;
    // img.image_data = "data:image/jpeg;base64,/9j/4AAQ..."
});
```

### Database Storage

**Structure:**

| Column | Type | Example |
|--------|------|---------|
| `id` | INTEGER | 1 |
| `tote_id` | INTEGER | 42 |
| `image_data` | BLOB | `\xFF\xD8\xFF\xE0...` (binary JPEG) |
| `image_type` | TEXT | `image/jpeg` |
| `display_order` | INTEGER | 0 |

**Binary Data:**
- Stored as raw bytes in BLOB column
- Supported formats: JPEG, PNG, GIF, WebP, etc.
- No file size limits (SQLite supports BLOBs up to 1GB+)

### Benefits of Database-Embedded Images

✅ **No Orphaned Files**
- CASCADE DELETE removes images when tote is deleted
- No manual cleanup needed
- Database integrity guaranteed

✅ **Fully Portable**
- Single `.db` file contains everything
- Copy database = copy all data + all images
- Example: `totetrax.db` is 45MB with 100 images

✅ **Simplified Backups**
- Backup 1 file instead of DB + image folder
- Network database already includes images
- Export JSON contains embedded images (base64)

✅ **Network Storage Ready**
- Set database to NAS: `\\NAS\share\totetrax.db`
- All images automatically stored in same database
- Multiple computers access same data + images

✅ **Import/Export with Images**
- Export creates JSON with base64-encoded images
- Import recreates exact totes with all images
- No separate image files to transfer

**Example Database Organization:**

```
\\NAS\Inventory\
└── totetrax.db         (contains everything: data + images)
```

Compared to old file-based approach:
```
\\NAS\Inventory\
├── totetrax.db         (data only)
└── tote_images\        (separate folder, can get orphaned)
    ├── 1739497234567890123.jpg
    ├── 1739497234567890124.png
    └── orphaned_file.jpg  (tote deleted but file remains)
```
- Old images in `web/static/images/uploads/` still accessible
- No migration required

✅ **Transparent to Frontend**
- Frontend code unchanged
- Always uses `/images/` URLs
- Server handles the complexity

### Image Deletion

**Delete Behavior:**
1. User clicks 🗑️ button on detail page
2. JavaScript calls `DELETE /api/tote-image/{id}`
3. Backend deletes database record
4. **Physical file NOT deleted** (intentional)
5. Page reloads to show updated gallery

**Why not delete file?**
- Files may be referenced by multiple totes (if user uploads same image)
- Prevents accidental data loss
- Disk space is cheap
- Consider adding cleanup script for orphaned files if needed

### Backward Compatibility

**Legacy `image_path` field:**
- Populated with first image from `image_paths` array
- Ensures old clients still see an image
- Not used by current UI (uses `images` array)

---

## Settings & Configuration

### Settings Model

**File:** `internal/models/settings.go`

```go
type Settings struct {
    Port         int    `json:"port"`          // Server port (default: 3818)
    Theme        string `json:"theme"`         // "light" or "dark"
    DatabasePath string `json:"database_path"` // Path to SQLite database file
}

func DefaultSettings() *Settings {
    return &Settings{
        Port:         3818,
        Theme:        "light",
        DatabasePath: "totetrax.db",
    }
}
```

### settings.json File

**Location:** Created in current working directory on first run

**Format:**
```json
{
  "port": 3818,
  "theme": "light",
  "database_path": "totetrax.db"
}
```

**Network Database Example (Windows UNC path):**
```json
{
  "port": 3818,
  "theme": "dark",
  "database_path": "\\\\192.168.1.100\\storage\\totetrax.db"
}
```

**Network Database Example (Linux mount):**
```json
{
  "port": 3818,
  "theme": "dark",
  "database_path": "/mnt/nas/totetrax.db"
}
```

### Settings Page UI

**URL:** http://localhost:3818/settings

**Features:**

1. **Server Port Configuration**
   - Range: 1024-65535
   - Validated on input
   - **Requires app restart** to take effect
   - Default: 3818

2. **Database Path Configuration**
   - Supports local paths: `totetrax.db`, `C:\data\totetrax.db`
   - Supports network paths: `\\NAS\share\totetrax.db`
   - Supports Linux paths: `/mnt/nas/totetrax.db`
   - **Requires app restart** to take effect
   - Enables database on NAS or network drive

3. **Theme Toggle**
   - Options: Light or Dark
   - **Applies immediately** (no restart needed)
   - Saved to both `settings.json` and `localStorage`
   - Persists across sessions

**JavaScript Implementation:** `web/static/js/settings-page.js`

```javascript
async function saveSettings() {
    const settings = {
        port: parseInt(document.getElementById('port').value),
        database_path: document.getElementById('database_path').value.trim(),
        theme: document.getElementById('theme').value
    };
    
    // Validate and save
    await fetch('/api/settings', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(settings)
    });
    
    // Apply theme immediately
    applyTheme(settings.theme);
    
    // Save to localStorage for instant UI update
    localStorage.setItem('totetrax_settings', JSON.stringify({theme: settings.theme}));
}
```

### How Settings Are Loaded

**On Application Startup:**

```go
// cmd/kontainer/main.go
func main() {
    // 1. Load settings from settings.json
    settingsService := service.NewSettingsService()
    settings, err := settingsService.LoadSettings()
    
    // 2. Use configured database path
    dbPath := settings.DatabasePath
    if dbPath == "" {
        dbPath = "totetrax.db" // Fallback
    }
    db, err := database.InitDB(dbPath)
    
    // 3. Start server on configured port
    port := fmt.Sprintf(":%d", settings.Port)
    http.ListenAndServe(port, router)
}
```

### Settings Usage Guide

#### Changing Server Port

1. Open Settings page: http://localhost:3818/settings
2. Enter new port (e.g., 8080)
3. Click "Save Settings"
4. **Restart the application**
5. Access at new port: http://localhost:8080

#### Using Network Database (NAS)

**Scenario:** Store database on network drive for multi-device access

**Steps:**

1. Ensure NAS is accessible and mounted
2. Open Settings page
3. Enter database path:
   - **Windows:** `\\192.168.1.100\storage\totetrax.db`
   - **Linux:** `/mnt/nas/totetrax.db`
4. Click "Save Settings"
5. **Confirm the migration** when prompted
6. Wait for migration to complete (database automatically copied to new location)
7. **Restart the application**
8. Database now stored on NAS

**What Happens During Migration:**
- Current database file is copied to new location
- Associated SQLite files (-journal, -wal, -shm) are also copied
- Destination directory is created if it doesn't exist
- Original database remains as backup
- Settings updated to point to new location
- Migration fails safely if destination file already exists

**Benefits:**
- Access from multiple computers
- Centralized backup
- Shared inventory across devices
- Network redundancy (if NAS has RAID)
- Automatic database migration (no manual file copying)

**Requirements:**
- Network path must be accessible
- Write permissions required
- Stable network connection
- SQLite supports network files (but performance may vary)

**Safety Features:**
- User confirmation required before migration
- Checks for existing files to prevent overwriting
- Original database preserved as backup
- Clear error messages if migration fails

#### Switching Theme

1. Open Settings page
2. Select "Dark" or "Light" from theme dropdown
3. Click "Save Settings"
4. Theme applies **immediately** (no restart)
5. Preference saved for next session

**Technical Details:**
- Theme saved to `settings.json` (server-side)
- Theme also saved to `localStorage` (client-side)
- On page load, theme applied before content renders (no flash)
- CSS variables used for easy theme switching

### Settings API Endpoints

#### GET /api/settings
Get current settings.

**Response:**
```json
{
  "port": 3818,
  "theme": "dark",
  "database_path": "\\\\NAS\\share\\totetrax.db"
}
```

#### PUT /api/settings
Update settings.

**Request:**
```json
{
  "port": 8080,
  "theme": "light",
  "database_path": "/mnt/nas/totetrax.db"
}
```

**Response:** Same as request (updated settings)

**Validation:**
- Port must be 1024-65535
- Database path cannot be empty
- Theme must be "light" or "dark"

### Data Management & Backup

ToteTrax provides comprehensive data management features accessible from both the **Dashboard** and **Settings** page.

#### Dashboard Quick Access

Located in the header alongside "Scan QR", "Add Tote", and "Settings":

```html
<button class="btn btn-secondary" onclick="exportData()">
    📥 Export
</button>
<button class="btn btn-secondary" onclick="document.getElementById('import-file').click()">
    📤 Import
</button>
```

**Features:**
- **Export**: One-click backup of all totes and images to JSON
- **Import**: Upload JSON backup file to restore/add totes

#### Settings Page Data Management

Located at the bottom of the Settings page in a dedicated "Data Management" section:

```html
<div class="settings-section">
    <h3>Data Management</h3>
    
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
```

**Features:**
- **Export All Data**: Same as dashboard export
- **Import Data**: Same as dashboard import
- **Delete All Totes**: Nuclear option with triple confirmation

#### Export Workflow

1. User clicks "Export" button
2. JavaScript executes: `window.location.href = '/api/export'`
3. Server responds with JSON file download
4. Browser downloads: `kontainer-export-2026-02-14.json`

**Export File Contents:**
- All totes with complete metadata
- Full images array with paths
- QR codes
- Timestamps

**Use Cases:**
- Daily/weekly backups
- Before major changes
- Database migration
- Sharing data between devices

#### Import Workflow

1. User clicks "Import" button
2. Hidden file input opens: `<input type="file" accept=".json">`
3. User selects JSON file
4. Confirmation dialog: "Import data from backup.json? This will ADD..."
5. JavaScript reads file and POSTs to `/api/import`
6. Server creates new totes with all images
7. Success message: "Successfully imported 15 tote(s)!"
8. Page reloads to show imported totes

**Important Notes:**
- **Additive operation**: Does NOT replace existing totes
- **Creates new totes**: Imported totes get new IDs and QR codes
- **Image files must exist**: Import assumes image files exist at specified paths
- **Cross-system imports**: Copy `/web/static/images/uploads/` folder first

**Image Path Handling:**
```json
{
  "images": [
    {"image_path": "/static/images/uploads/1771058813249973700.jpg"},
    {"image_path": "/static/images/uploads/1771058813255059400.png"}
  ]
}
```
- Paths are relative to web server root
- Import recreates all images in `tote_images` table
- Physical files must exist in `web/static/images/uploads/`

#### Delete All Workflow

**⚠️ Available ONLY in Settings page** (intentionally not on dashboard)

1. User clicks "Delete All Totes" button
2. **First confirmation**: Warning dialog with "CANNOT be undone" message
3. **Second confirmation**: "Are you ABSOLUTELY sure?"
4. **Third confirmation**: Prompt requiring user to type "YES"
5. If confirmed: DELETE request to `/api/totes/delete-all`
6. Server deletes all totes (CASCADE deletes images)
7. Success message: "Successfully deleted 42 tote(s)."
8. Redirect to dashboard (empty state)

**Safety Features:**
```javascript
if (!confirm('⚠️ WARNING: Delete ALL totes?...')) return;
if (!confirm('Are you ABSOLUTELY sure?...')) return;
const confirmation = prompt('Type YES to delete all data:');
if (confirmation !== 'YES') {
    alert('Deletion cancelled.');
    return;
}
```

**What Gets Deleted:**
- All totes from `totes` table
- All images from `tote_images` table (CASCADE)
- Database records only (physical image files remain)

**What Does NOT Get Deleted:**
- Physical image files in `web/static/images/uploads/`
- Settings in `settings.json`
- Database file itself (just emptied)

#### Backup Best Practices

1. **Regular Exports**: Weekly or before major changes
2. **Store exports externally**: Copy JSON to cloud/USB
3. **Copy image folder**: Backup `web/static/images/uploads/` separately
4. **Network database**: Use NAS for automatic redundancy
5. **Test imports**: Periodically verify backups work

#### Cross-System Migration

To move ToteTrax to another computer:

1. **Export data** on old system
2. **Copy image folder**: `web/static/images/uploads/`
3. **Install ToteTrax** on new system
4. **Paste image folder** to new installation
5. **Import data** on new system
6. Verify all images display correctly

**Alternative (Network Database):**
1. Configure Settings → Database Path: `\\NAS\share\totetrax.db`
2. All systems share same database
3. All systems access same image folder on network
4. No export/import needed

### Firewall Configuration

**Windows:**
```powershell
New-NetFirewallRule -DisplayName "ToteTrax Web Server" -Direction Inbound -LocalPort 3818 -Protocol TCP -Action Allow
```

**Linux (ufw):**
```bash
sudo ufw allow 3818/tcp
```

---

## Development Guide

### Prerequisites

- Go 1.25 or later
- No other dependencies (pure Go SQLite driver)

### Building

```bash
# Download dependencies
go mod download

# Build for current OS
go build -o totetrax cmd/kontainer/main.go

# Build for Windows (from any OS)
GOOS=windows GOARCH=amd64 go build -o totetrax.exe cmd/kontainer/main.go

# Build for Linux (from any OS)
GOOS=linux GOARCH=amd64 go build -o totetrax cmd/kontainer/main.go
```

### Running in Development

```bash
go run cmd/kontainer/main.go
```

**Or with live reload (install `air` first):**
```bash
go install github.com/cosmtrek/air@latest
air
```

### Database Reset

```bash
# Stop server
# Delete database
rm totetrax.db

# Restart server - fresh database created
go run cmd/kontainer/main.go
```

### Adding New Endpoints

1. **Define route** in `internal/api/router.go`:
   ```go
   mux.HandleFunc("/api/new-endpoint", handler.NewEndpointHandler)
   ```

2. **Create handler** in `internal/api/handlers.go`:
   ```go
   func (h *Handler) NewEndpointHandler(w http.ResponseWriter, r *http.Request) {
       // Implementation
   }
   ```

3. **Add service method** if needed in `internal/service/tote_service.go`

4. **Update model** if needed in `internal/models/tote.go`

### Testing APIs

**Using curl:**
```bash
# Get all totes
curl http://localhost:3818/api/totes

# Create tote
curl -X POST http://localhost:3818/api/tote \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","items":"Item 1\nItem 2","image_paths":["/static/images/uploads/test.jpg"]}'

# Get by ID
curl http://localhost:3818/api/tote/1

# Get by QR code
curl http://localhost:3818/api/tote/qr/TOTE-00001

# Delete image
curl -X DELETE http://localhost:3818/api/tote-image/5
```

---

## Deployment

### Single Binary Deployment

**Advantages:**
- No runtime dependencies
- No installation required
- Cross-platform (Windows, Linux, macOS)

**Files needed:**
```
totetrax.exe          (or totetrax on Linux)
web/                  (entire directory with static assets)
```

### Directory Structure on Server

```
/opt/totetrax/        (or C:\Apps\ToteTrax on Windows)
├── totetrax          (executable)
├── web/              (copied from source)
│   └── static/
├── totetrax.db       (created on first run)
└── settings.json     (created on first run)
```

### Running as Service

**Linux (systemd):**

Create `/etc/systemd/system/totetrax.service`:
```ini
[Unit]
Description=ToteTrax Storage Inventory
After=network.target

[Service]
Type=simple
User=totetrax
WorkingDirectory=/opt/totetrax
ExecStart=/opt/totetrax/totetrax
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable:
```bash
sudo systemctl enable totetrax
sudo systemctl start totetrax
```

**Windows (NSSM):**
```powershell
nssm install ToteTrax "C:\Apps\ToteTrax\totetrax.exe"
nssm set ToteTrax AppDirectory "C:\Apps\ToteTrax"
nssm start ToteTrax
```

### Reverse Proxy (Optional)

**Nginx:**
```nginx
server {
    listen 80;
    server_name totetrax.example.com;

    location / {
        proxy_pass http://localhost:3818;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

**Caddy:**
```
totetrax.example.com {
    reverse_proxy localhost:3818
}
```

### Backup Strategy

**What to backup:**
1. `totetrax.db` - All data
2. `web/static/images/uploads/` - User-uploaded images
3. `settings.json` - Configuration

**Backup script (Linux):**
```bash
#!/bin/bash
BACKUP_DIR="/backups/totetrax/$(date +%Y%m%d)"
mkdir -p "$BACKUP_DIR"
cp /opt/totetrax/totetrax.db "$BACKUP_DIR/"
cp -r /opt/totetrax/web/static/images/uploads "$BACKUP_DIR/"
cp /opt/totetrax/settings.json "$BACKUP_DIR/"
```

**Restore:**
```bash
# Stop service
sudo systemctl stop totetrax

# Restore files
cp backup/totetrax.db /opt/totetrax/
cp -r backup/uploads/* /opt/totetrax/web/static/images/uploads/

# Start service
sudo systemctl start totetrax
```

---

## Common Tasks

### Change Port

1. Edit `settings.json`:
   ```json
   {
     "port": 8080
   }
   ```

2. Restart application

3. Update firewall if needed

### Add Default Totes (Seed Data)

No built-in seeding. Use import endpoint:

```bash
curl -X POST http://localhost:3818/api/import \
  -H "Content-Type: application/json" \
  -d '[
    {
      "name": "Example Tote",
      "description": "Sample data",
      "items": "Item 1\nItem 2",
      "image_paths": []
    }
  ]'
```

### Cleanup Orphaned Images

Create maintenance script:

```go
// Find images not referenced in database
rows := db.Query("SELECT image_path FROM tote_images")
var usedPaths []string
// ... collect paths

// List all files in uploads dir
files, _ := ioutil.ReadDir("web/static/images/uploads")

// Delete files not in usedPaths
for _, file := range files {
    path := "/static/images/uploads/" + file.Name()
    if !contains(usedPaths, path) {
        os.Remove(filepath.Join("web/static/images/uploads", file.Name()))
    }
}
```

---

## Troubleshooting

### Server won't start

**Check port availability:**
```bash
# Windows
netstat -ano | findstr :3818

# Linux
lsof -i :3818
```

**Change port in settings.json if needed**

### Images not displaying

**Check paths:**
- Database stores: `/static/images/uploads/12345.jpg`
- File must exist at: `web/static/images/uploads/12345.jpg`
- Browser requests: `http://localhost:3818/static/images/uploads/12345.jpg`

**Check permissions:**
```bash
# Linux - ensure web/static/images/uploads is writable
chmod 755 web/static/images/uploads
```

### QR codes not generating

**Check JavaScript console for errors**

**Verify qrcode.min.js loaded:**
```html
<script src="/static/js/qrcode.min.js"></script>
```

### Database locked

**Symptoms:** `database is locked` error

**Causes:**
- Multiple instances running
- Unclosed database connection
- Crash left lock file

**Solutions:**
```bash
# Check for .db-journal file
ls -la totetrax.db*

# Remove if present and no instances running
rm totetrax.db-journal

# Ensure only one instance running
ps aux | grep totetrax
```

---

## Future Enhancements

### Planned Features

1. ~~**Dark Mode Toggle UI**~~ ✅ **IMPLEMENTED**
   - Settings page with instant theme switching
   - Saved to localStorage and settings.json
   - No restart required

2. **Image Reordering**
   - Drag-and-drop in gallery
   - Update `display_order` field

3. **Bulk Operations**
   - Select multiple totes
   - Bulk delete, export

4. **Advanced Search**
   - Filter by date range
   - Filter by image count
   - Sort by various fields

5. ~~**Settings Page**~~ ✅ **IMPLEMENTED**
   - Port configuration (1024-65535)
   - Database path (local or network/NAS)
   - Theme toggle with instant preview

6. **Location Tracking**
   - Add `storage_location` field (room, shelf, etc.)
   - Group totes by location

6. **Location Tracking**
   - Add `storage_location` field (room, shelf, etc.)
   - Group totes by location

7. **Barcode Support**
   - Alternative to QR codes
   - Code 128 or Code 39

7. **Mobile App**
   - React Native or Flutter
   - Native camera/NFC for better scanning

8. **Statistics Dashboard**
   - Total containers
   - Total items
   - Storage utilization charts

9. **Print Multiple Labels**
   - Select multiple totes
   - Generate PDF with all labels

10. **CSV Export**
    - Alternative to JSON
    - Better for spreadsheet users

### Technical Debt

1. **Error Handling**
   - Add structured logging (logrus or zap)
   - Better error messages to user

2. **Validation**
   - Input sanitization
   - Image file type validation
   - Max file size enforcement in UI

3. **Testing**
   - Unit tests for services
   - API integration tests
   - Frontend E2E tests (Playwright)

4. **Security**
   - Add authentication (optional)
   - Rate limiting on uploads
   - CSRF protection

5. **Performance**
   - Image thumbnails (ImageMagick/libvips)
   - Lazy loading in gallery
   - Pagination for large inventories

---

## Code Patterns & Conventions

### Error Handling

```go
// Services return (value, error)
func (s *ToteService) GetByID(id int) (*models.Tote, error) {
    // ...
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("tote not found")
    }
    return &tote, nil
}

// Handlers convert to HTTP errors
func (h *Handler) ToteHandler(w http.ResponseWriter, r *http.Request) {
    tote, err := h.toteService.GetByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(tote)
}
```

### JSON Responses

```go
// Success responses
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(data)

// Created responses (201)
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(data)

// No content responses (204)
w.WriteHeader(http.StatusNoContent)

// Error responses
http.Error(w, "error message", http.StatusBadRequest)
```

### SQL Queries

```go
// Use prepared statements
query := "SELECT id, name FROM totes WHERE id = ?"
row := db.QueryRow(query, id)

// Use placeholders, not string concatenation
// ❌ BAD: query := "SELECT * FROM totes WHERE name = '" + name + "'"
// ✅ GOOD: query := "SELECT * FROM totes WHERE name = ?"
```

### Naming Conventions

- **Handlers:** `{Action}Handler` (e.g., `TotesHandler`, `AddToteHandler`)
- **Services:** `{Entity}Service` (e.g., `ToteService`, `SettingsService`)
- **Models:** Singular nouns (e.g., `Tote`, `ToteImage`, not `Totes`)
- **API Endpoints:** 
  - Collections: plural (`/api/totes`)
  - Single resource: singular (`/api/tote`, `/api/tote/1`)
  - Actions: verb-based (`/api/tote/1/add-image`)

---

## Architecture Decisions

### Why SQLite?

**Pros:**
- Zero configuration
- Single file database
- Perfect for single-user/small team use
- No separate server process
- Built-in with Go (pure Go driver)
- ACID compliant

**Cons:**
- Not suitable for high concurrency (not needed for this use case)
- Single writer at a time (fine for web app)

**Alternative considered:** PostgreSQL (overkill for this use case)

### Why Embedded HTML in Go?

**Pros:**
- Single binary deployment
- No template file dependencies
- Fast startup (no file I/O)

**Cons:**
- Harder to edit (Go strings)
- No syntax highlighting for HTML

**Alternative:** Could use `embed` package with separate HTML files:
```go
//go:embed templates/*.html
var templates embed.FS
```

### Why Multiple Image Upload Immediately?

**Pros:**
- Better UX (see uploads progress)
- Prevents form loss if user navigates away
- Can validate images before form submit

**Cons:**
- Images uploaded even if user cancels
- No transaction rollback if tote creation fails

**Mitigation:** Could add cleanup job for orphaned images

### Why Additive Image Upload in Edit Mode?

**Pros:**
- Safer (never lose existing images)
- More intuitive (add, don't replace)
- Individual delete gives fine control

**Cons:**
- Can't replace all images in one action
- User must delete individually first

**Design Choice:** Safety over convenience. Losing images is worse than extra clicks.

---

## Performance Considerations

### Current Performance

**Database:**
- Small dataset (<10,000 totes): milliseconds
- Indexes on `name`, `qr_code`, `tote_id`
- No N+1 query issues (images loaded per tote, not per image)

**File I/O:**
- Images served directly by Go's `http.FileServer`
- No processing on read (just file serving)
- Upload capped at 10MB per file

**Frontend:**
- Vanilla JavaScript (no framework overhead)
- QR libraries ~100KB total
- CSS ~8KB

### Optimization Opportunities

1. **Image Thumbnails:**
   ```go
   // Generate on upload
   thumb := resize.Thumbnail(200, 200, img, resize.Lanczos3)
   jpeg.Encode(thumbFile, thumb, nil)
   ```

2. **Pagination:**
   ```sql
   SELECT * FROM totes LIMIT 50 OFFSET ?
   ```

3. **Caching:**
   - Add ETag headers for static files
   - Cache tote list in memory (invalidate on change)

4. **Lazy Loading:**
   ```javascript
   // Intersection Observer for images
   observer.observe(imgElement);
   ```

---

## Security Considerations

### Current Security

**Good:**
- No eval() or innerHTML with user data
- SQL prepared statements (no injection)
- File uploads to dedicated directory
- No authentication needed (local use)

**Needs Improvement:**
1. **File Type Validation:**
   ```go
   // Check magic bytes, not just extension
   contentType := http.DetectContentType(fileHeader)
   if !strings.HasPrefix(contentType, "image/") {
       return errors.New("not an image")
   }
   ```

2. **Path Traversal Prevention:**
   ```go
   // Sanitize filenames
   filename = filepath.Base(filename)
   ```

3. **Rate Limiting:**
   ```go
   // Use golang.org/x/time/rate
   limiter := rate.NewLimiter(rate.Every(time.Second), 10)
   ```

4. **HTTPS:**
   ```go
   // Add TLS support
   http.ListenAndServeTLS(":3818", "cert.pem", "key.pem", router)
   ```

### For Multi-User Deployment

**Add Authentication:**
```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check session/JWT
        if !isAuthenticated(r) {
            http.Error(w, "Unauthorized", 401)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

**Add User Context:**
```sql
ALTER TABLE totes ADD COLUMN user_id INTEGER;
```

---

## Git Repository

### Current Commits

```
675fad7 - Clarify additive image upload behavior in UI
75c8b2b - Add multiple images support per tote
8cd8b36 - Fix gitignore and add cmd/kontainer/main.go
59f5644 - Add quick start guide and main.go
8bb0a22 - Initial commit: ToteTrax storage container inventory management
```

### .gitignore

```
# Binaries
*.exe
*.dll
*.so
*.dylib
/totetrax
/totetrax.exe

# Database
*.db
*.db-journal

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Settings
settings.json

# Uploaded images
web/static/images/uploads/
```

---

## Summary

ToteTrax is a lightweight, self-contained storage inventory system built with:
- **Go backend** with SQLite (pure Go driver)
- **Vanilla JavaScript frontend**
- **Multiple images per container** with additive uploads
- **QR code system** for quick lookup
- **No dependencies** beyond Go standard library
- **Single binary deployment**

**Key architectural choices:**
1. SQLite for simplicity and portability
2. Embedded HTML for single-binary deployment
3. Additive image uploads for data safety
4. REST API for future extensibility
5. Pure Go SQLite driver for cross-platform builds

**Critical implementation details:**
- Images upload immediately on file selection
- Edit mode ADDS images (never replaces)
- Individual image deletion from gallery
- QR codes auto-generated sequentially
- Cascade delete removes images with tote
- Display order field for future drag-and-drop

**Perfect for:**
- Home organization
- Moving/packing
- Garage/workshop storage
- Warehouse (small scale)
- Anywhere you need to track "what's in which box"
- Multi-device access via network database

---

## Version History

### v1.8.0 (2026-02-17)
- ✅ **Automatic database migration** - Database automatically moved when path is changed in settings
- ✅ **Smart file handling** - Copies database and associated files (-journal, -wal, -shm)
- ✅ **Directory creation** - Automatically creates destination directories if needed
- ✅ **Safety checks** - Prevents overwriting existing database at new location
- ✅ **User confirmation** - Warns user before migrating database
- ✅ **Comprehensive tests** - Full test coverage for migration functionality
- ✅ **Error handling** - Clear error messages if migration fails

### v1.7.0 (2026-02-17)
- ✅ **Location field** - Added optional location tracking for containers
- ✅ **Database migration** - Automatic schema update adds location column to existing databases
- ✅ **UI integration** - Location field in add/edit forms and detail pages
- ✅ **Dashboard display** - Location shown on container cards with 📍 icon
- ✅ **Import/Export support** - Location field included in JSON export/import
- ✅ **Backward compatible** - NULL location values supported for existing containers

### v1.6.1 (2026-02-17)
- ✅ **UI update** - Changed "Scan QR" button to "Look up" with magnifying glass icon
- ✅ **Page title update** - Changed "Scan Tote QR Code" to "Look up Tote"
- ✅ **Removed camera scan** - Removed web-based camera scanning (reserved for mobile app)
- ✅ **Fixed image upload** - QR code image scanning now works properly with PNG, JPEG, GIF
- ✅ **Two lookup methods** - Image upload and manual entry (camera scanning via mobile app only)
- ✅ **API preserved** - All QR lookup API endpoints remain for mobile app integration

### v1.3.0 (2026-02-14)
- ✅ **Smart image storage location** - Images now saved next to database file
- ✅ **Portable backups** - Database folder contains both DB and images
- ✅ **Network storage ready** - NAS database = NAS images automatically
- ✅ **Organized structure** - Custom DB creates `tote_images/` subfolder
- ✅ **Image serving endpoint** - New `/images/` route serves from any location
- ✅ **Automatic path conversion** - Backend converts absolute paths to web URLs
- ✅ **Backward compatible** - Default DB still uses `web/static/images/uploads/`
- ✅ **Transparent to frontend** - No JavaScript changes required
- ✅ **Examples documented** - Clear scenarios for local/custom/network paths

### v1.2.0 (2026-02-14)
- ✅ **Import/Export UI** - Added user-friendly buttons to dashboard and settings
- ✅ **Dashboard quick access** - Export and Import buttons in header
- ✅ **Settings data management section** - Export, Import, and Delete All
- ✅ **Fixed image import** - Import now properly handles multiple images array
- ✅ **Backward compatibility** - Import supports old single image_path format
- ✅ **Enhanced safety** - Delete All requires triple confirmation with typed "YES"
- ✅ **Confirmation dialogs** - All destructive operations require confirmation
- ✅ **Auto-reload after import** - Page refreshes to show imported totes
- ✅ **Comprehensive documentation** - Complete data management and backup guide

### v1.1.0 (2026-02-14)
- ✅ Added comprehensive settings page
- ✅ Configurable server port (1024-65535)
- ✅ Configurable database path (supports NAS/network drives)
- ✅ Light/Dark theme toggle with instant application
- ✅ Settings UI with validation and reset functionality
- ✅ Updated main.go to use configurable database path
- ✅ Theme persistence in localStorage and settings.json

### v1.0.0 (2026-02-14)
- ✅ Initial release
- ✅ Multiple images per tote
- ✅ Additive image uploads
- ✅ QR code generation and scanning
- ✅ Import/Export API endpoints
- ✅ Complete REST API
- ✅ Web UI with dark mode CSS support

---

**End of Technical Documentation**

*Last Updated: 2026-02-17*  
*Version: 1.8.0*  
*Maintained at: D:\projects\totetrax\*
