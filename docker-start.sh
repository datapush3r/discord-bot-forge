#!/bin/bash

# DiscordBotForge Docker Startup Script

set -e

echo "🔥 DiscordBotForge Docker Startup Script"
echo "========================================"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found. Creating from template..."
    if [ -f env.docker ]; then
        cp env.docker .env
        echo "✅ Created .env from env.docker template"
        echo "📝 Please edit .env with your Discord bot token before continuing"
        echo ""
        read -p "Press Enter after you've configured .env file..."
    else
        echo "❌ env.docker template not found. Please create .env file manually."
        exit 1
    fi
fi

# Check if Discord bot token is set
if ! grep -q "DISCORD_BOT_TOKEN=your_discord_bot_token_here" .env; then
    echo "✅ Discord bot token appears to be configured"
else
    echo "❌ Please set your Discord bot token in .env file"
    exit 1
fi

echo ""
echo "🚀 Starting DiscordBotForge with Docker..."
echo ""

# Choose deployment mode
echo "Select deployment mode:"
echo "1) Development (with hot reload)"
echo "2) Production (full stack with monitoring)"
echo ""
read -p "Enter choice (1 or 2): " choice

case $choice in
    1)
        echo "🔧 Starting development environment..."
        docker-compose -f docker-compose.dev.yml up -d
        echo ""
        echo "✅ Development environment started!"
        echo "🌐 Web interface: http://localhost:8080"
        echo "📊 View logs: make docker-logs-dev"
        ;;
    2)
        echo "🏭 Starting production environment..."
        docker-compose up -d
        echo ""
        echo "✅ Production environment started!"
        echo "🌐 Web interface: http://localhost"
        echo "📊 Grafana: http://localhost:3000 (admin/admin)"
        echo "📈 Prometheus: http://localhost:9090"
        echo "📊 View logs: make docker-logs"
        ;;
    *)
        echo "❌ Invalid choice. Please run the script again."
        exit 1
        ;;
esac

echo ""
echo "🎉 DiscordBotForge is now running!"
echo ""
echo "Useful commands:"
echo "  View logs: make docker-logs (or docker-logs-dev)"
echo "  Stop: make docker-stop (or docker-stop-dev)"
echo "  Shell access: make docker-shell (or docker-shell-dev)"
echo "  Clean up: make docker-clean (or docker-clean-dev)"
echo ""
echo "Happy botting! 🔥"
