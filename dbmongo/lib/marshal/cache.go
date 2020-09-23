package marshal

import "errors"

// Cache saves values in memory
type Cache map[string]interface{}

// Get gets a value from the cache
func (ca Cache) Get(name string) (interface{}, error) {
	if ca == nil {
		return nil, errors.New("Entry not found: " + name)
	}
	if _, ok := ca[name]; !ok {
		return nil, errors.New("Entry not found: " + name)
	}
	return ca[name], nil
}

// Set writes a value to the Cache
func (ca Cache) Set(name string, value interface{}) {
	ca[name] = value
}

// NewCache returns a new cache object
func NewCache() Cache {
	return make(map[string]interface{})
}
