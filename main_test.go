package main

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
)

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

func TestRunDelegatesToCronModeWhenEnabled(t *testing.T) {
	originalStart := startWeeklySunday20
	defer func() {
		startWeeklySunday20 = originalStart
	}()

	expectedErr := errors.New("cron path used")
	startWeeklySunday20 = func(func()) error {
		return expectedErr
	}

	err := run(config.CLIOptions{CronMode: true})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected run to delegate to cron mode, got %v", err)
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
