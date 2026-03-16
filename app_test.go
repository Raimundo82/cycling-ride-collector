package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/application/usecase"
	"github.com/raimundo82/cycling-ride-collector/internal/application/usecase/input"
	"github.com/raimundo82/cycling-ride-collector/internal/config"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

var athlete = domain.NewAthlete(75, 135, 240)

const testWorkoutReportPath = "/tmp/workout-report.xlsx"

type spyDailyWorkoutPolicy struct {
	Called                int
	MinWorkoutDuration    int
	ReceivedDailyWorkouts []*domain.Workout
}

var _ contracts.DailyWorkoutPolicy = (*spyDailyWorkoutPolicy)(nil)

func (s *spyDailyWorkoutPolicy) GetDailyWorkout(dailyWorkouts []*domain.Workout, minWorkoutDuration int) *domain.Workout {
	s.Called++
	s.MinWorkoutDuration = minWorkoutDuration
	s.ReceivedDailyWorkouts = dailyWorkouts

	if len(dailyWorkouts) == 0 {
		return nil
	}

	return dailyWorkouts[0]
}

type spyWorkoutRepository struct {
	Called   int
	Workouts []*domain.Workout
	Err      error
}

var _ contracts.WorkoutRepository = (*spyWorkoutRepository)(nil)

func (s *spyWorkoutRepository) SaveAll(workouts []*domain.Workout, athlete *domain.Athlete) error {
	s.Called++
	s.Workouts = workouts
	return s.Err
}

type spyReportSender struct {
	Called     int
	ReportPath string
	Err        error
}

func (s *spyReportSender) Send(reportPath string) error {
	s.Called++
	s.ReportPath = reportPath
	return s.Err
}

var _ contracts.WorkoutRepository = (*spyWorkoutRepository)(nil)

type stubWorkoutProvider struct {
	Called int
	Period domain.Period
	Result []*domain.Workout
	Err    error
}

var _ contracts.WorkoutProvider = (*stubWorkoutProvider)(nil)

func (s *stubWorkoutProvider) GetWorkoutsByPeriod(period domain.Period) ([]*domain.Workout, error) {
	s.Called++
	s.Period = period
	if s.Err != nil {
		return nil, s.Err
	}

	return s.Result, nil
}

type stubAthleteProvider struct {
	Athlete *domain.Athlete
	Called  int
	Err     error
}

// GetAthleteData implements [contracts.AthleteDataProvider].
func (s *stubAthleteProvider) GetAthleteData() (*domain.Athlete, error) {
	s.Called++
	return s.Athlete, s.Err
}

var _ contracts.AthleteDataProvider = (*stubAthleteProvider)(nil)

func TestBuildDailyWorkoutPolicyMerge(t *testing.T) {
	Convey("Given the merge daily workout policy", t, func() {
		workouts := []*domain.Workout{
			domain.NewWorkout(domain.WorkoutParams{ID: 1, DurationInMin: 30, DistanceInKm: 10}),
			domain.NewWorkout(domain.WorkoutParams{ID: 2, DurationInMin: 40, DistanceInKm: 20}),
		}

		Convey("When building and selecting the daily workout", func() {
			policy := buildDailyWorkoutPolicy("merge")
			dailyWorkout := policy.GetDailyWorkout(workouts, 30)

			Convey("Then workouts are merged", func() {
				So(dailyWorkout, ShouldNotBeNil)
				So(dailyWorkout.ID, ShouldEqual, int64(1))
				So(dailyWorkout.DistanceInKm, ShouldEqual, 30.0)
				So(dailyWorkout.DurationInMin, ShouldEqual, 70)
			})
		})
	})
}

func TestBuildDailyWorkoutPolicyDefaultLongest(t *testing.T) {
	Convey("Given an invalid daily workout policy", t, func() {
		workouts := []*domain.Workout{
			domain.NewWorkout(domain.WorkoutParams{ID: 1, DurationInMin: 30, DistanceInKm: 10}),
			domain.NewWorkout(domain.WorkoutParams{ID: 2, DurationInMin: 40, DistanceInKm: 20}),
		}

		Convey("When building and selecting the daily workout", func() {
			policy := buildDailyWorkoutPolicy("invalid")
			dailyWorkout := policy.GetDailyWorkout(workouts, 30)

			Convey("Then it uses longest workout behavior", func() {
				So(dailyWorkout, ShouldNotBeNil)
				So(dailyWorkout.ID, ShouldEqual, int64(2))
				So(dailyWorkout.DistanceInKm, ShouldEqual, 20.0)
				So(dailyWorkout.DurationInMin, ShouldEqual, 40)
			})
		})
	})
}

