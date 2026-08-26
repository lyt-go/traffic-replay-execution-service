package service

import (
	"testing"

	"trafficreplay/internal/config"
	"trafficreplay/internal/model"
	"trafficreplay/internal/store"
	"trafficreplay/pkg/logger"
)

func bugService008() *Service {
	return New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
}

func TestRecordTaskReadIsSnapshot(t *testing.T) {
	svc := bugService008()
	task, err := svc.CreateRecordTask(model.RecordTask{Name: "original", SourceURL: "http://source"})
	if err != nil { panic(err) }
	view, err := svc.GetRecordTask(task.ID)
	if err != nil { panic(err) }
	view.Name = "mutated-without-update"
	got, err := svc.GetRecordTask(task.ID)
	if err != nil { panic(err) }
	if got.Name != "original" { t.Fatalf("read exposed mutable store state: %q", got.Name) }
}
