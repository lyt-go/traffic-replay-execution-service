package model

import (
	"strings"
	"time"
)

type TrafficSample struct {
	ID           string    `json:"id"`
	RecordTaskID string    `json:"record_task_id"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	Headers      string    `json:"headers"`
	Body         string    `json:"body"`
	StatusCode   int       `json:"status_code"`
	CapturedAt   time.Time `json:"captured_at"`
}

// Clone 返回 TrafficSample 的深拷贝副本。
// 读取层返回快照，使调用方对返回值的修改不会污染已持久化的数据。
func (s *TrafficSample) Clone() *TrafficSample {
	if s == nil {
		return nil
	}
	cp := *s
	return &cp
}

func (s *TrafficSample) Validate() error {
	s.RecordTaskID = strings.TrimSpace(s.RecordTaskID)
	s.Method = strings.TrimSpace(s.Method)
	s.Path = strings.TrimSpace(s.Path)
	if s.RecordTaskID == "" {
		return NewValidationError("record_task_id", "录制任务 ID 不能为空")
	}
	if s.Method == "" {
		return NewValidationError("method", "请求方法不能为空")
	}
	if s.Path == "" {
		return NewValidationError("path", "请求路径不能为空")
	}
	if s.CapturedAt.IsZero() {
		s.CapturedAt = time.Now()
	}
	return nil
}

type TrafficSampleFilter struct {
	RecordTaskID string
	Method       string
	Path         string
	StatusCode   int
}

func (f TrafficSampleFilter) Match(s *TrafficSample) bool {
	if f.RecordTaskID != "" && s.RecordTaskID != f.RecordTaskID {
		return false
	}
	if f.Method != "" && s.Method != f.Method {
		return false
	}
	if f.Path != "" && !strings.Contains(s.Path, f.Path) {
		return false
	}
	if f.StatusCode > 0 && s.StatusCode != f.StatusCode {
		return false
	}
	return true
}
