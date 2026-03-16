package usecase

import (
	"errors"
	"testing"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	. "github.com/smartystreets/goconvey/convey"
)

type spyWorkoutReportSender struct {
	Called     int
	ReportPath string
	Err        error
}

const testReportPath = "/tmp/report.xlsx"

var _ contracts.WorkoutReportSender = (*spyWorkoutReportSender)(nil)

func (s *spyWorkoutReportSender) Send(reportPath string) error {
	s.Called++
	s.ReportPath = reportPath
	return s.Err
}

func TestSendWorkoutReportExecuteSuccess(t *testing.T) {
	Convey("Given a report sender and a report path", t, func() {
		sender := &spyWorkoutReportSender{}
		useCase := NewSendWorkoutReport(sender)

		Convey("When Execute is called", func() {
			err := useCase.Execute(testReportPath)

			Convey("Then it delegates to the sender", func() {
				So(err, ShouldBeNil)
				So(sender.Called, ShouldEqual, 1)
				So(sender.ReportPath, ShouldEqual, testReportPath)
			})
		})
	})
}

func TestSendWorkoutReportExecutePropagatesSenderError(t *testing.T) {
	Convey("Given a sender that returns an error", t, func() {
		sender := &spyWorkoutReportSender{Err: errors.New("gmail failed")}
		useCase := NewSendWorkoutReport(sender)

		Convey("When Execute is called", func() {
			err := useCase.Execute(testReportPath)

			Convey("Then it returns the wrapped sender error", func() {
				So(err, ShouldNotBeNil)
				So(sender.Called, ShouldEqual, 1)
				So(err.Error(), ShouldContainSubstring, "gmail failed")
			})
		})
	})
}
