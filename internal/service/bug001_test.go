package service

import (
	"testing"

	"trafficreplay/internal/config"
	"trafficreplay/internal/model"
	"trafficreplay/internal/store"
	"trafficreplay/pkg/logger"
)

func bugService001() *Service {
	return New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
}

func TestReplayExecutionCompletesTask(t *testing.T) {
	svc := bugService001()
	record, err := svc.CreateRecordTask(model.RecordTask{Name: "capture", SourceURL: "http://source"})
	if err != nil { panic(err) }
	if _, err = svc.CreateTrafficSample(model.TrafficSample{RecordTaskID: record.ID, Method: "GET", Path: "/health", StatusCode: 200}); err != nil { panic(err) }
	task, err := svc.CreateReplayTask(model.ReplayTask{Name: "replay", TargetURL: "http://target"})
	if err != nil { panic(err) }
	if _, err = svc.UpdateReplayTask(task.ID, model.ReplayTask{Name: task.Name, TargetURL: task.TargetURL, Status: model.ReplayTaskRunning}); err != nil { panic(err) }
	if _, err = svc.ExecuteReplay(task.ID); err != nil { panic(err) }
	got, err := svc.GetReplayTask(task.ID)
	if err != nil { panic(err) }
	if got.Status != model.ReplayTaskCompleted || got.EndedAt.IsZero() {
		t.Fatalf("replay task remained %q with end time %v", got.Status, got.EndedAt)
	}
}
