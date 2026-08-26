package store

import (
	"trafficreplay/internal/model"
)

func (s *MemoryStore) CreateReplayResult(r *model.ReplayResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.replayResults[r.ID] = r
	return nil
}

// GetReplayResult 返回存储记录的副本（快照），调用方对返回值的修改不会污染存储数据。
func (s *MemoryStore) GetReplayResult(id string) (*model.ReplayResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.replayResults[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *r
	return &cp, nil
}

// ListReplayResults 返回所有记录的副本（快照），与存储数据隔离。
func (s *MemoryStore) ListReplayResults() []*model.ReplayResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ReplayResult, 0, len(s.replayResults))
	for _, r := range s.replayResults {
		cp := *r
		list = append(list, &cp)
	}
	return list
}

func (s *MemoryStore) UpdateReplayResult(r *model.ReplayResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.replayResults[r.ID]; !ok {
		return ErrNotFound
	}
	s.replayResults[r.ID] = r
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

// ListReplayResultsByTask 返回匹配记录的副本（快照），与存储数据隔离。
func (s *MemoryStore) ListReplayResultsByTask(replayTaskID string) []*model.ReplayResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ReplayResult, 0)
	for _, r := range s.replayResults {
		if r.ReplayTaskID == replayTaskID {
			cp := *r
			list = append(list, &cp)
		}
	}
	return list
}
