package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/anyshake/observer/internal/dao/action"
	"tailscale.com/ipn"
)

const (
	stateKey     = "tsnet_state"
	stateVersion = 0
)

type Database struct {
	mu            sync.RWMutex
	actionHandler *action.Handler
	namespace     string
	state         map[string][]byte
}

func New(actionHandler *action.Handler, namespace string) (*Database, error) {
	if actionHandler == nil {
		return nil, errors.New("tailscale state database handler is not available")
	}
	if namespace == "" {
		return nil, errors.New("tailscale state database namespace is not available")
	}

	if _, err := actionHandler.SettingsInit(
		namespace,
		stateKey,
		action.String,
		stateVersion,
		"{}",
	); err != nil {
		return nil, fmt.Errorf("failed to initialize tailscale state setting: %w", err)
	}

	value, valueType, _, err := actionHandler.SettingsGet(namespace, stateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load tailscale state setting: %w", err)
	}
	if valueType != action.String {
		return nil, errors.New("tailscale state setting must be a string")
	}
	encoded, ok := value.(string)
	if !ok {
		return nil, errors.New("tailscale state setting contains an invalid value")
	}

	state := make(map[string][]byte)
	if err := json.Unmarshal([]byte(encoded), &state); err != nil {
		return nil, fmt.Errorf("failed to decode tailscale state setting: %w", err)
	}

	return &Database{
		actionHandler: actionHandler,
		namespace:     namespace,
		state:         state,
	}, nil
}

func (s *Database) ReadState(id ipn.StateKey) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.state[string(id)]
	if !ok {
		return nil, ipn.ErrStateNotExist
	}

	return bytes.Clone(value), nil
}

func (s *Database) WriteState(id ipn.StateKey, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.actionHandler == nil {
		return errors.New("tailscale state database handler is not available")
	}

	current, exists := s.state[string(id)]
	if value == nil {
		if !exists {
			return nil
		}
	} else if exists && bytes.Equal(current, value) {
		return nil
	}

	nextState := make(map[string][]byte, len(s.state)+1)
	for key, stateValue := range s.state {
		nextState[key] = bytes.Clone(stateValue)
	}
	if value == nil {
		delete(nextState, string(id))
	} else {
		nextState[string(id)] = bytes.Clone(value)
	}

	encoded, err := json.Marshal(nextState)
	if err != nil {
		return fmt.Errorf("failed to encode tailscale state: %w", err)
	}
	if err := s.actionHandler.SettingsSet(
		s.namespace,
		stateKey,
		action.String,
		stateVersion,
		string(encoded),
	); err != nil {
		return fmt.Errorf("failed to persist tailscale state: %w", err)
	}

	s.state = nextState
	return nil
}

func (s *Database) String() string {
	return "databaseStateStore"
}
