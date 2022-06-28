package main

import (
	"fmt"
	"time"
)

type Task struct {
	Scheduler *TaskScheduler
	Runner    TaskRunner
	Logger    *TaskLogger
	Config    TaskConfig
}

func NewTask(config TaskConfig) *Task {
	return &Task{
		Scheduler: nil,
		Runner:    nil,
		Logger:    NewLogger(config.Logs),
		Config:    config,
	}
}

func (t *Task) Initialize() error {
	t.Scheduler = NewTaskScheduler(t.Config.Interval)

	newRunner, ok := taskRunnerMap[t.Config.Execute.Type]

	if !ok {
		return fmt.Errorf("unknown task runner type: %s", t.Config.Execute.Type)
	}

	t.Runner = newRunner(t.Config.Execute)

	if err := t.Scheduler.Initialize(t.Config.Label); err != nil {
		return err
	}

	if err := t.Logger.Start(); err != nil {
		return err
	}

	return t.Runner.Initialize(t.Logger)
}

func (t *Task) Start() {
	t.Scheduler.Start()

	go func() {
		for now := range t.Scheduler.ExecuteTrigger {
			t.Execute(now)
		}
	}()
}

func (t *Task) Execute(now time.Time) {
	start := time.Now()

	go func() {
		logFileMutex.Lock()
		logFile.WriteString(fmt.Sprintf("[%s] Running task '%s'\n", start.Format("2006-01-02 15:04:05"), t.Config.Label))
		logFileMutex.Unlock()
	}()

	t.Runner.Execute(now)

	end := time.Now()

	go func() {
		logFileMutex.Lock()
		logFile.WriteString(fmt.Sprintf("[%s] Finished task '%s' (%s)\n", start.Format("2006-01-02 15:04:05"), t.Config.Label, end.Sub(start)))
		logFileMutex.Unlock()
	}()
}

func (t *Task) Stop() error {
	t.Scheduler.Stop()

	return t.Runner.Stop()
}

func (t *Task) Close() error {
	t.Scheduler.Close()

	return t.Logger.Close()
}