func TestNewAppSuccess(t *testing.T) {
	Convey("Given a valid config with an existing valid token file", t, func() {
		tokenFile := createTempTokenFile(t, `{"access_token":"token","refresh_token":"refresh","expires_at":"2030-01-01T00:00:00Z"}`)
		defer func() { _ = os.Remove(tokenFile) }()

		cfg := &config.Config{
			TokenFilePath:      tokenFile,
			OutputFilePath:     filepath.Join(t.TempDir(), "workouts.csv"),
			StravaApiBaseUrl:   "https://example.com/api/v3",
			StravaOauthBaseUrl: "https://example.com/oauth",
		}

		Convey("When NewApp is called", func() {
			app, err := NewApp(cfg, "longest")

			Convey("Then the app is initialized", func() {
				So(err, ShouldBeNil)
				So(app, ShouldNotBeNil)
			})
		})
	})
}

func TestAppRunDelegatesToUseCase(t *testing.T) {
	Convey("Given an App with SaveWorkoutPeriod use case dependencies", t, func() {
		spyPolicy := &spyDailyWorkoutPolicy{}
		spyRepo := &spyWorkoutRepository{}
		provider := &stubWorkoutProvider{
			Result: []*domain.Workout{
				domain.NewWorkout(domain.WorkoutParams{ID: 10, DurationInMin: 90, StartTime: time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC)}),
			},
		}
		athleteProvider := &stubAthleteProvider{Athlete: athlete}
		uc := usecase.NewSaveWorkoutPeriod(spyPolicy, spyRepo, provider, athleteProvider)

		period, _ := domain.NewPeriod(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
		request := &input.SaveWorkoutPeriodRequest{Period: period, MinimalWorkoutDuration: 45}

		sender := &spyReportSender{}
		sendUseCase := usecase.NewSendWorkoutReport(sender)
		app := &App{SaveWorkoutPeriod: *uc, SendWorkoutReport: *sendUseCase}

		Convey("When Run is called", func() {
			err := app.Run(request, testWorkoutReportPath)

			Convey("Then it executes SaveWorkoutPeriod with request values", func() {
				So(err, ShouldBeNil)
				So(provider.Called, ShouldEqual, 1)
				So(spyPolicy.Called, ShouldEqual, 1)
				So(spyPolicy.MinWorkoutDuration, ShouldEqual, 45)
				So(spyRepo.Called, ShouldEqual, 1)
				So(spyRepo.Workouts, ShouldHaveLength, 1)
				So(spyRepo.Workouts[0].ID, ShouldEqual, int64(10))
			})
		})
	})
}

func TestAppRunPropagatesError(t *testing.T) {
	Convey("Given an App where the workout provider returns an error", t, func() {
		spyPolicy := &spyDailyWorkoutPolicy{}
		spyRepo := &spyWorkoutRepository{}
		provider := &stubWorkoutProvider{Err: errors.New("provider failure")}
		athleteProvider := &stubAthleteProvider{Athlete: athlete}

		uc := usecase.NewSaveWorkoutPeriod(spyPolicy, spyRepo, provider, athleteProvider)
		app := &App{SaveWorkoutPeriod: *uc}

		period, _ := domain.NewPeriod(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
		request := &input.SaveWorkoutPeriodRequest{Period: period, MinimalWorkoutDuration: 45}

		Convey("When Run is called", func() {
			err := app.Run(request, testWorkoutReportPath)

			Convey("Then it propagates the use case error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "provider failure")
				So(provider.Called, ShouldEqual, 1)
				So(spyPolicy.Called, ShouldEqual, 0)
				So(spyRepo.Called, ShouldEqual, 0)
			})
		})
	})
}

