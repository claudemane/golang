package main

import "sync"

type CachedResponse struct {
	StatusCode int
	Body       []byte
	Completed  bool
}

type MemoryStore struct {
	mu   sync.Mutex
	data map[string]*CachedResponse
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]*CachedResponse),
	}
}

func (m *MemoryStore) Get(key string) (*CachedResponse, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	resp, ok := m.data[key]
	return resp, ok
}

func (m *MemoryStore) StartProcessing(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[key]; exists {
		return false
	}

	m.data[key] = &CachedResponse{Completed: false}
	return true
}

func (m *MemoryStore) Finish(key string, status int, body []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = &CachedResponse{
		StatusCode: status,
		Body:       body,
		Completed:  true,
	}
}