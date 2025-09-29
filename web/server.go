package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"discord-bot-forge/core"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// WebServer handles the web interface for DiscordBotForge
type WebServer struct {
	bot        *core.Bot
	server     *http.Server
	router     *mux.Router
	upgrader   websocket.Upgrader
	clients    map[*websocket.Conn]bool
	clientsMux sync.RWMutex
	templates  *template.Template
}

// BotStatus represents the current status of the bot
type BotStatus struct {
	Running     bool                   `json:"running"`
	Version     string                 `json:"version"`
	Uptime      string                 `json:"uptime"`
	Commands    int                    `json:"commands"`
	Modules     int                    `json:"modules"`
	Middleware  int                    `json:"middleware"`
	Stats       map[string]interface{} `json:"stats"`
	LastUpdate  time.Time              `json:"last_update"`
}

// CommandInfo represents command information for the web interface
type CommandInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Usage       string   `json:"usage"`
	Category    string   `json:"category"`
	Cooldown    int      `json:"cooldown"`
	Permissions []string `json:"permissions"`
}

// ModuleInfo represents module information for the web interface
type ModuleInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

// NewWebServer creates a new web server instance
func NewWebServer(bot *core.Bot, port string) *WebServer {
	router := mux.NewRouter()
	
	ws := &WebServer{
		bot:      bot,
		router:   router,
		clients:  make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}
	
	// Load templates
	ws.loadTemplates()
	
	// Setup routes
	ws.setupRoutes()
	
	ws.server = &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	
	return ws
}

// Start starts the web server
func (ws *WebServer) Start() error {
	log.Printf("🌐 Starting DiscordBotForge web interface on port %s", ws.server.Addr[1:])
	return ws.server.ListenAndServe()
}

// loadTemplates loads HTML templates
func (ws *WebServer) loadTemplates() {
	tmpl := template.New("")
	
	// Parse all templates
	tmpl = template.Must(tmpl.ParseGlob("web/templates/*.html"))
	
	ws.templates = tmpl
}

// setupRoutes configures all HTTP routes
func (ws *WebServer) setupRoutes() {
	// Static files
	ws.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))
	
	// Web pages
	ws.router.HandleFunc("/", ws.handleDashboard).Methods("GET")
	ws.router.HandleFunc("/commands", ws.handleCommands).Methods("GET")
	ws.router.HandleFunc("/modules", ws.handleModules).Methods("GET")
	ws.router.HandleFunc("/logs", ws.handleLogs).Methods("GET")
	ws.router.HandleFunc("/settings", ws.handleSettings).Methods("GET")
	
	// API endpoints
	api := ws.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/status", ws.handleAPIStatus).Methods("GET")
	api.HandleFunc("/commands", ws.handleAPICommands).Methods("GET")
	api.HandleFunc("/modules", ws.handleAPIModules).Methods("GET")
	api.HandleFunc("/logs", ws.handleAPILogs).Methods("GET")
	api.HandleFunc("/restart", ws.handleAPIRestart).Methods("POST")
	api.HandleFunc("/stop", ws.handleAPIStop).Methods("POST")
	
	// WebSocket endpoint
	ws.router.HandleFunc("/ws", ws.handleWebSocket)
}

// handleDashboard serves the main dashboard page
func (ws *WebServer) handleDashboard(w http.ResponseWriter, r *http.Request) {
	status := ws.getBotStatus()
	
	data := map[string]interface{}{
		"Title": "DiscordBotForge Dashboard",
		"Bot":   status,
	}
	
	ws.templates.ExecuteTemplate(w, "dashboard.html", data)
}

// handleCommands serves the commands management page
func (ws *WebServer) handleCommands(w http.ResponseWriter, r *http.Request) {
	commands := ws.getCommandsInfo()
	
	data := map[string]interface{}{
		"Title":    "Commands Management",
		"Commands": commands,
	}
	
	ws.templates.ExecuteTemplate(w, "commands.html", data)
}

// handleModules serves the modules management page
func (ws *WebServer) handleModules(w http.ResponseWriter, r *http.Request) {
	modules := ws.getModulesInfo()
	
	data := map[string]interface{}{
		"Title":   "Modules Management",
		"Modules": modules,
	}
	
	ws.templates.ExecuteTemplate(w, "modules.html", data)
}

// handleLogs serves the logs viewing page
func (ws *WebServer) handleLogs(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Bot Logs",
	}
	
	ws.templates.ExecuteTemplate(w, "logs.html", data)
}

// handleSettings serves the settings page
func (ws *WebServer) handleSettings(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Bot Settings",
		"Config": ws.bot.Config,
	}
	
	ws.templates.ExecuteTemplate(w, "settings.html", data)
}

// API handlers
func (ws *WebServer) handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	status := ws.getBotStatus()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (ws *WebServer) handleAPICommands(w http.ResponseWriter, r *http.Request) {
	commands := ws.getCommandsInfo()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commands)
}

func (ws *WebServer) handleAPIModules(w http.ResponseWriter, r *http.Request) {
	modules := ws.getModulesInfo()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modules)
}

func (ws *WebServer) handleAPILogs(w http.ResponseWriter, r *http.Request) {
	// For now, return empty logs - in a real implementation, you'd read from log files
	logs := []string{"Bot started", "Command executed: ping", "Module loaded: logging"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (ws *WebServer) handleAPIRestart(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, you'd restart the bot
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "restarting"})
}

func (ws *WebServer) handleAPIStop(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, you'd stop the bot
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "stopping"})
}

// WebSocket handler
func (ws *WebServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()
	
	// Add client
	ws.clientsMux.Lock()
	ws.clients[conn] = true
	ws.clientsMux.Unlock()
	
	log.Printf("WebSocket client connected")
	
	// Send initial status
	status := ws.getBotStatus()
	conn.WriteJSON(status)
	
	// Keep connection alive and send updates
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			status := ws.getBotStatus()
			if err := conn.WriteJSON(status); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}
}

// Helper methods
func (ws *WebServer) getBotStatus() BotStatus {
	return BotStatus{
		Running:    ws.bot.Session != nil,
		Version:    ws.bot.Version,
		Uptime:     "2h 15m", // This would be calculated from start time
		Commands:   len(ws.bot.Commands),
		Modules:    len(ws.bot.Modules),
		Middleware: len(ws.bot.Middleware),
		Stats:      map[string]interface{}{"messages": 150, "commands_executed": 25},
		LastUpdate: time.Now(),
	}
}

func (ws *WebServer) getCommandsInfo() []CommandInfo {
	var commands []CommandInfo
	for _, cmd := range ws.bot.Commands {
		commands = append(commands, CommandInfo{
			Name:        cmd.Name(),
			Description: cmd.Description(),
			Usage:       cmd.Usage(),
			Category:    cmd.Category(),
			Cooldown:    cmd.Cooldown(),
			Permissions: cmd.Permissions(),
		})
	}
	return commands
}

func (ws *WebServer) getModulesInfo() []ModuleInfo {
	var modules []ModuleInfo
	for _, mod := range ws.bot.Modules {
		modules = append(modules, ModuleInfo{
			Name:    mod.Name(),
			Version: mod.Version(),
			Status:  "Running",
		})
	}
	return modules
}

// BroadcastToClients sends data to all connected WebSocket clients
func (ws *WebServer) BroadcastToClients(data interface{}) {
	ws.clientsMux.RLock()
	defer ws.clientsMux.RUnlock()
	
	for client := range ws.clients {
		if err := client.WriteJSON(data); err != nil {
			log.Printf("Error broadcasting to client: %v", err)
			delete(ws.clients, client)
			client.Close()
		}
	}
}
