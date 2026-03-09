package scheduler

import "github.com/robfig/cron/v3"

type schedulerRunner interface {
	AddFunc(spec string, cmd func()) (cron.EntryID, error)
	Start()
}
