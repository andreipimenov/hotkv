package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Storage - Hot Key-Value In-Memory Storage.
type Storage struct {
	mu          sync.Mutex                    // mu - mutex for thread-safe operations on maps.
	data        map[string]interface{}        // data contains key->value pairs.
	cancelFuncs map[string]context.CancelFunc // cancelFuncs contains cancel functions for keys contexts.
	timeout     time.Duration                 // timeout - time to live for keys.
}

// New creates new storage with specific timeout.
func New(timeout time.Duration) (*Storage, error) {
	if timeout <= 0 {
		return nil, errors.New("timeout must be positive")
	}
	return &Storage{
		data:        make(map[string]interface{}),
		cancelFuncs: make(map[string]context.CancelFunc),
		timeout:     timeout,
	}, nil
}

// Set sets key with value.
func (s *Storage) Set(key string, value interface{}) {
	s.mu.Lock()
	s.data[key] = value
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	s.cancelFuncs[key] = cancel
	s.mu.Unlock()
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			s.mu.Lock()
			delete(s.data, key)
			delete(s.cancelFuncs, key)
			s.mu.Unlock()
		}
	}(ctx)
}

// Get returns value by key and deletes key from storage.
func (s *Storage) Get(key string) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value := s.data[key]
	cancel, ok := s.cancelFuncs[key]
	if ok {
		cancel()
		return value, nil
	}
	return nil, fmt.Errorf("key %s not found", key)
}
