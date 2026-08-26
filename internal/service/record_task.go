package service

import (
	"sort"
	"time"

	"trafficreplay/internal/model"
	"trafficreplay/pkg/idgen"
)

func (s *Service) CreateRecordTask(input model.RecordTask) (*model.RecordTask, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	t := &model.RecordTask{
		ID:         idgen.Hex(),
		Name:       input.Name,
		SourceURL:  input.SourceURL,
		FilterPath: input.FilterPath,
		SampleRate: input.SampleRate,
		Status:     model.RecordTaskIdle,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateRecordTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetRecordTask(id string) (*model.RecordTask, error) {
	return s.store.GetRecordTask(id)
}

func (s *Service) ListRecordTasks(filter model.RecordTaskFilter, page, size int) ([]*model.RecordTask, int, error) {
	all := s.store.ListRecordTasks()
	matched := make([]*model.RecordTask, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RecordTask{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRecordTask(id string, input model.RecordTask) (*model.RecordTask, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	t, err := s.store.GetRecordTask(id)
	if err != nil {
		return nil, err
	}
	if input.Status != "" && input.Status != t.Status {
		if !model.CanTransitionRecordTask(t.Status, input.Status) {
			return nil, model.NewValidationError("status", "录制任务状态流转不合法")
		}
		t.Status = input.Status
		if input.Status == model.RecordTaskRecording {
			t.StartedAt = time.Now()
		}
		if input.Status == model.RecordTaskCompleted || input.Status == model.RecordTaskPaused {
			t.EndedAt = time.Now()
		}
	}
	t.Name = input.Name
	t.SourceURL = input.SourceURL
	t.FilterPath = input.FilterPath
	t.SampleRate = input.SampleRate
	t.UpdatedAt = time.Now()
	if err := s.store.UpdateRecordTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteRecordTask(id string) error {
	return s.store.DeleteRecordTask(id)
}
