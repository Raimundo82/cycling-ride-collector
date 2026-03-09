package scheduler

import "github.com/robfig/cron/v3"

type Scheduler struct {
	s schedulerRunner
}

var _ schedulerRunner = (*cron.Cron)(nil)

func New() *Scheduler {
	return &Scheduler{s: cron.New()}
}

func (s *Scheduler) StartWeeklySunday20(job func()) error {
	_, err := s.s.AddFunc("0 20 * * 0", job)
	if err != nil {
		return err
	}
	s.s.Start()
	return nil
}
