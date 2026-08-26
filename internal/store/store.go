// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"trafficreplay/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// RecordTask
	CreateRecordTask(t *model.RecordTask) error
	GetRecordTask(id string) (*model.RecordTask, error)
	ListRecordTasks() []*model.RecordTask
	UpdateRecordTask(t *model.RecordTask) error
	DeleteRecordTask(id string) error

	// TrafficSample
	CreateTrafficSample(s *model.TrafficSample) error
	GetTrafficSample(id string) (*model.TrafficSample, error)
	ListTrafficSamples() []*model.TrafficSample
	UpdateTrafficSample(s *model.TrafficSample) error
	DeleteTrafficSample(id string) error
	BatchCreateTrafficSamples(samples []*model.TrafficSample) error
	ListTrafficSamplesByTask(recordTaskID string) []*model.TrafficSample

	// ReplayTask
	CreateReplayTask(t *model.ReplayTask) error
	GetReplayTask(id string) (*model.ReplayTask, error)
	ListReplayTasks() []*model.ReplayTask
	UpdateReplayTask(t *model.ReplayTask) error
	DeleteReplayTask(id string) error

	// ReplayResult
	CreateReplayResult(r *model.ReplayResult) error
	GetReplayResult(id string) (*model.ReplayResult, error)
	ListReplayResults() []*model.ReplayResult
	UpdateReplayResult(r *model.ReplayResult) error
	DeleteReplayResult(id string) error
	ListReplayResultsByTask(replayTaskID string) []*model.ReplayResult

	// ReplayConfig
	CreateReplayConfig(c *model.ReplayConfig) error
	GetReplayConfig(id string) (*model.ReplayConfig, error)
	ListReplayConfigs() []*model.ReplayConfig
	UpdateReplayConfig(c *model.ReplayConfig) error
	DeleteReplayConfig(id string) error

	// Schedule
	CreateSchedule(s *model.Schedule) error
	GetSchedule(id string) (*model.Schedule, error)
	ListSchedules() []*model.Schedule
	UpdateSchedule(s *model.Schedule) error
	DeleteSchedule(id string) error
}
