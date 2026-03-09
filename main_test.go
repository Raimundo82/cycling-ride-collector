package main

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
)

func TestRunByModeCallsCronModeWhenCronEnabled(t *testing.T) {
	cronCalled := 0
	onceCalled := 0

	err := runByMode(
		config.CLIOptions{CronMode: true},
		func() error {
			cronCalled++
			return nil
		},
		func(config.CLIOptions) error {
			onceCalled++
			return nil
		},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cronCalled != 1 {
		t.Fatalf("expected cron path to be called once, got %d", cronCalled)
	}
	if onceCalled != 0 {
		t.Fatalf("expected once path not to be called, got %d", onceCalled)
	}
}

func TestRunByModeCallsOnceModeWhenCronDisabled(t *testing.T) {
	cronCalled := 0
	onceCalled := 0
	input := config.CLIOptions{CronMode: false, StartDate: "01/01/2026"}

	err := runByMode(
		input,
		func() error {
			cronCalled++
			return nil
		},
		func(opts config.CLIOptions) error {
			onceCalled++
			if opts.StartDate != input.StartDate {
				t.Fatalf("expected options to be forwarded, got %+v", opts)
			}
			return nil
		},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if onceCalled != 1 {
		t.Fatalf("expected once path to be called once, got %d", onceCalled)
	}
	if cronCalled != 0 {
		t.Fatalf("expected cron path not to be called, got %d", cronCalled)
	}
}

func TestRunByModePropagatesCronError(t *testing.T) {
	expectedErr := errors.New("cron failed")
	err := runByMode(
		config.CLIOptions{CronMode: true},
		func() error { return expectedErr },
		func(config.CLIOptions) error { return nil },
	)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected cron error to propagate, got %v", err)
	}
}

func TestRunByModePropagatesOnceError(t *testing.T) {
	expectedErr := errors.New("once failed")
	err := runByMode(
		config.CLIOptions{CronMode: false},
		func() error { return nil },
		func(config.CLIOptions) error { return expectedErr },
	)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected once error to propagate, got %v", err)
	}
}

func TestRunCronModePropagatesSchedulerError(t *testing.T) {
	originalStart := startWeeklySunday20
	defer func() {
		startWeeklySunday20 = originalStart
	}()

	expectedErr := errors.New("scheduler failed")
	startWeeklySunday20 = func(func()) error {
		return expectedErr
	}

	err := runCronMode()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected scheduler error to propagate, got %v", err)
	}
}

func TestBuildRuntimeConfigOverridesOutputPath(t *testing.T) {
	cfg := buildRuntimeConfig(config.CLIOptions{OutputFilePath: "custom.csv"})
	if cfg.OutputFilePath != "custom.csv" {
		t.Fatalf("expected output path override, got %q", cfg.OutputFilePath)
	}
}

