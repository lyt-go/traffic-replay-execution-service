package service

import (
	"testing"

	"trafficreplay/internal/config"
	"trafficreplay/internal/model"
	"trafficreplay/internal/store"
	"trafficreplay/pkg/logger"
)

func bugService005() *Service {
	return New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
}

func TestScheduleUpdateKeepsExistingConfigOnMissingReference(t *testing.T) {
	svc := bugService005()
	cfg, err := svc.CreateReplayConfig(model.ReplayConfig{Name: "valid", TargetHost: "target"})
	if err != nil { panic(err) }
	sch, err := svc.CreateSchedule(model.Schedule{Name: "daily", ConfigID: cfg.ID, CronExpr: "* * * * *"})
	if err != nil { panic(err) }
	if _, err = svc.UpdateSchedule(sch.ID, model.Schedule{Name: sch.Name, ConfigID: "missing-config", CronExpr: sch.CronExpr}); !model.IsValidationError(err) {
		t.Fatalf("expected missing config rejection, got %v", err)
	}
	got, err := svc.GetSchedule(sch.ID)
	if err != nil { panic(err) }
	if got.ConfigID != cfg.ID { t.Fatalf("schedule reference changed to %q", got.ConfigID) }
}
