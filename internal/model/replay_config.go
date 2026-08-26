package model

import (
	"strings"
	"time"
)

type ReplayConfig struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	TargetHost string    `json:"target_host"`
	TimeoutMs  int       `json:"timeout_ms"`
	Retries    int       `json:"retries"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Clone 返回 ReplayConfig 的深拷贝副本。
// 读取层返回快照，使调用方对返回值的修改不会污染已持久化的数据。
func (c *ReplayConfig) Clone() *ReplayConfig {
	if c == nil {
		return nil
	}
	cp := *c
	return &cp
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
