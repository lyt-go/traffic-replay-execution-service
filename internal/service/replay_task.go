package service

import (
	"sort"
	"time"

	"trafficreplay/internal/model"
	"trafficreplay/pkg/idgen"
)

func (s *Service) CreateReplayTask(input model.ReplayTask) (*model.ReplayTask, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	t := &model.ReplayTask{
		ID:          idgen.Hex(),
		Name:        input.Name,
		TargetURL:   input.TargetURL,
		Concurrency: input.Concurrency,
		TimeoutMs:   input.TimeoutMs,
		SampleCount: input.SampleCount,
		Status:      model.ReplayTaskPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateReplayTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetReplayTask(id string) (*model.ReplayTask, error) {
	return s.store.GetReplayTask(id)
}

func (s *Service) ListReplayTasks(filter model.ReplayTaskFilter, page, size int) ([]*model.ReplayTask, int, error) {
	all := s.store.ListReplayTasks()
	matched := make([]*model.ReplayTask, 0, len(all))
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
		return []*model.ReplayTask{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateReplayTask(id string, input model.ReplayTask) (*model.ReplayTask, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	t, err := s.store.GetReplayTask(id)
	if err != nil {
		return nil, err
	}
	if input.Status != "" && input.Status != t.Status {
		if !model.CanTransitionReplayTask(t.Status, input.Status) {
			return nil, model.NewValidationError("status", "回放任务状态流转不合法")
		}
		t.Status = input.Status
		if input.Status == model.ReplayTaskRunning {
			t.StartedAt = time.Now()
		}
		if input.Status == model.ReplayTaskCompleted || input.Status == model.ReplayTaskFailed {
			t.EndedAt = time.Now()
		}
	}
	t.Name = input.Name
	t.TargetURL = input.TargetURL
	t.Concurrency = input.Concurrency
	t.TimeoutMs = input.TimeoutMs
	t.SampleCount = input.SampleCount
	t.UpdatedAt = time.Now()
	if err := s.store.UpdateReplayTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteReplayTask(id string) error {
	return s.store.DeleteReplayTask(id)
}
