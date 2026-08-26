package service

import (
	"sort"
	"time"

	"trafficreplay/internal/model"
	"trafficreplay/pkg/idgen"
)

func (s *Service) CreateReplayConfig(input model.ReplayConfig) (*model.ReplayConfig, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	c := &model.ReplayConfig{
		ID:         idgen.Hex(),
		Name:       input.Name,
		TargetHost: input.TargetHost,
		TimeoutMs:  input.TimeoutMs,
		Retries:    input.Retries,
		Enabled:    input.Enabled,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateReplayConfig(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetReplayConfig(id string) (*model.ReplayConfig, error) {
	return s.store.GetReplayConfig(id)
}

func (s *Service) ListReplayConfigs(filter model.ReplayConfigFilter, page, size int) ([]*model.ReplayConfig, int, error) {
	all := s.store.ListReplayConfigs()
	matched := make([]*model.ReplayConfig, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ReplayConfig{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateReplayConfig(id string, input model.ReplayConfig) (*model.ReplayConfig, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	c, err := s.store.GetReplayConfig(id)
	if err != nil {
		return nil, err
	}
	c.Name = input.Name
	c.TargetHost = input.TargetHost
	c.TimeoutMs = input.TimeoutMs
	c.Retries = input.Retries
	c.Enabled = input.Enabled
	c.UpdatedAt = time.Now()
	if err := s.store.UpdateReplayConfig(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteReplayConfig(id string) error {
	return s.store.DeleteReplayConfig(id)
}
