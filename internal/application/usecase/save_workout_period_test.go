package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
	. "github.com/raimundo82/go-strava-weekly/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

type mockWorkoutPeriodSaver struct {
	Workouts      []*Workout
	SaveAllCalled int
}

var _ contracts.WorkoutPeriodSaver = (*mockWorkoutPeriodSaver)(nil)

type mockPeriodWorkoutProvider struct {
	Workouts                  []*Workout
	GetWorkoutsByPeriodCalled int
	Err                       error
}

var _ contracts.PeriodWorkoutProvider = (*mockPeriodWorkoutProvider)(nil)

// SaveAll implements [contracts.WorkoutPeriodSaver].
func (m *mockWorkoutPeriodSaver) SaveAll(workouts []*Workout) error {
	m.Workouts = append(m.Workouts, workouts...)
	m.SaveAllCalled++
	return nil
}

// GetWorkoutsByPeriod implements [contracts.PeriodWorkoutProvider].
func (m *mockPeriodWorkoutProvider) GetWorkoutsByPeriod(period Period) ([]*Workout, error) {
	m.GetWorkoutsByPeriodCalled++
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Workouts, nil
}

func TestSaveWorkoutPeriod_GivenWorkoutsForPeriod_WhenExecute_ThenWorkoutsSaved(t *testing.T) {
	Convey("Given workouts for the period from the provider", t, func() {
		workoutProvider := &mockPeriodWorkoutProvider{
			Workouts: []*Workout{
				{ID: 1, StartTime: time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC), DurationInMin: 90},
				{ID: 2, StartTime: time.Date(2024, 6, 1, 15, 0, 0, 0, time.UTC), DurationInMin: 60},
				{ID: 3, StartTime: time.Date(2024, 6, 2, 9, 0, 0, 0, time.UTC), DurationInMin: 30},
				{ID: 4, StartTime: time.Date(2024, 6, 3, 18, 0, 0, 0, time.UTC), DurationInMin: 60},
				{ID: 5, StartTime: time.Date(2024, 6, 3, 20, 0, 0, 0, time.UTC), DurationInMin: 25},
			},
			GetWorkoutsByPeriodCalled: 0,
		}

		workoutRepo := &mockWorkoutPeriodSaver{Workouts: []*Workout{}, SaveAllCalled: 0}

		saveWorkoutPeriod := NewSaveWorkoutPeriod(NewLongestWorkout(), workoutRepo, workoutProvider)

		period, _ := NewPeriod(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 3, 0, 0, 0, 0, time.UTC))

		minWorkoutDuration := 30

		Convey("When execute", func() {
			err := saveWorkoutPeriod.Execute(period, minWorkoutDuration)

			Convey("Then 3 workouts are saved", func() {
				So(workoutProvider.GetWorkoutsByPeriodCalled, ShouldEqual, 1)
				So(err, ShouldBeNil)
				So(workoutRepo.SaveAllCalled, ShouldEqual, 1)
				So(workoutRepo.Workouts, ShouldHaveLength, 3)
				So(workoutRepo.Workouts[0].ID, ShouldEqual, 1)
				So(workoutRepo.Workouts[1].ID, ShouldEqual, 3)
				So(workoutRepo.Workouts[2].ID, ShouldEqual, 4)
			})
		})
	})
}

func TestSaveWorkoutPeriod_GivenWorkoutsForOneDayPeriod_WhenExecute_ThenWorkoutIsSavedForOneDay(t *testing.T) {
	Convey("Given workouts for one day period from the provider", t, func() {
		workoutProvider := &mockPeriodWorkoutProvider{
			Workouts: []*Workout{
				{ID: 1, StartTime: time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC), DurationInMin: 90},
				{ID: 2, StartTime: time.Date(2024, 6, 2, 9, 0, 0, 0, time.UTC), DurationInMin: 30},
			},
			GetWorkoutsByPeriodCalled: 0,
		}

		workoutRepo := &mockWorkoutPeriodSaver{Workouts: []*Workout{}, SaveAllCalled: 0}

		saveWorkoutPeriod := NewSaveWorkoutPeriod(NewLongestWorkout(), workoutRepo, workoutProvider)

		period, _ := NewPeriod(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))

		minWorkoutDuration := 30

		Convey("When execute", func() {
			err := saveWorkoutPeriod.Execute(period, minWorkoutDuration)

			Convey("Then 1 workout is saved", func() {
				So(workoutProvider.GetWorkoutsByPeriodCalled, ShouldEqual, 1)
				So(err, ShouldBeNil)
				So(workoutRepo.SaveAllCalled, ShouldEqual, 1)
				So(workoutRepo.Workouts, ShouldHaveLength, 1)
				So(workoutRepo.Workouts[0].ID, ShouldEqual, 1)
			})
		})
	})
}

func TestSaveWorkoutPeriod_GivenNoWorkouts_WhenExecute_ThenRestWorkoutsAreSaved(t *testing.T) {
	Convey("Given no workouts for the period from the provider", t, func() {
		workoutProvider := &mockPeriodWorkoutProvider{Workouts: []*Workout{}, GetWorkoutsByPeriodCalled: 0}
		workoutRepo := &mockWorkoutPeriodSaver{Workouts: []*Workout{}, SaveAllCalled: 0}
		saveWorkoutPeriod := NewSaveWorkoutPeriod(NewLongestWorkout(), workoutRepo, workoutProvider)
		period, _ := NewPeriod(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 7, 0, 0, 0, 0, time.UTC))
		minWorkoutDuration := 30

		Convey("When execute", func() {
			err := saveWorkoutPeriod.Execute(period, minWorkoutDuration)

			Convey("Then rest workouts are saved", func() {
				So(workoutProvider.GetWorkoutsByPeriodCalled, ShouldEqual, 1)
				So(err, ShouldBeNil)
				So(workoutRepo.SaveAllCalled, ShouldEqual, 1)
				So(workoutRepo.Workouts, ShouldHaveLength, 7)
			})
		})
	})
}

func TestSaveWorkoutPeriod_GivenProviderError_WhenExecute_ThenReturnsError(t *testing.T) {
	Convey("Given error from the provider", t, func() {
		workoutProvider := &mockPeriodWorkoutProvider{Workouts: []*Workout{}, GetWorkoutsByPeriodCalled: 0, Err: errors.New("provider error")}
		workoutRepo := &mockWorkoutPeriodSaver{Workouts: []*Workout{}, SaveAllCalled: 0}
		saveWorkoutPeriod := NewSaveWorkoutPeriod(NewLongestWorkout(), workoutRepo, workoutProvider)
		period, _ := NewPeriod(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 7, 0, 0, 0, 0, time.UTC))
		minWorkoutDuration := 30

		Convey("When execute", func() {
			err := saveWorkoutPeriod.Execute(period, minWorkoutDuration)

			Convey("Then Error is returned", func() {
				So(workoutProvider.GetWorkoutsByPeriodCalled, ShouldEqual, 1)
				So(err, ShouldNotBeNil)
				So(workoutRepo.SaveAllCalled, ShouldEqual, 0)
				So(workoutRepo.Workouts, ShouldHaveLength, 0)
			})
		})
	})
}
