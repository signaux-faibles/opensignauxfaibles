package marshal

import (
	"errors"
	"sync"
)

type Cache interface {
	Get(name string) (any, error)
	Set(name string, value any)
}

// cache stores values in memory, for future retrieval
// It is thread-safe
type cache struct {
	data map[string]any
	mu   sync.RWMutex
}

// Get gets a value from the cache
func (ca *cache) Get(name string) (any, error) {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	if ca == nil {
		return nil, errors.New("Entry not found: " + name)
	}
	if _, ok := ca.data[name]; !ok {
		return nil, errors.New("Entry not found: " + name)
	}
	return ca.data[name], nil
}

// Set writes a value to the Cache
func (ca *cache) Set(name string, value any) {
	ca.mu.Lock()
	defer ca.mu.Unlock()

	ca.data[name] = value
}

// NewCache returns a new cache object
func NewCache(data map[string]any) Cache {
	return &cache{data: data}
}

func NewEmptyCache() Cache {
	return &cache{data: map[string]any{}}
}
