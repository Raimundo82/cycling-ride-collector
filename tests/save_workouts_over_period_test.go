package test

import (
	"errors"
	"testing"
	"time"

	. "github.com/raimundo82/go-strava-weekly/internal/application/contracts"
	. "github.com/raimundo82/go-strava-weekly/internal/application/orchestration"
	. "github.com/smartystreets/goconvey/convey"
)

type MockSaveWorkoutUseCase struct {
	ExecuteCalled int
	ExecuteError  error
}

// Execute implements [contracts.SaveWorkoutUseCase].
func (m *MockSaveWorkoutUseCase) Execute(date time.Time, minWorkoutDuration int) error {
	m.ExecuteCalled++
	return m.ExecuteError
}

var _ SaveWorkoutUseCase = (*MockSaveWorkoutUseCase)(nil)

func TestSaveWorkoutsOverPeriodOf3DaysOnTheSameMonth(t *testing.T) {
	Convey("Given a SaveWorkoutsOrchestrator with a mock SaveWorkoutUseCase", t, func() {
		mockSaveWorkoutUseCase := &MockSaveWorkoutUseCase{
			ExecuteCalled: 0,
			ExecuteError:  nil,
		}
		orchestrator := &SaveWorkoutsOrchestrator{
			SaveWorkoutUseCase: mockSaveWorkoutUseCase,
		}
		Convey("When SaveWorkoutsOverPeriod is called for a period of 3 days", func() {
			startDate := NewDate(2024, 1, 1)
			endDate := NewDate(2024, 1, 3)
			period := NewPeriod(startDate, endDate)
			minWorkoutDuration := 30
			err := orchestrator.SaveWorkoutsOverPeriod(period, minWorkoutDuration)
			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
			})
			Convey("And Period should be valid", func() {
				So(period.IsValid(), ShouldBeTrue)
			})
			Convey("And Execute should be called 3 times", func() {
				So(mockSaveWorkoutUseCase.ExecuteCalled, ShouldEqual, 3)
			})
		})
	})
}

func TestSaveWorkoutsOverPeriodOfSingleDay(t *testing.T) {
	Convey("Given a SaveWorkoutsOrchestrator with a mock SaveWorkoutUseCase", t, func() {
		mockSaveWorkoutUseCase := &MockSaveWorkoutUseCase{
			ExecuteCalled: 0,
			ExecuteError:  nil,
		}
		orchestrator := &SaveWorkoutsOrchestrator{
			SaveWorkoutUseCase: mockSaveWorkoutUseCase,
		}
		Convey("When SaveWorkoutsOverPeriod is called for a period of a single day", func() {
			startDate := NewDate(2024, 1, 1)
			endDate := NewDate(2024, 1, 1)
			period := NewPeriod(startDate, endDate)
			minWorkoutDuration := 30
			err := orchestrator.SaveWorkoutsOverPeriod(period, minWorkoutDuration)
			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
			})
			Convey("And Period should be valid", func() {
				So(period.IsValid(), ShouldBeTrue)
			})
			Convey("And Execute should be called 1 times", func() {
				So(mockSaveWorkoutUseCase.ExecuteCalled, ShouldEqual, 1)
			})
		})
	})
}

func TestSaveWorkoutsOverPeriodCrossingMonths(t *testing.T) {
	Convey("Given a SaveWorkoutsOrchestrator with a mock SaveWorkoutUseCase", t, func() {
		mockSaveWorkoutUseCase := &MockSaveWorkoutUseCase{
			ExecuteCalled: 0,
			ExecuteError:  nil,
		}
		orchestrator := &SaveWorkoutsOrchestrator{
			SaveWorkoutUseCase: mockSaveWorkoutUseCase,
		}
		Convey("When SaveWorkoutsOverPeriod is called for a period crossing months", func() {
			startDate := NewDate(2024, 1, 31)
			endDate := NewDate(2024, 2, 2)
			period := NewPeriod(startDate, endDate)
			minWorkoutDuration := 30
			err := orchestrator.SaveWorkoutsOverPeriod(period, minWorkoutDuration)
			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
			})
			Convey("And Period should be valid", func() {
				So(period.IsValid(), ShouldBeTrue)
			})
			Convey("And Execute should be called 3 times", func() {
				So(mockSaveWorkoutUseCase.ExecuteCalled, ShouldEqual, 3)
			})
		})
	})
}

