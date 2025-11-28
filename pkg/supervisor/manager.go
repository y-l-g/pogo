package supervisor

import (
	"sync"
	"sync/atomic"
)

// Manager holds the state of all active worker pools.
type Manager struct {
	pools     sync.Map // map[int64]*Pool
	idCounter int64
}

// NewManager creates a new Supervisor Manager.
func NewManager() *Manager {
	return &Manager{}
}

// CreatePool creates a new Pool with a unique ID and registers it.
func (m *Manager) CreatePool() *Pool {
	id := atomic.AddInt64(&m.idCounter, 1)
	p := NewPool(id)
	m.pools.Store(id, p)
	return p
}

// RegisterPool registers a specific pool (e.g. default pool 0).
func (m *Manager) RegisterPool(p *Pool) {
	m.pools.Store(p.ID, p)
}

// GetPool retrieves a pool by its ID.
func (m *Manager) GetPool(id int64) *Pool {
	if val, ok := m.pools.Load(id); ok {
		return val.(*Pool)
	}
	return nil
}

// RemovePool shuts down and unregisters a pool.
func (m *Manager) RemovePool(id int64) {
	if val, ok := m.pools.LoadAndDelete(id); ok {
		p := val.(*Pool)
		p.Shutdown()
	}
}

// Shutdown stops all managed pools.
func (m *Manager) Shutdown() {
	m.pools.Range(func(key, value any) bool {
		p := value.(*Pool)
		p.Shutdown()
		return true
	})
}

// Range iterates over all pools.
// logic should return true to continue iteration, false to stop.
func (m *Manager) Range(f func(*Pool) bool) {
	m.pools.Range(func(key, value any) bool {
		p, ok := value.(*Pool)
		if !ok {
			return true
		}
		return f(p)
	})
}
