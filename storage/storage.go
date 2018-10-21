package storage

import (
	"context"
	"fmt"
	"sync"
)

// Storage - Hot Key-Value In-Memory Storage.
type Storage struct {
	mu   sync.Mutex             // mu - mutex for thread-safe operations on map.
	data map[string]interface{} // data contains key->value pairs.
}

// New creates new storage with specific timeout.
func New() *Storage {
	return &Storage{
		data: make(map[string]interface{}),
	}
}

// Set sets key with value.
func (s *Storage) Set(ctx context.Context, key string, value interface{}) {
	s.mu.Lock()
	s.data[key] = value
	s.mu.Unlock()
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			s.mu.Lock()
			delete(s.data, key)
			s.mu.Unlock()
		}
	}(ctx)
}

// Get returns value by key and deletes key from storage.
func (s *Storage) Get(ctx context.Context, key string) (interface{}, error) {
	select {
	case <-ctx.Done():
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		value, ok := s.data[key]
		if ok {
			delete(s.data, key)
			return value, nil
		}
	}
	return nil, fmt.Errorf("key %s not found", key)
}
