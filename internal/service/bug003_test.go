package service

import (
	"testing"

	"trafficreplay/internal/config"
	"trafficreplay/internal/model"
	"trafficreplay/internal/store"
	"trafficreplay/pkg/logger"
)

func bugService003() *Service {
	return New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
}

func TestReplayResultRequiresRunningTask(t *testing.T) {
	svc := bugService003()
	record, err := svc.CreateRecordTask(model.RecordTask{Name: "capture", SourceURL: "http://source"})
	if err != nil {
		panic(err)
	}
	sample, err := svc.CreateTrafficSample(model.TrafficSample{RecordTaskID: record.ID, Method: "GET", Path: "/data"})
	if err != nil {
		panic(err)
	}
	task, err := svc.CreateReplayTask(model.ReplayTask{Name: "pending-replay", TargetURL: "http://target"})
	if err != nil {
		panic(err)
	}
	if _, err = svc.UpdateReplayTask(task.ID, model.ReplayTask{Name: task.Name, TargetURL: task.TargetURL, Status: model.ReplayTaskRunning}); err != nil {
		panic(err)
	}
	if _, err = svc.UpdateReplayTask(task.ID, model.ReplayTask{Name: task.Name, TargetURL: task.TargetURL, Status: model.ReplayTaskCompleted}); err != nil {
		panic(err)
	}
	if _, err = svc.CreateReplayResult(model.ReplayResult{ReplayTaskID: task.ID, SampleID: sample.ID}); !model.IsValidationError(err) {
		t.Fatalf("expected completed replay result to be rejected, got %v", err)
	}
}
