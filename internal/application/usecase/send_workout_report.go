package usecase

import (
	"fmt"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
)

type SendWorkoutReport struct {
	reportSender contracts.WorkoutReportSender
}

func NewSendWorkoutReport(reportSender contracts.WorkoutReportSender) *SendWorkoutReport {
	return &SendWorkoutReport{reportSender: reportSender}
}

func (s *SendWorkoutReport) Execute(reportPath string) error {
	if err := s.reportSender.Send(reportPath); err != nil {
		return fmt.Errorf("failed to send workout report: %w", err)
	}
	return nil
}
