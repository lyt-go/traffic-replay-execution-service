package service

import (
	"errors"
	"testing"

	"trafficreplay/internal/config"
	"trafficreplay/internal/model"
	"trafficreplay/internal/store"
	"trafficreplay/pkg/logger"
)

type failingReplayTaskStore struct { store.Store }

func (f failingReplayTaskStore) UpdateReplayTask(*model.ReplayTask) error { return errors.New("write failed") }

func TestFailedTaskWritePreservesPriorState(t *testing.T) {
	base := store.NewMemoryStore()
	cfg := &config.Config{MaxPageSize: 100}
	creator := New(base, logger.NewLevel(logger.LevelError), cfg)
	task, err := creator.CreateReplayTask(model.ReplayTask{Name: "before", TargetURL: "http://target"})
	if err != nil { panic(err) }
	failing := New(failingReplayTaskStore{Store: base}, logger.NewLevel(logger.LevelError), cfg)
	if _, err = failing.UpdateReplayTask(task.ID, model.ReplayTask{Name: "after", TargetURL: task.TargetURL, Status: model.ReplayTaskRunning}); err == nil {
		t.Fatal("expected backing store failure")
	}
	got, err := creator.GetReplayTask(task.ID)
	if err != nil { panic(err) }
	if got.Name != "before" || got.Status != model.ReplayTaskPending {
		t.Fatalf("failed update leaked into stored task: %#v", got)
	}
}
