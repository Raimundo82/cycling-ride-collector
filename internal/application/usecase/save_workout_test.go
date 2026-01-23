package usecase

import (
	"fmt"
	"testing"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

type MockWorkoutRepository struct {
	Workouts   []*domain.Workout
	SaveCalled int
}

type MockWorkoutProvider struct {
	Workouts         []*domain.Workout
	GetWorkoutCalled int
	Err              error
}

func (m *MockWorkoutRepository) Save(workout *domain.Workout) error {
	m.Workouts = append(m.Workouts, workout)
	m.SaveCalled++
	return nil
}

func (m *MockWorkoutProvider) GetWorkoutsByDate(unixDate int64) ([]*domain.Workout, error) {
	m.GetWorkoutCalled++
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Workouts, nil
}

var (
	testWorkoutShort     = &domain.Workout{Duration: 15, Distance: 5.0}
	testWorkoutLong      = &domain.Workout{Duration: 120, Distance: 20.0}
	otherTestWorkoutLong = &domain.Workout{Duration: 60, Distance: 10.0}
	emptyWorkouts        = []*domain.Workout{}
)

func TestSaveWorkout(t *testing.T) {
	testCases := []struct {
		name      string
		workouts  []*domain.Workout
		wantCalls int
		wantSaved int
	}{
		{"NoWorkouts", emptyWorkouts, 0, 0},
		{"SingleShort", []*domain.Workout{testWorkoutShort}, 0, 0},
		{"SingleLong", []*domain.Workout{testWorkoutLong}, 1, 1},
		{"MultipleShorts", []*domain.Workout{testWorkoutShort, testWorkoutShort}, 0, 0},
		{"MultipleLongs", []*domain.Workout{testWorkoutLong, otherTestWorkoutLong}, 1, 1},
		{"MultipleMixed", []*domain.Workout{testWorkoutShort, testWorkoutLong, otherTestWorkoutLong}, 1, 1},
	}

	for _, tc := range testCases {
		Convey(tc.name, t, func() {
			mockRepo := &MockWorkoutRepository{
				Workouts:   []*domain.Workout{},
				SaveCalled: 0,
			}

			mockProvider := &MockWorkoutProvider{
				Workouts:         tc.workouts,
				GetWorkoutCalled: 0,
			}

			useCase := &SaveWorkout{
				WorkoutRepo:     mockRepo,
				WorkoutProvider: mockProvider,
			}

			Convey("When Execute is called", func() {
				unixDate := int64(1622505600)
				err := useCase.Execute(unixDate, 30)

				Convey("Then no error occurs and no workout is saved", func() {
					So(err, ShouldBeNil)
					So(mockRepo.SaveCalled, ShouldEqual, tc.wantCalls)
					So(len(mockRepo.Workouts), ShouldEqual, tc.wantSaved)
				})
			})
		})
	}
}

func TestSaveWorkout_WithNoWorkouts(t *testing.T) {
	Convey("Given a SaveWorkout use case with no workouts from provider", t, func() {
		mockRepo := &MockWorkoutRepository{
			Workouts:   []*domain.Workout{},
			SaveCalled: 0,
		}

		mockProvider := &MockWorkoutProvider{
			Workouts:         []*domain.Workout{},
			GetWorkoutCalled: 0,
		}

		useCase := &SaveWorkout{
			WorkoutRepo:     mockRepo,
			WorkoutProvider: mockProvider,
		}

		Convey("When Execute is called", func() {
			unixDate := int64(1622505600)
			err := useCase.Execute(unixDate, 30)

			Convey("Then no error occurs and no workout is saved", func() {
				So(err, ShouldBeNil)
				So(mockRepo.SaveCalled, ShouldEqual, 0)
				So(len(mockRepo.Workouts), ShouldEqual, 0)
			})
		})
	})
}

func TestSaveWorkout_WithLongWorkout(t *testing.T) {
	Convey("Given a SaveWorkout use case with a long workout from provider", t, func() {
		mockRepo := &MockWorkoutRepository{
			Workouts:   []*domain.Workout{},
			SaveCalled: 0,
		}

		mockProvider := &MockWorkoutProvider{
			Workouts: []*domain.Workout{
				domain.NewWorkout(domain.WorkoutParams{
					Duration:     45,
					Id:           1,
					WorkoutType:  domain.Estrada,
					StartTime:    "10:00",
					Distance:     100,
					Elevation:    500,
					AvgCadence:   90,
					MaxHeartRate: 180,
					AvgHeartRate: 140,
					AvgPower:     225,
				}),
			},
			GetWorkoutCalled: 0,
		}

		useCase := &SaveWorkout{
			WorkoutRepo:     mockRepo,
			WorkoutProvider: mockProvider,
		}

		Convey("When Execute is called", func() {
			unixDate := int64(1622505600)
			err := useCase.Execute(unixDate, 30)

			Convey("Then no error occurs and a workout is saved", func() {
				So(err, ShouldBeNil)
				So(mockProvider.GetWorkoutCalled, ShouldEqual, 1)
				So(mockRepo.SaveCalled, ShouldEqual, 1)
				So(len(mockRepo.Workouts), ShouldEqual, 1)
				So(mockRepo.Workouts[0], ShouldEqual, mockProvider.Workouts[0])
			})
		})
	})
}

func TestSaveWorkout_WithShortWorkout(t *testing.T) {
	Convey("Given a SaveWorkout use case with a short workout from provider", t, func() {
		mockRepo := &MockWorkoutRepository{
			Workouts:   []*domain.Workout{},
			SaveCalled: 0,
		}

		mockProvider := &MockWorkoutProvider{
			Workouts: []*domain.Workout{
				domain.NewWorkout(domain.WorkoutParams{
					Duration:     15,
					Id:           1,
					WorkoutType:  domain.Estrada,
					StartTime:    "10:00",
					Distance:     100,
					Elevation:    500,
					AvgCadence:   90,
					MaxHeartRate: 180,
					AvgHeartRate: 140,
					AvgPower:     225,
				}),
			},
			GetWorkoutCalled: 0,
		}

		useCase := &SaveWorkout{
			WorkoutRepo:     mockRepo,
			WorkoutProvider: mockProvider,
		}

		Convey("When Execute is called", func() {
			unixDate := int64(1622505600)
			err := useCase.Execute(unixDate, 30)

			Convey("Then no error occurs and no workout is saved", func() {
				So(err, ShouldBeNil)
				So(mockProvider.GetWorkoutCalled, ShouldEqual, 1)
				So(mockRepo.SaveCalled, ShouldEqual, 0)
				So(len(mockRepo.Workouts), ShouldEqual, 0)
			})
		})
	})
}

func TestSaveWorkout_WithProviderError(t *testing.T) {
	Convey("Given a SaveWorkout use case with a provider that returns an error", t, func() {
		mockRepo := &MockWorkoutRepository{
			Workouts:   []*domain.Workout{},
			SaveCalled: 0,
		}

		providerErr := fmt.Errorf("provider connection failed")
		mockProvider := &MockWorkoutProvider{
			Workouts:         nil,
			GetWorkoutCalled: 0,
			Err:              providerErr,
		}

		useCase := &SaveWorkout{
			WorkoutRepo:     mockRepo,
			WorkoutProvider: mockProvider,
		}

		Convey("When Execute is called", func() {
			unixDate := int64(1622505600)
			err := useCase.Execute(unixDate, 30)

			Convey("Then the error is propagated and no workout is saved", func() {
				So(err, ShouldEqual, providerErr)
				So(mockProvider.GetWorkoutCalled, ShouldEqual, 1)
				So(mockRepo.SaveCalled, ShouldEqual, 0)
				So(len(mockRepo.Workouts), ShouldEqual, 0)
			})
		})
	})
}
