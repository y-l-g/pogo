package supervisor

// PoolObserver defines the interface for monitoring pool lifecycle events.
// Implementations must be thread-safe as methods are called from multiple goroutines.
type PoolObserver interface {
	OnWorkerStart(workerID int)
	OnWorkerExit(workerID int)
}

// NoOpObserver is a default implementation that does nothing.
// It is used when no observer is provided in the configuration.
type NoOpObserver struct{}

func (n *NoOpObserver) OnWorkerStart(workerID int) {}
func (n *NoOpObserver) OnWorkerExit(workerID int)  {}
