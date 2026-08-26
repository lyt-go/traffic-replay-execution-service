package store

import (
	"sync"

	"trafficreplay/internal/model"
)

type MemoryStore struct {
	mu            sync.RWMutex
	recordTasks   map[string]*model.RecordTask
	trafficSamples map[string]*model.TrafficSample
	replayTasks   map[string]*model.ReplayTask
	replayResults map[string]*model.ReplayResult
	replayConfigs map[string]*model.ReplayConfig
	schedules     map[string]*model.Schedule
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		recordTasks:    make(map[string]*model.RecordTask),
		trafficSamples: make(map[string]*model.TrafficSample),
		replayTasks:    make(map[string]*model.ReplayTask),
		replayResults:  make(map[string]*model.ReplayResult),
		replayConfigs:  make(map[string]*model.ReplayConfig),
		schedules:      make(map[string]*model.Schedule),
	}
}

var _ Store = (*MemoryStore)(nil)
