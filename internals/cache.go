package internals

import (
	"fmt"
)

// NewMemoryCache creates a new MemoryCache instance
func NewMemoryCache[T any]() *MemoryCache[T] {
	return &MemoryCache[T]{
		data: make(map[string]T),
	}
}

func (c *MemoryCache[T]) Save(key string, game T) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = game
	return nil
}

func (c *MemoryCache[T]) Load(key string) (T, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	game, exists := c.data[key]

	if !exists {
		var zero T
		return zero, fmt.Errorf("game not found")
	}

	return game, nil
}

func (c *MemoryCache[T]) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}
