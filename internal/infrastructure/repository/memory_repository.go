package repository

import (
	"errors"
	"sync"
	"urlChecker/internal/domain/monitor"
)

type MemoryRepository struct {
	mu      sync.RWMutex
	storage map[string]*monitor.URLMonitor
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		storage: make(map[string]*monitor.URLMonitor),
	}
}

func (r *MemoryRepository) Save(m *monitor.URLMonitor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.storage[m.ID] = m
	return nil
}

func (r *MemoryRepository) FindByID(id string) (*monitor.URLMonitor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, exists := r.storage[id]
	if !exists {
		return nil, errors.New("monitor not found")
	}
	return m, nil
}

func (r *MemoryRepository) FindAll() ([]*monitor.URLMonitor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*monitor.URLMonitor, 0, len(r.storage))
	for _, m := range r.storage {
		result = append(result, m)
	}
	return result, nil
}

func (r *MemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.storage, id)
	return nil
}

func (r *MemoryRepository) Update(m *monitor.URLMonitor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.storage[m.ID] = m
	return nil
}
