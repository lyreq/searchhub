package cron

import (
	"github.com/robfig/cron/v3"
)

type CronJob struct {
	Name     string
	Schedule string
	Func     func()
}

type CronService struct {
	cron *cron.Cron
}

func NewService() *CronService {
	return &CronService{
		cron: cron.New(
			cron.WithParser(cron.NewParser(
				cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
			)),
		),
	}
}

func (cs *CronService) RegisterJob(job CronJob) error {
	_, err := cs.cron.AddFunc(job.Schedule, job.Func)

	if err != nil {
		return err
	}

	return nil
}

func (cs *CronService) Start() {
	cs.cron.Start()
}

func (cs *CronService) Stop() {
	cs.cron.Stop()
}
