package catalog

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

type MethodConfig struct {
	Command string            `json:"command"`
	Params  map[string]string `json:"params"`
}

type Store struct {
	mu      sync.RWMutex
	path    string
	methods map[string]MethodConfig
}

func NewStore(path string) (*Store, error) {
	s := &Store{
		path:    path,
		methods: make(map[string]MethodConfig),
	}
	if err := s.Reload(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) Reload() error {
	raw, err := os.ReadFile(s.path)
	if err != nil {
		return fmt.Errorf("read methods config: %w", err)
	}

	var m map[string]MethodConfig
	if err := json.Unmarshal(raw, &m); err != nil {
		return fmt.Errorf("parse methods config: %w", err)
	}

	norm := make(map[string]MethodConfig, len(m))
	for k, v := range m {
		norm[strings.ToUpper(strings.TrimSpace(k))] = v
	}

	s.mu.Lock()
	s.methods = norm
	s.mu.Unlock()
	return nil
}

func (s *Store) Get(name string) (MethodConfig, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cfg, ok := s.methods[strings.ToUpper(strings.TrimSpace(name))]
	return cfg, ok
}

func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.methods)
}
