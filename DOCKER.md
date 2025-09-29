# 🐳 DiscordBotForge Docker Guide

This guide covers running DiscordBotForge using Docker and Docker Compose for both development and production environments.

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose installed
- Discord bot token from [Discord Developer Portal](https://discord.com/developers/applications)

### 1. Clone and Setup
```bash
git clone <repository-url>
cd discord-bot-forge
```

### 2. Configure Environment
```bash
# Copy Docker environment template
cp env.docker .env

# Edit .env with your Discord bot token
nano .env
```

### 3. Start DiscordBotForge
```bash
# Production deployment
make docker-prod

# Or development deployment
make docker-dev
```

### 4. Access Web Interface
- **Production**: http://localhost
- **Development**: http://localhost:8080
- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090

## 🏗️ Architecture

### Production Stack
- **DiscordBotForge**: Main bot application

### Development Stack
- **DiscordBotForge**: Development bot with hot reload

## 📁 Docker Files Structure

```
discord-bot-forge/
├── Dockerfile              # Production Docker image
├── Dockerfile.dev          # Development Docker image
├── docker-compose.yml      # Production stack
├── docker-compose.dev.yml  # Development stack
├── env.docker              # Docker environment template
├── docker/
│   ├── nginx/
│   │   └── nginx.conf      # Nginx configuration
│   ├── postgres/
│   │   └── init.sql        # Database initialization
│   ├── prometheus/
│   │   └── prometheus.yml  # Prometheus configuration
│   └── grafana/
│       └── provisioning/   # Grafana provisioning
└── Makefile               # Docker commands
```

## 🔧 Docker Commands

### Build Commands
```bash
# Build production image
make docker-build

# Build development image
make docker-build-dev
```

### Run Commands
```bash
# Start production stack
make docker-run

# Start development stack
make docker-run-dev

# Quick start (copies env.docker to .env)
make docker-quickstart
```

### Management Commands
```bash
# View logs
make docker-logs          # Production
make docker-logs-dev      # Development

# Access container shell
make docker-shell         # Production
make docker-shell-dev     # Development

# Stop services
make docker-stop          # Production
make docker-stop-dev      # Development

# Clean up
make docker-clean         # Production
make docker-clean-dev      # Development
```

## 🌐 Service Ports

| Service | Port | Description |
|---------|------|-------------|
| DiscordBotForge | 8080 | Web interface |

## 🔒 Environment Variables

### Required Variables
```bash
DISCORD_BOT_TOKEN=your_bot_token_here
BOT_OWNER_ID=your_user_id_here
```

### Optional Variables
```bash
DEBUG=false
WEB_PORT=8080
LOG_LEVEL=info
RATE_LIMIT_ENABLED=true
PERMISSION_CHECKS_ENABLED=true
LOGGING_ENABLED=true
```

## 🔧 Development

### Hot Reload Development
```bash
# Start development environment
make docker-dev

# View logs
make docker-logs-dev

# Access container
make docker-shell-dev
```

### Code Changes
The development container mounts your local code directory, so changes are reflected immediately without rebuilding.

## 🚀 Production Deployment

### Scaling
```bash
# Scale bot instances
docker-compose up -d --scale discord-bot-forge=3
```

### Backup
```bash
# Backup logs
docker run --rm -v discord-bot-forge_logs:/data -v $(pwd):/backup alpine tar czf /backup/logs-backup.tar.gz -C /data .
```

## 🐛 Troubleshooting

### Common Issues

1. **Bot not connecting**
   - Check Discord bot token in `.env`
   - Verify bot permissions in Discord Developer Portal

2. **Web interface not accessible**
   - Check if containers are running: `docker-compose ps`
   - Check logs: `make docker-logs`
   - Verify port 8080 is not blocked by firewall

3. **Container won't start**
   - Check Docker logs: `docker logs discord-bot-forge`
   - Verify environment variables are set correctly

### Debug Commands
```bash
# Check container status
docker-compose ps

# View all logs
docker-compose logs

# Restart specific service
docker-compose restart discord-bot-forge

# Rebuild and restart
docker-compose up -d --build
```

## 📈 Performance Tuning

### Resource Limits
Add to `docker-compose.yml`:
```yaml
services:
  discord-bot-forge:
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '0.5'
```

## 🔐 Security

### Production Security Checklist
- [ ] Use strong Discord bot token
- [ ] Configure proper firewall rules
- [ ] Regular security updates
- [ ] Monitor logs for suspicious activity
- [ ] Use environment variables for sensitive data

### Network Security
- Use Docker networks for service isolation
- Configure proper firewall rules
- Use your existing reverse proxy for SSL/TLS termination

## 📚 Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Discord Developer Portal](https://discord.com/developers/applications)
