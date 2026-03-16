package config

import "testing"

func TestParseCLIShouldReturnDefaultsWhenNoFlagsAreProvided(t *testing.T) {
	opts, err := ParseCLI([]string{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if opts.CronMode != false {
		t.Fatalf("expected CronMode to default to false")
	}
	if opts.StartDate != "" {
		t.Fatalf("expected StartDate to default to empty string")
	}
	if opts.EndDate != "" {
		t.Fatalf("expected EndDate to default to empty string")
	}
	if opts.DailyWorkoutPolicy != "longest" {
		t.Fatalf("expected DailyWorkoutPolicy to default to longest")
	}
	if opts.MinimalWorkoutDuration != 30 {
		t.Fatalf("expected MinimalWorkoutDuration to default to 30")
	}
	if opts.OutputFilePath != "" {
		t.Fatalf("expected OutputFilePath to default to empty string")
	}
}

func TestParseCLIShouldParseAllFlagsWhenValidArgsAreProvided(t *testing.T) {
	args := []string{
		"--cron",
		"--start-date", "01/01/2026",
		"--end-date", "01/07/2026",
		"--daily-workout-policy", "merge",
		"--min-duration", "45",
		"--output-file", "workouts-jan",
	}

	opts, err := ParseCLI(args)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if opts.CronMode != true {
		t.Fatalf("expected CronMode to be true")
	}
	if opts.StartDate != "01/01/2026" {
		t.Fatalf("expected StartDate to match input")
	}
	if opts.EndDate != "01/07/2026" {
		t.Fatalf("expected EndDate to match input")
	}
	if opts.DailyWorkoutPolicy != "merge" {
		t.Fatalf("expected DailyWorkoutPolicy to match input")
	}
	if opts.MinimalWorkoutDuration != 45 {
		t.Fatalf("expected MinimalWorkoutDuration to match input")
	}
	if opts.OutputFilePath != "workouts-jan" {
		t.Fatalf("expected OutputFilePath to match input")
	}
}

func TestParseCLIShouldReturnErrorWhenUnknownFlagIsProvided(t *testing.T) {
	opts, err := ParseCLI([]string{"--unknown-flag"})
	if err == nil {
		t.Fatalf("expected error for unknown flag")
	}

	if opts != (CLIOptions{}) {
		t.Fatalf("expected zero-value options when parse fails")
	}
}