func TestSaveWorkoutsOverPeriodCrossingYears(t *testing.T) {
	Convey("Given a SaveWorkoutsOrchestrator with a mock SaveWorkoutUseCase", t, func() {
		mockSaveWorkoutUseCase := &MockSaveWorkoutUseCase{
			ExecuteCalled: 0,
			ExecuteError:  nil,
		}
		orchestrator := &SaveWorkoutsOrchestrator{
			SaveWorkoutUseCase: mockSaveWorkoutUseCase,
		}
		Convey("When SaveWorkoutsOverPeriod is called for a period crossing years", func() {
			startDate := NewDate(2024, 12, 31)
			endDate := NewDate(2025, 1, 2)
			period := NewPeriod(startDate, endDate)
			minWorkoutDuration := 30
			err := orchestrator.SaveWorkoutsOverPeriod(period, minWorkoutDuration)
			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
			})
			Convey("And Period should be valid", func() {
				So(period.IsValid(), ShouldBeTrue)
			})
			Convey("And Execute should be called 3 times", func() {
				So(mockSaveWorkoutUseCase.ExecuteCalled, ShouldEqual, 3)
			})
		})
	})
}

func TestSaveWorkoutsOverInvalidPeriod(t *testing.T) {
	Convey("Given a SaveWorkoutsOrchestrator with a mock SaveWorkoutUseCase", t, func() {
		mockSaveWorkoutUseCase := &MockSaveWorkoutUseCase{
			ExecuteCalled: 0,
			ExecuteError:  nil,
		}
		orchestrator := &SaveWorkoutsOrchestrator{
			SaveWorkoutUseCase: mockSaveWorkoutUseCase,
		}
		Convey("When SaveWorkoutsOverPeriod is called for an invalid period", func() {
			startDate := NewDate(2025, 1, 2)
			endDate := NewDate(2024, 12, 31)
			period := NewPeriod(startDate, endDate)
			minWorkoutDuration := 30
			err := orchestrator.SaveWorkoutsOverPeriod(period, minWorkoutDuration)
			Convey("Then Period should be invalid", func() {
				So(period.IsValid(), ShouldBeFalse)
			})
			Convey("And invalid period error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
			Convey("And Execute should be called 0 times", func() {
				So(mockSaveWorkoutUseCase.ExecuteCalled, ShouldEqual, 0)
			})
		})
	})
}

func TestSaveWorkoutsOverPeriodWithPropagatedErrors(t *testing.T) {
	Convey("Given a SaveWorkoutsOrchestrator with a mock SaveWorkoutUseCase", t, func() {
		mockSaveWorkoutUseCase := &MockSaveWorkoutUseCase{
			ExecuteCalled: 0,
			ExecuteError:  errors.New("mock error"),
		}
		orchestrator := &SaveWorkoutsOrchestrator{
			SaveWorkoutUseCase: mockSaveWorkoutUseCase,
		}
		Convey("When SaveWorkoutsOverPeriod is called for an invalid period", func() {
			startDate := NewDate(2025, 1, 1)
			endDate := NewDate(2025, 1, 3)
			period := NewPeriod(startDate, endDate)
			minWorkoutDuration := 30
			err := orchestrator.SaveWorkoutsOverPeriod(period, minWorkoutDuration)
			Convey("Then Period should be invalid", func() {
				So(period.IsValid(), ShouldBeTrue)
			})
			Convey("And multiple errors should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "mock error\nmock error\nmock error")
			})
			Convey("And Execute should be called 3 times", func() {
				So(mockSaveWorkoutUseCase.ExecuteCalled, ShouldEqual, 3)
			})
		})
	})
}
