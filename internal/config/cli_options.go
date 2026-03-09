package config

import "flag"

type CLIOptions struct {
	CronMode               bool
	StartDate              string
	EndDate                string
	DailyWorkoutPolicy     string
	MinimalWorkoutDuration int
	OutputFilePath         string
}

func ParseCLI(args []string) (CLIOptions, error) {
	fs := flag.NewFlagSet("cycling-ride-collector", flag.ContinueOnError)

	var opts CLIOptions
	fs.BoolVar(&opts.CronMode, "cron", false, "Run in cron mode (weekly execution)")
	fs.StringVar(&opts.StartDate, "start-date", "", "Start date in MM/DD/YYYY format")
	fs.StringVar(&opts.EndDate, "end-date", "", "End date in MM/DD/YYYY format")
	fs.StringVar(&opts.DailyWorkoutPolicy, "daily-workout-policy", "longest", "Policy: 'longest' or 'merge'")
	fs.IntVar(&opts.MinimalWorkoutDuration, "min-duration", 30, "Minimal workout duration in minutes")
	fs.StringVar(&opts.OutputFilePath, "output-file", "", "Output CSV path")

	if err := fs.Parse(args); err != nil {
		return CLIOptions{}, err
	}

	return opts, nil
}
