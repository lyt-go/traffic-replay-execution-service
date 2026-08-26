package store

import (
	"trafficreplay/internal/model"
)

func (s *MemoryStore) CreateReplayConfig(c *model.ReplayConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.replayConfigs {
		if exist.Name == c.Name {
			return ErrConflict
		}
	}
	s.replayConfigs[c.ID] = c
	return nil
}

func (s *MemoryStore) GetReplayConfig(id string) (*model.ReplayConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.replayConfigs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) ListReplayConfigs() []*model.ReplayConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ReplayConfig, 0, len(s.replayConfigs))
	for _, c := range s.replayConfigs {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateReplayConfig(c *model.ReplayConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.replayConfigs[c.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.replayConfigs {
		if exist.ID != c.ID && exist.Name == c.Name {
			return ErrConflict
		}
	}
	s.replayConfigs[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteReplayConfig(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.replayConfigs[id]; !ok {
		return ErrNotFound
	}
	delete(s.replayConfigs, id)
	return nil
}
