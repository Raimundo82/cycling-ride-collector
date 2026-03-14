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
	authProvider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/provider"
	stravaAuthProvider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/provider/strava"
	auth_repository "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/repository/file"
	authService "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/service"
	custom_http "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/http"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/notification/email"
)

type App struct {
	SaveWorkoutPeriod usecase.SaveWorkoutPeriod
	SendWorkoutReport usecase.SendWorkoutReport
}

func NewApp(cfg *config.Config, dailyWorkoutPolicy string) (*App, error) {
	policy := buildDailyWorkoutPolicy(dailyWorkoutPolicy)

	tokenRepo, err := auth_repository.NewTokenRepository(cfg.TokenFilePath)
	if err != nil {
		return nil, err
	}

	oauthClient := stravaAuthProvider.NewOAuthHttpClient(&http.Client{Timeout: 10 * time.Second}, cfg)
	oauthProvider := authProvider.NewOAuthProvider(oauthClient, cfg)
	stravaTokenProvider := authService.NewStravaTokenService(oauthProvider, tokenRepo)
	googleTokenProvider := authService.NewGoogleTokenService(oauthProvider, tokenRepo)

	authTransport := custom_http.NewAuthTransport(stravaTokenProvider)

	apiHttpClient := &http.Client{Timeout: 10 * time.Second, Transport: authTransport}
	workoutProvider := activityProvider.NewWorkoutProvider(
		stravaActivityProvider.NewActivityProvider(apiHttpClient, cfg),
	)

	httpAthleteStatsProvider := athlete_strava.NewHttpAthleteStatsProvider(apiHttpClient, cfg.StravaApiBaseUrl)
	athleteProvider := athlete_provider.NewAthleteProvider(httpAthleteStatsProvider)

	saveWorkoutPeriodUseCase := usecase.NewSaveWorkoutPeriod(
		policy,
		activity_excel.NewExcelWorkoutPeriodSaverWithOptions(
			cfg.OutputFilePath,
			cfg.ExcelTemplate.TemplatePath,
			cfg.ExcelTemplate.SheetName,
			cfg.ExcelTemplate.StartCell,
		),
		workoutProvider,
		athleteProvider,
	)

	sendWorkoutReportUseCase := usecase.NewSendWorkoutReport(
		email.NewGmailWorkoutReportSender(&http.Client{Timeout: 10 * time.Second}, googleTokenProvider, cfg),
	)

	return &App{
		SaveWorkoutPeriod: *saveWorkoutPeriodUseCase,
		SendWorkoutReport: *sendWorkoutReportUseCase,
	}, nil
}

func (a *App) Run(request *input.SaveWorkoutPeriodRequest, reportPath string) error {
	err := a.SaveWorkoutPeriod.Execute(request.Period, request.MinimalWorkoutDuration)
	if err != nil {
		return err
	}
	if err := a.SendWorkoutReport.Execute(reportPath); err != nil {
		return err
	}
	return nil
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
