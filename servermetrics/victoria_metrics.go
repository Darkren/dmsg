package servermetrics

import (
	"fmt"

	"github.com/VictoriaMetrics/metrics"

	"github.com/skycoin/dmsg/metricsutil"
)

// Metrics collects metrics for metrics tracking system.
type Metrics interface {
	RecordSession(delta DeltaType)
	RecordStream(delta DeltaType)
	SetClientsCount(val int64)
}

// VictoriaMetrics implements `Metrics` using `VictoriaMetrics`.
type VictoriaMetrics struct {
	packetsPerSecond   *metricsutil.VictoriaMetricsGaugeWrapper
	clientsCount       *metricsutil.VictoriaMetricsGaugeWrapper
	activeSessions     *metricsutil.VictoriaMetricsGaugeWrapper
	successfulSessions *metrics.Counter
	failedSessions     *metrics.Counter
	activeStreams      *metricsutil.VictoriaMetricsGaugeWrapper
	successfulStreams  *metrics.Counter
	failedStreams      *metrics.Counter
}

// NewVictoriaMetrics returns the Victoria Metrics implementation of Metrics.
func NewVictoriaMetrics() *VictoriaMetrics {
	return &VictoriaMetrics{
		packetsPerSecond:   metricsutil.NewVictoriaMetricsGauge("packets_per_second"),
		clientsCount:       metricsutil.NewVictoriaMetricsGauge("clients_count"),
		activeSessions:     metricsutil.NewVictoriaMetricsGauge("active_sessions_count"),
		successfulSessions: metrics.GetOrCreateCounter("session_success_total"),
		failedSessions:     metrics.GetOrCreateCounter("session_fail_total"),
		activeStreams:      metricsutil.NewVictoriaMetricsGauge("active_streams_count"),
		successfulStreams:  metrics.GetOrCreateCounter("stream_success_total"),
		failedStreams:      metrics.GetOrCreateCounter("stream_fail_total"),
	}
}

func (m *VictoriaMetrics) SetPacketsPerSecond(val int64)

// SetClientsCount implements `Metrics`.
func (m *VictoriaMetrics) SetClientsCount(val int64) {
	m.clientsCount.Set(val)
}

// RecordSession implements `Metrics`.
func (m *VictoriaMetrics) RecordSession(delta DeltaType) {
	switch delta {
	case 0:
		m.failedSessions.Inc()
	case 1:
		m.successfulSessions.Inc()
		m.activeSessions.Inc()
	case -1:
		m.activeSessions.Dec()
	default:
		panic(fmt.Errorf("invalid delta: %d", delta))
	}
}

// RecordStream implements Metrics.
func (m *VictoriaMetrics) RecordStream(delta DeltaType) {
	switch delta {
	case 0:
		m.failedStreams.Inc()
	case 1:
		m.successfulStreams.Inc()
		m.activeStreams.Inc()
	case -1:
		m.activeStreams.Dec()
	default:
		panic(fmt.Errorf("invalid delta: %d", delta))
	}
}