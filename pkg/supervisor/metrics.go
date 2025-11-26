package supervisor

import (
	"log"
	"net/http"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Namespace for all metrics
	namespace = "pogo"

	// Singleton enforcement
	metricsInitialized atomic.Bool

	// Metric Definitions
	workersActive = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "workers_active",
		Help:      "Number of workers currently processing a job",
	}, []string{"pool_id"})

	workersTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "workers_total",
		Help:      "Total number of worker processes (Active + Idle)",
	}, []string{"pool_id"})

	queueDepth = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "ipc_queue_depth",
		Help:      "Number of pending tasks in the Go channel",
	}, []string{"pool_id"})

	goGoroutines = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "go_goroutines",
		Help:      "Number of Go routines",
	})

	goHeapBytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "go_heap_bytes",
		Help:      "Bytes of allocated heap objects",
	})
)

// InitMetrics registers the metrics and starts the HTTP server.
// It is idempotent.
func InitMetrics(addr string) {
	if metricsInitialized.Swap(true) {
		return
	}

	// Register metrics
	prometheus.MustRegister(workersActive)
	prometheus.MustRegister(workersTotal)
	prometheus.MustRegister(queueDepth)
	prometheus.MustRegister(goGoroutines)
	prometheus.MustRegister(goHeapBytes)

	// Start Runtime Stats Collector
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for range ticker.C {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			goHeapBytes.Set(float64(m.Alloc))
			goGoroutines.Set(float64(runtime.NumGoroutine()))
		}
	}()

	// Start HTTP Server
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		log.Printf("[Pogo] Metrics server listening on %s/metrics", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Printf("[Pogo] Metrics server failed: %v", err)
		}
	}()
}

// Helper functions for Pool to update metrics

func MetricWorkerSpawn(poolID int64) {
	workersTotal.WithLabelValues(strconv.FormatInt(poolID, 10)).Inc()
}

func MetricWorkerKill(poolID int64) {
	workersTotal.WithLabelValues(strconv.FormatInt(poolID, 10)).Dec()
}

func MetricWorkerBusy(poolID int64) {
	workersActive.WithLabelValues(strconv.FormatInt(poolID, 10)).Inc()
}

func MetricWorkerIdle(poolID int64) {
	workersActive.WithLabelValues(strconv.FormatInt(poolID, 10)).Dec()
}

func MetricQueueDepth(poolID int64, depth int) {
	queueDepth.WithLabelValues(strconv.FormatInt(poolID, 10)).Set(float64(depth))
}
