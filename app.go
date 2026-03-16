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
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/notification/email"
)

type App struct {
	SaveWorkoutPeriod usecase.SaveWorkoutPeriod
	SendWorkoutReport usecase.SendWorkoutReport
}

func NewApp(cfg *config.Config, dailyWorkoutPolicy string) (*App, error) {
	policy := buildDailyWorkoutPolicy(dailyWorkoutPolicy)
	stravaCfg := cfg.Strava

	stravaTokenClient := token_client.NewTokenClient(stravaCfg.OAuthBaseUrl, &http.Client{Timeout: 10 * time.Second})
	stravaTokenInput := &token_model.RefreshTokenInput{
		ClientID:     stravaCfg.ClientId,
		ClientSecret: stravaCfg.ClientSecret,
		GrantType:    "refresh_token",
		RefreshToken: stravaCfg.RefreshToken,
	}
	stravaTokenProvider := token_provider.NewTokenProvider(stravaTokenInput, stravaTokenClient)

	authTransport := custom_http.NewAuthTransport(stravaTokenProvider)

	apiHttpClient := &http.Client{Timeout: 10 * time.Second, Transport: authTransport}
	workoutProvider := activityProvider.NewWorkoutProvider(
		stravaActivityProvider.NewActivityProvider(apiHttpClient, cfg),
	)

	httpAthleteStatsProvider := athlete_strava.NewHttpAthleteStatsProvider(apiHttpClient, stravaCfg.ApiBaseUrl)
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

	googleCfg := cfg.GoogleOAuth
	googleTokenClient := token_client.NewTokenClient(googleCfg.OAuthBaseUrl, &http.Client{Timeout: 10 * time.Second})
	googleTokenInput := &token_model.RefreshTokenInput{
		ClientID:     googleCfg.ClientID,
		ClientSecret: googleCfg.ClientSecret,
		GrantType:    "refresh_token",
		RefreshToken: googleCfg.RefreshToken,
	}
	googleTokenProvider := token_provider.NewTokenProvider(googleTokenInput, googleTokenClient)
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
