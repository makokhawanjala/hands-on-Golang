package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemStats represents system resource statistics
type SystemStats struct {
	Timestamp     time.Time `json:"timestamp"`
	CPUPercent    float64   `json:"cpu_percent"`
	MemoryTotal   uint64    `json:"memory_total"`
	MemoryUsed    uint64    `json:"memory_used"`
	MemoryPercent float64   `json:"memory_percent"`
	DiskTotal     uint64    `json:"disk_total"`
	DiskUsed      uint64    `json:"disk_used"`
	DiskPercent   float64   `json:"disk_percent"`
	Goroutines    int       `json:"goroutines"`
}

// SystemInfo represents static system information
type SystemInfo struct {
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformFamily  string `json:"platform_family"`
	PlatformVersion string `json:"platform_version"`
	CPUCores        int    `json:"cpu_cores"`
	Uptime          uint64 `json:"uptime"`
}

// Monitor handles system monitoring
type Monitor struct {
	stats     []SystemStats
	statsMu   sync.RWMutex
	clients   map[*websocket.Conn]bool
	clientsMu sync.RWMutex
	upgrader  websocket.Upgrader
}

// NewMonitor creates a new system monitor
func NewMonitor() *Monitor {
	return &Monitor{
		stats:   make([]SystemStats, 0),
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for demo
			},
		},
	}
}

// collectStats gathers current system statistics
func (m *Monitor) collectStats() SystemStats {
	// CPU usage
	cpuPercent, _ := cpu.Percent(time.Second, false)
	var cpuUsage float64
	if len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	// Memory usage
	memInfo, _ := mem.VirtualMemory()

	// Disk usage (root partition)
	diskInfo, _ := disk.Usage("/")

	return SystemStats{
		Timestamp:     time.Now(),
		CPUPercent:    cpuUsage,
		MemoryTotal:   memInfo.Total,
		MemoryUsed:    memInfo.Used,
		MemoryPercent: memInfo.UsedPercent,
		DiskTotal:     diskInfo.Total,
		DiskUsed:      diskInfo.Used,
		DiskPercent:   diskInfo.UsedPercent,
		Goroutines:    runtime.NumGoroutine(),
	}
}

// startMonitoring begins collecting system stats
func (m *Monitor) startMonitoring() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := m.collectStats()

			// Store stats (keep last 100 entries)
			m.statsMu.Lock()
			m.stats = append(m.stats, stats)
			if len(m.stats) > 100 {
				m.stats = m.stats[1:]
			}
			m.statsMu.Unlock()

			// Broadcast to WebSocket clients
			m.broadcastStats(stats)
		}
	}
}

// broadcastStats sends stats to all connected WebSocket clients
func (m *Monitor) broadcastStats(stats SystemStats) {
	m.clientsMu.RLock()
	defer m.clientsMu.RUnlock()

	data, _ := json.Marshal(stats)

	for client := range m.clients {
		err := client.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			client.Close()
			delete(m.clients, client)
		}
	}
}

// getSystemInfo returns static system information
func getSystemInfo() SystemInfo {
	hostInfo, _ := host.Info()

	return SystemInfo{
		Hostname:        hostInfo.Hostname,
		OS:              hostInfo.OS,
		Platform:        hostInfo.Platform,
		PlatformFamily:  hostInfo.PlatformFamily,
		PlatformVersion: hostInfo.PlatformVersion,
		CPUCores:        runtime.NumCPU(),
		Uptime:          hostInfo.Uptime,
	}
}

// HTTP Handlers

func (m *Monitor) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	// Add client
	m.clientsMu.Lock()
	m.clients[conn] = true
	m.clientsMu.Unlock()

	// Send current stats immediately
	m.statsMu.RLock()
	if len(m.stats) > 0 {
		data, _ := json.Marshal(m.stats[len(m.stats)-1])
		conn.WriteMessage(websocket.TextMessage, data)
	}
	m.statsMu.RUnlock()

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	// Remove client
	m.clientsMu.Lock()
	delete(m.clients, conn)
	m.clientsMu.Unlock()
}