func TestAppRunCallsSenderWithProvidedReportPath(t *testing.T) {
	Convey("Given an App with successful save and report sender use cases", t, func() {
		spyPolicy := &spyDailyWorkoutPolicy{}
		spyRepo := &spyWorkoutRepository{}
		provider := &stubWorkoutProvider{
			Result: []*domain.Workout{
				domain.NewWorkout(domain.WorkoutParams{ID: 10, DurationInMin: 90, StartTime: time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC)}),
			},
		}
		athleteProvider := &stubAthleteProvider{Athlete: athlete}
		saveUseCase := usecase.NewSaveWorkoutPeriod(spyPolicy, spyRepo, provider, athleteProvider)

		sender := &spyReportSender{}
		sendUseCase := usecase.NewSendWorkoutReport(sender)

		app := &App{SaveWorkoutPeriod: *saveUseCase, SendWorkoutReport: *sendUseCase}

		period, _ := domain.NewPeriod(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
		request := &input.SaveWorkoutPeriodRequest{Period: period, MinimalWorkoutDuration: 45}

		Convey("When Run is called", func() {
			err := app.Run(request, testWorkoutReportPath)

			Convey("Then it should send the report using the provided path", func() {
				So(err, ShouldBeNil)
				So(spyRepo.Called, ShouldEqual, 1)
				So(sender.Called, ShouldEqual, 1)
				So(sender.ReportPath, ShouldEqual, testWorkoutReportPath)
			})
		})
	})
}

func TestAppRunDoesNotCallSenderWhenSaveFails(t *testing.T) {
	Convey("Given an App where save use case returns an error", t, func() {
		spyPolicy := &spyDailyWorkoutPolicy{}
		spyRepo := &spyWorkoutRepository{}
		provider := &stubWorkoutProvider{Err: errors.New("provider failure")}
		athleteProvider := &stubAthleteProvider{Athlete: athlete}
		saveUseCase := usecase.NewSaveWorkoutPeriod(spyPolicy, spyRepo, provider, athleteProvider)

		sender := &spyReportSender{}
		sendUseCase := usecase.NewSendWorkoutReport(sender)

		app := &App{SaveWorkoutPeriod: *saveUseCase, SendWorkoutReport: *sendUseCase}

		period, _ := domain.NewPeriod(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
		request := &input.SaveWorkoutPeriodRequest{Period: period, MinimalWorkoutDuration: 45}

		Convey("When Run is called", func() {
			err := app.Run(request, testWorkoutReportPath)

			Convey("Then sender should not be called", func() {
				So(err, ShouldNotBeNil)
				So(sender.Called, ShouldEqual, 0)
			})
		})
	})
}

func TestAppRunPropagatesSenderError(t *testing.T) {
	Convey("Given an App where sender use case returns an error", t, func() {
		spyPolicy := &spyDailyWorkoutPolicy{}
		spyRepo := &spyWorkoutRepository{}
		provider := &stubWorkoutProvider{
			Result: []*domain.Workout{
				domain.NewWorkout(domain.WorkoutParams{ID: 10, DurationInMin: 90, StartTime: time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC)}),
			},
		}
		athleteProvider := &stubAthleteProvider{Athlete: athlete}
		saveUseCase := usecase.NewSaveWorkoutPeriod(spyPolicy, spyRepo, provider, athleteProvider)

		sender := &spyReportSender{Err: errors.New("send failed")}
		sendUseCase := usecase.NewSendWorkoutReport(sender)

		app := &App{SaveWorkoutPeriod: *saveUseCase, SendWorkoutReport: *sendUseCase}

		period, _ := domain.NewPeriod(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
		request := &input.SaveWorkoutPeriodRequest{Period: period, MinimalWorkoutDuration: 45}

		Convey("When Run is called", func() {
			err := app.Run(request, testWorkoutReportPath)

			Convey("Then it should propagate the sender error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "send failed")
			})
		})
	})
}

func createTempTokenFile(t *testing.T, content string) string {
	t.Helper()

	file, err := os.CreateTemp("", "token-*.json")
	if err != nil {
		t.Fatalf("failed to create temp token file: %v", err)
	}

	if _, err := file.WriteString(content); err != nil {
		_ = file.Close()
		t.Fatalf("failed to write temp token file: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("failed to close temp token file: %v", err)
	}

	return file.Name()
}
