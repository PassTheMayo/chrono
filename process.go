package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

func init() {
	taskRunnerMap["process"] = NewProcessRunner
}

type ProcessRunner struct {
	IsRunning bool
	Config    TaskRunnerConfig
	Logger    *TaskLogger
}

func NewProcessRunner(config TaskRunnerConfig) TaskRunner {
	return &ProcessRunner{
		IsRunning: false,
		Config:    config,
		Logger:    nil,
	}
}

func (r *ProcessRunner) Initialize(logger *TaskLogger) error {
	r.Logger = logger

	return nil
}

func (r *ProcessRunner) Start() error {
	return nil
}

func (r *ProcessRunner) Execute(t time.Time) {
	cmd := exec.Command(r.Config.Process, r.Config.Args...)

	for key, value := range r.Config.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	cmd.Stdout = r.Logger.Stdout
	cmd.Stderr = r.Logger.Stderr

	if err := cmd.Start(); err != nil {
		log.Println(err)

		return
	}

	r.IsRunning = true

	if err := cmd.Wait(); err != nil {
		log.Println(err)

		return
	}

	r.IsRunning = false
}

func (r *ProcessRunner) Stop() error {
	return nil
}

var _ TaskRunner = &ProcessRunner{}
