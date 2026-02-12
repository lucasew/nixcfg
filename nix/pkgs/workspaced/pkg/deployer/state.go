package deployer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// StateStore é a interface para persistência de estado
type StateStore interface {
	// Load carrega estado atual
	Load() (*State, error)

	// Save persiste estado
	Save(state *State) error

	// Path retorna caminho/identificador do store (para logging)
	Path() string
}

// FileStateStore implementa StateStore usando arquivo JSON
type FileStateStore struct {
	path string
}

// NewFileStateStore cria um FileStateStore
func NewFileStateStore(path string) (*FileStateStore, error) {
	// Expand env vars e ~
	expanded := os.ExpandEnv(path)
	if len(expanded) > 0 && expanded[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		expanded = filepath.Join(home, expanded[1:])
	}

	// Garantir que diretório pai existe
	dir := filepath.Dir(expanded)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create state directory: %w", err)
	}

	return &FileStateStore{path: expanded}, nil
}

func (s *FileStateStore) Load() (*State, error) {
	state := &State{Files: make(map[string]ManagedInfo)}

	// Se arquivo não existe, retorna estado vazio
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		return state, nil
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	if err := json.Unmarshal(data, state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	// Garantir mapa não é nil
	if state.Files == nil {
		state.Files = make(map[string]ManagedInfo)
	}

	return state, nil
}

func (s *FileStateStore) Save(state *State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

func (s *FileStateStore) Path() string {
	return s.path
}

// MemoryStateStore implementa StateStore em memória (útil para testes)
type MemoryStateStore struct {
	state *State
	id    string
}

// NewMemoryStateStore cria um MemoryStateStore
func NewMemoryStateStore(id string) *MemoryStateStore {
	return &MemoryStateStore{
		state: &State{Files: make(map[string]ManagedInfo)},
		id:    id,
	}
}

func (s *MemoryStateStore) Load() (*State, error) {
	if s.state == nil {
		s.state = &State{Files: make(map[string]ManagedInfo)}
	}
	return s.state, nil
}

func (s *MemoryStateStore) Save(state *State) error {
	s.state = state
	return nil
}

func (s *MemoryStateStore) Path() string {
	return fmt.Sprintf("memory:%s", s.id)
}
