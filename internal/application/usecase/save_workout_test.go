package usecase

import (
	"fmt"
	"testing"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
	. "github.com/raimundo82/go-strava-weekly/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

type mockWorkoutRepository struct {
	Workouts   []*Workout
	SaveCalled int
}

var _ contracts.WorkoutSaver = (*mockWorkoutRepository)(nil)

type mockWorkoutProvider struct {
	Workouts                []*Workout
	GetWorkoutsByDateCalled int
	Err                     error
}

var _ contracts.SingleWorkoutProvider = (*mockWorkoutProvider)(nil)

func (m *mockWorkoutRepository) Save(workout *Workout) error {
	m.Workouts = append(m.Workouts, workout)
	m.SaveCalled++
	return nil
}

func (m *mockWorkoutProvider) GetWorkoutsByDate(date time.Time) ([]*Workout, error) {
	m.GetWorkoutsByDateCalled++
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Workouts, nil
}

var (
	testWorkoutShort     = &Workout{DurationInMin: 15, DistanceInKm: 5.0}
	testWorkoutLong      = &Workout{DurationInMin: 120, DistanceInKm: 20.0}
	otherTestWorkoutLong = &Workout{DurationInMin: 60, DistanceInKm: 10.0}
	emptyWorkouts        = []*Workout{}
)

func TestSaveWorkout(t *testing.T) {
	testCases := []struct {
		name                   string
		workouts               []*Workout
		saveCalls              int
		getWorkoutByDateCalled int
		getDailyWorkoutCalled  int
	}{
		{"NoWorkouts", emptyWorkouts, 1, 1, 0},
		{"SingleShort", []*Workout{testWorkoutShort}, 1, 1, 1},
		{"SingleLong", []*Workout{testWorkoutLong}, 1, 1, 1},
		{"MultipleShorts", []*Workout{testWorkoutShort, testWorkoutShort}, 1, 1, 1},
		{"MultipleLongs", []*Workout{testWorkoutLong, otherTestWorkoutLong}, 1, 1, 1},
		{"MultipleMixed", []*Workout{testWorkoutShort, testWorkoutLong, otherTestWorkoutLong}, 1, 1, 1},
	}

	for _, tc := range testCases {
		Convey(tc.name, t, func() {
			mockRepo := &mockWorkoutRepository{
				Workouts:   []*Workout{},
				SaveCalled: 0,
			}

			mockProvider := &mockWorkoutProvider{
				Workouts:                tc.workouts,
				GetWorkoutsByDateCalled: 0,
			}

			mockDailyWorkoutPolicy := &mockDailyWorkoutPolicy{
				GetDailyWorkoutCalled: 0,
				Workout:               nil,
			}

			useCase := NewSaveWorkout(mockDailyWorkoutPolicy, mockRepo, mockProvider)
			Convey("When Execute is called", func() {
				date := time.Date(2023, time.June, 1, 0, 0, 0, 0, time.UTC)
				err := useCase.Execute(date, 30)

				Convey("Then asserting saves match expectations", func() {
					So(err, ShouldBeNil)
					So(mockProvider.GetWorkoutsByDateCalled, ShouldEqual, tc.getWorkoutByDateCalled)
					So(mockDailyWorkoutPolicy.GetDailyWorkoutCalled, ShouldEqual, tc.getDailyWorkoutCalled)
					So(mockRepo.SaveCalled, ShouldEqual, tc.saveCalls)
				})
			})
		})
	}
}

func TestSaveWorkout_WithNoWorkouts(t *testing.T) {
	Convey("Given a SaveWorkout use case with no workouts from provider", t, func() {
		mockRepo := &mockWorkoutRepository{
			Workouts:   []*Workout{},
			SaveCalled: 0,
		}

		mockProvider := &mockWorkoutProvider{
			Workouts:                []*Workout{},
			GetWorkoutsByDateCalled: 0,
		}

		mockDailyWorkoutPolicy := &mockDailyWorkoutPolicy{
			GetDailyWorkoutCalled: 0,
		}

		useCase := NewSaveWorkout(mockDailyWorkoutPolicy, mockRepo, mockProvider)

		Convey("When Execute is called", func() {
			date := time.Date(2023, time.June, 1, 0, 0, 0, 0, time.UTC)
			err := useCase.Execute(date, 30)

			Convey("Then no error occurs and no workout is saved", func() {
				So(err, ShouldBeNil)
				So(mockRepo.SaveCalled, ShouldEqual, 1)
				So(mockDailyWorkoutPolicy.GetDailyWorkoutCalled, ShouldEqual, 0)
				So(mockProvider.GetWorkoutsByDateCalled, ShouldEqual, 1)
				So(len(mockRepo.Workouts), ShouldEqual, 1)
			})
		})
	})
}

func TestSaveWorkout_WithLongWorkout(t *testing.T) {
	Convey("Given a SaveWorkout use case with a long workout from provider", t, func() {
		date := time.Date(2023, time.June, 1, 10, 0, 0, 0, time.UTC)

		workout := NewWorkout(&WorkoutParams{
			DurationInMin:     45,
			ID:                1,
			WorkoutType:       Estrada,
			StartTime:         date,
			DistanceInKm:      100,
			ElevationInMeters: 500,
			AvgCadenceInRpm:   90,
			MaxHeartRateInBpm: 180,
			AvgHeartRateInBpm: 140,
			AvgPowerInWatts:   225,
		})

		mockRepo := &mockWorkoutRepository{
			Workouts:   []*Workout{},
			SaveCalled: 0,
		}

		mockProvider := &mockWorkoutProvider{
			Workouts:                []*Workout{workout},
			GetWorkoutsByDateCalled: 0,
		}

		mockDailyWorkoutPolicy := &mockDailyWorkoutPolicy{
			GetDailyWorkoutCalled: 0,
			Workout:               workout,
		}

		useCase := NewSaveWorkout(mockDailyWorkoutPolicy, mockRepo, mockProvider)
		Convey("When Execute is called", func() {
			err := useCase.Execute(date, 30)

			Convey("Then no error occurs and a workout is saved", func() {
				So(err, ShouldBeNil)
				So(mockProvider.GetWorkoutsByDateCalled, ShouldEqual, 1)
				So(mockDailyWorkoutPolicy.GetDailyWorkoutCalled, ShouldEqual, 1)
				So(mockRepo.SaveCalled, ShouldEqual, 1)
				So(len(mockRepo.Workouts), ShouldEqual, 1)
				So(mockRepo.Workouts[0], ShouldEqual, mockProvider.Workouts[0])
			})
		})
	})
}

func TestSaveWorkout_WithShortWorkout(t *testing.T) {
	Convey("Given a SaveWorkout use case with a short workout from provider", t, func() {
		date := time.Date(2023, time.June, 1, 10, 0, 0, 0, time.UTC)
		workout := NewWorkout(&WorkoutParams{
			DurationInMin:     15,
			ID:                1,
			WorkoutType:       Estrada,
			StartTime:         date,
			DistanceInKm:      100,
			ElevationInMeters: 500,
			AvgCadenceInRpm:   90,
			MaxHeartRateInBpm: 180,
			AvgHeartRateInBpm: 140,
			AvgPowerInWatts:   225,
		})

		mockRepo := &mockWorkoutRepository{
			Workouts:   []*Workout{},
			SaveCalled: 0,
		}
		mockDailyWorkoutPolicy := &mockDailyWorkoutPolicy{
			GetDailyWorkoutCalled: 0,
			Workout:               workout,
		}
		mockProvider := &mockWorkoutProvider{
			Workouts:                []*Workout{workout},
			GetWorkoutsByDateCalled: 0,
		}

		useCase := NewSaveWorkout(mockDailyWorkoutPolicy, mockRepo, mockProvider)

		Convey("When Execute is called", func() {
			err := useCase.Execute(date, 30)

			Convey("Then no error occurs and no workout is saved", func() {
				So(err, ShouldBeNil)
				So(mockProvider.GetWorkoutsByDateCalled, ShouldEqual, 1)
				So(mockRepo.SaveCalled, ShouldEqual, 1)
				So(mockDailyWorkoutPolicy.GetDailyWorkoutCalled, ShouldEqual, 1)
				So(len(mockRepo.Workouts), ShouldEqual, 1)
			})
		})
	})
}

func TestSaveWorkout_WithProviderError(t *testing.T) {
	Convey("Given a SaveWorkout use case with a provider that returns an error", t, func() {
		date := time.Date(2023, time.June, 1, 10, 0, 0, 0, time.UTC)
		mockRepo := &mockWorkoutRepository{
			Workouts:   []*Workout{},
			SaveCalled: 0,
		}

		mockDailyWorkoutPolicy := &mockDailyWorkoutPolicy{
			GetDailyWorkoutCalled: 0,
			Workout:               nil,
		}

		providerErr := fmt.Errorf("provider connection failed")
		mockProvider := &mockWorkoutProvider{
			Workouts:                nil,
			GetWorkoutsByDateCalled: 0,
			Err:                     providerErr,
		}

		useCase := NewSaveWorkout(mockDailyWorkoutPolicy, mockRepo, mockProvider)

		Convey("When Execute is called", func() {
			err := useCase.Execute(date, 30)

			Convey("Then the error is propagated and no workout is saved", func() {
				So(err, ShouldEqual, providerErr)
				So(mockProvider.GetWorkoutsByDateCalled, ShouldEqual, 1)
				So(mockDailyWorkoutPolicy.GetDailyWorkoutCalled, ShouldEqual, 0)
				So(mockRepo.SaveCalled, ShouldEqual, 0)
				So(len(mockRepo.Workouts), ShouldEqual, 0)
			})
		})
	})
}
