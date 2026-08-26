package model

import (
	"strings"
	"time"
)

const (
	ReplayTaskPending   = "pending"
	ReplayTaskRunning   = "running"
	ReplayTaskCompleted = "completed"
	ReplayTaskFailed    = "failed"
)

var replayTaskTransitions = map[string]map[string]bool{
	ReplayTaskPending:   {ReplayTaskRunning: true},
	ReplayTaskRunning:   {ReplayTaskCompleted: true, ReplayTaskFailed: true},
	ReplayTaskCompleted: {},
	ReplayTaskFailed:    {},
}

func CanTransitionReplayTask(from, to string) bool {
	if m, ok := replayTaskTransitions[from]; ok {
		return m[to]
	}
	return false
}

type ReplayTask struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	TargetURL    string    `json:"target_url"`
	Concurrency  int       `json:"concurrency"`
	TimeoutMs    int       `json:"timeout_ms"`
	SampleCount  int       `json:"sample_count"`
	Status       string    `json:"status"`
	StartedAt    time.Time `json:"started_at"`
	EndedAt      time.Time `json:"ended_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (t *ReplayTask) Validate() error {
	t.Name = strings.TrimSpace(t.Name)
	t.TargetURL = strings.TrimSpace(t.TargetURL)
	if t.Name == "" {
		return NewValidationError("name", "回放任务名称不能为空")
	}
	if t.TargetURL == "" {
		return NewValidationError("target_url", "目标地址不能为空")
	}
	if t.Concurrency < 1 {
		t.Concurrency = 1
	}
	if t.TimeoutMs < 1 {
		t.TimeoutMs = 5000
	}
	if t.Status == "" {
		t.Status = ReplayTaskPending
	}
	if t.Status != ReplayTaskPending && t.Status != ReplayTaskRunning && t.Status != ReplayTaskCompleted && t.Status != ReplayTaskFailed {
		return NewValidationError("status", "回放任务状态不合法")
	}
	return nil
}

type ReplayTaskFilter struct {
	Status  string
	Keyword string
}

func (f ReplayTaskFilter) Match(t *ReplayTask) bool {
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.Name), k) && !strings.Contains(strings.ToLower(t.TargetURL), k) {
			return false
		}
	}
	return true
}
