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
	activity_csv "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/repository/csv"
	athlete_provider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/provider"
	athlete_strava "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/provider/strava"
	authProvider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/provider"
	stravaAuthProvider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/provider/strava"
	auth_repository "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/repository/file"
	authService "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/service"
	custom_http "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/http"
)

type App struct {
	SaveWorkoutPeriod usecase.SaveWorkoutPeriod
}

func NewApp(cfg *config.Config, dailyWorkoutPolicy string) (*App, error) {
	policy := buildDailyWorkoutPolicy(dailyWorkoutPolicy)

	tokenRepo, err := auth_repository.NewTokenRepository(cfg.TokenFilePath)
	if err != nil {
		return nil, err
	}

	oauthClient := stravaAuthProvider.NewOAuthHttpClient(&http.Client{Timeout: 10 * time.Second}, cfg)
	oauthProvider := authProvider.NewOAuthProvider(oauthClient, cfg)
	tokenProvider := authService.NewTokenService(oauthProvider, tokenRepo)
	authTransport := custom_http.NewAuthTransport(tokenProvider)

	apiHttpClient := &http.Client{Timeout: 10 * time.Second, Transport: authTransport}
	workoutProvider := activityProvider.NewWorkoutProvider(
		stravaActivityProvider.NewActivityProvider(apiHttpClient, cfg),
	)

	httpAthleteStatsProvider := athlete_strava.NewHttpAthleteStatsProvider(apiHttpClient, cfg.StravaApiBaseUrl)
	athleteProvider := athlete_provider.NewAthleteProvider(httpAthleteStatsProvider)

	useCase := usecase.NewSaveWorkoutPeriod(
		policy,
		activity_csv.NewCSVWorkoutPeriodSaver(cfg.OutputFilePath),
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
