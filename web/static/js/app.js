// DiscordBotForge Web Interface JavaScript

class DiscordBotForge {
    constructor() {
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.init();
    }

    init() {
        this.connectWebSocket();
        this.setupEventListeners();
        this.startStatusUpdates();
    }

    connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;
        
        try {
            this.ws = new WebSocket(wsUrl);
            
            this.ws.onopen = () => {
                console.log('WebSocket connected');
                this.reconnectAttempts = 0;
                this.updateConnectionStatus(true);
            };
            
            this.ws.onmessage = (event) => {
                const data = JSON.parse(event.data);
                this.handleWebSocketMessage(data);
            };
            
            this.ws.onclose = () => {
                console.log('WebSocket disconnected');
                this.updateConnectionStatus(false);
                this.attemptReconnect();
            };
            
            this.ws.onerror = (error) => {
                console.error('WebSocket error:', error);
                this.updateConnectionStatus(false);
            };
        } catch (error) {
            console.error('Failed to connect WebSocket:', error);
            this.updateConnectionStatus(false);
        }
    }

    attemptReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            console.log(`Attempting to reconnect... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
            setTimeout(() => {
                this.connectWebSocket();
            }, 2000 * this.reconnectAttempts);
        }
    }

    updateConnectionStatus(connected) {
        const statusElement = document.getElementById('connection-status');
        if (statusElement) {
            if (connected) {
                statusElement.className = 'badge bg-success';
                statusElement.innerHTML = '<i class="fas fa-circle"></i> Connected';
            } else {
                statusElement.className = 'badge bg-danger';
                statusElement.innerHTML = '<i class="fas fa-circle"></i> Disconnected';
            }
        }
    }

    handleWebSocketMessage(data) {
        // Update bot status
        if (data.running !== undefined) {
            this.updateBotStatus(data);
        }
        
        // Update activity log
        if (data.activity) {
            this.addActivityLog(data.activity);
        }
        
        // Update logs
        if (data.logs) {
            this.addLogEntry(data.logs);
        }
    }

    updateBotStatus(data) {
        const elements = {
            'bot-status': data.running ? 'Running' : 'Stopped',
            'bot-uptime': data.uptime || 'Unknown',
            'last-update': new Date().toLocaleTimeString(),
            'messages-count': data.stats?.messages || 0,
            'commands-count': data.stats?.commands_executed || 0
        };

        Object.entries(elements).forEach(([id, value]) => {
            const element = document.getElementById(id);
            if (element) {
                element.textContent = value;
            }
        });
    }

    addActivityLog(activity) {
        const container = document.getElementById('activity-log');
        if (!container) return;

        const logEntry = document.createElement('div');
        logEntry.className = 'activity-item';
        logEntry.innerHTML = `
            <i class="fas fa-info-circle text-info"></i>
            <span class="ms-2">${activity.message}</span>
            <small class="text-muted ms-2">${new Date().toLocaleTimeString()}</small>
        `;

        container.insertBefore(logEntry, container.firstChild);
        
        // Keep only last 10 entries
        while (container.children.length > 10) {
            container.removeChild(container.lastChild);
        }
    }

    addLogEntry(logData) {
        const container = document.getElementById('logs-container');
        if (!container) return;

        const logEntry = document.createElement('div');
        logEntry.className = 'log-entry';
        logEntry.innerHTML = `
            <span class="log-time">[${new Date().toLocaleTimeString()}]</span>
            <span class="log-level ${logData.level}">[${logData.level.toUpperCase()}]</span>
            <span class="log-message">${logData.message}</span>
        `;

        container.appendChild(logEntry);
        
        // Auto-scroll if enabled
        const autoScroll = document.getElementById('auto-scroll');
        if (autoScroll && autoScroll.checked) {
            container.scrollTop = container.scrollHeight;
        }
    }

    setupEventListeners() {
        // Settings form
        const settingsForm = document.getElementById('settings-form');
        if (settingsForm) {
            settingsForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.saveSettings();
            });
        }

        // Add command form
        const addCommandForm = document.getElementById('add-command-form');
        if (addCommandForm) {
            addCommandForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.addCommand();
            });
        }

        // Log filters
        const logSearch = document.getElementById('log-search');
        if (logSearch) {
            logSearch.addEventListener('input', () => {
                this.filterLogs();
            });
        }
    }

    startStatusUpdates() {
        // Update status every 30 seconds
        setInterval(() => {
            this.fetchBotStatus();
        }, 30000);
    }

    async fetchBotStatus() {
        try {
            const response = await fetch('/api/status');
            const data = await response.json();
            this.updateBotStatus(data);
        } catch (error) {
            console.error('Failed to fetch bot status:', error);
        }
    }

    async saveSettings() {
        const formData = new FormData(document.getElementById('settings-form'));
        const settings = Object.fromEntries(formData);
        
        try {
            const response = await fetch('/api/settings', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(settings)
            });
            
            if (response.ok) {
                this.showAlert('Settings saved successfully!', 'success');
            } else {
                this.showAlert('Failed to save settings', 'danger');
            }
        } catch (error) {
            console.error('Error saving settings:', error);
            this.showAlert('Error saving settings', 'danger');
        }
    }

    async addCommand() {
        const formData = new FormData(document.getElementById('add-command-form'));
        const command = Object.fromEntries(formData);
        
        try {
            const response = await fetch('/api/commands', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(command)
            });
            
            if (response.ok) {
                this.showAlert('Command added successfully!', 'success');
                document.getElementById('add-command-form').reset();
                this.refreshCommands();
            } else {
                this.showAlert('Failed to add command', 'danger');
            }
        } catch (error) {
            console.error('Error adding command:', error);
            this.showAlert('Error adding command', 'danger');
        }
    }

    async restartBot() {
        try {
            const response = await fetch('/api/restart', {
                method: 'POST'
            });
            
            if (response.ok) {
                this.showAlert('Bot restart initiated', 'info');
            } else {
                this.showAlert('Failed to restart bot', 'danger');
            }
        } catch (error) {
            console.error('Error restarting bot:', error);
            this.showAlert('Error restarting bot', 'danger');
        }
    }

    async stopBot() {
        if (confirm('Are you sure you want to stop the bot?')) {
            try {
                const response = await fetch('/api/stop', {
                    method: 'POST'
                });
                
                if (response.ok) {
                    this.showAlert('Bot stop initiated', 'warning');
                } else {
                    this.showAlert('Failed to stop bot', 'danger');
                }
            } catch (error) {
                console.error('Error stopping bot:', error);
                this.showAlert('Error stopping bot', 'danger');
            }
        }
    }

    refreshStatus() {
        this.fetchBotStatus();
        this.showAlert('Status refreshed', 'info');
    }

    clearLogs() {
        const container = document.getElementById('logs-container');
        if (container) {
            container.innerHTML = '';
        }
        this.showAlert('Logs cleared', 'info');
    }

    downloadLogs() {
        // In a real implementation, this would download the log file
        this.showAlert('Log download started', 'info');
    }

    filterLogs() {
        const searchTerm = document.getElementById('log-search').value.toLowerCase();
        const logLevel = document.getElementById('log-level').value;
        const entries = document.querySelectorAll('.log-entry');
        
        entries.forEach(entry => {
            const message = entry.textContent.toLowerCase();
            const level = entry.querySelector('.log-level').textContent.toLowerCase();
            
            const matchesSearch = !searchTerm || message.includes(searchTerm);
            const matchesLevel = logLevel === 'all' || level.includes(logLevel);
            
            entry.style.display = matchesSearch && matchesLevel ? 'block' : 'none';
        });
    }

    applyFilters() {
        this.filterLogs();
        this.showAlert('Filters applied', 'info');
    }

    showCommandDetails(commandName) {
        // In a real implementation, this would fetch and display command details
        this.showAlert(`Showing details for command: ${commandName}`, 'info');
    }

    showModuleDetails(moduleName) {
        // In a real implementation, this would fetch and display module details
        this.showAlert(`Showing details for module: ${moduleName}`, 'info');
    }

    restartModule(moduleName) {
        if (confirm(`Are you sure you want to restart the ${moduleName} module?`)) {
            this.showAlert(`Restarting module: ${moduleName}`, 'info');
        }
    }

    refreshCommands() {
        // In a real implementation, this would refresh the commands table
        this.showAlert('Commands refreshed', 'info');
    }

    showAlert(message, type) {
        const alertDiv = document.createElement('div');
        alertDiv.className = `alert alert-${type} alert-dismissible fade show`;
        alertDiv.innerHTML = `
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        `;
        
        // Insert at the top of the container
        const container = document.querySelector('.container-fluid');
        if (container) {
            container.insertBefore(alertDiv, container.firstChild);
            
            // Auto-dismiss after 5 seconds
            setTimeout(() => {
                if (alertDiv.parentNode) {
                    alertDiv.remove();
                }
            }, 5000);
        }
    }
}

// Global functions for HTML onclick handlers
function restartBot() {
    window.discordBotForge.restartBot();
}

function stopBot() {
    window.discordBotForge.stopBot();
}

function refreshStatus() {
    window.discordBotForge.refreshStatus();
}

function clearLogs() {
    window.discordBotForge.clearLogs();
}

function downloadLogs() {
    window.discordBotForge.downloadLogs();
}

function applyFilters() {
    window.discordBotForge.applyFilters();
}

function showCommandDetails(commandName) {
    window.discordBotForge.showCommandDetails(commandName);
}

function showModuleDetails(moduleName) {
    window.discordBotForge.showModuleDetails(moduleName);
}

function restartModule(moduleName) {
    window.discordBotForge.restartModule(moduleName);
}

function refreshCommands() {
    window.discordBotForge.refreshCommands();
}

// Initialize the application when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.discordBotForge = new DiscordBotForge();
});
