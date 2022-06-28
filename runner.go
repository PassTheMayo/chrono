package main

import "time"

type NewTaskRunnerFunc func(TaskRunnerConfig) TaskRunner

var (
	taskRunnerMap map[string]NewTaskRunnerFunc = make(map[string]NewTaskRunnerFunc)
)

type TaskRunner interface {
	Initialize(*TaskLogger) error
	Start() error
	Execute(time.Time)
	Stop() error
}
