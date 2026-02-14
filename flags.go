package main

import "flag"

type flags struct {
	cronMode               bool
	startDate              string
	endDate                string
	dailyWorkoutPolicy     string
	minimalWorkoutDuration int
	outputFilePath         string
	accessToken            string
}

var parsedFlags *flags

func getFlags() *flags {
	if parsedFlags != nil {
		return parsedFlags
	}

	parsedFlags = &flags{}
	flag.BoolVar(&parsedFlags.cronMode, "cron", false, "Run in cron mode (weekly execution)")
	flag.StringVar(&parsedFlags.startDate, "start-date", "", "Start date in MM/DD/YYYY format")
	flag.StringVar(&parsedFlags.endDate, "end-date", "", "End date in MM/DD/YYYY format")
	flag.StringVar(&parsedFlags.dailyWorkoutPolicy, "daily-workout-policy", "longest", "Policy for handling multiple workouts in a day: 'longest' or 'merge'")
	flag.IntVar(&parsedFlags.minimalWorkoutDuration, "min-workout-duration", 30, "Minimum workout duration in minutes (default: 30)")
	flag.StringVar(&parsedFlags.outputFilePath, "output-file-path", "", "Path to save the output CSV file (default: generated based on date range)")
	flag.StringVar(&parsedFlags.accessToken, "access-token", "", "Strava API access token (optional, can also be set via STRAVA_ACCESS_TOKEN environment variable)")

	flag.Parse()

	return parsedFlags
}
