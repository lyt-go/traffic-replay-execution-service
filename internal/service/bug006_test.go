package service

import (
	"testing"

	"trafficreplay/internal/config"
	"trafficreplay/internal/model"
	"trafficreplay/internal/store"
	"trafficreplay/pkg/logger"
)

func bugService006() *Service {
	return New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
}

func TestDeletingReplayTaskRemovesResults(t *testing.T) {
	svc := bugService006()
	record, err := svc.CreateRecordTask(model.RecordTask{Name: "capture", SourceURL: "http://source"})
	if err != nil { panic(err) }
	sample, err := svc.CreateTrafficSample(model.TrafficSample{RecordTaskID: record.ID, Method: "GET", Path: "/data"})
	if err != nil { panic(err) }
	task, err := svc.CreateReplayTask(model.ReplayTask{Name: "replay", TargetURL: "http://target"})
	if err != nil { panic(err) }
	result, err := svc.CreateReplayResult(model.ReplayResult{ReplayTaskID: task.ID, SampleID: sample.ID})
	if err != nil { panic(err) }
	if err = svc.DeleteReplayTask(task.ID); err != nil { panic(err) }
	if _, err = svc.GetReplayResult(result.ID); err != store.ErrNotFound { t.Fatalf("orphan result still available: %v", err) }
}
