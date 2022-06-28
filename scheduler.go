package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

type TaskScheduler struct {
	Config         TaskIntervalConfig
	ExecuteTrigger chan time.Time
	Scheduler      *gocron.Scheduler
}

func NewTaskScheduler(config TaskIntervalConfig) *TaskScheduler {
	return &TaskScheduler{
		Config:         config,
		ExecuteTrigger: make(chan time.Time),
		Scheduler:      gocron.NewScheduler(time.Local),
	}
}

func (s *TaskScheduler) Initialize(label string) error {
	switch s.Config.Type {
	case "cron":
		s.Scheduler.Cron(s.Config.Value).Do(s.Execute)
	case "interval":
		s.Scheduler.Every(s.Config.Value).Do(s.Execute)
	default:
		return fmt.Errorf("unknown interval type for '%s': %s", label, s.Config.Type)
	}

	return nil
}

func (s *TaskScheduler) Start() {
	s.Scheduler.StartAsync()
}

func (s *TaskScheduler) Execute() {
	s.ExecuteTrigger <- time.Now()
}

func (s *TaskScheduler) Stop() {
	s.Scheduler.Stop()
}

func (s *TaskScheduler) Close() {
	close(s.ExecuteTrigger)
}
