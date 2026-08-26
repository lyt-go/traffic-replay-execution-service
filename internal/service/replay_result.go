package service

import (
	"sort"
	"time"

	"trafficreplay/internal/model"
	"trafficreplay/pkg/idgen"
)

func (s *Service) CreateReplayResult(input model.ReplayResult) (*model.ReplayResult, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetReplayTask(input.ReplayTaskID); err != nil {
		return nil, model.NewValidationError("replay_task_id", "关联的回放任务不存在")
	}
	if _, err := s.store.GetTrafficSample(input.SampleID); err != nil {
		return nil, model.NewValidationError("sample_id", "关联的样本不存在")
	}
	result := &model.ReplayResult{
		ID:             idgen.Hex(),
		ReplayTaskID:   input.ReplayTaskID,
		SampleID:       input.SampleID,
		ResponseStatus: input.ResponseStatus,
		LatencyMs:      input.LatencyMs,
		Matched:        input.Matched,
		Diff:           input.Diff,
		ReplayedAt:     time.Now(),
	}
	if err := s.store.CreateReplayResult(result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) GetReplayResult(id string) (*model.ReplayResult, error) {
	return s.store.GetReplayResult(id)
}

func (s *Service) ListReplayResults(filter model.ReplayResultFilter, page, size int) ([]*model.ReplayResult, int, error) {
	all := s.store.ListReplayResults()
	matched := make([]*model.ReplayResult, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].ReplayedAt.After(matched[j].ReplayedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ReplayResult{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateReplayResult(id string, input model.ReplayResult) (*model.ReplayResult, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	r, err := s.store.GetReplayResult(id)
	if err != nil {
		return nil, err
	}
	r.ResponseStatus = input.ResponseStatus
	r.LatencyMs = input.LatencyMs
	r.Matched = input.Matched
	r.Diff = input.Diff
	if err := s.store.UpdateReplayResult(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteReplayResult(id string) error {
	return s.store.DeleteReplayResult(id)
}

func (s *Service) ExecuteReplay(replayTaskID string) ([]*model.ReplayResult, error) {
	task, err := s.store.GetReplayTask(replayTaskID)
	if err != nil {
		return nil, err
	}
	if task.Status != model.ReplayTaskRunning {
		return nil, model.NewValidationError("status", "回放任务未处于 running 状态")
	}
	// 回放任务只处理与之关联的录制样本；未绑定录制任务的旧任务仍处理全部样本。
	var samples []*model.TrafficSample
	if task.RecordTaskID != "" {
		samples = s.store.ListTrafficSamplesByTask(task.RecordTaskID)
	} else {
		samples = s.store.ListTrafficSamples()
	}
	if task.SampleCount > 0 && len(samples) > task.SampleCount {
		samples = samples[:task.SampleCount]
	}
	results := make([]*model.ReplayResult, 0, len(samples))
	for _, sample := range samples {
		latency := 10 + int(time.Now().UnixNano()%200)
		matched := sample.StatusCode == 200
		diff := ""
		if !matched {
			diff = "status_code mismatch"
		}
		result := &model.ReplayResult{
			ID:             idgen.Hex(),
			ReplayTaskID:   replayTaskID,
			SampleID:       sample.ID,
			ResponseStatus: sample.StatusCode,
			LatencyMs:      latency,
			Matched:        matched,
			Diff:           diff,
			ReplayedAt:     time.Now(),
		}
		if err := s.store.CreateReplayResult(result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
