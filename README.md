# ToteTrax - Storage Container Inventory Management

A cross-platform (Windows/Linux) storage container tracking application with web interface and REST API.

## Features

- 📦 Track storage totes/boxes/containers by name
- 🖼️ Add images of container contents
- 📝 Maintain text lists of items in each container
- 🏷️ QR code generation and scanning for quick access
- 🖨️ Printable labels with QR codes
- 🌐 Web-based user interface
- 🔌 REST API for external integrations
- 💾 Local SQLite database (no cloud dependencies)
- 📥 Import/Export inventory data

## Technology Stack

- **Language:** Go 1.25+
- **Database:** SQLite (pure Go driver - no CGO)
- **Web Server:** Go standard library net/http
- **Architecture:** Modular, API-first design

## Installation

### Requirements

- **None!** Just download and run the executable

### Quick Start

1. Download the `totetrax` executable for your platform
2. Run it: `./totetrax` (Linux) or `totetrax.exe` (Windows)
3. Open your browser to `http://localhost:3818`

## Building from Source

### Prerequisites

- Go 1.25 or later

### Build Commands

```bash
# Clone the repository
git clone <your-repo-url>
cd totetrax

# Install dependencies
go mod download

# Build for current platform
go build -o totetrax cmd/totetrax/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o totetrax.exe cmd/totetrax/main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o totetrax cmd/totetrax/main.go
```

## API Documentation

### Endpoints

#### List All Totes
```
GET /api/totes
```

#### Create Tote
```
POST /api/tote
Content-Type: application/json

{
  "name": "Kitchen Supplies",
  "description": "Extra kitchen items",
  "items": "4x Dish towels\n2x Pot holders\n1x Apron",
  "image_path": "/static/images/uploads/123456.jpg"
}
```

#### Get Tote by ID
```
GET /api/tote/{id}
```

#### Get Tote by QR Code
```
GET /api/tote/qr/{qr_code}
```

#### Update Tote
```
PUT /api/tote/{id}
Content-Type: application/json

{
  "name": "Updated name",
  "items": "Updated list"
}
```

#### Delete Tote
```
DELETE /api/tote/{id}
```

#### Upload Image
```
POST /api/upload-image
Content-Type: multipart/form-data

FormData with 'image' field
```

#### Export Inventory
```
GET /api/export
```

#### Import Inventory
```
POST /api/import
Content-Type: application/json

[{tote objects}]
```

#### Delete All Totes
```
DELETE /api/totes/delete-all
```

## Project Structure

```
totetrax/
├── cmd/
│   └── totetrax/
│       └── main.go           # Application entry point
├── internal/
│   ├── api/                  # HTTP handlers and routing
│   │   ├── handlers.go
│   │   └── router.go
│   ├── database/             # Database initialization
│   │   └── database.go
│   ├── models/               # Data models
│   │   ├── tote.go
│   │   └── settings.go
│   └── service/              # Business logic
│       ├── tote_service.go
│       └── settings_service.go
├── web/
│   └── static/               # CSS, JS, images
│       ├── css/
│       ├── js/
│       └── images/
├── go.mod
└── README.md
```

## Development

### Running in Development Mode

```bash
go run cmd/totetrax/main.go
```

### Adding Dependencies

```bash
go get <package-name>
```

## Database

The application uses SQLite for local storage. The database file (`totetrax.db`) is created automatically in the current directory on first run.

### Schema

**totes** table:
- id (INTEGER PRIMARY KEY)
- name (TEXT)
- description (TEXT)
- items (TEXT)
- image_path (TEXT)
- qr_code (TEXT UNIQUE)
- created_at (DATETIME)
- updated_at (DATETIME)

## Usage Tips

### Organizing Your Storage

1. **Name Your Totes Clearly**: Use descriptive names like "Holiday Decorations", "Winter Clothes", "Garage Tools"
2. **Add Images**: Take photos of your packed containers to quickly see contents
3. **List Items**: Maintain detailed item lists for easy searching
4. **Print Labels**: Generate and print QR code labels for each tote
5. **Scan to Find**: Use the QR scanner to quickly locate specific containers

### QR Code Scanning

- **Method 1**: Use your device camera (requires HTTPS or localhost)
- **Method 2**: Upload a photo of the QR code
- **Method 3**: Manually type the code (e.g., TOTE-00001)

## Configuration

- **Default Port:** 3818 (configurable in settings.json)
- **Database Location:** `./totetrax.db` (created automatically on first run)
- **Web Interface:** `http://localhost:3818`
- **Image Uploads:** Stored in `web/static/images/uploads/`

### Network Access

To access from other devices on your local network:
1. Allow port 3818 in your firewall
2. Access via `http://<your-computer-ip>:3818` from other devices
3. Firewall rule example (Windows): 
   ```powershell
   New-NetFirewallRule -DisplayName "ToteTrax Web Server" -Direction Inbound -LocalPort 3818 -Protocol TCP -Action Allow
   ```

## Future Enhancements

- [ ] Dark mode theme toggle
- [ ] Mobile companion app
- [ ] Barcode scanning support
- [ ] Location/room tracking
- [ ] Statistics and reports
- [ ] CSV export/import
- [ ] Multi-user support

## License

Open Source - License TBD

## Contributing

Contributions welcome! Please open an issue or submit a pull request.

## Inspiration

Built using the same technology stack as [Filatrax](https://github.com/yourusername/filatrax) - a 3D printer filament inventory management system.
