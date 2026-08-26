package model

import (
	"strings"
	"time"
)

const (
	ScheduleActive = "active"
	SchedulePaused = "paused"
)

var scheduleTransitions = map[string]map[string]bool{
	ScheduleActive: {SchedulePaused: true},
	SchedulePaused: {ScheduleActive: true},
}

func CanTransitionSchedule(from, to string) bool {
	if m, ok := scheduleTransitions[from]; ok {
		return m[to]
	}
	return false
}

type Schedule struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ConfigID  string    `json:"config_id"`
	CronExpr  string    `json:"cron_expr"`
	Status    string    `json:"status"`
	LastRunAt time.Time `json:"last_run_at"`
	NextRunAt time.Time `json:"next_run_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Clone 返回 Schedule 的深拷贝副本。
// 读取层返回快照，使调用方对返回值的修改不会污染已持久化的数据。
func (s *Schedule) Clone() *Schedule {
	if s == nil {
		return nil
	}
	cp := *s
	return &cp
}

func (s *Schedule) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.ConfigID = strings.TrimSpace(s.ConfigID)
	s.CronExpr = strings.TrimSpace(s.CronExpr)
	if s.Name == "" {
		return NewValidationError("name", "调度计划名称不能为空")
	}
	if s.ConfigID == "" {
		return NewValidationError("config_id", "配置 ID 不能为空")
	}
	if s.CronExpr == "" {
		return NewValidationError("cron_expr", "Cron 表达式不能为空")
	}
	if s.Status == "" {
		s.Status = ScheduleActive
	}
	if s.Status != ScheduleActive && s.Status != SchedulePaused {
		return NewValidationError("status", "调度计划状态不合法")
	}
	return nil
}

type ScheduleFilter struct {
	Status  string
	Keyword string
}

func (f ScheduleFilter) Match(s *Schedule) bool {
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Name), k) && !strings.Contains(strings.ToLower(s.CronExpr), k) {
			return false
		}
	}
	return true
}
