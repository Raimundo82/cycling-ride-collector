//go:build integration

package integration

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/raimundo82/cycling-ride-collector/internal/application/usecase"
	"github.com/raimundo82/cycling-ride-collector/internal/config"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastucture/csv"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastucture/strava"
)

func TestStravaToCSVIntegration(t *testing.T) {
	Convey("Given a mock Strava API server", t, func() {
		// Setup mock Strava server
		startDate, _ := time.Parse("2006-01-02", "2026-02-10")
		endDate, _ := time.Parse("2006-01-02", "2026-02-10")
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check authentication header
			if r.Header.Get("Authorization") == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.Header().Set("Content-Type", "application/json")

			// Mock athlete activities endpoint
			if strings.HasPrefix(r.URL.Path, "/athlete/activities") && !strings.Contains(r.URL.Path, "/activities/") {
				activities := `[{
					"id": 1,
					"sport_type": "Ride",
					"commute": false,
					"trainer": false,
					"workout_type": 0,
					"start_date_local": "2026-02-10T08:00:00Z",
					"distance": 25000.0,
					"moving_time": 3600,
					"total_elevation_gain": 250.0,
					"average_watts": 180.0,
					"has_heartrate": true,
					"device_watts": true,
					"weighted_average_watts": 185.0,
					"average_heartrate": 145.0,
					"max_heartrate": 170.0,
					"average_cadence": 85.0
				}]`
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(activities))
				return
			}

			// Mock detailed activity endpoint (for leg sensations)
			if strings.Contains(r.URL.Path, "/activities") && !strings.Contains(r.URL.Path, "/streams") {
				detailedActivity := `{
					"id": 1,
					"private_note": "Boas"
				}`
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(detailedActivity))
				return
			}

			// Mock watts stream endpoint
			if strings.HasSuffix(r.URL.Path, "/streams") {
				wattsStream := `{
					"watts": {
						"data": [
							200,200,200,200,200,200,200,200,200,200,200,200,200,200,200,200,200,200,
							200,200,200,200,200,200,200,200,200,200,200,200,200,200,200,200,200,200,
							200,200,200,200,200,200,200,200,200,200,200,200,200,200,200,200,200,200,
							200,200,200,200,200,200,200,200,200,200,200,200,200,200,200,200,200,200
						]
					}
				}`
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(wattsStream))
				return
			}

			w.WriteHeader(http.StatusNotFound)
		}))
		defer mockServer.Close()

		// Setup temporary directory for CSV output
		tmpDir, err := os.MkdirTemp("", "integration-test-*")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tmpDir)

		Convey("When fetching activities and saving to CSV", func() {
			// Setup config pointing to mock server
			cfg := &config.Config{
				StravaApiBaseUrl:  mockServer.URL,
				StravaAccessToken: "test-token",
				OutputFilePath:    filepath.Join(tmpDir, "workouts.csv"),
			}

			// Setup components
			httpClient := &http.Client{Timeout: 10 * time.Second}
			stravaClient := strava.NewHttpClient(httpClient, cfg)
			workoutProvider := strava.NewProvider(stravaClient)
			workoutRepository := csv.NewCSVWorkoutPeriodSaver(cfg.OutputFilePath)
			dailyWorkoutPolicy := usecase.NewLongestWorkout()

			// Setup use case
			saveWorkoutUseCase := usecase.NewSaveWorkoutPeriod(
				dailyWorkoutPolicy,
				workoutRepository,
				workoutProvider,
			)

			// Create period for our test data
			period, _ := domain.NewPeriod(startDate, endDate)

			// Execute use case
			err := saveWorkoutUseCase.Execute(period, 30)

			Convey("Then the CSV file should be created with correct data", func() {
				So(err, ShouldBeNil)

				// Verify CSV file was created
				_, statErr := os.Stat(cfg.OutputFilePath)
				So(statErr, ShouldBeNil)

				// Read and verify CSV content
				csvContent, readErr := os.ReadFile(cfg.OutputFilePath)
				So(readErr, ShouldBeNil)
				So(len(csvContent), ShouldBeGreaterThan, 1)

				csvStr := string(csvContent)
				// Verify it contains the activity data
				So(csvStr, ShouldEqual, "2/10/2026,Estrada,08:00,1h0m,25.00,250,180,200,145,170,85,Boas\n")
			})
		})
	})
}
