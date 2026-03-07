package platform

import (
	"fmt"
	"sync"
)

var (
	registry = make(map[Type]Factory)
	mu       sync.RWMutex
)

// Register registers a platform factory
func Register(factory Factory) {
	mu.Lock()
	defer mu.Unlock()

	registry[factory.Platform()] = factory
}

// Get returns the factory for the specified platform
func Get(platformType Type) (Factory, error) {
	mu.RLock()
	defer mu.RUnlock()

	factory, ok := registry[platformType]
	if !ok {
		return nil, fmt.Errorf("unsupported platform: %s", platformType)
	}

	return factory, nil
}

// GetClient is a convenience method to create a client for the specified platform
func GetClient(platformType Type) (Client, error) {
	factory, err := Get(platformType)
	if err != nil {
		return nil, err
	}

	return factory.CreateClient()
}
