package observ

import (
	"sync"
	"time"
)

type Metrics struct {
	mu sync.RWMutex

	httpRequests      map[string]int64
	httpDurations     map[string][]time.Duration
	jobsProcessed     map[string]int64
	jobErrors         map[string]int64
	deliveryAttempts  int64
	deliverySuccesses int64
	deliveryFailures  int64
}

func NewMetrics() *Metrics {
	return &Metrics{
		httpRequests:  make(map[string]int64),
		httpDurations: make(map[string][]time.Duration),
		jobsProcessed: make(map[string]int64),
		jobErrors:     make(map[string]int64),
	}
}

func (m *Metrics) RecordHTTPRequest(path string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.httpRequests[path]++
	m.httpDurations[path] = append(m.httpDurations[path], duration)

	if len(m.httpDurations[path]) > 1000 {
		m.httpDurations[path] = m.httpDurations[path][len(m.httpDurations[path])-1000:]
	}
}

func (m *Metrics) RecordJobProcessed(jobType string, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.jobsProcessed[jobType]++
	if !success {
		m.jobErrors[jobType]++
	}
}

func (m *Metrics) RecordDelivery(success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.deliveryAttempts++
	if success {
		m.deliverySuccesses++
	} else {
		m.deliveryFailures++
	}
}

func (m *Metrics) GetHTTPStats() map[string]HTTPStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]HTTPStats)
	for path, count := range m.httpRequests {
		durations := m.httpDurations[path]
		if len(durations) == 0 {
			continue
		}

		stats[path] = HTTPStats{
			Count: count,
			P50:   percentile(durations, 0.5),
			P95:   percentile(durations, 0.95),
			P99:   percentile(durations, 0.99),
		}
	}

	return stats
}

func (m *Metrics) GetJobStats() map[string]JobStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]JobStats)
	for jobType, count := range m.jobsProcessed {
		stats[jobType] = JobStats{
			Processed: count,
			Errors:    m.jobErrors[jobType],
		}
	}

	return stats
}

func (m *Metrics) GetDeliveryStats() DeliveryStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return DeliveryStats{
		Attempts:  m.deliveryAttempts,
		Successes: m.deliverySuccesses,
		Failures:  m.deliveryFailures,
	}
}

type HTTPStats struct {
	Count int64
	P50   time.Duration
	P95   time.Duration
	P99   time.Duration
}

type JobStats struct {
	Processed int64
	Errors    int64
}

type DeliveryStats struct {
	Attempts  int64
	Successes int64
	Failures  int64
}

func percentile(durations []time.Duration, p float64) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)

	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	index := int(float64(len(sorted)) * p)
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}
