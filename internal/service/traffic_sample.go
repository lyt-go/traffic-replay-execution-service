package service

import (
	"math/rand"
	"sort"
	"strings"
	"time"

	"trafficreplay/internal/model"
	"trafficreplay/pkg/idgen"
)

func (s *Service) CreateTrafficSample(input model.TrafficSample) (*model.TrafficSample, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetRecordTask(input.RecordTaskID); err != nil {
		return nil, model.NewValidationError("record_task_id", "关联的录制任务不存在")
	}
	sample := &model.TrafficSample{
		ID:           idgen.Hex(),
		RecordTaskID: input.RecordTaskID,
		Method:       input.Method,
		Path:         input.Path,
		Headers:      input.Headers,
		Body:         input.Body,
		StatusCode:   input.StatusCode,
		CapturedAt:   time.Now(),
	}
	if err := s.store.CreateTrafficSample(sample); err != nil {
		return nil, err
	}
	return sample, nil
}

func (s *Service) GetTrafficSample(id string) (*model.TrafficSample, error) {
	return s.store.GetTrafficSample(id)
}

func (s *Service) ListTrafficSamples(filter model.TrafficSampleFilter, page, size int) ([]*model.TrafficSample, int, error) {
	all := s.store.ListTrafficSamples()
	matched := make([]*model.TrafficSample, 0, len(all))
	for _, sample := range all {
		if filter.Match(sample) {
			matched = append(matched, sample)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CapturedAt.After(matched[j].CapturedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.TrafficSample{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateTrafficSample(id string, input model.TrafficSample) (*model.TrafficSample, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	sample, err := s.store.GetTrafficSample(id)
	if err != nil {
		return nil, err
	}
	sample.Method = input.Method
	sample.Path = input.Path
	sample.Headers = input.Headers
	sample.Body = input.Body
	sample.StatusCode = input.StatusCode
	if err := s.store.UpdateTrafficSample(sample); err != nil {
		return nil, err
	}
	return sample, nil
}

func (s *Service) DeleteTrafficSample(id string) error {
	return s.store.DeleteTrafficSample(id)
}

func (s *Service) BatchCreateTrafficSamples(inputs []model.TrafficSample) ([]*model.TrafficSample, error) {
	samples := make([]*model.TrafficSample, 0, len(inputs))
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, err
		}
		if _, err := s.store.GetRecordTask(input.RecordTaskID); err != nil {
			return nil, model.NewValidationError("record_task_id", "关联的录制任务不存在")
		}
		sample := &model.TrafficSample{
			ID:           idgen.Hex(),
			RecordTaskID: input.RecordTaskID,
			Method:       input.Method,
			Path:         input.Path,
			Headers:      input.Headers,
			Body:         input.Body,
			StatusCode:   input.StatusCode,
			CapturedAt:   time.Now(),
		}
		samples = append(samples, sample)
	}
	if err := s.store.BatchCreateTrafficSamples(samples); err != nil {
		return nil, err
	}
	return samples, nil
}

func (s *Service) CaptureTraffic(recordTaskID string, input model.TrafficSample) (*model.TrafficSample, error) {
	task, err := s.store.GetRecordTask(recordTaskID)
	if err != nil {
		return nil, err
	}
	if task.Status != model.RecordTaskRecording {
		return nil, model.NewValidationError("status", "录制任务未处于 recording 状态")
	}
	if task.FilterPath != "" && !strings.Contains(input.Path, task.FilterPath) {
		return nil, model.NewValidationError("path", "请求路径不匹配过滤规则")
	}
	if task.SampleRate > 0 && rand.Intn(100) >= task.SampleRate {
		return nil, nil
	}
	return s.CreateTrafficSample(input)
}
