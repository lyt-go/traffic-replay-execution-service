package service

import (
	"testing"

	"trafficreplay/internal/config"
	"trafficreplay/internal/model"
	"trafficreplay/internal/store"
	"trafficreplay/pkg/logger"
)

func bugService007() *Service {
	return New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
}

func TestRecordRemovalClearsCapturedSamples(t *testing.T) {
	svc := bugService007()
	task, err := svc.CreateRecordTask(model.RecordTask{Name: "capture", SourceURL: "http://source"})
	if err != nil { panic(err) }
	sample, err := svc.CreateTrafficSample(model.TrafficSample{RecordTaskID: task.ID, Method: "GET", Path: "/data"})
	if err != nil { panic(err) }
	if err = svc.DeleteRecordTask(task.ID); err != nil { panic(err) }
	if _, err = svc.GetTrafficSample(sample.ID); err != store.ErrNotFound { t.Fatalf("orphan sample still available: %v", err) }
}
