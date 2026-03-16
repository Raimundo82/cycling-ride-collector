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
	activity_excel "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/repository/excel"
	athlete_provider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/provider"
	athlete_strava "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/provider/strava"
	token_client "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/client"
	token_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
	token_provider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/provider"
	custom_http "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/http"
)

type App struct {
	SaveWorkoutPeriod usecase.SaveWorkoutPeriod
}

func NewApp(cfg *config.Config, dailyWorkoutPolicy string) (*App, error) {
	policy := buildDailyWorkoutPolicy(dailyWorkoutPolicy)

	stravaTokenClient := token_client.NewTokenClient(cfg.StravaOauthBaseUrl)
	stravaTokenInput := &token_model.RefreshTokenInput{
		ClientID:     cfg.StravaClientId,
		ClientSecret: cfg.StravaClientSecret,
		GrantType:    "refresh_token",
		RefreshToken: cfg.StravaRefreshToken,
	}
	stravaTokenProvider := token_provider.NewTokenProvider(stravaTokenInput, stravaTokenClient)

	authTransport := custom_http.NewAuthTransport(stravaTokenProvider)

	apiHttpClient := &http.Client{Timeout: 10 * time.Second, Transport: authTransport}
	workoutProvider := activityProvider.NewWorkoutProvider(
		stravaActivityProvider.NewActivityProvider(apiHttpClient, cfg),
	)

	httpAthleteStatsProvider := athlete_strava.NewHttpAthleteStatsProvider(apiHttpClient, cfg.StravaApiBaseUrl)
	athleteProvider := athlete_provider.NewAthleteProvider(httpAthleteStatsProvider)

	useCase := usecase.NewSaveWorkoutPeriod(
		policy,
		activity_excel.NewExcelWorkoutPeriodSaverWithOptions(
			cfg.OutputFilePath+".xlsx",
			cfg.ExcelTemplate.TemplatePath,
			cfg.ExcelTemplate.SheetName,
			cfg.ExcelTemplate.StartCell,
		),
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