func TestBuildSaveWorkoutRequestSuccess(t *testing.T) {
	opts := config.CLIOptions{
		StartDate:              "01/01/2026",
		EndDate:                "01/02/2026",
		DailyWorkoutPolicy:     "longest",
		MinimalWorkoutDuration: 45,
	}

	request, err := buildSaveWorkoutRequest(opts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if request.DailyWorkoutPolicy != "longest" {
		t.Fatalf("expected policy longest, got %q", request.DailyWorkoutPolicy)
	}
	if request.MinimalWorkoutDuration != 45 {
		t.Fatalf("expected duration 45, got %d", request.MinimalWorkoutDuration)
	}
	if request.Period.StartDate().Format("2006-01-02") != "2026-01-01" {
		t.Fatalf("unexpected start date: %s", request.Period.StartDate().Format("2006-01-02"))
	}
	if request.Period.EndDate().Format("2006-01-02") != "2026-01-02" {
		t.Fatalf("unexpected end date: %s", request.Period.EndDate().Format("2006-01-02"))
	}
}

func TestBuildSaveWorkoutRequestMinDurationFloor(t *testing.T) {
	opts := config.CLIOptions{
		StartDate:              "01/01/2026",
		EndDate:                "01/02/2026",
		DailyWorkoutPolicy:     "longest",
		MinimalWorkoutDuration: 10,
	}

	request, err := buildSaveWorkoutRequest(opts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if request.MinimalWorkoutDuration != 30 {
		t.Fatalf("expected duration floor 30, got %d", request.MinimalWorkoutDuration)
	}
}

func TestBuildSaveWorkoutRequestMissingRequiredDates(t *testing.T) {
	_, err := buildSaveWorkoutRequest(config.CLIOptions{})
	if err == nil {
		t.Fatal("expected error for missing dates")
	}
	if err.Error() != "flags --start-date and --end-date are required" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestBuildSaveWorkoutRequestInvalidStartDate(t *testing.T) {
	opts := config.CLIOptions{
		StartDate:              "invalid",
		EndDate:                "01/02/2026",
		DailyWorkoutPolicy:     "longest",
		MinimalWorkoutDuration: 30,
	}

	_, err := buildSaveWorkoutRequest(opts)
	if err == nil {
		t.Fatal("expected invalid start date error")
	}
	if !strings.Contains(err.Error(), "invalid start date:") {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestBuildSaveWorkoutRequestInvalidEndDate(t *testing.T) {
	opts := config.CLIOptions{
		StartDate:              "01/01/2026",
		EndDate:                "invalid",
		DailyWorkoutPolicy:     "longest",
		MinimalWorkoutDuration: 30,
	}

	_, err := buildSaveWorkoutRequest(opts)
	if err == nil {
		t.Fatal("expected invalid end date error")
	}
	if !strings.Contains(err.Error(), "invalid end date:") {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestBuildSaveWorkoutRequestInvalidPeriod(t *testing.T) {
	opts := config.CLIOptions{
		StartDate:              "01/03/2026",
		EndDate:                "01/01/2026",
		DailyWorkoutPolicy:     "longest",
		MinimalWorkoutDuration: 30,
	}

	_, err := buildSaveWorkoutRequest(opts)
	if err == nil {
		t.Fatal("expected invalid period error")
	}
	if err.Error() != "invalid period: start date must be before end date" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestBuildSaveWorkoutRequestInvalidDailyWorkoutPolicy(t *testing.T) {
	opts := config.CLIOptions{
		StartDate:              "01/01/2026",
		EndDate:                "01/02/2026",
		DailyWorkoutPolicy:     "invalid",
		MinimalWorkoutDuration: 30,
	}

	_, err := buildSaveWorkoutRequest(opts)
	if err == nil {
		t.Fatal("expected invalid policy error")
	}

	expected := "invalid daily workout policy: invalid. Allowed values are 'longest' or 'merge'"
	if err.Error() != expected {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestRunReturnsInitializationError(t *testing.T) {
	t.Setenv("TOKEN_FILE_PATH", "token-file-does-not-exist.json")

	err := run(config.CLIOptions{
		StartDate:              "01/01/2026",
		EndDate:                "01/02/2026",
		DailyWorkoutPolicy:     "longest",
		MinimalWorkoutDuration: 30,
	})
	if err == nil {
		t.Fatal("expected run initialization error")
	}
	if !strings.Contains(err.Error(), "failed to initialize app: failed to open token file") {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}

func TestBuildCronOptionsCoversMondayToSunday(t *testing.T) {
	// Sunday, March 8, 2026
	now := time.Date(2026, 3, 8, 20, 0, 0, 0, time.UTC)
	opts := buildCronOptions(now)

	if opts.StartDate != "03/02/2026" {
		t.Fatalf("expected Monday start date 03/02/2026, got %q", opts.StartDate)
	}
	if opts.EndDate != "03/08/2026" {
		t.Fatalf("expected Sunday end date 03/08/2026, got %q", opts.EndDate)
	}
	if !opts.CronMode {
		t.Fatal("expected cron mode true")
	}
	if opts.DailyWorkoutPolicy != "longest" {
		t.Fatalf("expected longest policy, got %q", opts.DailyWorkoutPolicy)
	}
	if opts.MinimalWorkoutDuration != 30 {
		t.Fatalf("expected minimal duration 30, got %d", opts.MinimalWorkoutDuration)
	}
}
