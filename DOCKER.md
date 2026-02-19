# Kontainer Docker Deployment Guide

Complete guide for running Kontainer in Docker on Windows, Linux, NAS devices, and cloud servers.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Docker Compose Deployment](#docker-compose-deployment)
3. [NAS Deployment](#nas-deployment-synology-qnap)
4. [Configuration](#configuration)
5. [Troubleshooting](#troubleshooting)

---

## Quick Start

The fastest way to get Kontainer running with Docker:

```bash
# Clone the repository
git clone https://github.com/yourusername/kontainer.git
cd kontainer

# Start with Docker Compose
docker-compose up -d

# Access at http://localhost:3818
```

---

## Docker Compose Deployment

### Prerequisites
- Docker 20.10+
- Docker Compose 1.29+

### Basic Setup

1. **Create docker-compose.yml** (included in repo):
```yaml
version: '3.8'

services:
  kontainer:
    build: .
    image: kontainer:latest
    container_name: kontainer
    ports:
      - "3818:3818"
    environment:
      - PORT=3818
      - DATABASE_PATH=/data/kontainer.db
      - THEME=dark
    volumes:
      - kontainer-data:/data
    restart: unless-stopped

volumes:
  kontainer-data:
    driver: local
```

2. **Start the container**:
```bash
docker-compose up -d
```

3. **Verify it's running**:
```bash
docker-compose ps
docker-compose logs -f
```

4. **Access the application**:
   - Local: `http://localhost:3818`
   - Network: `http://<your-ip>:3818`

### Custom Port

Change the port by editing `docker-compose.yml`:
```yaml
ports:
  - "8080:3818"  # Access on port 8080
```

### Persistent Data

Data is stored in the `kontainer-data` Docker volume. To backup:
```bash
# Export data via web UI (recommended)
# Or backup the volume
docker run --rm -v kontainer-data:/data -v $(pwd):/backup alpine tar czf /backup/kontainer-backup.tar.gz /data
```

---

## NAS Deployment (Synology, QNAP)

Kontainer works great on NAS devices for centralized home inventory management.

### Synology NAS

#### Via Docker Package

1. **Install Docker** from Package Center

2. **Create project folder**:
   - File Station → Create folder `/docker/kontainer`
   - Upload files: `Dockerfile`, `docker-compose.yml`, and source code

3. **Open Docker app** → Project tab

4. **Create project**:
   - Path: `/docker/kontainer`
   - Source: `docker-compose.yml`
   - Click "Build"

5. **Access**:
   - `http://<nas-ip>:3818`

#### Via SSH (Advanced)

```bash
# SSH into your NAS
ssh admin@your-nas-ip

# Navigate to shared folder
cd /volume1/docker/kontainer

# Build and start
sudo docker-compose up -d
```

### QNAP NAS

1. **Install Container Station** from App Center

2. **Create container**:
   - Container Station → Create → Create Application
   - Upload `docker-compose.yml`
   - Click "Create"

3. **Access**: `http://<nas-ip>:3818`

### Why Run on NAS?

✅ **Centralized** - Access from any device on your network  
✅ **Always available** - NAS runs 24/7  
✅ **Reliable storage** - NAS RAID protects your inventory data  
✅ **Low power** - More efficient than dedicated PC  

### Network Access

To access from other devices:

1. **Find your NAS IP**: Check router or NAS settings
2. **Allow port 3818**: Usually no firewall changes needed on LAN
3. **Access**: `http://<nas-ip>:3818` from any browser

For internet access (not recommended without VPN):
- Set up reverse proxy (Synology: Application Portal)
- Use HTTPS with Let's Encrypt certificate
- Consider Tailscale or WireGuard VPN instead

---

## Configuration

### Environment Variables

Set these in `docker-compose.yml`:

```yaml
environment:
  - PORT=3818                        # Web server port
  - DATABASE_PATH=/data/kontainer.db # Database location
  - THEME=dark                       # Default theme (light/dark)
```

### Volume Mapping

Mount external storage:

```yaml
volumes:
  - /path/on/host:/data              # Linux/Mac
  - /volume1/docker/kontainer:/data  # Synology
  - C:\kontainer-data:/data          # Windows
```

### Resource Limits

For NAS with limited resources:

```yaml
deploy:
  resources:
    limits:
      cpus: '0.5'
      memory: 512M
    reservations:
      cpus: '0.25'
      memory: 256M
```

---

## Troubleshooting

### Container won't start

```bash
# Check logs
docker-compose logs

# Rebuild from scratch
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

### Can't access web interface

1. **Check container is running**:
```bash
docker-compose ps
```

2. **Check port binding**:
```bash
docker-compose port kontainer 3818
```

3. **Test locally**:
```bash
curl http://localhost:3818
```

4. **Firewall**: Ensure port 3818 is allowed

### Database errors

If database becomes corrupted:

1. **Export data** via web UI if possible
2. **Stop container**: `docker-compose down`
3. **Remove volume**: `docker volume rm kontainer-data`
4. **Restart**: `docker-compose up -d`
5. **Import data** from backup

### Permission errors (Linux/NAS)

```bash
# Fix permissions on volume
sudo chown -R 1000:1000 /path/to/volume

# Or in docker-compose.yml
user: "1000:1000"
```

### Updates

```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose build --no-cache
docker-compose up -d
```

---

## Advanced Topics

### Reverse Proxy (NGINX)

```nginx
server {
    listen 80;
    server_name kontainer.yourdomain.com;
    
    location / {
        proxy_pass http://localhost:3818;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### HTTPS with Let's Encrypt

Use Traefik or nginx-proxy with letsencrypt-companion for automatic HTTPS.

### Multi-Instance Setup

Run multiple instances for different locations:

```yaml
# docker-compose.yml
services:
  kontainer-garage:
    image: kontainer:latest
    ports:
      - "3818:3818"
    volumes:
      - garage-data:/data
  
  kontainer-basement:
    image: kontainer:latest
    ports:
      - "3819:3818"
    volumes:
      - basement-data:/data
```

---

## Support

- **Issues**: https://github.com/yourusername/kontainer/issues
- **Discussions**: https://github.com/yourusername/kontainer/discussions
- **Documentation**: https://github.com/yourusername/kontainer

---

**Need help?** Open an issue on GitHub with:
- Your OS/NAS model
- Docker version (`docker --version`)
- Error logs (`docker-compose logs`)
