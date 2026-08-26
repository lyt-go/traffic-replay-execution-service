package store

import (
	"trafficreplay/internal/model"
)

func (s *MemoryStore) CreateReplayTask(t *model.ReplayTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.replayTasks[t.ID] = t
	return nil
}

func (s *MemoryStore) GetReplayTask(id string) (*model.ReplayTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.replayTasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) ListReplayTasks() []*model.ReplayTask {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ReplayTask, 0, len(s.replayTasks))
	for _, t := range s.replayTasks {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateReplayTask(t *model.ReplayTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.replayTasks[t.ID]; !ok {
		return ErrNotFound
	}
	s.replayTasks[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteReplayTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.replayTasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.replayTasks, id)
	return nil
}
