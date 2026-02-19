package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	. "github.com/raimundo82/cycling-ride-collector/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

type mockWorkoutPeriodSaver struct {
	Workouts      []*Workout
	Athlete       *Athlete
	SaveAllCalled int
}

var _ contracts.WorkoutRepository = (*mockWorkoutPeriodSaver)(nil)

type mockPeriodWorkoutProvider struct {
	Workouts                  []*Workout
	GetWorkoutsByPeriodCalled int
	Err                       error
}

var _ contracts.WorkoutProvider = (*mockPeriodWorkoutProvider)(nil)

type mockAthleteDataProvider struct {
	AthleteData          *Athlete
	GetAthleteDataCalled int
	Err                  error
}

var _ contracts.AthleteDataProvider = (*mockAthleteDataProvider)(nil)

// SaveAll implements [contracts.WorkoutRepository].
func (m *mockWorkoutPeriodSaver) SaveAll(workouts []*Workout, athlete *Athlete) error {
	m.Workouts = append(m.Workouts, workouts...)
	m.Athlete = athlete
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

// GetAthleteData implements [contracts.AthleteDataProvider].
func (m *mockAthleteDataProvider) GetAthleteData() (*Athlete, error) {
	m.GetAthleteDataCalled++
	if m.Err != nil {
		return nil, m.Err
	}
	return m.AthleteData, nil
}

func TestSaveWorkoutPeriodSavesWorkoutsAndAthleteData(t *testing.T) {
	Convey("Given workouts for the period and an Athlete", t, func() {
		a := NewAthlete(79, 135, 250)
		athleteDataProvider := &mockAthleteDataProvider{AthleteData: a, Err: nil}
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

		saveWorkoutPeriod := NewSaveWorkoutPeriod(NewLongestWorkout(), workoutRepo, workoutProvider, athleteDataProvider)

		period, _ := NewPeriod(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 3, 0, 0, 0, 0, time.UTC))

		minWorkoutDuration := 30

		Convey("When execute", func() {
			err := saveWorkoutPeriod.Execute(period, minWorkoutDuration)

			Convey("Then 3 workouts and athlete data are saved", func() {
				So(athleteDataProvider.GetAthleteDataCalled, ShouldEqual, 1)
				So(workoutRepo.Athlete, ShouldResemble, a)
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

func TestSaveWorkoutPeriodGivenWorkoutsForOneDayPeriodWhenExecuteThenWorkoutIsSavedForOneDay(t *testing.T) {
	Convey("Given workouts for one day and an Athlete", t, func() {
		a := NewAthlete(79, 135, 250)
		athleteDataProvider := &mockAthleteDataProvider{AthleteData: a, Err: nil}
		workoutProvider := &mockPeriodWorkoutProvider{
			Workouts: []*Workout{
				{ID: 1, StartTime: time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC), DurationInMin: 90},
				{ID: 2, StartTime: time.Date(2024, 6, 2, 9, 0, 0, 0, time.UTC), DurationInMin: 30},
			},
			GetWorkoutsByPeriodCalled: 0,
		}

		workoutRepo := &mockWorkoutPeriodSaver{Workouts: []*Workout{}, SaveAllCalled: 0}

		saveWorkoutPeriod := NewSaveWorkoutPeriod(NewLongestWorkout(), workoutRepo, workoutProvider, athleteDataProvider)

		period, _ := NewPeriod(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))

		minWorkoutDuration := 30

		Convey("When execute", func() {
			err := saveWorkoutPeriod.Execute(period, minWorkoutDuration)

			Convey("Then the workout and athlete data are saved", func() {
				So(athleteDataProvider.GetAthleteDataCalled, ShouldEqual, 1)
				So(workoutRepo.Athlete, ShouldResemble, a)
				So(workoutProvider.GetWorkoutsByPeriodCalled, ShouldEqual, 1)
				So(athleteDataProvider.GetAthleteDataCalled, ShouldEqual, 1)
				So(err, ShouldBeNil)
				So(workoutRepo.SaveAllCalled, ShouldEqual, 1)
				So(workoutRepo.Workouts, ShouldHaveLength, 1)
				So(workoutRepo.Workouts[0].ID, ShouldEqual, 1)
			})
		})
	})
}

func TestSaveWorkoutPeriodGivenNoWorkoutsWhenExecuteThenRestWorkoutsAreSaved(t *testing.T) {
	Convey("Given no workouts and an Athlete", t, func() {
		a := NewAthlete(79, 135, 250)
		athleteDataProvider := &mockAthleteDataProvider{AthleteData: a, Err: nil}
		workoutProvider := &mockPeriodWorkoutProvider{Workouts: []*Workout{}, GetWorkoutsByPeriodCalled: 0}
		workoutRepo := &mockWorkoutPeriodSaver{Workouts: []*Workout{}, SaveAllCalled: 0}
		saveWorkoutPeriod := NewSaveWorkoutPeriod(NewLongestWorkout(), workoutRepo, workoutProvider, athleteDataProvider)
		period, _ := NewPeriod(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 7, 0, 0, 0, 0, time.UTC))
		minWorkoutDuration := 30

		Convey("When execute", func() {
			err := saveWorkoutPeriod.Execute(period, minWorkoutDuration)

			Convey("Then rest workouts are saved", func() {
				So(athleteDataProvider.GetAthleteDataCalled, ShouldEqual, 1)
				So(workoutRepo.Athlete, ShouldResemble, a)
				So(workoutProvider.GetWorkoutsByPeriodCalled, ShouldEqual, 1)
				So(err, ShouldBeNil)
				So(workoutRepo.SaveAllCalled, ShouldEqual, 1)
				So(workoutRepo.Workouts, ShouldHaveLength, 7)
			})
		})
	})
}

func TestSaveWorkoutPeriodGivenActivityProviderErrorWhenExecuteThenReturnsError(t *testing.T) {
	Convey("Given error from the activity provider", t, func() {
		a := NewAthlete(79, 135, 250)
		athleteDataProvider := &mockAthleteDataProvider{AthleteData: a, Err: nil}
		workoutProvider := &mockPeriodWorkoutProvider{Workouts: []*Workout{}, GetWorkoutsByPeriodCalled: 0, Err: errors.New("provider error")}
		workoutRepo := &mockWorkoutPeriodSaver{Workouts: []*Workout{}, SaveAllCalled: 0}
		saveWorkoutPeriod := NewSaveWorkoutPeriod(NewLongestWorkout(), workoutRepo, workoutProvider, athleteDataProvider)
		period, _ := NewPeriod(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 7, 0, 0, 0, 0, time.UTC))
		minWorkoutDuration := 30

		Convey("When execute", func() {
			err := saveWorkoutPeriod.Execute(period, minWorkoutDuration)

			Convey("Then Error is returned", func() {
				So(workoutProvider.GetWorkoutsByPeriodCalled, ShouldEqual, 1)
				So(err, ShouldNotBeNil)
				So(workoutRepo.SaveAllCalled, ShouldEqual, 0)
				So(workoutRepo.Workouts, ShouldHaveLength, 0)
				So(athleteDataProvider.GetAthleteDataCalled, ShouldEqual, 1)
				So(workoutRepo.Athlete, ShouldBeNil)
				So(err.Error(), ShouldContainSubstring, "provider error")
			})
		})
	})
}

func TestSaveWorkoutPeriodGivenAthleteProviderErrorWhenExecuteThenReturnsError(t *testing.T) {
	Convey("Given error from the athlete provider", t, func() {
		a := NewAthlete(79, 135, 250)
		athleteDataProvider := &mockAthleteDataProvider{AthleteData: a, Err: errors.New("provider error")}
		workoutProvider := &mockPeriodWorkoutProvider{Workouts: []*Workout{}, GetWorkoutsByPeriodCalled: 0}
		workoutRepo := &mockWorkoutPeriodSaver{Workouts: []*Workout{}, SaveAllCalled: 0}
		saveWorkoutPeriod := NewSaveWorkoutPeriod(NewLongestWorkout(), workoutRepo, workoutProvider, athleteDataProvider)
		period, _ := NewPeriod(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 7, 0, 0, 0, 0, time.UTC))
		minWorkoutDuration := 30

		Convey("When execute", func() {
			err := saveWorkoutPeriod.Execute(period, minWorkoutDuration)

			Convey("Then Error is returned", func() {
				So(workoutProvider.GetWorkoutsByPeriodCalled, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				So(workoutRepo.SaveAllCalled, ShouldEqual, 0)
				So(workoutRepo.Workouts, ShouldHaveLength, 0)
				So(athleteDataProvider.GetAthleteDataCalled, ShouldEqual, 1)
				So(workoutRepo.Athlete, ShouldBeNil)
			})
		})
	})
}
