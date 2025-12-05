package metrics

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/remotemgmt/gobased-remote-mgmt/pkg/common"
)

// MetricsCollector collects system metrics
type MetricsCollector struct {
	mu              sync.RWMutex
	ctx             context.Context
	cancel          context.CancelFunc
	pollingInterval time.Duration
	metrics         []common.MetricsSnapshot
	anomalyHandler  AnomalyHandler
}

// AnomalyHandler is called when an anomaly is detected
type AnomalyHandler func(alert common.AnomalyAlert)

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(ctx context.Context) *MetricsCollector {
	childCtx, cancel := context.WithCancel(ctx)
	return &MetricsCollector{
		ctx:             childCtx,
		cancel:          cancel,
		pollingInterval: 5 * time.Second,
		metrics:         make([]common.MetricsSnapshot, 0),
	}
}

// SetAnomalyHandler sets the anomaly handler
func (mc *MetricsCollector) SetAnomalyHandler(handler AnomalyHandler) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.anomalyHandler = handler
}

// Start starts the metrics collection loop
func (mc *MetricsCollector) Start() error {
	go func() {
		ticker := time.NewTicker(mc.pollingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-mc.ctx.Done():
				return
			case <-ticker.C:
				snapshot, err := mc.collectMetrics()
				if err != nil {
					log.Printf("Error collecting metrics: %v", err)
					continue
				}

				mc.mu.Lock()
				mc.metrics = append(mc.metrics, snapshot)
				// Keep only last 1000 metrics
				if len(mc.metrics) > 1000 {
					mc.metrics = mc.metrics[1:]
				}
				mc.mu.Unlock()

				// Check for anomalies
				mc.checkAnomalies(snapshot)
			}
		}
	}()

	log.Println("Metrics collection started")
	return nil
}

// collectMetrics collects current system metrics
func (mc *MetricsCollector) collectMetrics() (common.MetricsSnapshot, error) {
	snapshot := common.MetricsSnapshot{
		Timestamp: time.Now(),
	}

	// TODO: Implement actual metrics collection using gopsutil
	// For now, return dummy values
	snapshot.CPUPercent = 0.0
	snapshot.MemPercent = 0.0
	snapshot.Temperature = 0.0
	snapshot.Uptime = 0

	return snapshot, nil
}

// checkAnomalies checks for anomalies in the metrics
func (mc *MetricsCollector) checkAnomalies(snapshot common.MetricsSnapshot) {
	if mc.anomalyHandler == nil {
		return
	}

	// Check temperature threshold
	if snapshot.Temperature > 80.0 {
		mc.anomalyHandler(common.AnomalyAlert{
			AlertType:  "temperature",
			Severity:   "critical",
			Message:    "System temperature critically high",
			Value:      snapshot.Temperature,
			Threshold:  80.0,
			Timestamp:  time.Now(),
		})
	}

	// Check memory threshold
	if snapshot.MemPercent > 90.0 {
		mc.anomalyHandler(common.AnomalyAlert{
			AlertType:  "memory",
			Severity:   "warning",
			Message:    "Memory usage high",
			Value:      snapshot.MemPercent,
			Threshold:  90.0,
			Timestamp:  time.Now(),
		})
	}

	// Check CPU threshold
	if snapshot.CPUPercent > 95.0 {
		mc.anomalyHandler(common.AnomalyAlert{
			AlertType:  "cpu",
			Severity:   "warning",
			Message:    "CPU usage very high",
			Value:      snapshot.CPUPercent,
			Threshold:  95.0,
			Timestamp:  time.Now(),
		})
	}
}

// GetLatestMetrics returns the latest metrics
func (mc *MetricsCollector) GetLatestMetrics() *common.MetricsSnapshot {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if len(mc.metrics) == 0 {
		return nil
	}

	latest := mc.metrics[len(mc.metrics)-1]
	return &latest
}

// GetMetricsHistory returns the metrics history
func (mc *MetricsCollector) GetMetricsHistory(limit int) []common.MetricsSnapshot {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if limit <= 0 || limit > len(mc.metrics) {
		limit = len(mc.metrics)
	}

	start := len(mc.metrics) - limit
	if start < 0 {
		start = 0
	}

	result := make([]common.MetricsSnapshot, limit)
	copy(result, mc.metrics[start:])
	return result
}

// Stop stops the metrics collection
func (mc *MetricsCollector) Stop() error {
	mc.cancel()
	log.Println("Metrics collection stopped")
	return nil
}
