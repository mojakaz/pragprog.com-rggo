package pomodoro_test

import (
	"fmt"
	"pragprog.com/rggo/interactiveTools/pomo/pomodoro"
	"testing"
	"time"
)

func TestDailySummary(t *testing.T) {
	repo, cleanup := getRepo(t, dbname)
	defer cleanup()
	const duration = 2 * time.Second
	config := pomodoro.NewConfig(repo, duration, duration, duration)

	ds, err := pomodoro.DailySummary(time.Now(), config)
	if err != nil {
		t.Errorf("Expected no error, got %q", err)
	}
	if ds[0] != duration {
		t.Errorf("Expected duration %q, got %q", config.PomodoroDuration, ds[0])
	}
	if ds[1] != duration/2 {
		t.Errorf("Expected duration %q, got %q", config.ShortBreakDuration, ds[1])
	}
}

func TestRangeSummary(t *testing.T) {
	repo, cleanup := getRepo(t, dbname)
	defer cleanup()
	const duration = 2 * time.Second
	config := pomodoro.NewConfig(repo, duration, duration, duration)

	ws, err := pomodoro.RangeSummary(time.Now(), 1, config)
	if err != nil {
		t.Errorf("Expected no error, got %q", err)
	}
	if ws[0].Name != "Pomodoro" {
		t.Errorf("Expected Name 'Pomodoro', got %q", ws[0].Name)
	}
	now := time.Now()
	expLabel := fmt.Sprintf("%02d/%s", now.Day(), now.Format("Jan"))
	if ws[0].Labels[0] != expLabel {
		t.Errorf("Expected Label %q, got %q", expLabel, ws[0].Labels[0])
	}
	expValue := float64(2)
	if ws[0].Values[0] != expValue {
		t.Errorf("Expected Value %f, got %f", expValue, ws[0].Values[0])
	}
	if ws[1].Name != "Break" {
		t.Errorf("Expected Name 'Break', got %q", ws[1].Name)
	}
	if ws[1].Labels[0] != expLabel {
		t.Errorf("Expected Label %q, got %q", expLabel, ws[1].Labels[0])
	}
	expValue = float64(1)
	if ws[1].Values[0] != expValue {
		t.Errorf("Expected Value %f, got %f", expValue, ws[1].Values[0])
	}

}
