package main

import (
	"net/http"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/application/usecase"
	"github.com/raimundo82/cycling-ride-collector/internal/application/usecase/input"
	"github.com/raimundo82/cycling-ride-collector/internal/config"
	activityProvider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/provider"
	stravaActivityProvider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/provider/strava"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/repository/csv"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/provider"
	athlete_strava "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/provider/strava"
	authProvider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/provider"
	stravaAuthProvider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/provider/strava"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/repository/file"
	authService "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/service"
)

type App struct {
	SaveWorkoutPeriod usecase.SaveWorkoutPeriod
}

func NewApp(cfg *config.Config, dailyWorkoutPolicy string) (*App, error) {
	policy := buildDailyWorkoutPolicy(dailyWorkoutPolicy)

	tokenRepo, err := file.NewTokenRepository(cfg.TokenFilePath)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}
	oauthClient := stravaAuthProvider.NewOAuthHttpClient(httpClient, cfg)
	oauthProvider := authProvider.NewOAuthProvider(oauthClient, cfg)
	tokenProvider := authService.NewTokenService(oauthProvider, tokenRepo)

	workoutProvider := activityProvider.NewWorkoutProvider(
		stravaActivityProvider.NewActivityProvider(httpClient, cfg, tokenProvider),
	)

	athleteProvider := provider.NewAthleteProvider(athlete_strava.NewHttpAthleteStatsProvider(httpClient, cfg))

	useCase := usecase.NewSaveWorkoutPeriod(
		policy,
		csv.NewCSVWorkoutPeriodSaver(cfg.OutputFilePath),
		workoutProvider,
		athleteProvider,
	)

	return &App{SaveWorkoutPeriod: *useCase}, nil
}

func (a *App) Run(request *input.SaveWorkoutPeriodRequest) error {
	return a.SaveWorkoutPeriod.Execute(request.Period, request.MinimalWorkoutDuration)
}

func buildDailyWorkoutPolicy(workoutPolicy string) contracts.DailyWorkoutPolicy {
	switch workoutPolicy {
	case "merge":
		return usecase.NewMergeWorkouts()
	case "longest":
		fallthrough
	default:
		return usecase.NewLongestWorkout()
	}
}
