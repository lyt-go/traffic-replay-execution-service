package model

import (
	"strings"
	"time"
)

type ReplayResult struct {
	ID             string    `json:"id"`
	ReplayTaskID   string    `json:"replay_task_id"`
	SampleID       string    `json:"sample_id"`
	ResponseStatus int       `json:"response_status"`
	LatencyMs      int       `json:"latency_ms"`
	Matched        bool      `json:"matched"`
	Diff           string    `json:"diff"`
	ReplayedAt     time.Time `json:"replayed_at"`
}

func (r *ReplayResult) Validate() error {
	r.ReplayTaskID = strings.TrimSpace(r.ReplayTaskID)
	r.SampleID = strings.TrimSpace(r.SampleID)
	if r.ReplayTaskID == "" {
		return NewValidationError("replay_task_id", "回放任务 ID 不能为空")
	}
	if r.SampleID == "" {
		return NewValidationError("sample_id", "样本 ID 不能为空")
	}
	if r.ReplayedAt.IsZero() {
		r.ReplayedAt = time.Now()
	}
	return nil
}

type ReplayResultFilter struct {
	ReplayTaskID string
	Matched      *bool
}

func (f ReplayResultFilter) Match(r *ReplayResult) bool {
	if f.ReplayTaskID != "" && r.ReplayTaskID != f.ReplayTaskID {
		return false
	}
	if f.Matched != nil && r.Matched != *f.Matched {
		return false
	}
	return true
}
