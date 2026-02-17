package main

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/raimundo82/cycling-ride-collector/internal/application/usecase/input"
	"github.com/raimundo82/cycling-ride-collector/internal/config"
)

func TestParseFlagsAndConfigDefaults(t *testing.T) {
	t.Setenv("TOKEN_FILE_PATH", "/tmp/token.json")

	cfg, request, err := parseWithArgs(t,
		"--start-date", "01/01/2026",
		"--end-date", "01/02/2026",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.TokenFilePath != "/tmp/token.json" {
		t.Fatalf("expected TOKEN_FILE_PATH from env, got %q", cfg.TokenFilePath)
	}
	if cfg.OutputFilePath != "workouts_summary_2026-01-01_to_2026-01-02.csv" {
		t.Fatalf("unexpected default output file path: %q", cfg.OutputFilePath)
	}
	if request.DailyWorkoutPolicy != "longest" {
		t.Fatalf("expected default daily workout policy 'longest', got %q", request.DailyWorkoutPolicy)
	}
	if request.MinimalWorkoutDuration != 30 {
		t.Fatalf("expected default minimal workout duration 30, got %d", request.MinimalWorkoutDuration)
	}

	if request.Period.StartDate().Format("2006-01-02") != "2026-01-01" {
		t.Fatalf("unexpected period start date: %s", request.Period.StartDate().Format("2006-01-02"))
	}
	if request.Period.EndDate().Format("2006-01-02") != "2026-01-02" {
		t.Fatalf("unexpected period end date: %s", request.Period.EndDate().Format("2006-01-02"))
	}
}

func TestParseFlagsAndConfigOverrides(t *testing.T) {
	cfg, request, err := parseWithArgs(t,
		"--start-date", "01/01/2026",
		"--end-date", "01/03/2026",
		"--output-file", "custom.csv",
		"--min-duration", "45",
		"--daily-workout-policy", "merge",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.OutputFilePath != "custom.csv" {
		t.Fatalf("expected output-file override, got %q", cfg.OutputFilePath)
	}
	if request.DailyWorkoutPolicy != "merge" {
		t.Fatalf("expected daily workout policy 'merge', got %q", request.DailyWorkoutPolicy)
	}
	if request.MinimalWorkoutDuration != 45 {
		t.Fatalf("expected minimal workout duration 45, got %d", request.MinimalWorkoutDuration)
	}
}

func TestParseFlagsAndConfigMinDurationFloor(t *testing.T) {
	_, request, err := parseWithArgs(t,
		"--start-date", "01/01/2026",
		"--end-date", "01/02/2026",
		"--min-duration", "10",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if request.MinimalWorkoutDuration != 30 {
		t.Fatalf("expected min-duration floor to 30, got %d", request.MinimalWorkoutDuration)
	}
}

func TestParseFlagsAndConfigMissingRequiredFlags(t *testing.T) {
	_, _, err := parseWithArgs(t)
	if err == nil {
		t.Fatal("expected error for missing required flags")
	}

	if err.Error() != "flags --start-date and --end-date are required" {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}

func TestParseFlagsAndConfigInvalidStartDate(t *testing.T) {
	_, _, err := parseWithArgs(t,
		"--start-date", "invalid-date",
		"--end-date", "01/02/2026",
	)
	if err == nil {
		t.Fatal("expected invalid start date error")
	}

	if !strings.Contains(err.Error(), "invalid start date:") {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}

func TestParseFlagsAndConfigInvalidEndDate(t *testing.T) {
	_, _, err := parseWithArgs(t,
		"--start-date", "01/01/2026",
		"--end-date", "invalid-date",
	)
	if err == nil {
		t.Fatal("expected invalid end date error")
	}

	if !strings.Contains(err.Error(), "invalid end date:") {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}

func TestParseFlagsAndConfigInvalidPeriod(t *testing.T) {
	_, _, err := parseWithArgs(t,
		"--start-date", "01/03/2026",
		"--end-date", "01/01/2026",
	)
	if err == nil {
		t.Fatal("expected invalid period error")
	}

	if err.Error() != "invalid period: start date must be before end date" {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}

func TestParseFlagsAndConfigInvalidDailyWorkoutPolicy(t *testing.T) {
	_, _, err := parseWithArgs(t,
		"--start-date", "01/01/2026",
		"--end-date", "01/02/2026",
		"--daily-workout-policy", "not-valid",
	)
	if err == nil {
		t.Fatal("expected invalid daily workout policy error")
	}

	expected := "invalid daily workout policy: not-valid. Allowed values are 'longest' or 'merge'"
	if err.Error() != expected {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}

func TestRunReturnsInitializationError(t *testing.T) {
	t.Setenv("TOKEN_FILE_PATH", "token-file-does-not-exist.json")

	err := runWithArgs(t,
		"--start-date", "01/01/2026",
		"--end-date", "01/02/2026",
	)
	if err == nil {
		t.Fatal("expected run initialization error")
	}

	if !strings.Contains(err.Error(), "failed to initialize app: failed to open token file") {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}

func parseWithArgs(t *testing.T, args ...string) (*config.Config, *input.SaveWorkoutPeriodRequest, error) {
	t.Helper()

	oldArgs := os.Args
	oldFlagSet := flag.CommandLine

	os.Args = append([]string{"cycling-ride-collector"}, args...)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldFlagSet
	}()

	return parseFlagsAndConfig()
}

func runWithArgs(t *testing.T, args ...string) error {
	t.Helper()

	oldArgs := os.Args
	oldFlagSet := flag.CommandLine

	os.Args = append([]string{"cycling-ride-collector"}, args...)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldFlagSet
	}()

	return run()
}
