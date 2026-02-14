# ToteTrax Project Information

**Project Name:** ToteTrax  
**Created:** 2026-02-14  
**Based On:** Filatrax (3D Filament Inventory Management)

## Project Overview
ToteTrax is a storage container inventory and tracking cross-platform application designed for Windows and Linux. It helps users organize and track items stored in boxes, totes, bins, and containers for home organization, moving, storage, and inventory management.

The application runs as a local service providing a web-based interface, making it accessible from any browser while keeping data local. It features a modular architecture and includes a REST API for integration with companion applications.

## Technology Stack
- **Language:** Go 1.25+ (single binary deployment, cross-platform)
- **Database:** SQLite (embedded, zero-configuration)
- **SQLite Driver:** modernc.org/sqlite (pure Go, no CGO dependencies)
- **Web Framework:** Go standard library net/http
- **Frontend:** HTML/CSS/JavaScript with QR code support
- **API:** REST API for external integrations
- **License:** Open source / Free components only

### Deployment Benefits
- Single executable binary (no runtime dependencies)
- Zero installation requirements for users
- Embedded SQLite (no separate database server)
- Cross-platform native compilation

## Project Structure
```
D:\projects\totetrax\
├── cmd\totetrax\              # Application entry point
│   └── main.go
├── internal\                  # Internal packages
│   ├── api\                   # HTTP handlers and routing
│   │   ├── handlers.go
│   │   └── router.go
│   ├── database\              # Database initialization
│   │   └── database.go
│   ├── models\                # Data models
│   │   ├── tote.go
│   │   └── settings.go
│   └── service\               # Business logic
│       ├── tote_service.go
│       └── settings_service.go
├── web\                       # Web interface assets
│   └── static\                # CSS, JS, images
│       ├── css\
│       │   └── style.css
│       ├── js\
│       │   ├── app.js
│       │   ├── form.js
│       │   ├── detail.js
│       │   ├── qrcode.min.js
│       │   └── html5-qrcode.min.js
│       └── images\
│           └── uploads\       # User-uploaded images
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
├── README.md                  # Project documentation
├── .gitignore                 # Git ignore rules
├── totetrax.exe              # Compiled executable (Windows)
└── totetrax-project-info.md  # This file
```

## Key Features

### Core Functionality
- **Container Tracking**: Track storage totes/boxes by name and description
- **Item Lists**: Maintain text-based lists of items in each container
- **Image Support**: Upload and attach images of container contents
- **QR Code System**: Auto-generated QR codes (TOTE-XXXXX format) for each container
- **QR Scanning**: Multi-method scanning (camera, image upload, manual entry)
- **Printable Labels**: Generate printer-friendly labels with QR codes
- **Search**: Quick search across tote names, descriptions, and item lists

### Technical Features
- **Web Interface**: Browser-based UI served by local service
- **REST API**: Full API access for external applications
- **Cross-Platform**: Windows and Linux support
- **Modular Architecture**: Easy to extend with new features
- **Import/Export**: JSON-based backup and restore
- **Image Storage**: Local file storage for container photos

## Data Model

Tote records include:
- Name (e.g., "Kitchen Supplies", "Holiday Decorations")
- Description (brief overview)
- Items (multi-line text list)
- Image path (uploaded photo reference)
- QR code (unique identifier, TOTE-XXXXX format)
- Created/Updated timestamps

## REST API Endpoints

**Collection:**
- `GET /api/totes` - List all totes

**Single Resource:**
- `POST /api/tote` - Create new tote
- `GET /api/tote/{id}` - Get tote by ID
- `GET /api/tote/qr/{qr_code}` - Get tote by QR code
- `PUT /api/tote/{id}` - Update tote by ID
- `DELETE /api/tote/{id}` - Delete tote

**Import/Export:**
- `GET /api/export` - Export all totes as JSON
- `POST /api/import` - Import totes from JSON
- `DELETE /api/totes/delete-all` - Remove all totes

**Utilities:**
- `POST /api/upload-image` - Upload container image
- `GET /api/settings` - Get app settings
- `PUT /api/settings` - Update app settings

## Configuration
- **Default Port:** 3818 (configurable in settings.json)
- **Database Location:** `./totetrax.db` (created automatically on first run)
- **Web Interface:** `http://localhost:3818`
- **Image Uploads:** `web/static/images/uploads/`

## Development Status
- [x] Initial setup (git repository initialized)
- [x] Project requirements defined
- [x] Technology stack selection (Go + SQLite)
- [x] Database schema design
- [x] REST API design
- [x] Backend service implementation (full CRUD)
- [x] Complete web interface
- [x] QR code generation and scanning
- [x] Image upload functionality
- [x] Printable labels
- [x] Import/Export
- [x] Project documentation
- [x] Executable built and tested
- [ ] Dark mode theme
- [ ] Settings page
- [ ] Testing suite
- [ ] Production deployment

## Design Principles
- Must use only open source or free components
- Modular architecture for easy feature additions
- Local-first: All data stored locally, no cloud dependencies
- API-first design to enable companion apps
- RESTful API with singular/plural resource naming conventions

## Target Users
- Home organizers managing storage containers
- People moving or packing
- Garage/workshop organization
- Warehouse/inventory management (small scale)
- Seasonal storage tracking

## Use Cases
- **Moving House**: Track which boxes contain which items
- **Garage Organization**: Find tools and supplies quickly
- **Seasonal Storage**: Know what's in holiday decoration boxes
- **Workshop**: Organize parts, components, and materials
- **General Storage**: Manage attic, basement, or storage unit contents

## Future Enhancements
- Dark/light theme toggle
- Settings page with preferences
- Location/room tracking
- Advanced search and filtering
- Usage statistics
- Barcode scanning support
- Mobile companion app (Android)
- CSV export option
- Bulk operations
- Print multiple labels at once
- Custom label sizes

## Recent Changes

### 2026-02-14 - Initial Build
- Project created based on Filatrax architecture
- Complete backend implementation:
  - SQLite database with totes schema
  - Full CRUD REST API
  - Pure Go SQLite driver (no CGO dependencies)
  - QR code system (TOTE-XXXXX format)
- Web UI implemented:
  - Dashboard with search
  - Add/edit forms with image upload
  - Detail view with QR code display
  - QR scanning page (3 methods)
  - Printable labels
- Features:
  - Image upload to local storage
  - QR code generation and scanning
  - Import/export JSON
  - Responsive design
  - Real-time search
- All API endpoints tested and functional
- README documentation complete
- Executable built successfully (15MB)

## Differences from Filatrax

### Simplified Data Model
- Single entity (Tote) vs multiple (Filament + FilamentTypes)
- No temperature or weight tracking
- No custom type management
- Focus on general storage vs specialized 3D printing use case

### Enhanced Features
- Image upload and display
- Multi-line item lists
- Simpler UI focused on quick lookup

### Similar Features
- QR code generation and scanning
- Printable labels
- Import/Export
- REST API architecture
- Pure Go SQLite implementation
