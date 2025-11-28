package supervisor

import (
	"sort"
	"sync"
	"time"
)

// LatencyTracker keeps a rolling window of response times to calculate P95.
type LatencyTracker struct {
	mu      sync.Mutex
	samples [100]int64
	idx     int
	count   int
}

func (l *LatencyTracker) Add(ms int64) {
	l.mu.Lock()
	l.samples[l.idx] = ms
	l.idx = (l.idx + 1) % 100
	if l.count < 100 {
		l.count++
	}
	l.mu.Unlock()
}

func (l *LatencyTracker) P95() int64 {
	l.mu.Lock()
	var snapshot [100]int64
	count := l.count
	copy(snapshot[:], l.samples[:])
	l.mu.Unlock()

	if count == 0 {
		return 0
	}

	activeSlice := snapshot[:count]
	sort.Slice(activeSlice, func(i, j int) bool { return activeSlice[i] < activeSlice[j] })

	p95Index := int(float64(count) * 0.95)
	if p95Index >= count {
		p95Index = count - 1
	}
	return activeSlice[p95Index]
}

type ScalingAction int

const (
	ActionNone ScalingAction = iota
	ActionScaleUp
	ActionScaleDown
)

// AutoScaler encapsulates the decision logic for expanding or shrinking the pool.
type AutoScaler struct {
	mu           sync.Mutex
	lastSpawn    time.Time
	scaleUpVotes int
	latency      LatencyTracker

	scaleLatencyThreshold int64
}

func NewAutoScaler(scaleLatencyThreshold int64) *AutoScaler {
	return &AutoScaler{
		scaleLatencyThreshold: scaleLatencyThreshold,
		lastSpawn:             time.Now(),
	}
}

func (s *AutoScaler) RecordLatency(d time.Duration) {
	s.latency.Add(d.Milliseconds())
}

func (s *AutoScaler) P95() int64 {
	return s.latency.P95()
}

// Assess returns the recommended scaling action based on current metrics.
func (s *AutoScaler) Assess(queueDepth, idleWorkers, totalWorkers, minWorkers, maxWorkers, availableWorkers int) ScalingAction {
	s.mu.Lock()
	defer s.mu.Unlock()

	latencyP95 := s.latency.P95()

	// Scale Up Logic
	// Condition: Queue has more work than idle workers can handle immediately,
	// AND latency is high, AND we haven't hit the cap.
	if queueDepth > idleWorkers && latencyP95 > s.scaleLatencyThreshold && totalWorkers < maxWorkers {
		s.scaleUpVotes++
	} else {
		s.scaleUpVotes = 0
	}

	// Vote mechanism to debounce spikes (requires 2 consecutive votes)
	// Time check to prevent flapping
	if s.scaleUpVotes >= 2 && time.Since(s.lastSpawn) > 2*time.Second {
		s.lastSpawn = time.Now()
		s.scaleUpVotes = 0
		return ActionScaleUp
	}

	// Scale Down Logic
	// Condition: Empty queue, more workers than min, and we have excess idle capacity.
	// We also check availableWorkers (channel length) to ensure we actually have someone to kill.
	if queueDepth == 0 && totalWorkers > minWorkers && idleWorkers > 1 && availableWorkers > minWorkers {
		return ActionScaleDown
	}

	return ActionNone
}
