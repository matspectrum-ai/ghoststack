package monitoring

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type MetricType int

const (
	MetricCounter MetricType = iota
	MetricGauge
	MetricHistogram
)

type MetricValue struct {
	Type  MetricType        `json:"type"`
	Name  string            `json:"name"`
	Value float64           `json:"value"`
	Tags  map[string]string `json:"tags,omitempty"`
}

type MetricsRegistry struct {
	mu       sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64
}

func NewMetricsRegistry() *MetricsRegistry {
	return &MetricsRegistry{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (r *MetricsRegistry) SetGauge(name string, value float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.gauges[name] = value
}

func (r *MetricsRegistry) IncCounter(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counters[name]++
}

func (r *MetricsRegistry) AddCounter(name string, delta int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counters[name] += delta
}

func (r *MetricsRegistry) GetGauge(name string) float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.gauges[name]
}

func (r *MetricsRegistry) GetCounter(name string) int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.counters[name]
}

func (r *MetricsRegistry) Snapshot() map[string]float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	snap := make(map[string]float64, len(r.gauges)+len(r.counters))
	for k, v := range r.gauges {
		snap[k] = v
	}
	for k, v := range r.counters {
		snap[k] = float64(v)
	}
	return snap
}

type ProviderMetrics struct {
	mu            sync.RWMutex
	RxBytes       int64
	TxBytes       int64
	Latency       time.Duration
	UptimeSeconds int64
	Peers         int
	Errors        int64
}

func (pm *ProviderMetrics) Snapshot() ProviderMetrics {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return ProviderMetrics{
		RxBytes:       pm.RxBytes,
		TxBytes:       pm.TxBytes,
		Latency:       pm.Latency,
		UptimeSeconds: pm.UptimeSeconds,
		Peers:         pm.Peers,
		Errors:        pm.Errors,
	}
}

func (pm *ProviderMetrics) String() string {
	snap := pm.Snapshot()
	return fmt.Sprintf("rx=%d tx=%d latency=%s uptime=%ds peers=%d errors=%d",
		snap.RxBytes, snap.TxBytes, snap.Latency, snap.UptimeSeconds, snap.Peers, snap.Errors)
}

type SystemMetrics struct {
	mu     sync.RWMutex
	CPU    float64
	Memory uint64
	Rx     int64
	Tx     int64
}

func (sm *SystemMetrics) Snapshot() SystemMetrics {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return SystemMetrics{
		CPU:    sm.CPU,
		Memory: sm.Memory,
		Rx:     sm.Rx,
		Tx:     sm.Tx,
	}
}

type MetricsCollector struct {
	system   *SystemMetrics
	provider *ProviderMetrics
	registry *MetricsRegistry
	interval time.Duration
}

func NewMetricsCollector(interval time.Duration) *MetricsCollector {
	return &MetricsCollector{
		system:   &SystemMetrics{},
		provider: &ProviderMetrics{},
		registry: NewMetricsRegistry(),
		interval: interval,
	}
}

func (mc *MetricsCollector) Run(ctx context.Context) error {
	ticker := time.NewTicker(mc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			mc.collect()
		}
	}
}

func (mc *MetricsCollector) collect() {
	mc.system.mu.Lock()
	defer mc.system.mu.Unlock()

	mc.system.CPU = readCPUUsage()
	mc.system.Memory = readMemoryUsage()
	mc.system.Rx = readNetStat("rx_bytes")
	mc.system.Tx = readNetStat("tx_bytes")
}

func readCPUUsage() float64 {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0
	}
	var user, nice, system, idle int64
	_, err = fmt.Sscanf(string(data[:]), "cpu  %d %d %d %d", &user, &nice, &system, &idle)
	if err != nil {
		return 0
	}
	total := user + nice + system + idle
	if total == 0 {
		return 0
	}
	return float64(user+system) / float64(total) * 100
}

func readMemoryUsage() uint64 {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}
	var total, free, buffers, cached uint64
	fmt.Sscanf(string(data[:]), "MemTotal: %d kB\nMemFree: %d kB\nMemAvailable: %d kB", &total, &free, &buffers)
	_ = cached
	if total == 0 {
		return 0
	}
	return (total - free) * 1024
}

func readNetStat(kind string) int64 {
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return 0
	}
	lines := strings.Split(string(data), "\n")
	idx := 1
	if kind == "tx_bytes" {
		idx = 9
	}
	var total int64
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		if !strings.Contains(fields[0], ":") {
			continue
		}
		rx, _ := fmt.Sscanf(fields[idx], "%d", &total)
		if rx == 1 {
			break
		}
	}
	return total
}

func (mc *MetricsCollector) System() *SystemMetrics {
	return mc.system
}

func (mc *MetricsCollector) Provider() *ProviderMetrics {
	return mc.provider
}

func (mc *MetricsCollector) Registry() *MetricsRegistry {
	return mc.registry
}

type AlertSeverity string

const (
	AlertInfo     AlertSeverity = "info"
	AlertWarning  AlertSeverity = "warning"
	AlertCritical AlertSeverity = "critical"
)

type Alert struct {
	ID           string        `json:"id"`
	Severity     AlertSeverity `json:"severity"`
	Title        string        `json:"title"`
	Message      string        `json:"message"`
	Timestamp    time.Time     `json:"timestamp"`
	Acknowledged bool          `json:"acknowledged"`
}

type AlertManager struct {
	mu     sync.RWMutex
	alerts []Alert
}

func NewAlertManager() *AlertManager {
	return &AlertManager{}
}

func (am *AlertManager) Add(severity AlertSeverity, title, message string) Alert {
	am.mu.Lock()
	defer am.mu.Unlock()

	alert := Alert{
		ID:        fmt.Sprintf("alert-%d", time.Now().UnixNano()),
		Severity:  severity,
		Title:     title,
		Message:   message,
		Timestamp: time.Now(),
	}

	am.alerts = append(am.alerts, alert)

	if len(am.alerts) > 100 {
		am.alerts = am.alerts[len(am.alerts)-100:]
	}

	return alert
}

func (am *AlertManager) List(limit int) []Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	if limit <= 0 || limit > len(am.alerts) {
		limit = len(am.alerts)
	}

	result := make([]Alert, limit)
	copy(result, am.alerts[len(am.alerts)-limit:])
	return result
}

func (am *AlertManager) Acknowledge(id string) bool {
	am.mu.Lock()
	defer am.mu.Unlock()

	for i := range am.alerts {
		if am.alerts[i].ID == id {
			am.alerts[i].Acknowledged = true
			return true
		}
	}
	return false
}

func (am *AlertManager) UnacknowledgedCount() int {
	am.mu.RLock()
	defer am.mu.RUnlock()

	count := 0
	for _, a := range am.alerts {
		if !a.Acknowledged {
			count++
		}
	}
	return count
}
