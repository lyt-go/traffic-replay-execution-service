package model

import (
	"strings"
	"time"
)

type ReplayConfig struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	TargetHost       string    `json:"target_host"`
	TimeoutMs        int       `json:"timeout_ms"`
	Retries          int       `json:"retries"`
	Enabled          bool      `json:"enabled"`
	DisabledByUpdate bool      `json:"disabled_by_update"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// IsScheduleEligible 判断挂在配置上的调度计划是否可执行。
// 仅当配置在创建后通过更新操作被停用（Enabled 由 true 变为 false，
// 即 DisabledByUpdate=true）时，调度计划不可执行；创建时即为停用
// 且从未更新过的旧配置仍可执行。
func (c *ReplayConfig) IsScheduleEligible() bool {
	return !c.DisabledByUpdate
}

func (c *ReplayConfig) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	c.TargetHost = strings.TrimSpace(c.TargetHost)
	if c.Name == "" {
		return NewValidationError("name", "配置名称不能为空")
	}
	if c.TargetHost == "" {
		return NewValidationError("target_host", "目标主机不能为空")
	}
	if c.TimeoutMs < 1 {
		c.TimeoutMs = 5000
	}
	if c.Retries < 0 {
		c.Retries = 0
	}
	return nil
}

type ReplayConfigFilter struct {
	Enabled *bool
	Keyword string
}

func (f ReplayConfigFilter) Match(c *ReplayConfig) bool {
	if f.Enabled != nil && c.Enabled != *f.Enabled {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.Name), k) && !strings.Contains(strings.ToLower(c.TargetHost), k) {
			return false
		}
	}
	return true
}
