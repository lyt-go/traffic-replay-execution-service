package service

import (
	"sort"
	"time"

	"trafficreplay/internal/model"
	"trafficreplay/pkg/idgen"
)

func (s *Service) CreateSchedule(input model.Schedule) (*model.Schedule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetReplayConfig(input.ConfigID); err != nil {
		return nil, model.NewValidationError("config_id", "关联的回放配置不存在")
	}
	now := time.Now()
	sch := &model.Schedule{
		ID:        idgen.Hex(),
		Name:      input.Name,
		ConfigID:  input.ConfigID,
		CronExpr:  input.CronExpr,
		Status:    model.ScheduleActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateSchedule(sch); err != nil {
		return nil, err
	}
	return sch, nil
}

func (s *Service) GetSchedule(id string) (*model.Schedule, error) {
	return s.store.GetSchedule(id)
}

func (s *Service) ListSchedules(filter model.ScheduleFilter, page, size int) ([]*model.Schedule, int, error) {
	all := s.store.ListSchedules()
	matched := make([]*model.Schedule, 0, len(all))
	for _, sch := range all {
		if filter.Match(sch) {
			matched = append(matched, sch)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Schedule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateSchedule(id string, input model.Schedule) (*model.Schedule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	sch, err := s.store.GetSchedule(id)
	if err != nil {
		return nil, err
	}
	if input.Status != "" && input.Status != sch.Status {
		if !model.CanTransitionSchedule(sch.Status, input.Status) {
			return nil, model.NewValidationError("status", "调度计划状态流转不合法")
		}
		sch.Status = input.Status
	}
	sch.Name = input.Name
	sch.ConfigID = input.ConfigID
	sch.CronExpr = input.CronExpr
	sch.UpdatedAt = time.Now()
	if err := s.store.UpdateSchedule(sch); err != nil {
		return nil, err
	}
	return sch, nil
}

func (s *Service) DeleteSchedule(id string) error {
	return s.store.DeleteSchedule(id)
}

func (s *Service) RunSchedule(id string) (*model.Schedule, error) {
	sch, err := s.store.GetSchedule(id)
	if err != nil {
		return nil, err
	}
	if sch.Status != model.ScheduleActive {
		return nil, model.NewValidationError("status", "调度计划未处于 active 状态")
	}
	now := time.Now()
	sch.LastRunAt = now
	sch.NextRunAt = now.Add(time.Minute * 5)
	sch.UpdatedAt = now
	if err := s.store.UpdateSchedule(sch); err != nil {
		return nil, err
	}
	return sch, nil
}
