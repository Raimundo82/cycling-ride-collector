package test

import (
	"testing"
	"time"

	. "github.com/raimundo82/go-strava-weekly/internal/application/contracts"
	. "github.com/raimundo82/go-strava-weekly/internal/application/orchestration"
	. "github.com/smartystreets/goconvey/convey"
)

type MockSaveWorkoutUseCase struct {
	ExecuteCalled int
}

// Execute implements [contracts.SaveWorkoutUseCase].
func (m *MockSaveWorkoutUseCase) Execute(date time.Time, minWorkoutDuration int) error {
	m.ExecuteCalled++
	return nil
}

var _ SaveWorkoutUseCase = (*MockSaveWorkoutUseCase)(nil)

func TestSaveWorkoutsOverPeriod(t *testing.T) {
	Convey("Given a SaveWorkoutsOrchestrator with a mock SaveWorkoutUseCase", t, func() {
		mockSaveWorkoutUseCase := &MockSaveWorkoutUseCase{
			ExecuteCalled: 0,
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
			Convey("And Execute should be called 3 times", func() {
				So(mockSaveWorkoutUseCase.ExecuteCalled, ShouldEqual, 3)
			})
		})
	})
}
