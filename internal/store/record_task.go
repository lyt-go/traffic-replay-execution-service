package store

import (
	"trafficreplay/internal/model"
)

func (s *MemoryStore) CreateRecordTask(t *model.RecordTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recordTasks[t.ID] = t
	return nil
}

func (s *MemoryStore) GetRecordTask(id string) (*model.RecordTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.recordTasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) ListRecordTasks() []*model.RecordTask {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RecordTask, 0, len(s.recordTasks))
	for _, t := range s.recordTasks {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateRecordTask(t *model.RecordTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.recordTasks[t.ID]; !ok {
		return ErrNotFound
	}
	s.recordTasks[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteRecordTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.recordTasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.recordTasks, id)
	return nil
}
