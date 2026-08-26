package store

import (
	"trafficreplay/internal/model"
)

func (s *MemoryStore) CreateTrafficSample(sample *model.TrafficSample) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.trafficSamples[sample.ID] = sample
	return nil
}

func (s *MemoryStore) GetTrafficSample(id string) (*model.TrafficSample, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sample, ok := s.trafficSamples[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sample, nil
}

func (s *MemoryStore) ListTrafficSamples() []*model.TrafficSample {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.TrafficSample, 0, len(s.trafficSamples))
	for _, sample := range s.trafficSamples {
		list = append(list, sample)
	}
	return list
}

func (s *MemoryStore) UpdateTrafficSample(sample *model.TrafficSample) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.trafficSamples[sample.ID]; !ok {
		return ErrNotFound
	}
	s.trafficSamples[sample.ID] = sample
	return nil
}

func (s *MemoryStore) DeleteTrafficSample(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.trafficSamples[id]; !ok {
		return ErrNotFound
	}
	delete(s.trafficSamples, id)
	// 级联清理：删除依赖该样本的全部回放结果，避免孤儿数据。
	s.deleteReplayResultsBySampleLocked(id)
	return nil
}

func (s *MemoryStore) BatchCreateTrafficSamples(samples []*model.TrafficSample) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sample := range samples {
		s.trafficSamples[sample.ID] = sample
	}
	return nil
}

func (s *MemoryStore) ListTrafficSamplesByTask(recordTaskID string) []*model.TrafficSample {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.TrafficSample, 0)
	for _, sample := range s.trafficSamples {
		if sample.RecordTaskID == recordTaskID {
			list = append(list, sample)
		}
	}
	return list
}
