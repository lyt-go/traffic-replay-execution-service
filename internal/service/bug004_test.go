package service

import (
	"testing"

	"trafficreplay/internal/config"
	"trafficreplay/internal/model"
	"trafficreplay/internal/store"
	"trafficreplay/pkg/logger"
)

func bugService004() *Service {
	return New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
}

func TestDisabledReplayConfigDoesNotRunSchedule(t *testing.T) {
	svc := bugService004()
	cfg, err := svc.CreateReplayConfig(model.ReplayConfig{Name: "disabled", TargetHost: "target", Enabled: true})
	if err != nil {
		panic(err)
	}
	if _, err = svc.UpdateReplayConfig(cfg.ID, model.ReplayConfig{Name: cfg.Name, TargetHost: cfg.TargetHost, Enabled: false}); err != nil {
		panic(err)
	}
	sch, err := svc.CreateSchedule(model.Schedule{Name: "nightly", ConfigID: cfg.ID, CronExpr: "* * * * *"})
	if err != nil {
		panic(err)
	}
	if _, err = svc.RunSchedule(sch.ID); !model.IsValidationError(err) {
		t.Fatalf("expected disabled config to block schedule, got %v", err)
	}
	got, err := svc.GetSchedule(sch.ID)
	if err != nil {
		panic(err)
	}
	if !got.LastRunAt.IsZero() {
		t.Fatal("disabled schedule recorded a run")
	}
	legacy, err := svc.CreateReplayConfig(model.ReplayConfig{Name: "legacy-default", TargetHost: "legacy", Enabled: false})
	if err != nil {
		panic(err)
	}
	legacySchedule, err := svc.CreateSchedule(model.Schedule{Name: "legacy-schedule", ConfigID: legacy.ID, CronExpr: "* * * * *"})
	if err != nil {
		panic(err)
	}
	if _, err = svc.RunSchedule(legacySchedule.ID); err != nil {
		t.Fatalf("untouched legacy config should remain runnable: %v", err)
	}
}