func (m *Monitor) handleCurrentStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	stats := m.collectStats()
	json.NewEncoder(w).Encode(stats)
}

func (m *Monitor) handleHistoricalStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get limit parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	m.statsMu.RLock()
	defer m.statsMu.RUnlock()

	start := 0
	if len(m.stats) > limit {
		start = len(m.stats) - limit
	}

	json.NewEncoder(w).Encode(m.stats[start:])
}

func handleSystemInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	info := getSystemInfo()
	json.NewEncoder(w).Encode(info)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := map[string]string{
		"status": "healthy",
		"time":   time.Now().UTC().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(response)
}

// Serve static files (dashboard)
func serveDashboard(w http.ResponseWriter, r *http.Request) {
	dashboard := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>System Monitor Dashboard</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: #333;
            min-height: 100vh;
        }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        .header { 
            background: rgba(255,255,255,0.95);
            padding: 20px;
            border-radius: 15px;
            margin-bottom: 20px;
            box-shadow: 0 8px 32px rgba(0,0,0,0.1);
        }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        .stat-card {
            background: rgba(255,255,255,0.95);
            padding: 25px;
            border-radius: 15px;
            box-shadow: 0 8px 32px rgba(0,0,0,0.1);
            backdrop-filter: blur(10px);
        }
        .stat-value { font-size: 2.5rem; font-weight: bold; color: #4f46e5; }
        .stat-label { font-size: 0.9rem; color: #6b7280; margin-top: 5px; }
        .chart-container {
            background: rgba(255,255,255,0.95);
            padding: 25px;
            border-radius: 15px;
            box-shadow: 0 8px 32px rgba(0,0,0,0.1);
            height: 400px;
        }
        .status { 
            display: inline-block;
            padding: 5px 12px;
            border-radius: 20px;
            font-size: 0.8rem;
            font-weight: bold;
        }
        .status.online { background: #dcfce7; color: #166534; }
        .progress-bar {
            width: 100%;
            height: 8px;
            background: #e5e7eb;
            border-radius: 4px;
            overflow: hidden;
            margin-top: 10px;
        }
        .progress-fill {
            height: 100%;
            background: linear-gradient(90deg, #10b981, #059669);
            transition: width 0.3s ease;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🖥️ System Monitor Dashboard</h1>
            <p>Real-time system resource monitoring</p>
            <span class="status online" id="connectionStatus">🟢 Connected</span>
        </div>

        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-value" id="cpuValue">0%</div>
                <div class="stat-label">CPU Usage</div>
                <div class="progress-bar">
                    <div class="progress-fill" id="cpuProgress" style="width: 0%"></div>
                </div>
            </div>
            
            <div class="stat-card">
                <div class="stat-value" id="memoryValue">0%</div>
                <div class="stat-label">Memory Usage</div>
                <div class="progress-bar">
                    <div class="progress-fill" id="memoryProgress" style="width: 0%"></div>
                </div>
            </div>
            
            <div class="stat-card">
                <div class="stat-value" id="diskValue">0%</div>
                <div class="stat-label">Disk Usage</div>
                <div class="progress-bar">
                    <div class="progress-fill" id="diskProgress" style="width: 0%"></div>
                </div>
            </div>
            
            <div class="stat-card">
                <div class="stat-value" id="goroutinesValue">0</div>
                <div class="stat-label">Active Goroutines</div>
            </div>
        </div>

        <div class="chart-container">
            <canvas id="systemChart"></canvas>
        </div>
    </div>

    <script>
        // WebSocket connection
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = protocol + '//' + window.location.host + '/ws';
        let ws;
        let chart;
        const maxDataPoints = 50;

        // Initialize chart
        const ctx = document.getElementById('systemChart').getContext('2d');
        chart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: [],
                datasets: [
                    {
                        label: 'CPU %',
                        data: [],
                        borderColor: '#ef4444',
                        backgroundColor: 'rgba(239, 68, 68, 0.1)',
                        tension: 0.4
                    },
                    {
                        label: 'Memory %',
                        data: [],
                        borderColor: '#3b82f6',
                        backgroundColor: 'rgba(59, 130, 246, 0.1)',
                        tension: 0.4
                    },
                    {
                        label: 'Disk %',
                        data: [],
                        borderColor: '#10b981',
                        backgroundColor: 'rgba(16, 185, 129, 0.1)',
                        tension: 0.4
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    title: {
                        display: true,
                        text: 'System Resources Over Time'
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        max: 100,
                        ticks: {
                            callback: function(value) {
                                return value + '%';
                            }
                        }
                    }
                },
                animation: {
                    duration: 750
                }
            }
        });

        function connectWebSocket() {
            ws = new WebSocket(wsUrl);
            
            ws.onopen = function() {
                document.getElementById('connectionStatus').innerHTML = '🟢 Connected';
                document.getElementById('connectionStatus').className = 'status online';
            };
            
            ws.onclose = function() {
                document.getElementById('connectionStatus').innerHTML = '🔴 Disconnected';
                document.getElementById('connectionStatus').className = 'status offline';
                // Reconnect after 3 seconds
                setTimeout(connectWebSocket, 3000);
            };
            
            ws.onmessage = function(event) {
                const data = JSON.parse(event.data);
                updateDashboard(data);
            };
        }

        function updateDashboard(stats) {
            // Update stat cards
            document.getElementById('cpuValue').textContent = stats.cpu_percent.toFixed(1) + '%';
            document.getElementById('memoryValue').textContent = stats.memory_percent.toFixed(1) + '%';
            document.getElementById('diskValue').textContent = stats.disk_percent.toFixed(1) + '%';
            document.getElementById('goroutinesValue').textContent = stats.goroutines;

            // Update progress bars
            document.getElementById('cpuProgress').style.width = stats.cpu_percent + '%';
            document.getElementById('memoryProgress').style.width = stats.memory_percent + '%';
            document.getElementById('diskProgress').style.width = stats.disk_percent + '%';

            // Update chart
            const time = new Date(stats.timestamp).toLocaleTimeString();
            
            chart.data.labels.push(time);
            chart.data.datasets[0].data.push(stats.cpu_percent);
            chart.data.datasets[1].data.push(stats.memory_percent);
            chart.data.datasets[2].data.push(stats.disk_percent);

            // Keep only last N data points
            if (chart.data.labels.length > maxDataPoints) {
                chart.data.labels.shift();
                chart.data.datasets.forEach(dataset => dataset.data.shift());
            }

            chart.update('none');
        }

        // Start WebSocket connection
        connectWebSocket();
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashboard))
}

func main() {
	monitor := NewMonitor()

	// Start monitoring in background
	go monitor.startMonitoring()

	// Setup routes
	router := mux.NewRouter()

	// API endpoints
	router.HandleFunc("/api/stats/current", monitor.handleCurrentStats).Methods("GET")
	router.HandleFunc("/api/stats/history", monitor.handleHistoricalStats).Methods("GET")
	router.HandleFunc("/api/system/info", handleSystemInfo).Methods("GET")
	router.HandleFunc("/api/health", handleHealth).Methods("GET")

	// WebSocket endpoint
	router.HandleFunc("/ws", monitor.handleWebSocket)

	// Dashboard
	router.HandleFunc("/", serveDashboard).Methods("GET")

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 System Monitor starting on port %s\n", port)
	fmt.Printf("📊 Dashboard: http://localhost:%s\n", port)
	fmt.Printf("🔗 API: http://localhost:%s/api/stats/current\n", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
