package scheduler

import (
	"errors"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	. "github.com/smartystreets/goconvey/convey"
)

type fakeRunner struct {
	spec          string
	addFuncCalled bool
	started       bool
	addErr        error
}

// AddFunc implements [schedulerRunner].
func (f *fakeRunner) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	f.addFuncCalled = true
	f.spec = spec
	if f.addErr != nil {
		return 0, f.addErr
	}
	return 1, nil
}

// Start implements [schedulerRunner].
func (f *fakeRunner) Start() {
	f.started = true
}

var _ schedulerRunner = (*fakeRunner)(nil)

func TestShouldCreateSchedulerWithCronRunnerWhenNewIsCalled(t *testing.T) {
	Convey("Given no dependencies", t, func() {
		Convey("When New is called", func() {
			s := New()

			Convey("Then it should return a scheduler with initialized runner", func() {
				So(s, ShouldNotBeNil)
				So(s.s, ShouldNotBeNil)
			})
		})
	})
}

func TestShouldStartSchedulerWhenAddFuncSucceeds(t *testing.T) {
	Convey("Given a scheduler with a valid runner", t, func() {
		runner := &fakeRunner{}
		s := &Scheduler{s: runner}

		Convey("When StartWeeklySunday20 is called", func() {
			err := s.StartWeeklySunday20(func() { /* no-op job */ })

			Convey("Then it should add Sunday-20 cron and start scheduler", func() {
				So(err, ShouldBeNil)
				So(runner.addFuncCalled, ShouldBeTrue)
				So(runner.spec, ShouldEqual, "0 20 * * 0")
				So(runner.started, ShouldBeTrue)
			})
		})
	})
}

func TestShouldReturnErrorWhenAddFuncFails(t *testing.T) {
	Convey("Given a scheduler with failing AddFunc", t, func() {
		runner := &fakeRunner{addErr: errors.New("invalid cron spec")}
		s := &Scheduler{s: runner}

		Convey("When StartWeeklySunday20 is called", func() {
			err := s.StartWeeklySunday20(func() { /* no-op job */ })

			Convey("Then it should return error and not start scheduler", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid cron spec")
				So(runner.started, ShouldBeFalse)
			})
		})
	})
}

func TestShouldScheduleNextRunOnSundayAt20WhenStartWeeklySunday20IsCalled(t *testing.T) {
	Convey("Given a scheduler with real cron runner in UTC", t, func() {
		c := cron.New(cron.WithLocation(time.UTC))
		s := &Scheduler{s: c}

		Convey("When StartWeeklySunday20 is called", func() {
			err := s.StartWeeklySunday20(func() { /* no-op job */ })
			So(err, ShouldBeNil)

			entries := c.Entries()
			Convey("Then next execution should be on Sunday at 20:00", func() {
				So(entries, ShouldHaveLength, 1)

				next := entries[0].Next.UTC()
				So(next.Weekday(), ShouldEqual, time.Sunday)
				So(next.Hour(), ShouldEqual, 20)
				So(next.Minute(), ShouldEqual, 0)
			})
		})
	})
}
