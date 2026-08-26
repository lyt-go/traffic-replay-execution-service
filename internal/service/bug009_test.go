package service

import (
	"testing"

	"trafficreplay/internal/config"
	"trafficreplay/internal/model"
	"trafficreplay/internal/store"
	"trafficreplay/pkg/logger"
)

func bugService009() *Service {
	return New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
}

func TestResultQueryKeepsStoredLatencyIsolated(t *testing.T) {
	svc := bugService009()
	record, err := svc.CreateRecordTask(model.RecordTask{Name: "capture", SourceURL: "http://source"})
	if err != nil { panic(err) }
	sample, err := svc.CreateTrafficSample(model.TrafficSample{RecordTaskID: record.ID, Method: "GET", Path: "/data"})
	if err != nil { panic(err) }
	task, err := svc.CreateReplayTask(model.ReplayTask{Name: "replay", TargetURL: "http://target"})
	if err != nil { panic(err) }
	result, err := svc.CreateReplayResult(model.ReplayResult{ReplayTaskID: task.ID, SampleID: sample.ID, LatencyMs: 10})
	if err != nil { panic(err) }
	view, err := svc.GetReplayResult(result.ID)
	if err != nil { panic(err) }
	view.LatencyMs = 999
	got, err := svc.GetReplayResult(result.ID)
	if err != nil { panic(err) }
	if got.LatencyMs != 10 { t.Fatalf("read exposed mutable result state: %d", got.LatencyMs) }
}
