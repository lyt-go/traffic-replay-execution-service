package model

import (
	"strings"
	"time"
)

const (
	RecordTaskIdle      = "idle"
	RecordTaskRecording = "recording"
	RecordTaskPaused    = "paused"
	RecordTaskCompleted = "completed"
)

var recordTaskTransitions = map[string]map[string]bool{
	RecordTaskIdle:      {RecordTaskRecording: true},
	RecordTaskRecording: {RecordTaskPaused: true, RecordTaskCompleted: true},
	RecordTaskPaused:    {RecordTaskRecording: true, RecordTaskCompleted: true},
	RecordTaskCompleted: {},
}

func CanTransitionRecordTask(from, to string) bool {
	if m, ok := recordTaskTransitions[from]; ok {
		return m[to]
	}
	return false
}

type RecordTask struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	SourceURL  string    `json:"source_url"`
	FilterPath string    `json:"filter_path"`
	SampleRate int       `json:"sample_rate"`
	Status     string    `json:"status"`
	StartedAt  time.Time `json:"started_at"`
	EndedAt    time.Time `json:"ended_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (t *RecordTask) Validate() error {
	t.Name = strings.TrimSpace(t.Name)
	t.SourceURL = strings.TrimSpace(t.SourceURL)
	t.FilterPath = strings.TrimSpace(t.FilterPath)
	if t.Name == "" {
		return NewValidationError("name", "录制任务名称不能为空")
	}
	if t.SourceURL == "" {
		return NewValidationError("source_url", "来源地址不能为空")
	}
	if t.SampleRate < 0 || t.SampleRate > 100 {
		return NewValidationError("sample_rate", "采样率必须在 0-100 之间")
	}
	if t.Status == "" {
		t.Status = RecordTaskIdle
	}
	if t.Status != RecordTaskIdle && t.Status != RecordTaskRecording && t.Status != RecordTaskPaused && t.Status != RecordTaskCompleted {
		return NewValidationError("status", "录制任务状态不合法")
	}
	return nil
}

type RecordTaskFilter struct {
	Status  string
	Keyword string
}

func (f RecordTaskFilter) Match(t *RecordTask) bool {
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.Name), k) && !strings.Contains(strings.ToLower(t.SourceURL), k) {
			return false
		}
	}
	return true
}
