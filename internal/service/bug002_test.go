package service

import (
	"testing"

	"trafficreplay/internal/config"
	"trafficreplay/internal/model"
	"trafficreplay/internal/store"
	"trafficreplay/pkg/logger"
)

func bugService002() *Service {
	return New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
}

func TestReplayUsesOnlyLinkedSamples(t *testing.T) {
	svc := bugService002()
	first, err := svc.CreateRecordTask(model.RecordTask{Name: "first", SourceURL: "http://source-a"})
	if err != nil {
		panic(err)
	}
	second, err := svc.CreateRecordTask(model.RecordTask{Name: "second", SourceURL: "http://source-b"})
	if err != nil {
		panic(err)
	}
	linked, err := svc.CreateTrafficSample(model.TrafficSample{RecordTaskID: first.ID, Method: "GET", Path: "/linked", StatusCode: 200})
	if err != nil {
		panic(err)
	}
	if _, err = svc.CreateTrafficSample(model.TrafficSample{RecordTaskID: second.ID, Method: "GET", Path: "/unrelated", StatusCode: 200}); err != nil {
		panic(err)
	}
	task, err := svc.CreateReplayTask(model.ReplayTask{Name: "replay-first", TargetURL: "http://target", RecordTaskID: first.ID})
	if err != nil {
		panic(err)
	}
	if _, err = svc.UpdateReplayTask(task.ID, model.ReplayTask{Name: task.Name, TargetURL: task.TargetURL, RecordTaskID: first.ID, Status: model.ReplayTaskRunning}); err != nil {
		panic(err)
	}
	results, err := svc.ExecuteReplay(task.ID)
	if err != nil {
		panic(err)
	}
	if len(results) != 1 || results[0].SampleID != linked.ID {
		t.Fatalf("replay returned samples from another recording task: %#v", results)
	}
	legacy, err := svc.CreateReplayTask(model.ReplayTask{Name: "legacy-replay", TargetURL: "http://target"})
	if err != nil {
		panic(err)
	}
	if _, err = svc.UpdateReplayTask(legacy.ID, model.ReplayTask{Name: legacy.Name, TargetURL: legacy.TargetURL, Status: model.ReplayTaskRunning}); err != nil {
		panic(err)
	}
	legacyResults, err := svc.ExecuteReplay(legacy.ID)
	if err != nil {
		panic(err)
	}
	if len(legacyResults) != 2 {
		t.Fatalf("unlinked replay should retain all samples, got %d", len(legacyResults))
	}
}
