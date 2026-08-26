package store

import (
	"trafficreplay/internal/model"
)

func (s *MemoryStore) CreateReplayResult(r *model.ReplayResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.replayResults[r.ID] = r.Clone()
	return nil
}

func (s *MemoryStore) GetReplayResult(id string) (*model.ReplayResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.replayResults[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r.Clone(), nil
}

func (s *MemoryStore) ListReplayResults() []*model.ReplayResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ReplayResult, 0, len(s.replayResults))
	for _, r := range s.replayResults {
		list = append(list, r.Clone())
	}
	return list
}

func (s *MemoryStore) UpdateReplayResult(r *model.ReplayResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.replayResults[r.ID]; !ok {
		return ErrNotFound
	}
	s.replayResults[r.ID] = r.Clone()
	return nil
}

func (s *MemoryStore) DeleteReplayResult(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.replayResults[id]; !ok {
		return ErrNotFound
	}
	delete(s.replayResults, id)
	return nil
}

func (s *MemoryStore) ListReplayResultsByTask(replayTaskID string) []*model.ReplayResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ReplayResult, 0)
	for _, r := range s.replayResults {
		if r.ReplayTaskID == replayTaskID {
			list = append(list, r.Clone())
		}
	}
	return list
}
