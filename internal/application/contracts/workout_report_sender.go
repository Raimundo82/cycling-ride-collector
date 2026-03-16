package contracts

type WorkoutReportSender interface {
	Send(reportPath string) error
}
