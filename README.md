# Kontainer - Storage Container Inventory Management

> 📦 Never lose track of what's in your storage boxes again!

A self-hosted, cross-platform inventory management system for tracking items in storage containers, boxes, and totes. Features QR code labels, image galleries, and a modern web interface.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com/)

## ✨ Features

- 📦 **Container Tracking** - Track unlimited storage containers with names, descriptions, and locations
- 🖼️ **Image Galleries** - Add multiple photos of container contents stored securely in the database
- 🏷️ **QR Code Labels** - Auto-generate unique QR codes for instant container lookup
- 📝 **Item Lists** - Maintain detailed inventories of what's in each container
- 🔍 **Quick Search** - Find items across all containers instantly
- 🖨️ **Printable Labels** - Generate printer-ready labels with QR codes
- 📥 **Import/Export** - Backup and restore your entire inventory as JSON
- 🌓 **Dark/Light Modes** - Choose your preferred theme
- 🐳 **Docker Support** - Deploy easily on NAS, Raspberry Pi, or cloud servers
- 💾 **Local-First** - All data stored locally in SQLite (no cloud dependencies)
- 🔌 **REST API** - Full API access for integrations and automation

## 📸 Screenshots

<!-- TODO: Add screenshots here -->
*Coming soon - screenshots of dashboard, detail view, and QR scanning*

## 🚀 Quick Start

### Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/yourusername/kontainer.git
cd kontainer

# Start with Docker Compose
docker-compose up -d

# Access at http://localhost:3818
```

See [DOCKER.md](DOCKER.md) for detailed Docker deployment, NAS setup, and configuration options.

### Pre-built Binaries

Download the latest release for your platform:

- [Windows (amd64)](https://github.com/yourusername/kontainer/releases/latest)
- [Linux (amd64)](https://github.com/yourusername/kontainer/releases/latest)
- [macOS (amd64/arm64)](https://github.com/yourusername/kontainer/releases/latest)

```bash
# Linux/macOS
chmod +x kontainer
./kontainer

# Windows
kontainer.exe

# Access at http://localhost:3818
```

### Build from Source

**Prerequisites:** Go 1.24 or later

```bash
# Clone the repository
git clone https://github.com/yourusername/kontainer.git
cd kontainer

# Build
go build -o kontainer cmd/kontainer/main.go

# Run
./kontainer
```

## 🎯 Use Cases

- 🏠 **Home Organization** - Track seasonal decorations, holiday items, camping gear
- 📦 **Moving Houses** - Know exactly which box contains what during a move
- 🔧 **Workshop/Garage** - Organize tools, parts, and supplies
- 📚 **Storage Units** - Manage items in off-site storage
- 🏢 **Small Business** - Track inventory, supplies, and equipment
- 🎨 **Hobby/Craft Supplies** - Organize materials across multiple containers

## 📖 Documentation

- [Docker Deployment Guide](DOCKER.md) - Docker, NAS, and cloud deployment
- [Technical Documentation](TECHNICAL-DOCS.md) - Architecture and development details
- [API Documentation](#api-endpoints) - REST API reference (see below)

## 🔌 API Endpoints

Kontainer provides a full REST API for automation and integrations:

```bash
# List all containers
GET /api/totes

# Get container by ID
GET /api/tote/{id}

# Look up by QR code
GET /api/tote/qr/{qr_code}

# Create container
POST /api/tote

# Update container
PUT /api/tote/{id}

# Delete container
DELETE /api/tote/{id}

# Export all data
GET /api/export

# Import data
POST /api/import
```

Full API documentation available in [TECHNICAL-DOCS.md](TECHNICAL-DOCS.md).

## 🛠️ Technology Stack

- **Backend:** Go 1.24+ (single binary, no runtime dependencies)
- **Database:** SQLite with pure Go driver ([modernc.org/sqlite](https://modernc.org/sqlite))
- **Frontend:** Vanilla HTML/CSS/JavaScript
- **QR Codes:** qrcode.min.js, html5-qrcode.min.js
- **Deployment:** Native binaries or Docker containers

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by the need for simple, self-hosted inventory management
- Built with ❤️ using Go and modern web standards
- QR code libraries: [davidshimjs/qrcodejs](https://github.com/davidshimjs/qrcodejs) and [mebjas/html5-qrcode](https://github.com/mebjas/html5-qrcode)

## 📧 Support

- **Issues:** [GitHub Issues](https://github.com/yourusername/kontainer/issues)
- **Discussions:** [GitHub Discussions](https://github.com/yourusername/kontainer/discussions)

---

**Made with 📦 for anyone tired of forgetting what's in their boxes!**
