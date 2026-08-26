package service

import (
	"trafficreplay/internal/model"
)

type OverviewStats struct {
	RecordTaskCount   int     `json:"record_task_count"`
	SampleTotal       int     `json:"sample_total"`
	ReplaySuccess     int     `json:"replay_success"`
	ReplayFailed      int     `json:"replay_failed"`
	ReplaySuccessRate float64 `json:"replay_success_rate"`
	AvgLatencyMs      float64 `json:"avg_latency_ms"`
	DiffRate          float64 `json:"diff_rate"`
}

func (s *Service) GetOverviewStats() (*OverviewStats, error) {
	recordTasks := s.store.ListRecordTasks()
	samples := s.store.ListTrafficSamples()
	results := s.store.ListReplayResults()

	var replaySuccess, replayFailed, totalLatency int
	var diffCount int
	for _, r := range results {
		if r.Matched {
			replaySuccess++
		} else {
			replayFailed++
		}
		totalLatency += r.LatencyMs
		if r.Diff != "" {
			diffCount++
		}
	}

	totalReplay := replaySuccess + replayFailed
	successRate := 0.0
	if totalReplay > 0 {
		successRate = float64(replaySuccess) / float64(totalReplay) * 100
	}
	avgLatency := 0.0
	if totalReplay > 0 {
		avgLatency = float64(totalLatency) / float64(totalReplay)
	}
	diffRate := 0.0
	if totalReplay > 0 {
		diffRate = float64(diffCount) / float64(totalReplay) * 100
	}

	return &OverviewStats{
		RecordTaskCount:   len(recordTasks),
		SampleTotal:       len(samples),
		ReplaySuccess:     replaySuccess,
		ReplayFailed:      replayFailed,
		ReplaySuccessRate: successRate,
		AvgLatencyMs:      avgLatency,
		DiffRate:          diffRate,
	}, nil
}

func (s *Service) GetRecordTaskStatusStats() map[string]int {
	all := s.store.ListRecordTasks()
	stats := make(map[string]int)
	for _, t := range all {
		stats[t.Status]++
	}
	return stats
}

func (s *Service) GetReplayTaskStatusStats() map[string]int {
	all := s.store.ListReplayTasks()
	stats := make(map[string]int)
	for _, t := range all {
		stats[t.Status]++
	}
	return stats
}

func (s *Service) GetSampleCountByTask() map[string]int {
	all := s.store.ListTrafficSamples()
	stats := make(map[string]int)
	for _, sample := range all {
		stats[sample.RecordTaskID]++
	}
	return stats
}

func (s *Service) GetResultCountByReplayTask() map[string]int {
	all := s.store.ListReplayResults()
	stats := make(map[string]int)
	for _, r := range all {
		stats[r.ReplayTaskID]++
	}
	return stats
}

func (s *Service) GetTopLatencyResults(limit int) []*model.ReplayResult {
	all := s.store.ListReplayResults()
	if len(all) == 0 {
		return []*model.ReplayResult{}
	}
	list := make([]*model.ReplayResult, len(all))
	copy(list, all)
	for i := 0; i < len(list)-1; i++ {
		for j := i + 1; j < len(list); j++ {
			if list[i].LatencyMs < list[j].LatencyMs {
				list[i], list[j] = list[j], list[i]
			}
		}
	}
	if limit > len(list) {
		limit = len(list)
	}
	return list[:limit]
}
